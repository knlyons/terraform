package ibmcloud

import (
	"fmt"
	"testing"

	"strings"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/IBM-Bluemix/bluemix-go/api/cf/cfv2"
)

func TestAccIBMCloudCFSharedDomain_Basic(t *testing.T) {
	var conf cfv2.SharedDomainFields
	name := fmt.Sprintf("terraform%d.com", acctest.RandInt())

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMCloudCFSharedDomain_basic(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMCloudCFSharedDomainExists("ibmcloud_cf_shared_domain.domain", &conf),
					resource.TestCheckResourceAttr("ibmcloud_cf_shared_domain.domain", "name", name),
				),
			},
		},
	})
}

func testAccCheckIBMCloudCFSharedDomainExists(n string, obj *cfv2.SharedDomainFields) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		cfClient, err := testAccProvider.Meta().(ClientSession).CFAPI()
		if err != nil {
			return err
		}
		sharedDomainGUID := rs.Primary.ID

		shdomain, err := cfClient.SharedDomains().Get(sharedDomainGUID)
		if err != nil {
			return err
		}

		*obj = *shdomain
		return nil
	}
}

func testAccCheckIBMCloudCFSharedDomainDestroy(s *terraform.State) error {
	cfClient, err := testAccProvider.Meta().(ClientSession).CFAPI()
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibmcloud_cf_shared_domain" {
			continue
		}

		sharedDomainGUID := rs.Primary.ID

		// Try to find the shared domain
		_, err := cfClient.SharedDomains().Get(sharedDomainGUID)

		if err != nil && !strings.Contains(err.Error(), "404") {
			return fmt.Errorf("Error waiting for CF shared domain (%s) to be destroyed: %s", rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckIBMCloudCFSharedDomain_basic(name string) string {
	return fmt.Sprintf(`
	
		resource "ibmcloud_cf_shared_domain" "domain" {
			name = "%s"
		}
	`, name)
}
