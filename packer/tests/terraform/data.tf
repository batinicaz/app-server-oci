data "http" "my_ip" {
  url = "https://ipinfo.io/ip"
}

data "oci_core_subnet" "selected" {
  subnet_id = var.subnet_ocid
}