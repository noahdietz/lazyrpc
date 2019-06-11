test:
	go install github.com/noahdietz/lazyrpc/cmd/protoc-gen-lazyrpc
	./cmd/protoc-gen-lazyrpc/test.sh

test-update:
	go install github.com/noahdietz/lazyrpc/cmd/protoc-gen-lazyrpc
	./cmd/protoc-gen-lazyrpc/test.sh -update_golden

annotations:
	./compile_protos.sh
	go install ./config

clean:
	rm -rf ./cmd/protoc-gen-lazyrpc/example/todo/todo_service.proto
	rm -rf ./protos/api-common-protos