---
layout: "ibm"
page_title: "IBM : service_key"
sidebar_current: "docs-ibm-resource-service-key"
description: |-
  Manages IBM Service Key.
---

# ibm\_service_key

Create, update, or delete service keys for IBM Bluemix.

## Example Usage

```hcl
data "ibm_service_instance" "service_instance" {
  name = "mycloudant"
}

resource "ibm_service_key" "serviceKey" {
  name                  = "mycloudantkey"
  service_instance_guid = "${data.ibm_service_instance.service_instance.id}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, string) A descriptive name used to identify a service key.
* `parameters` - (Optional, map) Arbitrary parameters to pass along to the service broker. Must be a JSON object.
* `service_instance_guid` - (Required, string) The GUID of the service instance that the service key needs to be associated with.



## Attributes Reference

The following attributes are exported:

* `credentials` - Credentials associated with the key.
