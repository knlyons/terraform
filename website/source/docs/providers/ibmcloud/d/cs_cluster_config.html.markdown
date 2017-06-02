---
layout: "ibmcloud"
page_title: "IBM Cloud: ibmcloud_cs_cluster_config"
sidebar_current: "docs-ibmcloud-datasource-cs-cluster-config"
description: |-
  Get the cluster configuration for Kubernetes on IBM Bluemix.
---

# ibmcloud\_cs_cluster_config


Download a configuration for Kubernetes clusters on IBM Bluemix.


## Example Usage

```hcl
data "ibmcloud_cs_cluster_config" "cluster_foo" {
  org_guid     = "test"
  space_guid   = "test_space"
  account_guid = "test_acc"
  name         = "FOO"
  config_dir   = "/home/foo_config"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_name_id` - (Required) Name or ID of the cluster.
* `config_dir` - (Required) The directory where you want the cluster configuration to download.
* `admin` - (Optional) Set it to true to download config for the admin. Default false.
* `download` - (Optional) Set it to false to skip downloading the config for the admin. Default true. Since it is part of a data source it means by default the config is downloaded for every `terraform` call. For a particular cluster name/id the config is guaranteed to be downloaded to the same path for a given `config_dir`
* `org_guid` - (Required) The GUID for the Bluemix organization that the cluster is associated with. The value can be retrieved from the `ibmcloud_cf_org` data source.
* `space_guid` - (Required) The GUID for the Bluemix space that the cluster is associated with. The value can be retrieved from the `ibmcloud_cf_space` data source.
* `account_guid` - (Required) The GUID for the Bluemix account that the cluster is associated with. The value can be retrieved from the `ibmcloud_cf_account` data source.


## Attributes Reference

The following attributes are exported:

* `id` - The unique identifier of the Cluster config 
* `config_file_path` -  The path to the cluster config file. Typically the Kubernetes yml config file.
