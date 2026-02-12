package breaker

import (
	"errors"
	"postificus/internal/metrics"
	"sync"
	"time"
)

var (
	ErrCircuitOpen = errors.New("circuit breaker is open")
)

type State int

const (
	StateClosed State = iota
	StateOpen
	StateHalfOpen
)

type Settings struct {
	FailureThreshold int
	ResetTimeout     time.Duration
}

type CircuitBreaker struct {
	mu           sync.Mutex
	name         string
	state        State
	failures     int
	failureLimit int
	resetTimeout time.Duration
	lastFailure  time.Time
}

func NewCircuitBreaker(failureLimit int, resetTimeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		name:         "default", // TODO: Pass name
		state:        StateClosed,
		failureLimit: failureLimit,
		resetTimeout: resetTimeout,
	}
}

func NewCircuitBreakerWithName(name string, failureLimit int, resetTimeout time.Duration) *CircuitBreaker {
	metrics.CircuitBreakerState.WithLabelValues(name).Set(0) // Closed
	return &CircuitBreaker{
		name:         name,
		state:        StateClosed,
		failureLimit: failureLimit,
		resetTimeout: resetTimeout,
	}
}

func (cb *CircuitBreaker) Execute(fn func() error) error {
	if !cb.Allow() {
		return ErrCircuitOpen
	}

	err := fn()
	if err != nil {
		cb.RecordFailure()
		return err
	}

	cb.RecordSuccess()
	return nil
}

func (cb *CircuitBreaker) Allow() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if cb.state == StateOpen {
		if time.Since(cb.lastFailure) > cb.resetTimeout {
			cb.state = StateHalfOpen
			metrics.CircuitBreakerState.WithLabelValues(cb.name).Set(2) // HalfOpen
			return true
		}
		return false
	}
	return true
}

func (cb *CircuitBreaker) RecordFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.failures++
	cb.lastFailure = time.Now()

	if cb.state == StateClosed && cb.failures >= cb.failureLimit {
		cb.state = StateOpen
		metrics.CircuitBreakerState.WithLabelValues(cb.name).Set(1) // Open
	} else if cb.state == StateHalfOpen {
		cb.state = StateOpen                                        // Fail fast in half-open
		metrics.CircuitBreakerState.WithLabelValues(cb.name).Set(1) // Open
	}
}

func (cb *CircuitBreaker) RecordSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if cb.state == StateHalfOpen {
		cb.state = StateClosed
		cb.failures = 0
		metrics.CircuitBreakerState.WithLabelValues(cb.name).Set(0) // Closed
	} else if cb.state == StateClosed {
		cb.failures = 0
	}
}
