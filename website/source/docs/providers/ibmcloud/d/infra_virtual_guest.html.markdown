---
layout: "ibmcloud"
page_title: "IBM Cloud: ibmcloud_virtual_guest"
sidebar_current: "docs-ibmcloud-datasource-infra-virtual_guest"
description: |-
  Get information on a IBM Cloud Infrastructure Virtual Guest resource
---

# ibmcloud\_infra_virtual_guest

Import the details of an existing virtual guest as a read-only data source. The fields of the data source can then be referenced by other resources within the same configuration using interpolation syntax.

## Example Usage

```hcl
data "ibmcloud_infra_virtual_guest" "virtual_guest" {
  hostname    = "jumpbox"
  domain      = "example.com"
  most_recent = true
}
```

## Argument Reference

The following arguments are supported:

* `hostname` - (Required) The hostname of the virtual guest.
* `domain` - (Required) The domain of the virtual guest.
* `most_recent` - (Optional) `True` or `False`. If `true` and multiple entries are found, the most recently created virtual guest is used. If `false`, an error is returned.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the virtual guest.
* `datacenter` - Data center in which the virtual guest is deployed.
* `cores` - Number of cpu cores.
* `status` - The VSI status.
* `last_known_power_state` - The last known power state of a virtual guest in the event the guest is turned off outside of IMS or has gone offline.
* `power_state` - The current power state of a virtual guest.
