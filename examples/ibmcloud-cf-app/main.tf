resource "null_resource" "prepare_app_zip" {
  triggers = {
    app_version = "${var.app_version}"
  }

  provisioner "local-exec" {
    command = "mkdir -p  ${var.dir_to_clone}; cd ${var.dir_to_clone}; git init; git remote add origin ${var.git_repo}; git fetch; git checkout -t origin/master; zip -r ${var.app_zip} *"
  }
}

data "ibmcloud_cf_space" "space" {
  org   = "${var.org}"
  space = "${var.space}"
}

data "ibmcloud_cf_shared_domain" "domain" {
  name = "mybluemix.net"
}

resource "ibmcloud_cf_route" "route" {
  domain_guid = "${data.ibmcloud_cf_shared_domain.domain.id}"
  space_guid  = "${data.ibmcloud_cf_space.space.id}"
  host        = "${var.route}"
}

resource "ibmcloud_cf_service_instance" "service" {
  name       = "${var.service_instance_name}"
  space_guid = "${data.ibmcloud_cf_space.space.id}"
  service    = "${var.service_offering}"
  plan       = "${var.plan}"
  tags       = ["my-service"]
}

resource "ibmcloud_cf_service_key" "key" {
    name = "%s"
    service_instance_guid = "${ibmcloud_cf_service_instance.service.id}"
}

resource "ibmcloud_cf_app" "app" {
  depends_on = ["ibmcloud_cf_service_key.key", "null_resource.prepare_app_zip"]
  name              = "${var.app_name}"
  space_guid        = "${data.ibmcloud_cf_space.space.id}"
  app_path          = "${var.app_zip}"
  wait_time_minutes = 2

  buildpack  = "${var.buildpack}"
  disk_quota = 512

  command               = "${var.command}"
  memory                = 128
  instances             = 1
  disk_quota            = 512
  route_guid            = ["${ibmcloud_cf_route.route.id}"]
  service_instance_guid = ["${ibmcloud_cf_service_instance.service.id}"]

  environment_json = {
    "somejson" = "somevalue"
  }

  app_version = "${var.app_version}"
}
