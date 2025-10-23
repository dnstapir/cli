PROG:=dnstapir-cli
# -----
VERSION:=`cat ./VERSION`
COMMIT:=`git describe --dirty=+WiP --always`
APPDATE=`date +"%Y-%m-%d-%H:%M"`
GOFLAGS:=-v -ldflags "-X app.version=$(VERSION)-$(COMMIT)"

GOOS ?= $(shell uname -s | tr A-Z a-z)

GO:=GOOS=$(GOOS) CGO_ENABLED=0 go
# GO:=GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=1 go

SPECFILE:=rpm/SPECS/dnstapir-cli.spec

default: ${PROG}

${PROG}: build

build:  
	/bin/sh make-version.sh $(VERSION)-$(COMMIT) $(APPDATE) $(PROG)
	$(GO) build $(GOFLAGS) -o ${PROG}

linux:  
	/bin/sh make-version.sh $(VERSION)-$(COMMIT) $(APPDATE) $(PROG)
	GOOS=linux GOARCH=amd64 go build $(GOFLAGS) -o ${PROG}.linux

netbsd:  
	/bin/sh make-version.sh $(VERSION)-$(COMMIT) $(APPDATE) $(PROG)
	GOOS=netbsd GOARCH=amd64 go build $(GOFLAGS) -o ${PROG}.netbsd

test:
	$(GO) test -v -cover

clean: SHELL:=/bin/bash
clean:
	@rm -f $(PROG) *~ cmd/*~
	@rm -f *.tar.gz
	@rm -f rpm/SOURCES/*.tar.gz
	@rm -rf rpm/{BUILD,BUILDROOT,SRPMS,RPMS}
	@rm -rf deb/usr
	@rm -rf deb/etc
	@rm -rf deb/var
	@rm -f *.deb

install:
	install -b -c -s ${PROG} /usr/local/bin/

tarball:
	git archive --format=tar.gz --prefix=$(PROG)/ -o $(PROG)-$(VERSION).tar.gz HEAD

srpm: SHELL:=/bin/bash
srpm: tarball
	test $$(rpmspec -q --qf '%{version}' $(SPECFILE) 2>/dev/null || grep '^Version:' $(SPECFILE) | awk '{print $$2}') == $(VERSION)
	mkdir -p rpm/{BUILD,RPMS,SRPMS}
	cp $(PROG)-$(VERSION).tar.gz rpm/SOURCES/
	rpmbuild -bs --define "%_topdir ./rpm" --undefine=dist $(SPECFILE)
	test -z "$(outdir)" || cp rpm/SRPMS/*.src.rpm "$(outdir)"

deb: build
	mkdir -p deb/usr/bin
	mkdir -p deb/etc/dnstapir/certs
	mkdir -p deb/usr/lib/systemd/system
	cp dnstapir-cli deb/usr/bin
	cp rpm/SOURCES/dnstapir-renew.service deb/usr/lib/systemd/system
	cp rpm/SOURCES/dnstapir-renew.timer deb/usr/lib/systemd/system
	dpkg-deb -b deb/ $(PROG)-$(VERSION).deb

.PHONY: build clean
