
PLATFORMS = linux/amd64 darwin/amd64 windows/amd64 linux/arm

VERSION = $(shell git describe --tags | cut -dv -f2)
LDFLAGS := -X github.com/schnoddelbotz/albutim/cmd.AppVersion=$(VERSION) -w
ASSETS := $(wildcard assets/*)

all: albutim

albutim: lib/assets.go
	CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)"

lib/assets.go: $(ASSETS)
	test -n "$(shell which esc)" || go get -v -u github.com/mjibson/esc
	go generate

clean:
	rm -f albutim albutim_* lib/assets.go

run: clean albutim
	./albutim serve --root testalbum

release: lib/assets.go
	for platform in $(PLATFORMS); do \
		echo "Building for $$platform..."; \
		export GOOS=`echo $$platform | cut -d/ -f1` GOARCH=`echo $$platform | cut -d/ -f2`; \
			export SUFFIX=`test $${GOOS} = windows && echo .exe || echo` ; \
			CGO_ENABLED=0 go build -o albutim_$${GOOS}-$${GOARCH}$${SUFFIX} -ldflags "$(LDFLAGS)"; \
	done

ziprelease: release
	for bin in albutim_darwin* albutim_linux* albutim_windows*; do \
		archive=`echo $${bin} | sed -e 's@.exe@@'` ; \
		zip $${archive}_v$(VERSION).zip $$bin; \
	done
