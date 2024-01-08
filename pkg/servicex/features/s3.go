package features

import (
	"io"
	"io/ioutil"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.startlite.cn/itapp/startlite/pkg/lines/errorx"
	"github.startlite.cn/itapp/startlite/pkg/lines/featurex"
	"github.startlite.cn/itapp/startlite/pkg/lines/typesx"
	"github.startlite.cn/itapp/startlite/pkg/lines/utilx/textx"
	"github.startlite.cn/itapp/startlite/pkg/servicex/types"
)

type S3 struct {
	types.S3Config

	*s3.S3
}

const pathStyle = "path"

func (impl *S3) Retrieve() (credentials.Value, error) {
	return credentials.Value{
		AccessKeyID:     impl.AccessKeyID,
		SecretAccessKey: impl.SecretAccessKey,
		SessionToken:    "",
		ProviderName:    "S3Config",
	}, nil
}

func (impl *S3) IsExpired() bool {
	return false
}

func MustNewS3Client(cl *featurex.ConfigLoader) *S3 {
	s3Client := S3{}
	cl.Load(&s3Client)

	sess := session.Must(session.NewSession())
	awsConfig := aws.Config{
		Credentials:      credentials.NewCredentials(&s3Client),
		Region:           aws.String(s3Client.Region),
		S3ForcePathStyle: typesx.BoolP(s3Client.Style == pathStyle),
	}
	if !textx.Blank(s3Client.Endpoint) {
		awsConfig.Endpoint = typesx.StringP(s3Client.Endpoint)
	}
	client := s3.New(sess, &awsConfig)

	s3Client.S3 = client

	return &s3Client
}

func (c *S3) Upload(fileKey string, body io.ReadSeeker) error {
	_, err := c.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(c.Bucket),
		Body:   body,
		Key:    &fileKey,
	})
	if err != nil {
		return errorx.Wrap(err, "S3 backend failed.")
	}

	return nil
}

func (c *S3) Download(fileKey string) (*s3.GetObjectOutput, error) {
	resp, err := c.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(c.Bucket),
		Key:    &fileKey,
	})
	if err != nil {
		return nil, errorx.Wrap(err, "S3 backend failed.")
	}

	return resp, nil
}

func (c *S3) DownloadBytes(fileKey string) ([]byte, error) {
	resp, err := c.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(c.Bucket),
		Key:    &fileKey,
	})
	if err != nil {
		return nil, errorx.Wrap(err, "S3 backend failed.")
	}
	defer resp.Body.Close()

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	return buf, nil
}
