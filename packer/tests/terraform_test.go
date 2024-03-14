package test

import (
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"golang.org/x/exp/maps"
)

const (
	robotsTXTPath       = "/robots.txt"
	robotsTag           = "X-Robots-Tag"
	expectedRobotsValue = "noindex, nofollow"
)

var (
	ingressTypes = map[string]string{
		"80":   "http",
		"3000": "http",
		"8080": "http",
		"8081": "http",
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
		validateHTTPConnectionWithRetry(t, protocol, publicIP, port+robotsTXTPath)
	}
}

func validateHTTPConnectionWithRetry(t *testing.T, protocol, instanceIP, port string) {
	maxRetries := 10
	sleepBetweenRetries := 10 * time.Second
	retry.DoWithRetry(t, "HTTP Request", maxRetries, sleepBetweenRetries, func() (string, error) {
		url := fmt.Sprintf("%s://%s:%s", protocol, instanceIP, port)
		req, err := http.NewRequest("GET", url, nil)

		if err != nil {
			return "", err
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return "", err
		}

		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return "", fmt.Errorf("Expected HTTP status 200, but got: %d", resp.StatusCode)
		}

		xRobotsTagValue := resp.Header.Get(robotsTag)
		if xRobotsTagValue == "" {
			t.Fatalf("Missing %s header", robotsTag)
		}

		if xRobotsTagValue != expectedRobotsValue {
			t.Fatalf("Expected %s to be '%s', but got: %s", robotsTag, expectedRobotsValue, xRobotsTagValue)
		}

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}

		return string(bodyBytes), nil
	})
}
