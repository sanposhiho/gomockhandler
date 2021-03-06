.PHONY:genjson
genjson:
	easyjson -all ./model/chunk.go

# test on playground
.PHONY:playground
playground:
	go build . ;\
	cd playground ;\
	go generate ./...; \
	rm ../gomockhandler
