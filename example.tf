resource "etcd_discovery" "test" {
   size = 1
}

output "etcd" {
    value = "${etcd_discovery.test.url}"
}
