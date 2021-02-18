run: build
	@./heimatcli

build: 
	@go build -o ./heimatcli ./cmd/*.go