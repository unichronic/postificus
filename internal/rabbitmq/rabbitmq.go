package rabbitmq

import (
	"context"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// -------------------------------------------------------
// CONNECTION
// -------------------------------------------------------

// Connection wraps the AMQP connection and channel
type Connection struct {
	conn *amqp.Connection
	uri  string
}

// NewConnection establishes a connection to RabbitMQ
func NewConnection(uri string) (*Connection, error) {
	conn, err := amqp.Dial(uri)
	if err != nil {
		return nil, err
	}

	return &Connection{
		conn: conn,
		uri:  uri,
	}, nil
}

// Close closes the connection
func (c *Connection) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}

// CreateChannel creates a new channel from the connection
func (c *Connection) CreateChannel() (*amqp.Channel, error) {
	return c.conn.Channel()
}

// WaitUntilReady blocks until RabbitMQ is reachable (for startup)
func WaitUntilReady(uri string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		conn, err := amqp.Dial(uri)
		if err == nil {
			conn.Close()
			return nil
		}
		log.Println("Waiting for RabbitMQ...", err)
		time.Sleep(2 * time.Second)
	}
	return amqp.ErrClosed
}

// -------------------------------------------------------
// PRODUCER
// -------------------------------------------------------

// Producer handles message publishing
type Producer struct {
	conn *Connection
}

// NewProducer creates a new producer attached to a connection
func NewProducer(conn *Connection) *Producer {
	return &Producer{conn: conn}
}

// Publish sends specific payload to a queue
func (p *Producer) Publish(queueName string, payload []byte) error {
	ch, err := p.conn.CreateChannel()
	if err != nil {
		return fmt.Errorf("failed to create channel: %w", err)
	}
	defer ch.Close()

	// Ensure queue exists before publishing (safety)
	_, err = ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		amqp.Table{
			"x-dead-letter-exchange":    "dlx",
			"x-dead-letter-routing-key": queueName,
		}, // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return ch.PublishWithContext(ctx,
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         payload,
			DeliveryMode: amqp.Persistent, // Persist messages
			Timestamp:    time.Now(),
		},
	)
}

// -------------------------------------------------------
// CONSUMER
// -------------------------------------------------------

// Consumer handles message consumption
type Consumer struct {
	conn *Connection
}

// NewConsumer creates a new consumer attached to a connection
func NewConsumer(conn *Connection) *Consumer {
	return &Consumer{conn: conn}
}

// Consume starts a worker for a specific queue
func (c *Consumer) Consume(queueName string, handler func([]byte) error) error {
	ch, err := c.conn.CreateChannel()
	if err != nil {
		return err
	}
	// Note: We do NOT close the channel here, passing it to the goroutine

	// Ensure DLX exists
	err = ch.ExchangeDeclare(
		"dlx",    // name
		"direct", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		return err
	}

	// Ensure DLQ exists
	_, err = ch.QueueDeclare(
		queueName+":dlq", // name
		true,             // durable
		false,            // delete when unused
		false,            // exclusive
		false,            // no-wait
		nil,              // arguments
	)
	if err != nil {
		return err
	}

	// Bind DLQ to DLX
	err = ch.QueueBind(
		queueName+":dlq", // queue name
		queueName,        // routing key (same as original queue for simplicity)
		"dlx",            // exchange
		false,
		nil,
	)
	if err != nil {
		return err
	}

	// Ensure Main Queue exists with DLQ arguments
	args := amqp.Table{
		"x-dead-letter-exchange":    "dlx",
		"x-dead-letter-routing-key": queueName,
	}

	_, err = ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		args,      // arguments
	)
	if err != nil {
		ch.Close()
		return err
	}

	// Set QoS (prefetch count) - process 1 message at a time per worker (fair dispatch)
	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		ch.Close()
		return err
	}

	msgs, err := ch.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack (WE WILL MANUALLY ACK)
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		ch.Close()
		return err
	}

	go func() {
		log.Printf("🐰 Consumer started for queue: %s", queueName)
		for d := range msgs {
			log.Printf("Received a message on %s", queueName)

			if err := handler(d.Body); err != nil {
				log.Printf("❌ Error processing message: %v", err)

				// RETRY LOGIC
				retryCount := 0
				if val, ok := d.Headers["x-retry-count"]; ok {
					if i, ok := val.(int32); ok {
						retryCount = int(i)
					}
				}

				if retryCount < 2 {
					log.Printf("🔄 Retrying message (Attempt %d/2)...", retryCount+1)
					// Verify Channel is open
					if ch.IsClosed() {
						log.Println("Channel closed, cannot retry")
						continue
					}

					// Publish Copy with Header
					err := ch.Publish(
						"",        // exchange
						queueName, // routing key
						false,     // mandatory
						false,     // immediate
						amqp.Publishing{
							ContentType:  d.ContentType,
							Body:         d.Body,
							Headers:      amqp.Table{"x-retry-count": retryCount + 1},
							DeliveryMode: d.DeliveryMode,
						},
					)
					if err != nil {
						log.Printf("❌ Failed to republish retry: %v", err)
						d.Nack(false, true) // Force requeue if republish failed
					} else {
						d.Ack(false) // Ack original, new one is in queue
					}
				} else {
					log.Printf("💀 Max retries reached. Moving to DLQ.")
					d.Nack(false, false) // Requeue=false -> DLQ
				}

			} else {
				// Ack
				d.Ack(false)
				log.Printf("✅ Message processed on %s", queueName)
			}
		}
		log.Printf("Consumer stopped for queue: %s", queueName)
	}()

	return nil
}
