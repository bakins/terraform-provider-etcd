PREFIX := /usr/local

all: bin/terraform-provider-etcd

bin/terraform-provider-etcd:
	gb build all

clean:
	rm bin/*

distclean: clean
	rm -rf pkg

install: bin/terraform-provider-etcd
	install -m 755 -d $(PREFIX)/bin
	install -m 755 $< $(PREFIX)/bin/terraform-provider-etcd
