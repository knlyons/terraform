package ibmcloud

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceIBMCloudCfServicePlan() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMCloudCfServicePlanRead,

		Schema: map[string]*schema.Schema{
			"service": {
				Description: "Service name for example, cleardb",
				Type:        schema.TypeString,
				Required:    true,
			},

			"plan": {
				Description: "The plan type ex- shared ",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func dataSourceIBMCloudCfServicePlanRead(d *schema.ResourceData, meta interface{}) error {

	service := d.Get("service").(string)
	srOff := meta.(ClientSession).CloudFoundryServiceOfferingClient()
	serviceOff, err := srOff.FindByLabel(service)
	if err != nil {
		return fmt.Errorf("Error retrieving service offering: %s", err)
	}

	srPlan := meta.(ClientSession).CloudFoundryServicePlanClient()
	plan := d.Get("plan").(string)

	servicePlan, err := srPlan.FindPlanInServiceOffering(serviceOff.GUID, plan)
	if err != nil {
		return fmt.Errorf("Error retrieving plan: %s", err)
	}

	d.SetId(servicePlan.GUID)

	return nil
}
