---
layout: "ibmcloud"
page_title: "IBM Cloud: ibmcloud_cf_space"
sidebar_current: "docs-ibmcloud-datasource-cf-space"
description: |-
  Get information about an IBM Bluemix space.
---

# ibmcloud\_cf_space

Import the details of an existing IBM Bluemix space as a read-only data source. The fields of the data source can then be referenced by other resources within the same configuration by using interpolation syntax. 

## Example Usage

```hcl
data "ibmcloud_cf_space" "spaceData" {
  space = "prod"
  org   = "someexample.com"
}
```

The following example shows how you can use the data source to reference the space ID in the `ibmcloud_cf_service_instance` resource.

```hcl
resource "ibmcloud_cf_service_instance" "service_instance" {
  name              = "test"
  space_guid        = "${data.ibmcloud_cf_space.spaceData.id}"
  service           = "cloudantNOSQLDB"
  plan              = "Lite"
  tags              = ["cluster-service", "cluster-bind"]
}

```

## Argument Reference

The following arguments are supported:

* `org` - (Required) The name of your Bluemix org. The value can be retrieved by running the `bx iam orgs` command in the [Bluemix CLI](https://console.ng.bluemix.net/docs/cli/reference/bluemix_cli/index.html#getting-started).
* `space` - (Required) The name of your space. The value can be retrieved by running the `bx iam spaces` command in the Bluemix CLI.

## Attributes Reference

The following attributes are exported:

* `id` - The unique identifier of the space.  
