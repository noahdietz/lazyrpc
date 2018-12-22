#!/bin/bash

TODO_DIR=$GOPATH/src/github.com/noahdietz/lazyrpc/cmd/protoc-gen-lazyrpc/example/todo

protoc -I protos -I $TODO_DIR --lazyrpc_out $TODO_DIR $TODO_DIR/todo.proto