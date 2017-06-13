---
layout: "ibm"
page_title: "IBM: lb_vpx"
sidebar_current: "docs-ibm-resource-lb-vpx"
description: |-
  Manages IBM VPX Load Balancer.
---

# ibm\_lb_vpx

Provides a resource for VPX load balancers. This allows VPX load balancers to be created, updated, and deleted.

**NOTE**: IBM VPX load balancers consist of Citrix NetScaler VPX devices (virtual), which are currently priced on a per-month basis. The cost for an entire month is incurred immediately upon creation, so use caution when creating the resource. [See the network appliance docs](http://www.softlayer.com/network-appliances) for more information about pricing. Under the Citrix log, click **see more pricing** for a current price matrix.

You can also use the following REST URL to get a listing of VPX choices along with version numbers, speed, and plan type:

```
https://{{userName}}:{{apiKey}}@api.softlayer.com/rest/v3/SoftLayer_Product_Package/192/getItems.json?objectMask=id;capacity;description;units;keyName;prices.id;prices.categories.id;prices.categories.name
```

## Example Usage

[SLDN reference](http://sldn.softlayer.com/reference/datatypes/SoftLayer_Network_Application_Delivery_Controller)

```hcl
resource "ibm_lb_vpx" "test_vpx" {
    datacenter = "dal06"
    speed = 10
    version = "10.1"
    plan = "Standard"
    ip_count = 2
    public_vlan_id = 1251234
    private_vlan_id = 1540786
    public_subnet = "23.246.226.248/29"
    private_subnet = "10.107.180.0/26"
}
```

## Argument Reference

The following arguments are supported:

* `datacenter` - (Required, string) The data center that the VPX load balancer is to be provisioned in. Accepted values can be found [in the data center docs](http://www.softlayer.com/data-centers).
* `speed` - (Required, integer) The speed in Mbps. Accepted values are `10`, `200`, and `1000`.
* `version` - (Required, string) The VPX load balancer version. Accepted values are `10.1` and `10.5`.
* `plan` - (Required, string) The VPX load balancer plan. Accepted values are `Standard` and `Platinum`.
* `ip_count` - (Required, integer) The number of static public IP addresses assigned to the VPX load balancer. Accepted values are `2`, `4`, `8`, and `16`.
* `public_vlan_id` - (Optional, integer) Public VLAN ID that is used for the public network interface of the VPX load balancer. Accepted values can be found [in the VLAN docs](https://control.softlayer.com/network/vlans). Click the desired VLAN and note the ID in the resulting URL. Or, you can also [refer to a VLAN by name using a data source](../d/network_vlan.html).
* `private_vlan_id` - (Optional, integer) Private VLAN ID that is used for the private network interface of the VPX load balancer. Accepted values can be found [in the VLAN docs](https://control.softlayer.com/network/vlans). Click  the desired VLAN and note the ID in the resulting URL. Or, you can also [refer to a VLAN by name using a data source](../d/network_vlan.html).
* `public_subnet` - (Optional, string) Public subnet that is used for the public network interface of the VPX load balancer. Accepted values are primary public networks and can be found [in the subnet docs](https://control.softlayer.com/network/subnets).
* `private_subnet` - (Optional, string) Public subnet that is used for the private network interface of the VPX load balancer. Accepted values are primary private networks and can be found [in the subnet docs](https://control.softlayer.com/network/subnets).

## Attributes Reference

The following attributes are exported:

* `id` - The internal identifier of a VPX load balancer
* `name` - The internal name of a VPX load balancer.
* `vip_pool` - List of virtual IP addresses for the VPX load balancer.
