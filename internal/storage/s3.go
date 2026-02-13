package storage

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var S3Client *minio.Client
var S3Bucket string
var S3PublicURL string

func InitS3() error {
	endpoint := os.Getenv("S3_ENDPOINT")
	accessKey := os.Getenv("S3_ACCESS_KEY")
	secretKey := os.Getenv("S3_SECRET_KEY")
	S3Bucket = os.Getenv("S3_BUCKET")
	if S3Bucket == "" {
		S3Bucket = "postificus-uploads" // Default bucket
	}
	// For Supabase/AWS, we construct a public URL base
	// S3_PUBLIC_URL can be set explicitly, or we try to guess
	S3PublicURL = os.Getenv("S3_PUBLIC_URL")

	if endpoint == "" || accessKey == "" || secretKey == "" {
		return fmt.Errorf("S3 configuration missing")
	}

	// Clean endpoint scheme
	useSSL := true
	if strings.HasPrefix(endpoint, "http://") {
		endpoint = strings.TrimPrefix(endpoint, "http://")
		useSSL = false
	} else if strings.HasPrefix(endpoint, "https://") {
		endpoint = strings.TrimPrefix(endpoint, "https://")
	}

	var err error
	S3Client, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return fmt.Errorf("failed to init S3 client: %w", err)
	}

	// Ensure bucket exists
	ctx := context.Background()
	exists, errBucket := S3Client.BucketExists(ctx, S3Bucket)
	if errBucket != nil {
		log.Printf("⚠️ Check bucket exists error: %v", errBucket)
	} else if !exists {
		err = S3Client.MakeBucket(ctx, S3Bucket, minio.MakeBucketOptions{})
		if err != nil {
			log.Printf("⚠️ Failed to create bucket %s: %v", S3Bucket, err)
		} else {
			log.Printf("✅ Created bucket: %s", S3Bucket)
			// Set policy public? For now just bucket.
			// Supabase buckets are usually public by policy.
		}
	}

	log.Printf("✅ S3 Client initialized for endpoint: %s", endpoint)
	return nil
}

// UploadFile uploads a multipart file to S3 and returns the public URL
func UploadFile(ctx context.Context, file *multipart.FileHeader) (string, error) {
	if S3Client == nil {
		return "", fmt.Errorf("S3 client not initialized")
	}

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%d-%s%s", time.Now().UnixNano(), "upload", ext)
	contentType := file.Header.Get("Content-Type")

	// Upload
	info, err := S3Client.PutObject(ctx, S3Bucket, filename, src, file.Size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("S3 upload failed: %w", err)
	}

	// Construct public URL
	if S3PublicURL != "" {
		// Use explicit public URL base (e.g. Supabase Public URL)
		// Supabase format: https://[project].supabase.co/storage/v1/object/public/[bucket]/[filename]
		// If S3PublicURL is just the base, we append bucket/filename
		base := strings.TrimSuffix(S3PublicURL, "/")
		return fmt.Sprintf("%s/%s/%s", base, S3Bucket, filename), nil
	} else {
		// MinIO / Generic S3 (Presigned or direct)
		// For now, let's assume direct access if no custom public URL
		// NOTE: This might be internal Docker URL (minio:9000), which isn't reachable by frontend.
		// Safe default: Return the filename and let frontend decide? No, frontend expects URL.
		// Standard MinIO Browser Access: http://localhost:9000/bucket/filename

		// Fallback: If running locally, we might need a workaround.
		// But for "Hosted S3", we usually have a public CDN URL.
		return fmt.Sprintf("https://%s/%s/%s", S3Client.EndpointURL().Host, S3Bucket, info.Key), nil
	}
}
