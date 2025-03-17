package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/packer"
)

type buildConfig struct {
	subnetOCID string
	version    string
}

func packerBuild(t *testing.T, providerConf ociProviderConfig, buildConf buildConfig) (string, error) {
	packerOptions := &packer.Options{
		Template: "../packer.pkr.hcl",

		Vars: map[string]string{
			"oci_fingerprint":        providerConf.fingerprint,
			"oci_key_file":           providerConf.privateKeyFilePath,
			"root_tenancy_ocid":      providerConf.rootTenancyOCID,
			"subnet_ocid":            buildConf.subnetOCID,
			"terraform_tenancy_ocid": providerConf.terraformTenancyOCID,
			"version":                buildConf.version,
		},
	}

	imageID, err := packer.BuildArtifactE(t, packerOptions)
	return imageID, err
}
