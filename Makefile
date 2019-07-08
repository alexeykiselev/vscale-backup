PROJECT=vscale-backup

VERSION=$(shell git describe --tags --always --dirty)

.PHONY: vetcheck fmtcheck clean gotest

all: vetcheck fmtcheck gotest

ver:
	@echo Building version: $(VERSION)

gotest:
	go test -cover ./...

fmtcheck:
	@gofmt -l -s . | grep ".*\.go"; if [ "$$?" = "0" ]; then exit 1; fi

clean:
	@rm -rf build

vetcheck:
	go vet ./...

build-linux:
	@CGO_ENABLE=0 GOOS=linux GOARCH=amd64 go build -o build/bin/linux-amd64/$(PROJECT) -ldflags="-X main.version=$(VERSION)" .
build-darwin:
	@CGO_ENABLE=0 GOOS=darwin GOARCH=amd64 go build -o build/bin/darwin-amd64/$(PROJECT) -ldflags="-X main.version=$(VERSION)" .
build-windows:
	@CGO_ENABLE=0 GOOS=windows GOARCH=amd64 go build -o build/bin/windows-amd64/$(PROJECT).exe -ldflags="-X main.version=$(VERSION)" .

release: ver build-linux build-darwin build-windows

dist: clean release
	@mkdir -p build/dist
	@cd ./build/; zip -j ./dist/$(PROJECT)_$(VERSION)_Windows-64bit.zip ./bin/windows-amd64/$(PROJECT)*
	@cd ./build/bin/linux-amd64/; tar pzcvf ../../dist/$(PROJECT)_$(VERSION)_Linux-64bit.tar.gz ./$(PROJECT)*
	@cd ./build/bin/darwin-amd64/; tar pzcvf ../../dist/$(PROJECT)_$(VERSION)_macOS-64bit.tar.gz ./$(PROJECT)*
