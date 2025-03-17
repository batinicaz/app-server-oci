variable "availability_domain" {
  type        = string
  default     = "xykJ:UK-LONDON-1-AD-3"
  description = "Availability domain where instance will be launched. Default is the always free domain."
}

variable "image_id" {
  description = "The ID of the built image"
  type        = string
}

variable "ingress_ports" {
  description = "Ports to open"
  type        = set(string)
}

variable "instance_shape" {
  type        = string
  default     = "VM.Standard.A1.Flex"
  description = "Instance type to use, default is the always free domain ARM option."
}

variable "oci_fingerprint" {
  description = "The fingerprint of the key used to authenticate with OCI"
  type        = string
}

variable "oci_private_key_path" {
  description = "The path to the private key to authenticate with OCI"
  type        = string
}

variable "oci_region" {
  description = "The region in which to create resources"
  type        = string
}

variable "oci_tenancy_id" {
  description = "The tenancy id where to resources are to be created"
  type        = string
}

variable "oci_user_id" {
  description = "The ID of user that terraform will use to create the resources"
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

locals {
  defined_tags = {
    "terraform.managed" = "packer"
    "terraform.name"    = "App Server Build"
    "terraform.repo"    = "https://github.com/batinicaz/app-server-oci"
  }
}