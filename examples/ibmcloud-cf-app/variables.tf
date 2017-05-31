variable "app_version" {
  default = "4"
}

variable "git_repo" {
  default = "https://github.com/ashishth09/nodejs-cloudantdb-crud-example"
}

variable "dir_to_clone" {
  default = "/tmp/my_cf_code"
}

variable "app_zip" {
  default = "/tmp/myzip.zip"
}

variable "org" {
  default = "ashishth@in.ibm.com"
}

variable "space" {
  default = "dev"
}

variable "route" {
  default = "my-app-route"
}

variable "service_instance_name" {
  default = "myservice"
}

variable "service_offering" {
 default = "cloudantNoSQLDB"
}

variable "plan" {
  default = "Lite"
}

variable "app_name" {
  default = "myapp"
}

variable "command" {
  default = ""
}

variable "buildpack" {
  default = ""
}
