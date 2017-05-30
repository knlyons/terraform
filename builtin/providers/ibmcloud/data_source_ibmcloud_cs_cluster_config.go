package ibmcloud

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mitchellh/go-homedir"
)

func dataSourceIBMCloudArmadaClusterConfig() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMCloudArmadaClusterConfigRead,

		Schema: map[string]*schema.Schema{

			"org_guid": {
				Description: "The bluemix organization guid this cluster belongs to",
				Type:        schema.TypeString,
				Required:    true,
			},
			"space_guid": {
				Description: "The bluemix space guid this cluster belongs to",
				Type:        schema.TypeString,
				Required:    true,
			},
			"account_guid": {
				Description: "The bluemix account guid this cluster belongs to",
				Type:        schema.TypeString,
				Required:    true,
			},
			"cluster_name_id": {
				Description: "The name/id of the cluster",
				Type:        schema.TypeString,
				Required:    true,
			},
			"config_dir": {
				Description: "The directory where the cluster config to be downloaded. Default is home directory ",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"config_file_path": {
				Description: "The path to the kubernetes yml file ",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceIBMCloudArmadaClusterConfigRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(ClientSession).ClusterClient()
	if err != nil {
		return err
	}

	name := d.Get("cluster_name_id").(string)

	targetEnv := getClusterTargetHeader(d)
	configDir := d.Get("config_dir").(string)
	if len(configDir) == 0 {
		configDir, err = homedir.Dir()
		if err != nil {
			return fmt.Errorf("Error fetching homedir: %s", err)
		}

	}

	configPath, err := client.GetClusterConfig(name, configDir, targetEnv)
	if err != nil {
		return fmt.Errorf("Error downloading the cluster config [%s]: %s", name, err)
	}

	d.SetId(name)
	d.Set("config_dir", configDir)
	d.Set("config_file_path", configPath)
	return nil
}
