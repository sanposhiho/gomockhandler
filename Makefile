.PHONY:build
build:
	go build .

# for playground
# build gomockhandler and generate gomockhandler.json with it.
.PHONY:playground
playground:
	go build . ;\
	cd playground ;\
	go generate ./...; \
	rm ../gomockhandler;

# for playground
# build gomockhandler and generate mocks with it.
# It use gomockhandler.json on /playground dir. So you may have to run make playground first.
.PHONY:playground-gen
playground-gen:
	go build . ;\
	cd playground ;\
	../gomockhandler -config=gomockhandler.json mockgen ;\
	rm ../gomockhandler;

# test on playground
# build gomockhandler and check mocks with it.
# It use gomockhandler.json on /playground dir. So you may have to run make playground first.
.PHONY:playground-check
playground-check:
	go build . ;\
	cd playground ;\
	../gomockhandler -config=gomockhandler.json check ;\
	rm ../gomockhandler;

# test on playground
# build gomockhandler and delete mocks with it.
# It use gomockhandler.json on /playground dir. So you may have to run make playground first.
.PHONY:playground-delete
playground-delete:
	go build . ;\
	cd playground ;\
	../gomockhandler -config=gomockhandler.json -destination=./mock/user.go deletemock ;\
	rm ../gomockhandler;


# clean up playground
.PHONY:clean
clean:
	cd playground ;\
	rm gomockhandler.json;\
	rm -rf mock/
