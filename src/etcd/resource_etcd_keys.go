package etcd

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"strconv"

	etcdErr "github.com/coreos/etcd/error"
	client "github.com/coreos/go-etcd/etcd"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceEtcdKeys() *schema.Resource {
	return &schema.Resource{
		Create: resourceEtcdKeysCreate,
		Update: resourceEtcdKeysCreate,
		Read:   resourceEtcdKeysRead,
		Delete: resourceEtcdKeysDelete,

		Schema: map[string]*schema.Schema{
			"key": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},

						"path": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},

						"value": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"default": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},

						"delete": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
				Set: resourceEtcdKeysHash,
			},

			"var": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
	}
}

func resourceEtcdKeysHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%s-", m["name"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["path"].(string)))
	return hashcode.String(buf.String())
}

func resourceEtcdKeysCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client.Client)

	// Store the computed vars
	vars := make(map[string]string)

	// Extract the keys
	keys := d.Get("key").(*schema.Set).List()
	for _, raw := range keys {
		key, path, sub, err := parseKey(raw)
		if err != nil {
			return err
		}

		value := sub["value"].(string)
		if value != "" {
			log.Printf("[DEBUG] setting etcd key '%s' to '%v", path, value, d)
			if _, err := c.Set(path, value, 0); err != nil {
				return fmt.Errorf("Failed to set etcd key '%s': %v", path, err)
			}
			vars[key] = value
			sub["value"] = value

		} else {
			log.Printf("[DEBUG] Getting etcd key '%s'", path)
			resp, err := c.Get(path, false, false)
			if err != nil && !IsKeyNotFound(err) {
				return fmt.Errorf("Failed to get etcd key '%s': %v", path, err)
			}
			value, err := attributeValue(sub, key, resp)
			if err != nil {
				return fmt.Errorf("Failed to get etcd value for '%s': %v", path, err)
			}
			vars[key] = value
		}
	}

	// Update the resource
	d.SetId("etcd")
	d.Set("key", keys)
	d.Set("var", vars)
	return nil
}

func resourceEtcdKeysRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client.Client)
	// Store the computed vars
	vars := make(map[string]string)

	// Extract the keys
	keys := d.Get("key").(*schema.Set).List()
	for _, raw := range keys {
		key, path, sub, err := parseKey(raw)
		if err != nil {
			return err
		}

		log.Printf("[DEBUG] Refreshing etcd value of key '%s", path)

		resp, err := c.Get(path, false, false)
		if err != nil && !IsKeyNotFound(err) {
			return fmt.Errorf("Failed to get value for path '%s' from etcd: %v", path, err)
		}

		value, err := attributeValue(sub, key, resp)
		if err != nil {
			return fmt.Errorf("Failed to get value for path '%s' from etcd: %v", path, err)
		}
		vars[key] = value
		sub["value"] = value
	}

	// Update the resource
	d.Set("key", keys)
	d.Set("var", vars)
	return nil
}

func resourceEtcdKeysDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client.Client)

	// Extract the keys
	keys := d.Get("key").(*schema.Set).List()
	for _, raw := range keys {
		_, path, sub, err := parseKey(raw)
		if err != nil {
			return err
		}

		// Ignore if the key is non-managed
		shouldDelete, ok := sub["delete"].(bool)
		if !ok || !shouldDelete {
			continue
		}

		log.Printf("[DEBUG] Deleting etcd key '%s'", path)

		if _, err := c.Delete(path, false); err != nil && !IsKeyNotFound(err) {
			return fmt.Errorf("Failed to delete etcd key '%s': %v", path, err)
		}
	}

	// Clear the ID
	d.SetId("")
	return nil
}

// parseKey is used to parse a key into a name, path, config or error
func parseKey(raw interface{}) (string, string, map[string]interface{}, error) {
	sub, ok := raw.(map[string]interface{})
	if !ok {
		return "", "", nil, fmt.Errorf("Failed to unroll: %#v", raw)
	}

	key, ok := sub["name"].(string)
	if !ok {
		return "", "", nil, fmt.Errorf("Failed to expand key '%#v'", sub)
	}

	path, ok := sub["path"].(string)
	if !ok {
		return "", "", nil, fmt.Errorf("Failed to get path for key '%s'", key)
	}
	return key, path, sub, nil
}

// attributeValue determines the value for a key, potentially
// using a default value if provided.
func attributeValue(sub map[string]interface{}, key string, resp *client.Response) (string, error) {
	if resp == nil {
		// Use a default if given
		if raw, ok := sub["default"]; ok {
			switch def := raw.(type) {
			case string:
				return def, nil
			case bool:
				return strconv.FormatBool(def), nil
			}
		}
	}

	n := resp.Node
	if n == nil {
		return "", errors.New("response has nil node")
	}

	if n.Dir {
		return "", fmt.Errorf("%s is a directory", key)
	}

	return n.Value, nil
}

func IsKeyNotFound(err error) bool {
	e, ok := err.(*client.EtcdError)
	return ok && e.ErrorCode == etcdErr.EcodeKeyNotFound
}
