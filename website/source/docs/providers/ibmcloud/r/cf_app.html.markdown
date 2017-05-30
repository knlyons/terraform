---
layout: "ibmcloud"
page_title: "IBM Cloud: cf_app"
sidebar_current: "docs-ibmcloud-resource-cf-app"
description: |-
  Manages IBM Cloud Cloud Foundry app.
---

# ibmcloud\_cf_app

Create, update, or delete Cloud Foundry application on IBM Bluemix.

## Example Usage

```hcl	
data "ibmcloud_cf_space" "space" {
  org   = "example.com"
  space = "dev"
}

resource "ibmcloud_cf_app" "app" {
  name         = "my-app"
  space_guid   = "${data.ibmcloud_cf_space.space.id}"
  app_path     = "hello.zip"
  wait_timeout = 90
  buildpack    = "sdk-for-nodejs"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, string) The name of the application.
* `memory` - (Optional, int) The amount of memory (in megabytes) each instance should have. If user doesn't specify then system assigns some pre-defined values based on the quota allocated to the application. You can check the default values by issuing `bluemix cf org <your-org>`. This will list the quotas defined on your org and space.
If space quotas are defined you can get them by `bluemix cf space-quota <space-quota-name>`. Otherwise you can check the organization quotas by `bluemix cf quota <quota-name>`
* `instances` - (Optional, int) The number of instances of the application.
* `disk_quota` - (Optional, int) The maximum amount of disk (in megabytes) available to an instance of an application. Default value is [1024 MB](http://bosh.io/jobs/cloud_controller_ng?source=github.com/cloudfoundry/cf-release&version=234#p=cc.default_app_disk_in_mb). Please check with your cloud provider if the value has been set differently.
* `space_guid` - (Required, string) Define space guid to which application belongs. The values can be retrieved from data source `ibmcloud_cf_space`.
* `buildpack` - (Optional, string) Buildpack to build the application. You can provide its values in the following ways
  * Blank value means autodetection
  * A Git URL pointing to a buildpack. Example - https://github.com/cloudfoundry/nodejs-buildpack.git
  * Name of an installed buildpack. Example - `go_buildpack`
* `environment_json` - (Optional, map) Key/value pairs of all the environment variables to run in your application. Does not include any system or service variables.
* `command` - (Optional, string) The initial command for the app.
* `route_guid` - (Optional, set) Define the route guids which should be bound to the application.Route should be in the same space as application.
* `service_instance_guid` - (Optional, set) Define the service instance guids that should be bound to this application.
* `wait_time_minutes` - (Optional, int) Define timeout to wait for the application to start.
* `app_path` - (Required, string) Define the path to the zip file of the application. The zip must contain all the application files directly within it and not inside some top-level folder. Typically you should go to the directory where your application files reside and issue `zip -r myapplication.zip *`.
* `app_version`	 - (Optional, string) Version of the application. If the application content in the file specified by _app_path_ changes then terraform can't detect that. So you should either change the application zip file name to let terraform know your zip content has changed or you can use this attribute to let the provider know that without changing the _app_path_

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the application.
