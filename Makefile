.PHONY: build restart

build:
	go build -o deployer
restart:
	systemctl restart deployer.service
