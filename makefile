build:
	go build -o out/hue-hook ./...

build_arm:
	GOARCH=arm
	GOARM=7
	make build