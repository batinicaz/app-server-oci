package test

import (
	"encoding/base64"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

var deleteImage bool
var fingerprint string
var hcpBucketName string
var hcpClientID string
var hcpClientSecret string
var imageID string
var privateKey string
var rootTenancyOCID string
var region string
var skipDestroy bool
var subnetOCID string
var terraformTenancyOCID string
var userOCID string
var version string

type ociProviderConfig struct {
	configDir            string
	configFilePath       string
	fingerprint          string
	privateKey           string
	privateKeyFilePath   string
	region               string
	rootTenancyOCID      string
	terraformTenancyOCID string
	userOCID             string
}

func TestMain(m *testing.M) {
	flag.BoolVar(&deleteImage, "deleteImage", false, "Delete the built image from OCI")
	flag.BoolVar(&skipDestroy, "skipDestroy", false, "Used for local dev, if set terraform stack will be left up after test run for debugging")
	flag.StringVar(&fingerprint, "fingerprint", "", "The fingerprint of the key used to authenticate with OCI")
	flag.StringVar(&hcpBucketName, "hcpBucketName", "dev", "The name of the HCP bucket in which to store build info")
	flag.StringVar(&hcpClientID, "hcpClientID", "", "The ID of the HCP Packer Registry where build metadata is to be stored")
	flag.StringVar(&hcpClientSecret, "hcpClientSecret", "", "The secret to authenticate with HCP Packer Registry where build metadata is stored")
	flag.StringVar(&imageID, "imageID", "", "Used for local dev, providing an existing image ID allows skipping the packer build and just running the terraform/infra tests")
	flag.StringVar(&privateKey, "privateKey", "", "The base64 encoded private key to authenticate with OCI")
	flag.StringVar(&region, "region", "", "The region in which to create resources")
	flag.StringVar(&rootTenancyOCID, "rootTenancyOCID", "", "The root tenancy id where images can be found")
	flag.StringVar(&subnetOCID, "subnetOCID", "", "The name of the subnet in which to launch the build instance")
	flag.StringVar(&terraformTenancyOCID, "terraformTenancyOCID", "", "The terraform tenancy id where resources will be created")
	flag.StringVar(&userOCID, "userOCID", "", "The ID of user that terraform will use to create the resources")
	flag.StringVar(&version, "version", "", "The git ref of the configuration used to build the image")
	flag.Parse()

	os.Exit(m.Run())
}

func TestAndBuildImage(t *testing.T) {
	buildConfig := &buildConfig{
		subnetOCID: subnetOCID,
		version:    version,
	}

	hcpConf := &hcpConfig{
		clientID:     hcpClientID,
		clientSecret: hcpClientSecret,
	}

	rawKey, err := base64.StdEncoding.DecodeString(privateKey)
	require.NoError(t, err, "Failed to base64 decode the provided private key")

	ociConfDir := filepath.Join(os.Getenv("HOME"), ".oci")
	providerConf := &ociProviderConfig{
		configDir:            ociConfDir,
		configFilePath:       filepath.Join(ociConfDir, "config"),
		fingerprint:          fingerprint,
		privateKey:           string(rawKey),
		privateKeyFilePath:   filepath.Join(ociConfDir, "oci_api_key.pem"),
		region:               region,
		rootTenancyOCID:      rootTenancyOCID,
		terraformTenancyOCID: terraformTenancyOCID,
		userOCID:             userOCID,
	}

	setupOCIConfig(t, *providerConf)
	t.Logf("OCI Config file created")

	if imageID == "" {
		var err error
		t.Logf("Running packer build")
		imageID, err = packerBuild(t, *hcpConf, *providerConf, *buildConfig)
		require.NoErrorf(t, err, "Packer build failed: %s", err)
	}

	t.Cleanup(func() {
		if !deleteImage {
			t.Logf("Image %s retained", imageID)
		} else {
			// TODO: Implement natively once OCI better supported in terratest
			cmd := exec.Command("oci", "compute", "image", "delete", "--force", "--image-id", imageID)
			_, err := cmd.Output()
			require.NoError(t, err, "Failed to delete oci image")
			t.Logf("Deleted image %s", imageID)
		}
	})

	t.Logf("Packer finished building: %s", imageID)
	runTerraform(t, imageID, *providerConf)

	t.Logf("Deployment testing complete")
}

// Terratest OCI module requires the config to be in place and it's also used by the packer builder
func setupOCIConfig(t *testing.T, conf ociProviderConfig) {
	err := os.MkdirAll(conf.configDir, 0700)
	require.NoError(t, err, "Failed to create directory to house OCI config")

	config := fmt.Sprintf(`[DEFAULT]
user=%s
fingerprint=%s
tenancy=%s
region=%s
key_file=%s`,
		conf.userOCID, conf.fingerprint, conf.rootTenancyOCID, conf.region, conf.privateKeyFilePath)

	err = os.WriteFile(conf.configFilePath, []byte(config), 0600)
	require.NoError(t, err, "Failed to create OCI config file")

	err = os.WriteFile(conf.privateKeyFilePath, []byte(conf.privateKey), 0400)
	require.NoError(t, err, "Failed to create OCI private key file")
}
