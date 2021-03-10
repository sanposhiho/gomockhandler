.PHONY:easyjson
easyjson:
	easyjson -all ./model/config.go

.PHONY:build
build:
	go build .

# test on playground
.PHONY:playground
playground:
	go build . ;\
	cd playground ;\
	go generate ./...; \
	rm ../gomockhandler;

# test on playground
.PHONY:playground-gen
playground-gen:
	go build . ;\
	cd playground ;\
	../gomockhandler -config=gomockhandler.json mockgen ;\
	rm ../gomockhandler;

# test on playground
.PHONY:playground-check
playground-check:
	go build . ;\
	cd playground ;\
	../gomockhandler -config=gomockhandler.json check ;\
	rm ../gomockhandler;

# clean playground
.PHONY:clean
clean:
	cd playground ;\
	rm gomockhandler.json;\
	rm -rf mock/
