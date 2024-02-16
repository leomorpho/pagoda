package filestorage

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io"
	"log"
	"time"

	"github.com/mikestefanello/pagoda/config"
	"github.com/mikestefanello/pagoda/ent"
	"github.com/mikestefanello/pagoda/ent/filestorage"
	"github.com/mikestefanello/pagoda/pkg/domain"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Bucket int

const (
	BucketMainApp Bucket = iota
	BucketStaticFiles
)

var ErrBucketDoesNotExist = errors.New("requested bucket does not exist")

// StorageClientInterface defines the interface for the storage client.
type StorageClientInterface interface {
	CreateBucket(bucketName string, location string) error
	UploadFile(bucket Bucket, objectName string, fileStream io.Reader) (*int, error)
	DeleteFile(bucket Bucket, objectName string) error
	GetPresignedURL(bucket Bucket, objectName string, expiry time.Duration) (string, error)
	GetPhotoObjectFromFile(file *ent.FileStorage) (*domain.Photo, error)
	GetPhotoObjectsFromFiles(files []*ent.FileStorage) ([]domain.Photo, error)
}

type StorageClient struct {
	config      *config.Config
	orm         *ent.Client
	minioClient *minio.Client
}

func NewStorageClient(cfg *config.Config, orm *ent.Client) *StorageClient {
	minioClient, err := minio.New(cfg.Storage.S3Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.Storage.S3AccessKey, cfg.Storage.S3SecretKey, ""),
		Secure: cfg.Storage.S3UseSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}

	return &StorageClient{
		config:      cfg,
		orm:         orm,
		minioClient: minioClient,
	}
}

func (sc *StorageClient) getBucketName(b Bucket) (string, error) {
	switch b {
	case BucketMainApp:
		return sc.config.Storage.AppBucketName, nil
	case BucketStaticFiles:
		return sc.config.Storage.StaticFilesBucketName, nil
	default:
		return "", ErrBucketDoesNotExist
	}
}

func (sc *StorageClient) CreateBucket(bucketName string, location string) error {
	ctx := context.Background()

	bucketName = bucketName + string(sc.config.App.Environment)

	err := sc.minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		exists, errBucketExists := sc.minioClient.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Printf("Bucket %s already exists\n", bucketName)
		} else {
			return err
		}
	} else {
		log.Printf("Successfully created bucket %s\n", bucketName)
	}

	return nil
}

func (sc *StorageClient) UploadFile(bucket Bucket, objectName string, fileStream io.Reader) (*int, error) {
	ctx := context.Background()

	bucketName, err := sc.getBucketName(bucket)
	if err != nil {
		return nil, err
	}

	// Calculate file size and hash
	hash := md5.New()
	size, err := io.Copy(hash, fileStream)
	if err != nil {
		return nil, err
	}
	fileHash := hex.EncodeToString(hash.Sum(nil))

	// Seek to the beginning of the file stream if possible
	if seeker, ok := fileStream.(io.Seeker); ok {
		_, err = seeker.Seek(0, io.SeekStart)
		if err != nil {
			return nil, err
		}
	}

	// Upload file to S3-compatible storage
	_, err = sc.minioClient.PutObject(ctx, bucketName, objectName, fileStream, size, minio.PutObjectOptions{})
	if err != nil {
		return nil, err
	}

	// Create a new entry in the filestorage table
	filestorageEntry, err := sc.orm.FileStorage.
		Create().
		SetBucketName(bucketName).
		SetObjectKey(objectName).
		SetFileSize(size).
		SetFileHash(fileHash).
		SetCreatedAt(time.Now()).
		SetUpdatedAt(time.Now()).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return &filestorageEntry.ID, nil
}

func (sc *StorageClient) GetPresignedURL(bucket Bucket, objectName string, expiry time.Duration) (string, error) {
	ctx := context.Background()

	bucketName, err := sc.getBucketName(bucket)
	if err != nil {
		return "", err
	}

	presignedURL, err := sc.minioClient.PresignedGetObject(ctx, bucketName, objectName, expiry, nil)
	if err != nil {
		return "", err
	}

	return presignedURL.String(), nil
}

func (sc *StorageClient) DeleteFile(bucket Bucket, objectName string) error {
	ctx := context.Background()
	bucketName, err := sc.getBucketName(bucket)
	if err != nil {
		return err
	}
	err = sc.minioClient.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{ForceDelete: true})
	if err != nil {
		return err
	}

	_, err = sc.orm.FileStorage.Delete().Where(filestorage.ObjectKeyEQ(objectName)).Exec(ctx)

	if err != nil {
		return err
	}
	return nil
}

// TODO: GetPhotoObjectFromFile and GetPhotoObjectsFromFiles can be standardized
// to return a specific file object. Expiration can also be parametrized if necessary.
// getPhotoObjectFromFile generates a signed URL for a single file.
func (sc *StorageClient) GetPhotoObjectFromFile(file *ent.FileStorage) (*domain.Photo, error) {
	if file == nil {
		return nil, nil
	}

	// Generate a presigned URL with a specified duration
	url, err := sc.GetPresignedURL(BucketMainApp, file.ObjectKey, 2*24*time.Hour) // Adjust duration as needed
	if err != nil {
		return nil, err
	}

	return &domain.Photo{
		Key: file.ObjectKey,
		Url: url,
	}, nil
}

func (sc *StorageClient) GetPhotoObjectsFromFiles(
	files []*ent.FileStorage,
) ([]domain.Photo, error) {
	if len(files) == 0 {
		return nil, errors.New("no files given")
	}
	var photos []domain.Photo
	for _, f := range files {
		photo, err := sc.GetPhotoObjectFromFile(f)
		if err != nil {
			return nil, err
		}
		photos = append(photos, *photo)
	}
	return photos, nil
}
