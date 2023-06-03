package helper

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gofiber/fiber/v2"
)

type MediaInterface interface {
	Save(c *fiber.Ctx) ([]string, error)
	UploadToAWSS3(fileName string) error
	GetObjectFromS3(fileName string) (*s3.GetObjectOutput, error)
}

type Media struct{}

func NewMediaHelper() MediaInterface {
	return &Media{}
}

func AWSConfig() (string, *session.Session) {
	bucket := os.Getenv("AWS_BUCKET_NAME")

	session := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_BUCKET_REGION")),
		Credentials: credentials.NewStaticCredentials(
			os.Getenv("AWS_ACCESS_KEY"),
			os.Getenv("AWS_SECRET_KEY"),
			""),
	}))

	return bucket, session
}

func (m *Media) FileName(name string) string {
	postfix := strconv.FormatInt(time.Now().Unix(), 10)
	fileExt := strings.Split(name, ".")[1]
	fileName := strings.Split(name, ".")[0]
	newFileName := strings.ReplaceAll(fileName, " ", "-")

	return newFileName + postfix + "." + fileExt
}

func (m *Media) Save(c *fiber.Ctx) ([]string, error) {
	form, formErr := c.MultipartForm()

	if formErr != nil {
		return []string{}, formErr
	}

	// get image files from form
	files := form.File["image"]

	if len(files) == 0 {
		return []string{}, fmt.Errorf("image is required")
	}

	fileNames := make([]string, len(files))

	for i, file := range files {
		fileNames[i] = m.FileName(file.Filename)
	}

	for i, file := range files {
		if err := c.SaveFile(file, fmt.Sprintf("./images/%s", fileNames[i])); err != nil {
			return []string{}, err
		}
	}

	return fileNames, nil
}

func (m *Media) UploadToAWSS3(filename string) error {
	timeout := time.Duration(5 * time.Second)
	key := os.Getenv("AWS_BUCKET_BASE_FOLDER") + "/" + m.FileName(filename)
	bucket, session := AWSConfig()

	svc := s3.New(session)

	ctx := context.Background()
	ctx, cancelFn := context.WithTimeout(ctx, timeout)

	defer cancelFn()

	file, openErr := os.Open(fmt.Sprintf("./images/%s", filename))

	if openErr != nil {
		log.Fatal(openErr)
		return openErr
	}

	// Uploads the object to S3. The Context will interrupt the request if the
	// timeout expires.
	_, err := svc.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   file,
	})
	if err != nil {
		return err
	}

	return nil
}

func (m *Media) GetObjectFromS3(key string) (*s3.GetObjectOutput, error) {
	bucket, session := AWSConfig()

	svc := s3.New(session)

	var objectKey string = fmt.Sprintf(os.Getenv("AWS_BUCKET_BASE_FOLDER") + "/" + key)

	// Downloads the object to a file
	obj, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(objectKey),
	})

	if err != nil {
		return nil, err
	}

	return obj, nil
}
