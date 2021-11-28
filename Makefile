.ONESHELL:
all:
	@mkdir -p bin
	@cd ./src/rouphc
	@go build .
	@cd ../../
	@mv src/rouphc/rouphc bin
	@cd ./src/rouph
	@go build .
	@cd ../../
	@mv src/rouph/rouph bin
clean:
	rm -f *.o *.s test
