package test

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/athena"
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
	table := terraform.Output(t, terraformOptions, "table")
	database := terraform.Output(t, terraformOptions, "database")

	// test sending an HTTP POST request and checking that the data arrived in successfully S3
	resp, err := sendEvent(functionUrl, "x-hub-signature-256", sharedSecret)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	assert.NoError(t, waitForS3Objects(bucket, region, 10, 30))

	// test sending an HTTP POST request with an invalid secret
	resp, err = sendEvent(functionUrl, "x-hub-signature-256", "not valid")
	assert.NoError(t, err)
	assert.Equal(t, 502, resp.StatusCode)

	// test running a query in AWS athena and checking that data is returned
	queryId := startAthenaQuery(region, table, database, bucket)
	err = waitForAthenaQuery(region, queryId, 10, 30)
	assert.NoError(t, err)
	rows, err := getAthenaResults(region, queryId)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(rows), 1)
}

// sendEvent sends an HTTP POST with a test json body to the Lambda function url
func sendEvent(functionUrl string, secretHeader string, sharedSecret string) (*http.Response, error) {
	fmt.Println("sending HTTP POST to: ", functionUrl)

	var jsonData = []byte(`{"action":"completed","number":"54"}`)
	req, err := http.NewRequest("POST", functionUrl, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set(secretHeader, "sha256="+GenerateSha256Hmac(string(jsonData), sharedSecret))

	client := &http.Client{
		Timeout: time.Second * 30,
	}
	response, err := client.Do(req)
	fmt.Println(response)
	fmt.Println(err)
	return response, err
}

// waitForS3Objects waits for any objects to be created in the given bucket
func waitForS3Objects(bucketName string, region string, delaySeconds, attempts int) error {
	for i := 0; ; i++ {
		output, err := listBucketObjects(bucketName, region)
		if err != nil {
			fmt.Println(err)
		} else {
			if len(output.Contents) > 0 {
				return nil
			}
		}

		if i >= (attempts - 1) {
			return fmt.Errorf("timed out while retrying")
		}

		fmt.Printf("Retrying in %d seconds...\n", delaySeconds)
		time.Sleep(time.Second * time.Duration(delaySeconds))
	}
}

// listBucketObjects lists all the objects in the given bucket
func listBucketObjects(bucketName string, region string) (*s3.ListObjectsV2Output, error) {
	fmt.Println("listing objects in bucket: ", bucketName)
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

// startAthenaQuery runs a query in AWS Athena and returns the query id
func startAthenaQuery(region, table, database, bucket string) string {
	fmt.Println("running Athena query for table: ", table)
	sess, _ := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	svc := athena.New(sess, aws.NewConfig().WithRegion(region))
	output, err := svc.StartQueryExecution(&athena.StartQueryExecutionInput{
		QueryExecutionContext: &athena.QueryExecutionContext{
			Database: aws.String(database),
		},
		QueryString: aws.String(fmt.Sprintf("SELECT * FROM \"%s\"", table)),
		ResultConfiguration: &athena.ResultConfiguration{
			OutputLocation: aws.String(fmt.Sprintf("s3://%s/", bucket)),
		},
	})
	fmt.Println(err)
	fmt.Println(output)
	return *output.QueryExecutionId
}

// getAthenaQueryState returns the status state of a AWS Athena query
func getAthenaQueryState(region, queryId string) (string, error) {
	fmt.Println("getting Athena query status for: ", queryId)
	sess, _ := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	svc := athena.New(sess, aws.NewConfig().WithRegion(region))
	output, err := svc.GetQueryExecution(&athena.GetQueryExecutionInput{
		QueryExecutionId: aws.String(queryId),
	})
	fmt.Println(err)
	fmt.Println(output)
	return *output.QueryExecution.Status.State, err
}

// waitForAthenaQuery waits for an AWS Athena query to complete
func waitForAthenaQuery(region, queryId string, delaySeconds, attempts int) error {
	for i := 0; ; i++ {
		state, err := getAthenaQueryState(region, queryId)
		if err != nil {
			fmt.Println(err)
		} else {
			if state == "SUCCEEDED" {
				return nil
			}
			if state != "RUNNING" && state != "QUEUED" {
				return fmt.Errorf("query failed")
			}
		}

		if i >= (attempts - 1) {
			return fmt.Errorf("timed out while retrying")
		}

		fmt.Printf("Retrying in %d seconds...\n", delaySeconds)
		time.Sleep(time.Second * time.Duration(delaySeconds))
	}
}

// getAthenaResults gets the result rows from an AWS Athena query
func getAthenaResults(region, queryId string) ([]*athena.Row, error) {
	fmt.Println("getting Athena query results for: ", queryId)
	sess, _ := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	svc := athena.New(sess, aws.NewConfig().WithRegion(region))
	output, err := svc.GetQueryResults(&athena.GetQueryResultsInput{
		QueryExecutionId: aws.String(queryId),
	})
	fmt.Println(err)
	fmt.Println(output)
	return output.ResultSet.Rows, err
}
