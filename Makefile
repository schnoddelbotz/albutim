
PLATFORMS = linux/amd64 darwin/amd64 windows/amd64 linux/arm

VERSION = $(shell git describe --tags | cut -dv -f2)
LDFLAGS := -X github.com/schnoddelbotz/albutim/cmd.AppVersion=$(VERSION) -w

all: albutim

albutim:
	go generate
	CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)"

dependencies:
	go get github.com/mjibson/esc

clean:
	rm -f albutim lib/assets.go

run: clean albutim
	./albutim serve --root testalbum

release:
	for platform in $(PLATFORMS); do \
		echo "Building for $$platform..."; \
		export GOOS=`echo $$platform | cut -d/ -f1` GOARCH=`echo $$platform | cut -d/ -f2`; \
			export SUFFIX=`test $${GOOS} = windows && echo .exe || echo` ; \
			CGO_ENABLED=0 go build -o dashboard-nerf_$${GOOS}-$${GOARCH}$${SUFFIX} -ldflags "$(LDFLAGS)"; \
	done

ziprelease: release
	for bin in dashboard-nerf_darwin* dashboard-nerf_linux* dashboard-nerf_windows*; do \
		archive=`echo $${bin} | sed -e 's@.exe@@'` ; \
		zip $${archive}_v$(VERSION).zip $$bin; \
	done
