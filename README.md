# etcd provider for Terraform

[Terraform](http://terraform.io) provider for etcd.
## Status

Development/Testing.


## TODO

## Install

This project used [gb](http://getgb.io), so you must have it
installed.

```shell
$ git clone https://github.com/bakins/terraform-provider-etcd
$ cd terraform-provider-etcd
$ make
$ sudo make install
```

will install to `/usr/local/bin/terraform-provider-etcd`. Set PREFIX
to change this:

```shell
$sudo make install PREFIX=/usr
```


Note: You may need to add something like the following to `~/.terraformrc` if you get an error about missing the etcd provider when running terraform:

```
providers {
  etcd = "/usr/local/bin/terraform-provider-etcd"
}
```

## Usage

### Discovery

Simple usage:

```
resource "etcd_discovery" "test" {
   size = 1
}

output "etcd" {
    value = "${etcd_discovery.test.url}"
}
```

The resource `etcd_discovery` has the following optional fields:

- `size` - cluster size. default is 3.
- `endpoint` - discovery endpoint. default is "https://discovery.etcd.io/new"

The resulting URL is availible in the `url` output of the resource -- `etcd_discovery.test.url` in this example.

### Keys

`etcd_keys` operates similar to
[consul_keys](https://www.terraform.io/docs/providers/consul/r/keys.html)

```
provider "etcd" {
    endpoint = "http://oneof.my.etcd.servers.or.proxies:port"
}

resource "etcd_keys" "ami" {
    # Read the launch AMI from etcd
    key {
        name = "ami"
        path = "service/app/launch_ami"
        default = "ami-1234"
    }

    # Set the CNAME of our load balancer as a key
    key {
        name = "elb_cname"
        path = "service/app/elb_address"
        value = "${aws_elb.app.dns_name}"
    }

# Start our instance with the dynamic ami value
resource "aws_instance" "app" {
    ami = "${etcd_keys.app.var.ami}"
    ...
}
```
