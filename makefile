run: build
	@./heimatcli

build: 
	@go build -o ./heimatcli ./src/cmd/*.go

install: build
	@cp ./heimatcli ~/go/bin/heimat_dev

release:
	goreleaser --rm-dist
