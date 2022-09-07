package test

import (
	"bytes"
	"encoding/json"
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

	terraformOptions := &terraform.Options{
		TerraformDir:  "../../examples/complete",
		BackendConfig: map[string]interface{}{},
		EnvVars:       map[string]string{},
		Vars:          map[string]interface{}{},
	}
	//defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApplyAndIdempotent(t, terraformOptions)

	functionUrl := terraform.Output(t, terraformOptions, "function_url")
	bucket := terraform.Output(t, terraformOptions, "bucket")
	region := terraform.Output(t, terraformOptions, "region")

	resp, err := sendEvent(functionUrl)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	assert.NoError(t, waitForS3Objects(bucket, region, 10, 90))
}

// sendEvent sends an HTTP POST with a test json body to the Lambda function url
func sendEvent(functionUrl string) (*http.Response, error) {
	fmt.Println("sending HTTP POST to: ", functionUrl)
	values := map[string]string{"test1": "value1"}
	jsonData, _ := json.Marshal(values)
	resp, err := http.Post(functionUrl, "application/json", bytes.NewBuffer(jsonData))
	fmt.Println(err)
	fmt.Println(resp)
	return resp, err
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
