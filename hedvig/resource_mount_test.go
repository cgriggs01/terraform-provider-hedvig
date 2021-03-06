package hedvig

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"testing"
	"os"
)

func testHedvigMount() error {
	return nil
}

func TestAccHedvigMount(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccHedvigMountConfig,
				Check:  testAccCheckHedvigMountExists("hedvig_mount.test-mount"),
			},
		},
	})
}

var testAccHedvigMountConfig = fmt.Sprintf(`
provider "hedvig" {
  node = "%s"
  username = "%s"
  password = "%s"
}

resource "hedvig_vdisk" "test-mount-vdisk" {
  cluster = "%s"
  name = "HedvigVdiskTest4"
  size = 11
  type = "NFS"
}

resource "hedvig_mount" "test-mount" {
  cluster = "%s"
  vdisk = "${hedvig_vdisk.test-mount-vdisk.name}"
  controller = "%s"
}
`, os.Getenv("HV_TESTNODE"), os.Getenv("HV_TESTUSER"), os.Getenv("HV_TESTPASS"), os.Getenv("HV_TESTCLUST"), os.Getenv("HV_TESTCLUST"), os.Getenv("HV_TESTCONT"))

func testAccCheckHedvigMountExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No lun ID is set")
		}

		return nil
	}
}
