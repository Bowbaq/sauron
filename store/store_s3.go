package store

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type s3Store struct {
	*jsonStore

	bucket, key string

	client *s3.S3
}

// NewS3 instantiates a new concrete Store backed by an amazon S3 object.
func NewS3(bucket, key string) Store {
	s := &s3Store{
		bucket: bucket,
		key:    key,

		client: s3.New(session.Must(session.NewSession())),
	}
	s.jsonStore = &jsonStore{
		read:  s.read,
		write: s.write,
	}

	return s
}

func (s3s *s3Store) read() (sauronState, error) {
	resp, err := s3s.client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s3s.bucket),
		Key:    aws.String(s3s.key),
	})
	if err != nil {
		if err.(awserr.Error).Code() == "NoSuchKey" {
			return make(sauronState), nil
		}
		return nil, fmt.Errorf("Failed to read state from S3: %v", err)
	}

	var s sauronState
	err = json.NewDecoder(resp.Body).Decode(&s)
	if err != nil {
		return nil, fmt.Errorf("Failed to decode state: %v", err)
	}

	return s, nil
}

func (s3s *s3Store) write(s sauronState) error {
	data, err := json.Marshal(s)
	if err != nil {
		return fmt.Errorf("Failed to encode state: %v", err)
	}

	_, err = s3s.client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(s3s.bucket),
		Key:    aws.String(s3s.key),
		Body:   bytes.NewReader(data),
	})
	if err != nil {
		return fmt.Errorf("Failed to write state to S3: %v", err)
	}

	return nil
}
