test:
	go install github.com/noahdietz/lazyrpc/cmd/protoc-gen-lazyrpc
	./cmd/protoc-gen-lazyrpc/test.sh

annotations:
	./compile_protos.sh

clean:
	rm -rf ./cmd/protoc-gen-lazyrpc/example/todo/todo_service.proto
	rm -rf ./protos/api-common-protos