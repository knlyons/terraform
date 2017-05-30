package ibmcloud

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceIBMCloudCfServiceKey() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMCloudCfServiceKeyRead,

		Schema: map[string]*schema.Schema{
			"credentials": {
				Description: "Credentials asociated with the key",
				Type:        schema.TypeMap,
				Computed:    true,
			},

			"name": {
				Description: "The name of the service key",
				Type:        schema.TypeString,
				Required:    true,
			},
			"service_instance_name": {
				Description: "Service instance name for example, cleardbinstance",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func dataSourceIBMCloudCfServiceKeyRead(d *schema.ResourceData, meta interface{}) error {
	sr, err := meta.(ClientSession).CloudFoundryServiceInstanceClient()
	if err != nil {
		return nil
	}
	serviceInstanceName := d.Get("service_instance_name").(string)
	inst, err := sr.FindByName(serviceInstanceName)
	if err != nil {
		return err
	}

	serviceInstance, err := sr.Get(inst.GUID)
	if err != nil {
		return fmt.Errorf("Error retrieving service: %s", err)
	}

	name := d.Get("name").(string)
	srKey, _ := meta.(ClientSession).CloudFoundryServiceKeyClient()

	serviceKey, err := srKey.FindByName(serviceInstance.Metadata.GUID, name)
	if err != nil {
		return fmt.Errorf("Error retrieving service key: %s", err)
	}

	d.SetId(serviceKey.GUID)
	d.Set("credentials", serviceKey.Credentials)

	return nil
}
