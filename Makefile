test:
	go install github.com/noahdietz/lazyrpc/cmd/protoc-gen-lazy_rpc
	./cmd/protoc-gen-lazy_rpc/test.sh

annotations:
	./compile_protos.sh

clean:
	rm -rf ./cmd/protoc-gen-lazy_rpc/example/todo/todo_service.proto
	rm -rf ./protos/api-common-protos