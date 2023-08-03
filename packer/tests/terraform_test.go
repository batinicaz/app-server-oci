package test

import (
	"fmt"
	"testing"
	"time"

	http_helper "github.com/gruntwork-io/terratest/modules/http-helper"
	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"golang.org/x/exp/maps"
)

var (
	ingressTypes = map[string]string{
		"80":   "http",
		"8080": "http",
	}
)

func runTerraform(t *testing.T, imageID string, providerConf ociProviderConfig) {
	terraformOptions := &terraform.Options{
		TerraformDir: "./terraform",

		Vars: map[string]interface{}{
			"image_id":               imageID,
			"ingress_ports":          maps.Keys(ingressTypes),
			"oci_fingerprint":        providerConf.fingerprint,
			"oci_private_key_path":   providerConf.privateKeyFilePath,
			"oci_region":             providerConf.region,
			"oci_tenancy_id":         providerConf.rootTenancyOCID,
			"oci_user_id":            providerConf.userOCID,
			"subnet_ocid":            subnetOCID,
			"terraform_tenancy_ocid": providerConf.terraformTenancyOCID,
		},
	}

	terraform.Init(t, terraformOptions)
	if !skipDestroy {
		defer terraform.Destroy(t, terraformOptions)
	}
	terraform.Apply(t, terraformOptions)

	publicIP := terraform.Output(t, terraformOptions, "public_ip")

	for port, protocol := range ingressTypes {
		validateHTTPConnectionWithRetry(t, protocol, publicIP, port)
	}
}

func validateHTTPConnectionWithRetry(t *testing.T, protocol, instanceIP string, port string) {
	maxRetries := 10
	sleepBetweenRetries := 10 * time.Second
	retry.DoWithRetry(t, "HTTP Request", maxRetries, sleepBetweenRetries, func() (string, error) {
		url := fmt.Sprintf("%s://%s:%s", protocol, instanceIP, port)
		statusCode, body, err := http_helper.HTTPDoE(t, "GET", url, nil, nil, nil)

		if err != nil {
			return "", err
		}

		if statusCode != 200 {
			return "", fmt.Errorf("Expected HTTP status 200, but got: %d", statusCode)
		}

		return body, err
	})
}
