resource "oci_core_instance" "freshrss" {
  availability_domain = var.availability_domain
  compartment_id      = var.terraform_tenancy_ocid
  shape               = var.instance_shape

  create_vnic_details {
    subnet_id = var.subnet_ocid
  }

  defined_tags = {
    "terraform.managed" = "packer"
    "terraform.name"    = "Fresh RSS Build"
    "terraform.repo"    = "https://github.com/batinicaz/freshrss-oci"
  }

  instance_options {
    are_legacy_imds_endpoints_disabled = true
  }

  launch_options {
    network_type                        = "PARAVIRTUALIZED"
    is_pv_encryption_in_transit_enabled = true
  }

  source_details {
    source_id   = var.image_id
    source_type = "image"
  }
}