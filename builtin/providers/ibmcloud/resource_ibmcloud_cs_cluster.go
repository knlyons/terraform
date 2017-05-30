package ibmcloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	v1 "github.com/IBM-Bluemix/bluemix-go/api/k8scluster/k8sclusterv1"
	"github.com/IBM-Bluemix/bluemix-go/bmxerror"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

const (
	clusterNormal     = "normal"
	workerNormal      = "normal"
	workerReadyState  = "Ready"
	workerDeleteState = "deleted"

	clusterProvisioning = "provisioning"
	workerProvisioning  = "provisioning"
)

func resourceIBMCloudArmadaCluster() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMCloudArmadaClusterCreate,
		Read:     resourceIBMCloudArmadaClusterRead,
		Update:   resourceIBMCloudArmadaClusterUpdate,
		Delete:   resourceIBMCloudArmadaClusterDelete,
		Exists:   resourceIBMCloudArmadaClusterExists,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The cluster name",
			},
			"datacenter": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The datacenter where this cluster will be deployed",
			},
			"workers": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"action": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validateAllowedStringValue([]string{"add", "reboot", "reload"}),
						},
					},
				},
			},

			"machine_type": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"isolation": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},

			"billing": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Default:  "hourly",
			},

			"public_vlan_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  nil,
			},

			"private_vlan_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  nil,
			},
			"ingress_hostname": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ingress_secret": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"no_subnet": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},
			"server_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"worker_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"subnet_id": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"webhook": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"level": {
							Type:     schema.TypeString,
							Required: true,
						},
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateAllowedStringValue([]string{"slack"}),
						},
						"url": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"org_guid": {
				Description: "The bluemix organization guid this cluster belongs to",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"space_guid": {
				Description: "The bluemix space guid this cluster belongs to",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"account_guid": {
				Description: "The bluemix account guid this cluster belongs to",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"wait_time_minutes": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  90,
			},
		},
	}
}

func resourceIBMCloudArmadaClusterCreate(d *schema.ResourceData, meta interface{}) error {
	clusterClient, err := meta.(ClientSession).ClusterClient()
	if err != nil {
		return err
	}

	name := d.Get("name").(string)
	datacenter := d.Get("datacenter").(string)
	workers := d.Get("workers").([]interface{})
	billing := d.Get("billing").(string)
	machineType := d.Get("machine_type").(string)
	publicVlanID := d.Get("public_vlan_id").(string)
	privateVlanID := d.Get("private_vlan_id").(string)
	webhooks := d.Get("webhook").([]interface{})
	noSubnet := d.Get("no_subnet").(bool)
	isolation := d.Get("isolation").(string)

	params := &v1.ClusterCreateRequest{
		Name:        name,
		Datacenter:  datacenter,
		WorkerNum:   len(workers),
		Billing:     billing,
		MachineType: machineType,
		PublicVlan:  publicVlanID,
		PrivateVlan: privateVlanID,
		NoSubnet:    noSubnet,
		Isolation:   isolation,
	}

	targetEnv := getClusterTargetHeader(d)

	cls, err := clusterClient.Create(params, targetEnv)
	if err != nil {
		return err
	}
	d.SetId(cls.ID)
	//wait for cluster availability
	_, err = WaitForClusterAvailable(d, meta, targetEnv)
	//wait for worker  availability
	_, err = WaitForWorkerAvailable(d, meta, targetEnv)
	if err != nil {
		return fmt.Errorf(
			"Error waiting for workers of cluster (%s) to become ready: %s", d.Id(), err)
	}

	subnetClient, _ := meta.(ClientSession).ClusterSubnetClient()
	subnetIDs := d.Get("subnet_id").(*schema.Set)
	for _, subnetID := range subnetIDs.List() {
		if subnetID != "" {
			err = subnetClient.AddSubnet(cls.ID, subnetID.(string), targetEnv)
			if err != nil {
				return err
			}
		}
	}
	webhookClient, _ := meta.(ClientSession).ClusterWebHooksClient()
	for _, e := range webhooks {
		pack := e.(map[string]interface{})
		webhook := v1.WebHook{
			Level: pack["level"].(string),
			Type:  pack["type"].(string),
			URL:   pack["url"].(string),
		}

		webhookClient.Add(cls.ID, webhook, targetEnv)

	}

	workersInfo := []map[string]string{}
	workerClient, _ := meta.(ClientSession).ClusterWorkerClient()
	workerFields, err := workerClient.List(cls.ID, targetEnv)
	if err != nil {
		return err
	}
	//Create a map with worker name and id
	for i, e := range workers {
		pack := e.(map[string]interface{})
		var worker = map[string]string{
			"name":   pack["name"].(string),
			"id":     workerFields[i].ID,
			"action": pack["action"].(string),
		}
		workersInfo = append(workersInfo, worker)
	}
	d.Set("workers", workersInfo)

	if err != nil {
		return fmt.Errorf(
			"Error waiting for cluster (%s) to become ready: %s", d.Id(), err)
	}

	return resourceIBMCloudArmadaClusterRead(d, meta)
}

func resourceIBMCloudArmadaClusterRead(d *schema.ResourceData, meta interface{}) error {

	targetEnv := getClusterTargetHeader(d)

	client, err := meta.(ClientSession).ClusterClient()
	if err != nil {
		return err
	}

	clusterID := d.Id()
	cls, err := client.Find(clusterID, targetEnv)
	if err != nil {
		return fmt.Errorf("Error retrieving armada cluster: %s", err)
	}

	d.Set("name", cls.Name)
	d.Set("server_url", cls.ServerURL)
	d.Set("ingress_hostname", cls.IngressHostname)
	d.Set("ingress_secret", cls.IngressSecretName)
	d.Set("worker_num", cls.WorkerCount)
	d.Set("subnet_id", d.Get("subnet_id").(*schema.Set))
	return nil
}

func resourceIBMCloudArmadaClusterUpdate(d *schema.ResourceData, meta interface{}) error {

	targetEnv := getClusterTargetHeader(d)

	subnetClient, err := meta.(ClientSession).ClusterSubnetClient()
	if err != nil {
		return err
	}
	webhookClient, _ := meta.(ClientSession).ClusterWebHooksClient()
	workerClient, _ := meta.(ClientSession).ClusterWorkerClient()

	clusterID := d.Id()
	workersInfo := []map[string]string{}
	if d.HasChange("workers") {
		oldWorkers, newWorkers := d.GetChange("workers")
		oldWorker := oldWorkers.([]interface{})
		newWorker := newWorkers.([]interface{})
		for _, nW := range newWorker {
			newPack := nW.(map[string]interface{})
			exists := false
			for _, oW := range oldWorker {
				oldPack := oW.(map[string]interface{})
				if strings.Compare(newPack["name"].(string), oldPack["name"].(string)) == 0 {
					exists = true
					if strings.Compare(newPack["action"].(string), oldPack["action"].(string)) != 0 {
						params := v1.WorkerParam{
							Action: newPack["action"].(string),
						}
						workerClient.Update(clusterID, oldPack["id"].(string), params, targetEnv)
						var worker = map[string]string{
							"name":   newPack["name"].(string),
							"id":     newPack["id"].(string),
							"action": newPack["action"].(string),
						}
						workersInfo = append(workersInfo, worker)
					} else {
						var worker = map[string]string{
							"name":   oldPack["name"].(string),
							"id":     oldPack["id"].(string),
							"action": oldPack["action"].(string),
						}
						workersInfo = append(workersInfo, worker)
					}
				}
			}
			if !exists {
				params := v1.WorkerParam{
					Action: "add",
					Count:  1,
				}
				err := workerClient.Add(clusterID, params, targetEnv)
				if err != nil {
					return fmt.Errorf("Error adding worker to cluster")
				}
				id, err := getID(d, meta, clusterID, oldWorker, workersInfo)
				if err != nil {
					return fmt.Errorf("Error getting id of worker")
				}
				var worker = map[string]string{
					"name":   newPack["name"].(string),
					"id":     id,
					"action": newPack["action"].(string),
				}
				workersInfo = append(workersInfo, worker)
			}
		}
		for _, oW := range oldWorker {
			oldPack := oW.(map[string]interface{})
			exists := false
			for _, nW := range newWorker {
				newPack := nW.(map[string]interface{})
				exists = exists || (strings.Compare(oldPack["name"].(string), newPack["name"].(string)) == 0)
			}
			if !exists {
				workerClient.Delete(clusterID, oldPack["id"].(string), targetEnv)
			}

		}
		//wait for new workers to available
		//Done - Can we not put WaitForWorkerAvailable after all client.DeleteWorker
		WaitForWorkerAvailable(d, meta, targetEnv)
		d.Set("workers", workersInfo)
	}

	//TODO put webhooks can't deleted in the error message if such case is observed in the chnages
	if d.HasChange("webhook") {
		oldHooks, newHooks := d.GetChange("webhook")
		oldHook := oldHooks.([]interface{})
		newHook := newHooks.([]interface{})
		for _, nH := range newHook {
			newPack := nH.(map[string]interface{})
			exists := false
			for _, oH := range oldHook {
				oldPack := oH.(map[string]interface{})
				if (strings.Compare(newPack["level"].(string), oldPack["level"].(string)) == 0) && (strings.Compare(newPack["type"].(string), oldPack["type"].(string)) == 0) && (strings.Compare(newPack["url"].(string), oldPack["url"].(string)) == 0) {
					exists = true
				}
			}
			if !exists {
				webhook := v1.WebHook{
					Level: newPack["level"].(string),
					Type:  newPack["type"].(string),
					URL:   newPack["url"].(string),
				}

				webhookClient.Add(clusterID, webhook, targetEnv)
			}
		}
	}
	//TODO put subnet can't deleted in the error message if such case is observed in the chnages
	if d.HasChange("subnet_id") {
		oldSubnets, newSubnets := d.GetChange("subnet_id")
		oldSubnet := oldSubnets.(*schema.Set)
		newSubnet := newSubnets.(*schema.Set)
		for _, nS := range newSubnet.List() {
			exists := false
			for _, oS := range oldSubnet.List() {
				if strings.Compare(nS.(string), oS.(string)) == 0 {
					exists = true
				}
			}
			if !exists {
				err := subnetClient.AddSubnet(clusterID, nS.(string), targetEnv)
				if err != nil {
					return err
				}
			}
		}
	}
	return resourceIBMCloudArmadaClusterRead(d, meta)
}

func getID(d *schema.ResourceData, meta interface{}, clusterID string, oldWorkers []interface{}, workerInfo []map[string]string) (string, error) {
	targetEnv := getClusterTargetHeader(d)
	workerClient, err := meta.(ClientSession).ClusterWorkerClient()
	if err != nil {
		return "", err
	}
	workerFields, err := workerClient.List(clusterID, targetEnv)
	if err != nil {
		return "", err
	}
	for _, wF := range workerFields {
		exists := false
		for _, oW := range oldWorkers {
			oldPack := oW.(map[string]interface{})
			if strings.Compare(wF.ID, oldPack["id"].(string)) == 0 || strings.Compare(wF.State, "deleted") == 0 {
				exists = true
			}
		}
		if !exists {
			for i := 0; i < len(workerInfo); i++ {
				pack := workerInfo[i]
				exists = exists || (strings.Compare(wF.ID, pack["id"]) == 0)
			}
			if !exists {
				return wF.ID, nil
			}
		}
	}

	return "", fmt.Errorf("Unable to get ID of worker")
}

func resourceIBMCloudArmadaClusterDelete(d *schema.ResourceData, meta interface{}) error {
	targetEnv := getClusterTargetHeader(d)
	client, err := meta.(ClientSession).ClusterClient()
	if err != nil {
		return err
	}
	clusterID := d.Id()
	err = client.Delete(clusterID, targetEnv)
	if err != nil {
		return fmt.Errorf("Error deleting cluster: %s", err)
	}
	return nil
}

// WaitForClusterAvailable Waits for cluster creation
func WaitForClusterAvailable(d *schema.ResourceData, meta interface{}, target *v1.ClusterTargetHeader) (interface{}, error) {
	log.Printf("Waiting for cluster (%s) to be available.", d.Id())
	id := d.Id()

	client, _ := meta.(ClientSession).ClusterClient()

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", clusterProvisioning},
		Target:     []string{clusterNormal},
		Refresh:    clusterStateRefreshFunc(client, id, d, target),
		Timeout:    time.Duration(d.Get("wait_time_minutes").(int)) * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func clusterStateRefreshFunc(client v1.Clusters, instanceID string, d *schema.ResourceData, target *v1.ClusterTargetHeader) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		clusterFields, err := client.Find(instanceID, target)
		if err != nil {
			return nil, "", fmt.Errorf("Error retrieving cluster: %s", err)
		}
		// Check active transactions
		log.Println("Checking cluster")
		//Check for cluster state to be normal
		log.Println("Checking cluster state %s", strings.Compare(clusterFields.State, clusterNormal))
		if strings.Compare(clusterFields.State, clusterNormal) != 0 {
			return clusterFields, clusterProvisioning, nil
		}
		return clusterFields, clusterNormal, nil
	}
}

// WaitForWorkerAvailable Waits for worker creation
func WaitForWorkerAvailable(d *schema.ResourceData, meta interface{}, target *v1.ClusterTargetHeader) (interface{}, error) {
	log.Printf("Waiting for worker of the cluster (%s) to be available.", d.Id())
	id := d.Id()

	workerClient, _ := meta.(ClientSession).ClusterWorkerClient()
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", workerProvisioning},
		Target:     []string{workerNormal},
		Refresh:    workerStateRefreshFunc(workerClient, id, d, target),
		Timeout:    time.Duration(d.Get("wait_time_minutes").(int)) * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func workerStateRefreshFunc(client v1.Workers, instanceID string, d *schema.ResourceData, target *v1.ClusterTargetHeader) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		workerFields, err := client.List(instanceID, target)
		if err != nil {
			return nil, "", fmt.Errorf("Error retrieving workers for cluster: %s", err)
		}
		log.Println("Checking workers...")
		//Done worker has two fields State and Status , so check for those 2
		for _, e := range workerFields {
			if strings.Compare(e.State, workerNormal) != 0 || strings.Compare(e.Status, workerReadyState) != 0 {
				if strings.Compare(e.State, "deleted") != 0 {
					return workerFields, workerProvisioning, nil
				}
			}
		}
		return workerFields, workerNormal, nil
	}
}

func resourceIBMCloudArmadaClusterExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	targetEnv := getClusterTargetHeader(d)
	client, err := meta.(ClientSession).ClusterClient()
	if err != nil {
		return false, err
	}
	clusterID := d.Id()
	cls, err := client.Find(clusterID, targetEnv)
	if err != nil {
		if apiErr, ok := err.(bmxerror.RequestFailure); ok {
			if apiErr.StatusCode() == 404 {
				return false, nil
			}
		}
		return false, fmt.Errorf("Error communicating with the API: %s", err)
	}
	return cls.ID == clusterID, nil
}
