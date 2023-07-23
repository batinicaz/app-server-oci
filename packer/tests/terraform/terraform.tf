terraform {
  backend "local" {
    path = "terratest-run.tfstate"
  }
  required_providers {
    http = {
      source  = "hashicorp/http"
      version = "~> 3.0"
    }
    oci = {
      source  = "oracle/oci"
      version = "~> 5.0"
    }
  }

  required_version = "~> 1.0"
}

provider "oci" {
  fingerprint      = var.oci_fingerprint
  private_key_path = var.oci_private_key_path
  region           = var.oci_region
  tenancy_ocid     = var.oci_tenancy_id
  user_ocid        = var.oci_user_id
}