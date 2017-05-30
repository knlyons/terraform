package ibmcloud

import (
	"fmt"

	"github.com/IBM-Bluemix/bluemix-go/bmxerror"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceIBMCloudCfServiceKey() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMCloudCfServiceKeyCreate,
		Read:     resourceIBMCloudCfServiceKeyRead,
		Delete:   resourceIBMCloudCfServiceKeyDelete,
		Exists:   resourceIBMCloudCfServiceKeyExists,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: " The name of the service key ",
			},

			"service_instance_guid": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The guid of the service instance for which to create service key",
			},
			"parameters": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Description: "Arbitrary parameters to pass along to the service broker. Must be a JSON object",
			},
			"credentials": {
				Description: "Credentials asociated with the key",
				Type:        schema.TypeMap,
				Computed:    true,
			},
		},
	}
}

func resourceIBMCloudCfServiceKeyCreate(d *schema.ResourceData, meta interface{}) error {
	serviceClient, err := meta.(ClientSession).CloudFoundryServiceKeyClient()
	if err != nil {
		return err
	}
	name := d.Get("name").(string)
	serviceInstanceGUID := d.Get("service_instance_guid").(string)
	var parameters map[string]interface{}

	if parameters, ok := d.GetOk("parameters"); ok {
		parameters = parameters.(map[string]interface{})
	}

	serviceKey, err := serviceClient.Create(serviceInstanceGUID, name, parameters)
	if err != nil {
		return fmt.Errorf("Error creating service key: %s", err)
	}

	d.SetId(serviceKey.Metadata.GUID)

	return resourceIBMCloudCfServiceKeyRead(d, meta)
}

func resourceIBMCloudCfServiceKeyRead(d *schema.ResourceData, meta interface{}) error {
	serviceClient, err := meta.(ClientSession).CloudFoundryServiceKeyClient()
	if err != nil {
		return err
	}
	serviceKeyGUID := d.Id()

	serviceKey, err := serviceClient.Get(serviceKeyGUID)
	if err != nil {
		return fmt.Errorf("Error retrieving service key: %s", err)
	}
	d.Set("credentials", serviceKey.Entity.Credentials)

	return nil
}

func resourceIBMCloudCfServiceKeyDelete(d *schema.ResourceData, meta interface{}) error {
	serviceClient, err := meta.(ClientSession).CloudFoundryServiceKeyClient()
	if err != nil {
		return err
	}

	serviceKeyGUID := d.Id()

	err = serviceClient.Delete(serviceKeyGUID)
	if err != nil {
		return fmt.Errorf("Error deleting service key: %s", err)
	}

	d.SetId("")

	return nil
}

func resourceIBMCloudCfServiceKeyExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	serviceClient, err := meta.(ClientSession).CloudFoundryServiceKeyClient()
	if err != nil {
		return false, err
	}
	serviceKeyGUID := d.Id()

	serviceKey, err := serviceClient.Get(serviceKeyGUID)
	if err != nil {
		if apiErr, ok := err.(bmxerror.RequestFailure); ok {
			if apiErr.StatusCode() == 404 {
				return false, nil
			}
		}
		return false, fmt.Errorf("Error communicating with the API: %s", err)
	}

	return serviceKey.Metadata.GUID == serviceKeyGUID, nil
}
