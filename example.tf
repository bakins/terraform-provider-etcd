resource "etcd_discovery" "test" {
   size = 1
}

output "etcd" {
    value = "${etcd_discovery.test.url}"
}

provider "etcd" {
    endpoint = "http://localhost:2379"
}

# Access a key in Consul
resource "etcd_keys" "test" {
    key {
        name = "test"
        path = "/test/1/2/3"
        value = "testing"
    }
}

output "key" {
    value = "${etcd_keys.test.var.test}"
}
