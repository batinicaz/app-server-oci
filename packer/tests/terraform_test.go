package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
)

func runTerraform(t *testing.T, imageID string, providerConf ociProviderConfig) {
	terraformOptions := &terraform.Options{
		TerraformDir: "./terraform",

		Vars: map[string]interface{}{
			"image_id":               imageID,
			"hcp_packer_bucket_name": hcpBucketName,
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
	terraform.ApplyAndIdempotent(t, terraformOptions)
}
