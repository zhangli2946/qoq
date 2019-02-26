
all:
	ifeq($(PNAME), "ALL")
	else
		go build -buildmode=plugin -o=$(PNAME)/$(PNAME)  $(PNAME)/$(PNAME).go
	endif