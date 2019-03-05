package send_lambda_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func generateEnvironmentName() string {
	return "test-" + random.UniqueId()
}

func getAccountID(t *testing.T) string {
	session, err := session.NewSession()
	if err != nil {
		t.Fatalf("unable to create aws session: %v", err)
	}
	svc := sts.New(session)
	out, err := svc.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		t.Fatalf("unable to get caller identity: %v", err)
	}
	return *out.Account
}

func loadFile(t *testing.T, name string) []byte {
	path := filepath.Join("testdata", name) // relative path
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return bytes
}

func TestSendLambdaDeploy(t *testing.T) {
	// setup the  test fixture
	accountID := getAccountID(t)
	environmentName := generateEnvironmentName()

	preSetupOptions := &terraform.Options{
		TerraformDir: "pre_setup",

		Vars: map[string]interface{}{
			"environment": environmentName,
			"aws_region":  "eu-west-2",
		},

		NoColor: true,
	}

	defer terraform.Destroy(t, preSetupOptions)
	terraform.InitAndApply(t, preSetupOptions)

	// deploy the lambda
	setupOptions := &terraform.Options{
		TerraformDir: "../../src/send_lambda",

		Vars: map[string]interface{}{
			"environment": environmentName,
			"aws_region":  "eu-west-2",
			"lambda_zip":  "../../../lambda/send/lambda.zip",
		},

		BackendConfig: map[string]interface{}{
			"bucket": "prm-" + accountID + "-terraform-states",
			"key":    environmentName + "/terratest/send_lambda/terraform.tfstate",
			"region": "eu-west-2",
		},

		NoColor: true,
	}

	defer terraform.Destroy(t, setupOptions)
	terraform.InitAndApply(t, setupOptions)

	// setup the rest of the test fixture
	postSetupOptions := &terraform.Options{
		TerraformDir: "post_setup",

		Vars: map[string]interface{}{
			"environment": environmentName,
			"aws_region":  "eu-west-2",
		},

		NoColor: true,
	}

	defer terraform.Destroy(t, postSetupOptions)
	terraform.InitAndApply(t, postSetupOptions)

	// post data at the published endpoint and verify that the response is 200
	url := terraform.Output(t, postSetupOptions, "invoke_endpoint") + "/send"

	t.Logf("Posting data to %s", url)
	body := loadFile(t, "send.xml")
	resp, err := http.Post(url, "application/xml", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("unable to post data to URL %s: %v", url, err)
	}

	assert.Equal(t, "200 OK", resp.Status)
}
