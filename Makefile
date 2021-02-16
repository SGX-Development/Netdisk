.PHONY: all
all:
	rm -f enclave.signed.so
	$(MAKE) -C sgx
	cp sgx/bin/enclave.signed.so test/
	cp sgx/bin/enclave.signed.so netdisk/
	sync
	sudo $(MAKE) -C test

.PHONY: cleandb
cleandb:
	rm -rf test/idx
	sync

.PHONY: clean
clean:
	$(MAKE) -C sgx clean
	$(MAKE) -C test clean
	rm -f enclave.signed.so tantivy-sgx tantivy-sgx-part
	rm -rf idx
	sync

