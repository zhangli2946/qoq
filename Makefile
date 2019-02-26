
all: build
	@echo "Done"
build: clean 
	@go build -o agent
clean: 
	@go clean
plugin: 
	$(MAKE) -C handler/biz
	$(MAKE) -C handler/sys