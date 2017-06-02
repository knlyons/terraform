---
layout: "ibmcloud"
page_title: "IBM Cloud: ibmcloud_cf_account"
sidebar_current: "docs-ibmcloud-datasource-cf-account"
description: |-
  Get information about an IBM Bluemix account.
---

# ibmcloud\_cf_account

Import the details of an existing IBM Bluemix account as a read-only data source. The fields of the data source can then be referenced by other resources within the same configuration by using interpolation syntax. 

## Example Usage

```hcl
data "ibmcloud_cf_org" "orgData" {
  org = "example.com"
}

data "ibmcloud_cf_account" "accountData" {
  org_guid = "${data.ibmcloud_cf_org.orgData.id}"
}
```

## Argument Reference

The following arguments are supported:

* `org_guid` - (Required) The GUID of the Bluemix org. The value can be retrieved from the `ibmcloud_cf_org` data source, or by running the `bx iam orgs --guid` command in the [Bluemix CLI](https://console.ng.bluemix.net/docs/cli/reference/bluemix_cli/index.html#getting-started).

## Attributes Reference

The following attributes are exported:

* `id` - The unique identifier of the account.  
