.PHONY: gen build run clean get_result

gen:
	protoc --go_out=. static/*.proto

build: gen
	mkdir -p build
	cd src/serializer && go build -o ../../build/serializer
	cd src/proxy && go build -o ../../build/proxy

run-serializer:
	build/serializer

run-proxy:
	build/proxy

clean:
	rm -rf gen
	rm -rf build

get_result:
	echo get_result | nc localhost 8080

docker-build:
	docker build -t serializer-test .

docker-run:
	docker run -p 8081:8080 -e S_FORMAT=${S_FORMAT} serializer-test
