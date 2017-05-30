package ibmcloud

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	slsession "github.com/softlayer/softlayer-go/session"

	bluemix "github.com/IBM-Bluemix/bluemix-go"
	"github.com/IBM-Bluemix/bluemix-go/api/account/accountv2"
	"github.com/IBM-Bluemix/bluemix-go/api/cf/cfv2"
	"github.com/IBM-Bluemix/bluemix-go/api/k8scluster/k8sclusterv1"
	"github.com/IBM-Bluemix/bluemix-go/bmxerror"
	"github.com/IBM-Bluemix/bluemix-go/endpoints"
	bxsession "github.com/IBM-Bluemix/bluemix-go/session"
)

//SoftlayerRestEndpoint rest endpoint of SoftLayer
const SoftlayerRestEndpoint = "https://api.softlayer.com/rest/v3"

var (
	errEmptySoftLayerCredentials = errors.New("softlayer_username and softlayer_api_key must be provided. Please see the documentation on how to configure them")
	errEmptyBluemixCredentials   = errors.New("bluemix_api_key must be provided. Please see the documentation on how to configure it")
)

//Config stores user provider input
type Config struct {
	//BluemixAPIKey is the Bluemix api key
	BluemixAPIKey string
	//Bluemix region
	Region string
	//Bluemix API timeout
	BluemixTimeout time.Duration

	//Softlayer end point url
	SoftLayerEndpointURL string

	//Softlayer API timeout
	SoftLayerTimeout time.Duration

	// Softlayer User Name
	SoftLayerUserName string

	// Softlayer API Key
	SoftLayerAPIKey string

	//SkipServiceConfig is a set of services whose configuration is to be skipped. Valid values could be bluemix, softlayer etc
	SkipServiceConfig *schema.Set

	//Retry Count for API calls
	//Unexposed in the schema at this point as they are used only during session creation for a few calls
	//When sdk implements it we an expose them for expected behaviour
	//https://github.com/softlayer/softlayer-go/issues/41
	RetryCount int
	//Constant Retry Delay for API calls
	RetryDelay time.Duration
}

//Session stores the information required for communication with the SoftLayer and Bluemix API
type Session struct {
	// SoftLayerSesssion is the the SoftLayer session used to connect to the SoftLayer API
	SoftLayerSession *slsession.Session

	// BluemixSession is the the Bluemix session used to connect to the Bluemix API
	BluemixSession *bxsession.Session
}

// ClientSession  contains  Bluemix/SoftLayer session and clients
type ClientSession interface {
	SoftLayerSession() *slsession.Session
	BluemixSession() (*bxsession.Session, error)

	ClusterClient() (k8sclusterv1.Clusters, error)
	ClusterWorkerClient() (k8sclusterv1.Workers, error)
	ClusterSubnetClient() (k8sclusterv1.Subnets, error)
	ClusterWebHooksClient() (k8sclusterv1.Webhooks, error)

	CloudFoundryAppClient() (cfv2.Apps, error)
	CloudFoundryOrgClient() (cfv2.Organizations, error)
	CloudFoundryServiceBindingClient() (cfv2.ServiceBindings, error)
	CloudFoundryServiceInstanceClient() (cfv2.ServiceInstances, error)
	CloudFoundryServicePlanClient() (cfv2.ServicePlans, error)
	CloudFoundryServiceKeyClient() (cfv2.ServiceKeys, error)
	CloudFoundryServiceOfferingClient() (cfv2.ServiceOfferings, error)
	CloudFoundrySpaceClient() (cfv2.Spaces, error)
	CloudFoundrySpaceQuotaClient() (cfv2.SpaceQuotas, error)
	CloudFoundryRouteClient() (cfv2.Routes, error)
	CloudFoundrySharedDomainClient() (cfv2.SharedDomains, error)
	CloudFoundryPrivateDomainClient() (cfv2.PrivateDomains, error)

	BluemixAcccountClient() accountv2.Accounts
}

type clientSession struct {
	session *Session

	csConfigErr error
	csClient    k8sclusterv1.Clusters
	csWorker    k8sclusterv1.Workers
	csSubnet    k8sclusterv1.Subnets
	csWebHook   k8sclusterv1.Webhooks

	cfConfigErr              error
	cfAppClient              cfv2.Apps
	cfOrgClient              cfv2.Organizations
	cfServiceInstanceClient  cfv2.ServiceInstances
	cfSpaceClient            cfv2.Spaces
	cfSpaceQuotaClient       cfv2.SpaceQuotas
	cfServiceBindingClient   cfv2.ServiceBindings
	cfServicePlanClient      cfv2.ServicePlans
	cfServiceKeysClient      cfv2.ServiceKeys
	cfServiceOfferingsClient cfv2.ServiceOfferings
	cfRouteClient            cfv2.Routes
	cfSharedDomainClient     cfv2.SharedDomains
	cfPrivateDomainClient    cfv2.PrivateDomains

	accountConfigErr     error
	bluemixAccountClient accountv2.Accounts
}

// SoftLayerSession providers SoftLayer Session
func (sess clientSession) SoftLayerSession() *slsession.Session {
	return sess.session.SoftLayerSession
}

// CloudFoundryOrgClient providers Cloud Foundary org APIs
func (sess clientSession) CloudFoundryOrgClient() (cfv2.Organizations, error) {
	return sess.cfOrgClient, sess.cfConfigErr
}

// CloudFoundrySpaceClient providers Cloud Foundary space APIs
func (sess clientSession) CloudFoundrySpaceClient() (cfv2.Spaces, error) {
	return sess.cfSpaceClient, sess.cfConfigErr
}

// CloudFoundrySpaceQuotaClient providers Cloud Foundary space quota APIs
func (sess clientSession) CloudFoundrySpaceQuotaClient() (cfv2.SpaceQuotas, error) {
	return sess.cfSpaceQuotaClient, sess.cfConfigErr
}

// CloudFoundryAppClient providers Cloud Foundary app APIs
func (sess clientSession) CloudFoundryAppClient() (cfv2.Apps, error) {
	return sess.cfAppClient, sess.cfConfigErr
}

// CloudFoundryServiceBindingClient providers Cloud Foundary service binding APIs
func (sess clientSession) CloudFoundryServiceBindingClient() (cfv2.ServiceBindings, error) {
	return sess.cfServiceBindingClient, sess.cfConfigErr
}

// CloudFoundryServiceInstanceClient providers Cloud Foundary service APIs
func (sess clientSession) CloudFoundryServiceInstanceClient() (cfv2.ServiceInstances, error) {
	return sess.cfServiceInstanceClient, sess.cfConfigErr
}

// CloudFoundryServiceClient providers Cloud Foundary service APIs
func (sess clientSession) CloudFoundryServicePlanClient() (cfv2.ServicePlans, error) {
	return sess.cfServicePlanClient, sess.cfConfigErr
}

// CloudFoundryServiceKeyClient providers Cloud Foundary service APIs
func (sess clientSession) CloudFoundryServiceKeyClient() (cfv2.ServiceKeys, error) {
	return sess.cfServiceKeysClient, sess.cfConfigErr
}

// CloudFoundryServiceClient providers Cloud Foundary service APIs
func (sess clientSession) CloudFoundryServiceOfferingClient() (cfv2.ServiceOfferings, error) {
	return sess.cfServiceOfferingsClient, sess.cfConfigErr
}

// CloudFoundryRoute providers Cloud Foundary route APIs
func (sess clientSession) CloudFoundryRouteClient() (cfv2.Routes, error) {
	return sess.cfRouteClient, sess.cfConfigErr
}

// CloudFoundrySharedDomainClient providers Cloud Foundary shared domain APIs
func (sess clientSession) CloudFoundrySharedDomainClient() (cfv2.SharedDomains, error) {
	return sess.cfSharedDomainClient, sess.cfConfigErr
}

// CloudFoundryPrivateDomainClient providers Cloud Foundary private domain APIs
func (sess clientSession) CloudFoundryPrivateDomainClient() (cfv2.PrivateDomains, error) {
	return sess.cfPrivateDomainClient, sess.cfConfigErr
}

// BluemixAcccountClient providers Bluemix Account APIs
func (sess clientSession) BluemixAcccountClient() (accountv2.Accounts, error) {
	return sess.bluemixAccountClient, sess.accountConfigErr
}

// ClusterClient providers Bluemix Kubernetes Cluster APIs
func (sess clientSession) ClusterClient() (k8sclusterv1.Clusters, error) {
	return sess.csClient, sess.csConfigErr
}

// ClusterWorkerClient providers Bluemix Kubernetes Cluster APIs
func (sess clientSession) ClusterWorkerClient() (k8sclusterv1.Workers, error) {
	return sess.csWorker, sess.csConfigErr
}

// ClusterSubnetClient providers Bluemix Kubernetes Cluster APIs
func (sess clientSession) ClusterSubnetClient() (k8sclusterv1.Subnets, error) {
	return sess.csSubnet, sess.csConfigErr
}

// ClusterWebHooksClient providers Bluemix Kubernetes Cluster APIs
func (sess clientSession) ClusterWebHooksClient() (k8sclusterv1.Webhooks, error) {
	return sess.csWebHook, sess.csConfigErr
}

// BluemixSession to provide the Bluemix Session
func (sess clientSession) BluemixSession() (*bxsession.Session, error) {
	return sess.session.BluemixSession, sess.cfConfigErr
}

// ClientSession configures and returns a fully initialized ClientSession
func (c *Config) ClientSession() (interface{}, error) {

	sess, err := newSession(c)
	if err != nil {
		return nil, err
	}

	session := clientSession{
		session: sess,
	}

	if sess.BluemixSession == nil {
		log.Println("Skipping Bluemix Clients configuration")
		session.csConfigErr = errEmptyBluemixCredentials
		session.cfConfigErr = errEmptyBluemixCredentials
		session.accountConfigErr = errEmptyBluemixCredentials
		return session, nil
	}

	cfClient, err := cfv2.New(sess.BluemixSession)

	if err != nil {
		return nil, err
	}

	appAPI := cfClient.Apps()
	orgAPI := cfClient.Organizations()
	spaceAPI := cfClient.Spaces()
	serviceBindingAPI := cfClient.ServiceBindings()
	serviceInstanceAPI := cfClient.ServiceInstances()
	servicePlanAPI := cfClient.ServicePlans()
	serviceKeysAPI := cfClient.ServiceKeys()
	serviceOfferringAPI := cfClient.ServiceOfferings()
	routeAPI := cfClient.Routes()
	sharedDomainAPI := cfClient.SharedDomains()
	privateDomainAPI := cfClient.PrivateDomains()

	accClient, err := accountv2.New(sess.BluemixSession)
	if err != nil {
		return nil, err
	}
	accountAPI := accClient.Accounts()

	clusterClient, err := k8sclusterv1.New(sess.BluemixSession)
	var clusterConfigErr error
	if err != nil {
		if apiErr, ok := err.(bmxerror.Error); ok {
			if apiErr.Code() == endpoints.ErrCodeServiceEndpoint {
				clusterConfigErr = fmt.Errorf(`Cluster service doesn't exist in the region %q`, c.Region)
				session.csConfigErr = clusterConfigErr
			}
		}
		return nil, err
	}
	if clusterConfigErr == nil {
		clustersAPI := clusterClient.Clusters()
		clusterWorkerAPI := clusterClient.Workers()
		clusterSubnetsAPI := clusterClient.Subnets()
		clusterWebhookAPI := clusterClient.WebHooks()
		session.csClient = clustersAPI
		session.csSubnet = clusterSubnetsAPI
		session.csWorker = clusterWorkerAPI
		session.csWebHook = clusterWebhookAPI
	}

	session.cfAppClient = appAPI
	session.cfOrgClient = orgAPI
	session.cfServiceBindingClient = serviceBindingAPI
	session.cfServiceInstanceClient = serviceInstanceAPI
	session.cfServiceKeysClient = serviceKeysAPI
	session.cfServicePlanClient = servicePlanAPI
	session.cfServiceOfferingsClient = serviceOfferringAPI
	session.cfSpaceClient = spaceAPI
	session.cfRouteClient = routeAPI
	session.cfSharedDomainClient = sharedDomainAPI
	session.cfPrivateDomainClient = privateDomainAPI
	session.bluemixAccountClient = accountAPI

	return session, nil
}

func newSession(c *Config) (*Session, error) {
	ibmcloudSession := &Session{}

	if c.SoftLayerUserName != "" && c.SoftLayerAPIKey != "" {
		log.Println("Configuring SoftLayer Session ")
		softlayerSession := &slsession.Session{
			Endpoint: c.SoftLayerEndpointURL,
			Timeout:  c.SoftLayerTimeout,
			UserName: c.SoftLayerUserName,
			APIKey:   c.SoftLayerAPIKey,
			Debug:    os.Getenv("TF_LOG") != "",
		}
		ibmcloudSession.SoftLayerSession = softlayerSession
	}

	if c.BluemixAPIKey != "" {
		log.Println("Configuring Bluemix Session")
		var sess *bxsession.Session
		bmxConfig := &bluemix.Config{
			BluemixAPIKey: c.BluemixAPIKey,
			Debug:         os.Getenv("TF_LOG") != "",
			HTTPTimeout:   c.BluemixTimeout,
			Region:        c.Region,
			RetryDelay:    &c.RetryDelay,
			MaxRetries:    &c.RetryCount,
		}
		sess, err := bxsession.New(bmxConfig)
		if err != nil {
			return nil, err
		}
		ibmcloudSession.BluemixSession = sess
	}

	return ibmcloudSession, nil
}
