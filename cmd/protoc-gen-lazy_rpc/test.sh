#!/bin/bash

TODO_DIR=$GOPATH/src/github.com/noahdietz/lazyrpc/cmd/protoc-gen-lazy_rpc/example/todo

protoc -I $TODO_DIR --lazy_rpc_out $TODO_DIR $TODO_DIR/todo.proto