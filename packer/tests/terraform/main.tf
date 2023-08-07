resource "oci_core_instance" "freshrss" {
  availability_domain                 = var.availability_domain
  compartment_id                      = var.terraform_tenancy_ocid
  defined_tags                        = local.defined_tags
  display_name                        = "packer-test-freshrss"
  is_pv_encryption_in_transit_enabled = true
  shape                               = var.instance_shape

  create_vnic_details {
    nsg_ids   = [oci_core_network_security_group.freshrss.id]
    subnet_id = var.subnet_ocid
  }

  instance_options {
    are_legacy_imds_endpoints_disabled = true
  }

  launch_options {
    network_type                        = "PARAVIRTUALIZED"
    is_pv_encryption_in_transit_enabled = true
  }

  shape_config {
    ocpus         = 1
    memory_in_gbs = 6
  }

  source_details {
    source_id   = var.image_id
    source_type = "image"
  }
}

resource "oci_core_network_security_group" "freshrss" {
  compartment_id = var.terraform_tenancy_ocid
  display_name   = "Fresh RSS Packer Build Test Security Group"
  defined_tags   = local.defined_tags
  vcn_id         = data.oci_core_subnet.selected.vcn_id
}

resource "oci_core_network_security_group_security_rule" "freshrss" {
  // checkov:skip=CKV_OCI_21: Intentionally using stateful rules
  for_each                  = var.ingress_ports
  network_security_group_id = oci_core_network_security_group.freshrss.id
  description               = "Allow ${each.key} from test runner"
  direction                 = "INGRESS"
  protocol                  = "6" // TCP
  source                    = "${data.http.my_ip.response_body}/32"

  tcp_options {
    destination_port_range {
      max = each.key
      min = each.key
    }
  }
}
