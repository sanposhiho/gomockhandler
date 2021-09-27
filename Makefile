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
.PHONY:playground-clean
clean:
	cd playground ;\
	rm gomockhandler.json;\
	rm -rf mock/


# for playground/benchmark
# build gomockhandler and generate gomockhandler.json with it.
.PHONY:benchmark
benchmark:
	cd playground/benchmark ;\
	go generate ./...; \
	rm ../../gomockhandler ;


# for playground/benchmark
# build gomockhandler and generate mocks with it.
# It use gomockhandler.json on /playground dir. So you may have to run make playground first.
.PHONY:benchmark-gen
benchmark-gen:
	cd playground/benchmark ;\
	echo "gomockhandler";\
	time ../../gomockhandler -config=gomockhandler.json -f mockgen ;\
	sleep 30s;\
	echo "mockgen + go generate";\
	cd generator ;\
    time go generate ./...; \
	rm ../../../gomockhandler ;

# test on playground/benchmark
# build gomockhandler and check mocks with it.
# It use gomockhandler.json on /playground dir. So you may have to run make playground first.
.PHONY:benchmark-check
benchmark-check:
	go build . ;\
	cd playground/benchmark ;\
	time ../../gomockhandler -config=gomockhandler.json check ;\
	rm ../../gomockhandler ;

# clean up playground/benchmark
.PHONY:benchmark-clean
benchmark-clean:
	cd playground/benchmark ;\
	rm gomockhandler.json;\
	rm -rf mock*
