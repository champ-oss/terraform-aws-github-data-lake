package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gruntwork-io/terratest/modules/terraform"
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
	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApplyAndIdempotent(t, terraformOptions)

	functionUrl := terraform.Output(t, terraformOptions, "function_url")
	sendEvent(functionUrl)

	t.Log("Sleeping...")
	time.Sleep(5 * time.Minute)
}

func sendEvent(functionUrl string) {
	values := map[string]string{"test1": "value1"}
	jsonData, _ := json.Marshal(values)
	resp, err := http.Post(functionUrl, "application/json", bytes.NewBuffer(jsonData))
	fmt.Println(err)
	fmt.Println(resp)
}
