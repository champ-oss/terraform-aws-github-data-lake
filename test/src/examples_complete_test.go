package test

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/avast/retry-go"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

// TestExamplesComplete tests a typical deployment of this module
func TestExamplesComplete(t *testing.T) {
	t.Parallel()

	// Secret used to sign the payload when sending test events
	sharedSecret := "testing123"

	terraformOptions := &terraform.Options{
		TerraformDir:  "../../examples/complete",
		BackendConfig: map[string]interface{}{},
		EnvVars:       map[string]string{},
		Vars: map[string]interface{}{
			"shared_secret": sharedSecret,
		},
	}
	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApplyAndIdempotent(t, terraformOptions)

	functionUrl := terraform.Output(t, terraformOptions, "function_url")
	bucket := terraform.Output(t, terraformOptions, "bucket")
	region := terraform.Output(t, terraformOptions, "region")
	signatureHeaderKey := terraform.Output(t, terraformOptions, "signature_header_key")

	resp, err := sendEvent(functionUrl, signatureHeaderKey, sharedSecret)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	assert.NoError(t, waitForS3Objects(bucket, region, 10, 30))

	fmt.Println("sleeping before destroy")
	time.Sleep(5 * time.Minute)
}

// sendEvent sends an HTTP POST with a test json body to the Lambda function url
func sendEvent(functionUrl string, secretHeader string, sharedSecret string) (*http.Response, error) {
	fmt.Println("sending HTTP POST to: ", functionUrl)

	var jsonData = []byte(`{
		"test1": "value1",
		"test2": "value2"
	}`)
	req, err := http.NewRequest("POST", functionUrl, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set(secretHeader, GenerateSha256Hmac(string(jsonData), sharedSecret))

	client := &http.Client{
		Timeout: time.Second * 30,
	}
	response, err := client.Do(req)
	fmt.Println(response)
	fmt.Println(err)
	return response, err
}

// waitForS3Objects waits for any objects to be created in the given bucket
func waitForS3Objects(bucketName string, region string, delaySeconds, attempts uint) error {
	return retry.Do(func() error {
		output, err := listBucketObjects(bucketName, region)
		if err != nil {
			return err
		}
		if len(output.Contents) == 0 {
			return fmt.Errorf("bucket is empty")
		}
		return nil

	}, retry.Delay(time.Duration(delaySeconds)*time.Second), retry.Attempts(attempts))
}

// listBucketObjects lists all the objects in the given bucket
func listBucketObjects(bucketName string, region string) (*s3.ListObjectsV2Output, error) {
	fmt.Println("Listing objects in bucket: ", bucketName)
	sess, _ := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	svc := s3.New(sess, aws.NewConfig().WithRegion(region))
	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(bucketName)})
	fmt.Println("objects in bucket: ", len(resp.Contents))
	return resp, err
}

// GenerateSha256Hmac generates a HMAC digest of the given data string
func GenerateSha256Hmac(data, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}
