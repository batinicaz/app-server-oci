variable "availability_domain" {
  type        = string
  default     = "xykJ:UK-LONDON-1-AD-3"
  description = "Availability domain where instance will be launched. Default is the always free domain."
}

variable "instance_shape" {
  type        = string
  default     = "VM.Standard.A1.Flex"
  description = "Instance type to use, default is the always free domain ARM option."
}

variable "oci_key_file" {
  description = "Path to the private key used for authenticating with OCI"
  type        = string
}

variable "oci_fingerprint" {
  description = "The fingerprint of the key used to authenticate with OCI"
  type        = string
}

variable "root_tenancy_ocid" {
  description = "The tenancy id used to lookup the source image"
  type        = string
}

variable "subnet_ocid" {
  description = "The OCID of the subnet to use for the build instance"
  type        = string
}

variable "terraform_tenancy_ocid" {
  description = "The tenancy id where to resources are to be created"
  type        = string
}

variable "source_image_username" {
  description = "The username of the user in the source image to connect as"
  type        = string
  default     = "ubuntu"
}

variable "ubuntu_version" {
  description = "The version of ubuntu to use as a source image"
  type        = string
  default     = "24.04"
}

variable "version" {
  description = "The version of the configuration being used to build the image. Should be a git ref"
  type        = string
}

packer {
  required_plugins {
    ansible = {
      version = "~> 1.0"
      source  = "github.com/hashicorp/ansible"
    }
    oracle = {
      version = "~> 1.0"
      source  = "github.com/hashicorp/oracle"
    }
  }
}

source "oracle-oci" "ubuntu" {
  availability_domain    = var.availability_domain
  compartment_ocid       = var.terraform_tenancy_ocid
  fingerprint            = var.oci_fingerprint
  image_compartment_ocid = var.terraform_tenancy_ocid
  image_name             = "app-server-${var.version}"
  key_file               = var.oci_key_file
  ssh_username           = var.source_image_username
  shape                  = var.instance_shape
  subnet_ocid            = var.subnet_ocid
  tenancy_ocid           = var.root_tenancy_ocid

  base_image_filter {
    compartment_id           = var.root_tenancy_ocid
    operating_system         = "Canonical Ubuntu"
    operating_system_version = var.ubuntu_version
  }

  create_vnic_details {
    display_name = "packer-build-app-server-${var.version}"
  }

  defined_tags_json = jsonencode({
    terraform = {
      managed = "packer"
      name    = "app-server-${var.version}"
      repo    = "https://github.com/batinicaz/app-server-oci"
    }
  })

  tags = {
    build-time     = "${timestamp()}"
    ubuntu-version = var.ubuntu_version
    version        = var.version
  }

  instance_defined_tags_json = jsonencode({
    terraform = {
      managed = "packer"
      name    = "App Server Build"
      repo    = "https://github.com/batinicaz/app-server-oci"
    }
  })

  shape_config {
    ocpus         = 1
    memory_in_gbs = 6
  }
}

build {
  sources = ["source.oracle-oci.ubuntu"]

  provisioner "ansible" {
    galaxy_file   = "../ansible/requirements.yml"
    playbook_file = "../ansible/playbook.yml"
    user          = var.source_image_username
    use_sftp      = true

    extra_arguments = [
      "--vault-password-file=.vault-password"
    ]
  }
}