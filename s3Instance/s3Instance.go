package s3Instance

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

type innitialS3 struct {
	awsAccessKeyID     string
	awsSecretAccessKey string
	awsRegion          string
	Session            *session.Session
	BucketName         string
}

var S3Creds *innitialS3 = &innitialS3{}

func (s3Object *innitialS3) startS3Session() error {

	S3Creds.awsAccessKeyID = os.Getenv("AWS_ACCESS_KEY_ID")
	S3Creds.awsSecretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
	S3Creds.awsRegion = os.Getenv("AWS_REGION")
	S3Creds.BucketName = os.Getenv("BUCKET_NAME")

	if S3Creds.awsAccessKeyID == "" || S3Creds.awsSecretAccessKey == "" || S3Creds.awsRegion == "" {
		return fmt.Errorf("as credenciais da AWS não estão configuradas corretamente")
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(S3Creds.awsRegion),
		Credentials: credentials.NewStaticCredentials(
			S3Creds.awsAccessKeyID,
			S3Creds.awsSecretAccessKey,
			"",
		),
	})
	if err != nil {
		return err
	}
	S3Creds.Session = sess

	return nil
}

func InitS3() {
	S3Creds.startS3Session()
}
