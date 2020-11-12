package worker

import (
	"context"
	"github.com/google/uuid"
	minio "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/klog"
	"os"
	"strings"
	"sync"
)

type Uploader interface {
	Upload(string, context.Context) error
	GetUploadedFiles() sets.String
}

var defaultMinioUploader *MinioUploader
var defaultLock sync.Mutex

type MinioUploader struct {
	minioClient   *minio.Client
	uploadedFiles sets.String
	contentType   string
	region        string
	bucket        string
}

// NewMinioUploader use single-mode new minioUploader
func NewMinioUploader(endpoint, accessKeyID, secretAccessKey, bucket, contentType, region string, useSSL bool) (*MinioUploader, error) {
	defaultLock.Lock()
	defer defaultLock.Unlock()
	if defaultMinioUploader == nil {
		// Initialize minio client object.
		minioClient, err := minio.New(endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
			Secure: useSSL,
		})
		if err != nil {
			klog.Errorln(err)
			return nil, err
		}
		defaultMinioUploader = &MinioUploader{
			minioClient:   minioClient,
			uploadedFiles: sets.String{},
			contentType:   contentType,
			region:        region,
			bucket:        bucket + uuid.New().String(),
		}
	}
	return defaultMinioUploader, nil
}

// Upload use go-sdk update data
func (up *MinioUploader) Upload(filePath string, ctx context.Context) error {
	strs := strings.Split(filePath, "/")
	objectName := strs[len(strs)-1]
	bucketName := up.bucket
	err := up.minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: up.region})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := up.minioClient.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			klog.Warningf("We already own %s\n", bucketName)
		} else {
			klog.Errorln(err)
			return err
		}
	} else {
		klog.Infof("Successfully created %s\n", bucketName)
	}

	// Upload the zip file with FPutObject
	n, err := up.minioClient.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: up.contentType})
	if err != nil {
		klog.Infoln(err)
		return err
	}
	klog.Infof("Successfully uploaded file %s of size %d\n", filePath, n.Size)
	up.uploadedFiles.Insert(filePath)
	os.Remove(filePath)
	return nil
}

func (up *MinioUploader) GetUploadedFiles() sets.String {
	return up.uploadedFiles
}
