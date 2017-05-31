package ibmcloud

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceIBMCloudCfOrg() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMCloudCfOrgRead,

		Schema: map[string]*schema.Schema{
			"org": {
				Description: "Org name, for example myorg@domain",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func dataSourceIBMCloudCfOrgRead(d *schema.ResourceData, meta interface{}) error {
	cfAPI, err := meta.(ClientSession).CFAPI()
	if err != nil {
		return err
	}
	orgAPI := cfAPI.Organizations()
	org := d.Get("org").(string)
	orgFields, err := orgAPI.FindByName(org)
	if err != nil {
		return fmt.Errorf("Error retrieving organisation: %s", err)
	}

	d.SetId(orgFields.GUID)

	return nil
}
