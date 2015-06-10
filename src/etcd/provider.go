package etcd

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/mitchellh/mapstructure"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"endpoint": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "http://127.0.0.1:2379",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"etcd_discovery": resourceEtcdDiscovery(),
			"etcd_keys":      resourceEtcdKeys(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	var config Config
	configRaw := d.Get("").(map[string]interface{})
	if err := mapstructure.Decode(configRaw, &config); err != nil {
		return nil, err
	}
	log.Printf("[INFO] Initializing etcd client")
	return config.Client()
}
