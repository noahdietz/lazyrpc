#!/bin/bash

TODO_DIR=$GOPATH/src/github.com/noahdietz/lazyrpc/cmd/protoc-gen-lazyrpc/example/todo
UPDATE=false

if [ "$1" = "-update_golden" ]; then
    UPDATE=true
fi

protoc -I protos -I $TODO_DIR --lazyrpc_out $TODO_DIR $TODO_DIR/todo.proto

if $UPDATE ; then
    cp $TODO_DIR/todo_service.proto $TODO_DIR/todo_service.want
else
    diff $TODO_DIR/todo_service.proto $TODO_DIR/todo_service.want
fi