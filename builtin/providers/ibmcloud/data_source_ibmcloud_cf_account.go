package ibmcloud

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceIBMCloudCfAccount() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMCloudCfAccountRead,

		Schema: map[string]*schema.Schema{
			"org_guid": {
				Description: "The guid of the org",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func dataSourceIBMCloudCfAccountRead(d *schema.ResourceData, meta interface{}) error {
	bmxSess, err := meta.(ClientSession).BluemixSession()
	if err != nil {
		return err
	}
	accClient, err := meta.(ClientSession).BluemixAcccountAPI()
	if err != nil {
		return err
	}
	orgGUID := d.Get("org_guid").(string)
	account, err := accClient.Accounts().FindByOrg(orgGUID, bmxSess.Config.Region)
	if err != nil {
		return fmt.Errorf("Error retrieving organisation: %s", err)
	}
	d.SetId(account.GUID)
	return nil
}
