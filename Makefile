.PHONY: all
all:
	rm -f bin/*
	$(MAKE) -C sgx
	cp sgx/bin/enclave.signed.so test/
	cp sgx/bin/enclave.signed.so netdisk/
	sync
	sudo $(MAKE) -C test

.PHONY: cleandb
cleandb:
	rm -rf test/idx netdisk/idx
	sync

.PHONY: clean
clean:
	$(MAKE) -C sgx clean
	$(MAKE) -C test clean
	rm -f netdisk/enclave.signed.so  netdisk/netdisk
	rm -rf netdisk/idx
	sync

