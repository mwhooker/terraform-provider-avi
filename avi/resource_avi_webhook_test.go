package avi

import (
	"fmt"
	"strings"
	"testing"

	"github.com/avinetworks/sdk/go/clients"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAVIWebhookBasic(t *testing.T) {
	//updatedConfig := fmt.Sprintf(testAccAVIWebhookConfig, "abc")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAVIWebhookDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAVIWebhookConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAVIWebhookExists("avi_webhook.testwebhook"),
					resource.TestCheckResourceAttr(
						"avi_webhook.testwebhook", "name", "wb-test")),
			},
			{
				Config: updatetestAccAVIWebhookConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAVIWebhookExists("avi_webhook.testwebhook"),
					resource.TestCheckResourceAttr(
						"avi_webhook.testwebhook", "name", "wb-abc")),
			},
		},
	})

}

func testAccCheckAVIWebhookExists(resourcename string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*clients.AviClient).AviSession
		var obj interface{}
		rs, ok := s.RootModule().Resources[resourcename]
		if !ok {
			return fmt.Errorf("Not found: %s", resourcename)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Webhook ID is set")
		}
		url := strings.SplitN(rs.Primary.ID, "/api", 2)[1]
		uuid := strings.Split(url, "#")[0]
		path := "api" + uuid
		err := conn.Get(path, &obj)
		if err != nil {
			return err
		}
		return nil
	}

}

func testAccCheckAVIWebhookDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*clients.AviClient).AviSession
	var obj interface{}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "avi_webhook" {
			continue
		}
		url := strings.SplitN(rs.Primary.ID, "/api", 2)[1]
		uuid := strings.Split(url, "#")[0]
		path := "api" + uuid
		err := conn.Get(path, &obj)
		if err != nil {
			if strings.Contains(err.Error(), "404") {
				return nil
			}
			return err
		}
		if len(obj.(map[string]interface{})) > 0 {
			return fmt.Errorf("AVI Webhook still exists")
		}
	}
	return nil
}

const testAccAVIWebhookConfig = `
data "avi_tenant" "default_tenant"{
	name= "admin"
}

resource "avi_webhook" "testwebhook" {
	name = "wb-test"
	tenant_ref= "${data.avi_tenant.default_tenant.id}"
}
`

const updatetestAccAVIWebhookConfig = `
data "avi_tenant" "default_tenant"{
	name= "admin"
}

resource "avi_webhook" "testwebhook" {
	name = "wb-abc"
	tenant_ref= "${data.avi_tenant.default_tenant.id}"
}
`