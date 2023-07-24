package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/packer"
)

type buildConfig struct {
	subnetOCID string
	version    string
}

type hcpConfig struct {
	clientID     string
	clientSecret string
}

func packerBuild(t *testing.T, hcpConf hcpConfig, providerConf ociProviderConfig, buildConf buildConfig) (string, error) {
	packerOptions := &packer.Options{
		Template: "../packer.pkr.hcl",

		Vars: map[string]string{
			"hcp_packer_bucket_name": hcpBucketName,
			"oci_fingerprint":        providerConf.fingerprint,
			"oci_key_file":           providerConf.privateKeyFilePath,
			"root_tenancy_ocid":      providerConf.rootTenancyOCID,
			"subnet_ocid":            buildConf.subnetOCID,
			"terraform_tenancy_ocid": providerConf.terraformTenancyOCID,
			"version":                buildConf.version,
		},

		Env: map[string]string{
			"HCP_CLIENT_ID":     hcpConf.clientID,
			"HCP_CLIENT_SECRET": hcpConf.clientSecret,
		},
	}

	imageID, err := packer.BuildArtifactE(t, packerOptions)
	return imageID, err
}
