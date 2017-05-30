package ibmcloud

import (
	"fmt"

	v2 "github.com/IBM-Bluemix/bluemix-go/api/cf/cfv2"

	"github.com/IBM-Bluemix/bluemix-go/bmxerror"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceIBMCloudCfSharedDomain() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMCloudCfSharedDomainCreate,
		Read:     resourceIBMCloudCfSharedDomainRead,
		Delete:   resourceIBMCloudCfSharedDomainDelete,
		Exists:   resourceIBMCloudCfSharedDomainExists,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "The name of the domain",
				ValidateFunc: validateDomainName,
			},

			"router_group_guid": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The guid of the router group.",
			},
		},
	}
}

func resourceIBMCloudCfSharedDomainCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(ClientSession).CloudFoundrySharedDomainClient()
	if err != nil {
		return err
	}
	name := d.Get("name").(string)
	routerGroupGUID := d.Get("router_group_guid").(string)

	params := v2.SharedDomainRequest{
		Name:            name,
		RouterGroupGUID: routerGroupGUID,
	}

	shdomain, err := client.Create(params)
	if err != nil {
		return fmt.Errorf("Error creating shared domain: %s", err)
	}

	d.SetId(shdomain.Metadata.GUID)

	return resourceIBMCloudCfSharedDomainRead(d, meta)
}

func resourceIBMCloudCfSharedDomainRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(ClientSession).CloudFoundrySharedDomainClient()
	if err != nil {
		return err
	}
	shdomainGUID := d.Id()

	shdomain, err := client.Get(shdomainGUID)
	if err != nil {
		return fmt.Errorf("Error retrieving shared domain: %s", err)
	}
	d.Set("name", shdomain.Entity.Name)
	d.Set("router_group_guid", shdomain.Entity.RouterGroupGUID)

	return nil
}

func resourceIBMCloudCfSharedDomainDelete(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(ClientSession).CloudFoundrySharedDomainClient()
	if err != nil {
		return err
	}

	shdomainGUID := d.Id()

	err = client.Delete(shdomainGUID, true)
	if err != nil {
		return fmt.Errorf("Error deleting shared domain: %s", err)
	}

	d.SetId("")

	return nil
}

func resourceIBMCloudCfSharedDomainExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client, err := meta.(ClientSession).CloudFoundrySharedDomainClient()
	if err != nil {
		return false, err
	}
	shdomainGUID := d.Id()

	shdomain, err := client.Get(shdomainGUID)
	if err != nil {
		if apiErr, ok := err.(bmxerror.RequestFailure); ok {
			if apiErr.StatusCode() == 404 {
				return false, nil
			}
		}
		return false, fmt.Errorf("Error communicating with the API: %s", err)
	}

	return shdomain.Metadata.GUID == shdomainGUID, nil
}
