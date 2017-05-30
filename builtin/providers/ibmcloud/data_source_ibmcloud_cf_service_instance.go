package ibmcloud

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceIBMCloudCfServiceInstance() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMCloudCfServiceInstanceRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Service instance name for example, cleardb",
				Type:        schema.TypeString,
				Required:    true,
			},

			"credentials": {
				Description: "Credentials asociated with the key",
				Type:        schema.TypeMap,
				Computed:    true,
			},

			"service_plan_guid": {
				Description: "The uniquie identifier of the service offering plan type",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceIBMCloudCfServiceInstanceRead(d *schema.ResourceData, meta interface{}) error {
	sr, err := meta.(ClientSession).CloudFoundryServiceInstanceClient()
	if err != nil {
		return nil
	}
	name := d.Get("name").(string)
	inst, err := sr.FindByName(name)
	if err != nil {
		return err
	}

	serviceInstance, err := sr.Get(inst.GUID)
	if err != nil {
		return fmt.Errorf("Error retrieving service: %s", err)
	}

	d.SetId(serviceInstance.Metadata.GUID)
	d.Set("credentials", serviceInstance.Entity.Credentials)
	d.Set("service_plan_guid", serviceInstance.Entity.ServicePlanGUID)

	return nil
}
