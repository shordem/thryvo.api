package config

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/shordem/api.thryvo/dto"
	"github.com/shordem/api.thryvo/lib/constants"
	"github.com/shordem/api.thryvo/lib/helper"
)

type FileConfigInterface interface {
	UploadFile(userId string, file *multipart.FileHeader) (string, error)
	GetObject(path string) (dto.GetFileDTO, error)
	DeleteObject(key string) error
	GetObjectPath(userId string, key string) string
}

type file struct {
	bucket  string
	service *s3.S3
}

func NewFileConfig(env constants.Env) FileConfigInterface {
	return &file{
		bucket:  env.AWS_BUCKET,
		service: s3.New(AWSConfig(env.AWS_REGION, env.AWS_ACCESS_KEY, env.AWS_SECRET_KEY)),
	}
}

func AWSConfig(region string, accessKey string, secretKey string) *session.Session {
	return session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
	}))
}

func (m *file) UploadFile(userId string, file *multipart.FileHeader) (string, error) {
	key := m.FileKey(file.Filename)
	path := m.GetObjectPath(userId, key)
	fileOpen, openErr := file.Open()

	if openErr != nil {
		return "", openErr
	}

	defer fileOpen.Close()

	fileContent := new(bytes.Buffer)

	_, copyErr := io.Copy(fileContent, fileOpen)

	if copyErr != nil {
		return "", copyErr
	}

	// Uploads the object to S3
	_, err := m.service.PutObject(&s3.PutObjectInput{
		Bucket: helper.StringToPointer(m.bucket),
		Key:    helper.StringToPointer(path),
		Body:   bytes.NewReader(fileContent.Bytes()),
	})

	if err != nil {
		return "", err
	}

	return key, nil
}

func (m *file) GetObject(path string) (dto.GetFileDTO, error) {
	var media dto.GetFileDTO

	// Downloads the object to a file
	obj, err := m.service.GetObject(&s3.GetObjectInput{
		Bucket: helper.StringToPointer(m.bucket),
		Key:    helper.StringToPointer(path),
	})

	if err != nil {
		return dto.GetFileDTO{}, err
	}

	media.Body = obj.Body
	media.ContentType = obj.ContentType
	media.ContentLength = obj.ContentLength

	return media, nil
}

func (m *file) DeleteObject(key string) error {
	_, err := m.service.DeleteObject(&s3.DeleteObjectInput{
		Bucket: helper.StringToPointer(m.bucket),
		Key:    helper.StringToPointer(key),
	})

	return err
}

func (m *file) FileKey(name string) string {
	filename, _ := helper.GenerateSnowflakeID()
	fileExt := strings.Split(name, ".")[len(strings.Split(name, "."))-1]

	return fmt.Sprintf("%d.%s", filename, fileExt)
}

func (m *file) GetObjectPath(userId string, key string) string {
	return fmt.Sprintf(userId + "/" + key)
}
