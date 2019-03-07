all:
	@go build -buildmode=plugin -o=$(PNAME)/$(PNAME).so  $(PNAME)/$(PNAME).go
	@md5sum $(PNAME)/$(PNAME).so