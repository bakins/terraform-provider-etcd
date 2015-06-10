package etcd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceEtcdDiscovery() *schema.Resource {
	return &schema.Resource{
		Create: discoverCreate,
		Delete: discoverDelete,
		Exists: discoverExists,
		Read:   discoverRead,

		Schema: map[string]*schema.Schema{
			"endpoint": {
				Type:        schema.TypeString,
				Description: "discovery endpoint",
				Default:     "https://discovery.etcd.io/new",
				Optional:    true,
				ForceNew:    true,
			},
			"size": &schema.Schema{
				Type:        schema.TypeInt,
				Description: "size of cluster",
				Default:     3,
				Optional:    true,
				ForceNew:    true,
			},
			"url": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "url",
			},
		},
	}
}

func discoverCreate(d *schema.ResourceData, meta interface{}) error {
	resp, err := http.Get(fmt.Sprintf("%s?size=%d", d.Get("endpoint").(string), d.Get("size").(int)))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	url := strings.TrimSpace(string(data))
	d.Set("url", url)
	d.SetId(url)
	return nil
}

func discoverDelete(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")
	return nil
}

func discoverExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	return d.Id() != "", nil
}

func discoverRead(d *schema.ResourceData, meta interface{}) error {
	d.Set("url", d.Id())
	d.SetId(d.Id())
	return nil
}
