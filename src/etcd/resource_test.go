package etcd

import "github.com/hashicorp/terraform/terraform"

var testProviders = map[string]terraform.ResourceProvider{
	"etcd": Provider(),
}
