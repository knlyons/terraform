#service instance guid
output "guid" {
  value = "${ibmcloud_cf_service_instance.service-instance.id}"
}
