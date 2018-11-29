*Disclaimer: this is just a thing I slapped together, this is not supported and is just for fun. Try it out though!*

# lazyrpc

This is a protoc plugin that generates a basic CRUD service proto for messages of the input protos.
I made this because I was tired of writing the same boilerplate request/response messages and service methods for simple APIs.

I also wanted to play around with the rich descriptor facilities of [protoreflect](https://github.com/jhump/protoreflect)
in the context of building a protoc plugin.

## Using lazyrpc

### Install

You can either install the plugin directly:

```
$ go get github.com/noahdietz/lazyrpc/cmd/protoc-gen-lazy_rpc
```

Or build it from source:

```
$ cd $GOPATH/src
$ mkdir -p github.com/noahdietz
$ cd github.com/noahdietz
$ git clone https://github.com/noahdietz/lazyrpc.git
$ cd lazyrpc/cmd/protoc-gen-lazy_rpc
$ go get
$ go install 
```

### Configuration Comments

To keep things simple, but still somewhat configurable, `lazyrpc` takes message level comments as configuration.
The config phrases can be anywhere in the leading comments of a message and are as follows:

* Methods: a comma-delimited list of methods to generate for this message's service
  * Supported methods are: `create,get,list,update,delete` (this is valid example)
* Key: a single field name for use as a key in get & delete methods.
  * Omission of this phrase results in the entire message being used as input

Example:

```
// Methods: create,get,list,update,delete
// Key: id
message Todo {
  int32 id = 1;
  string task = 2;
  bool done = 3;
}
```

A message without a `Methods` phrase will be ignored entirely.

### Invocation

Invoke the plugin like any other protoc plugin:

```
$ protoc -I my/protos/ --lazy_rpc_out my/protos/ my/protos/a.proto
```

### Output

A proto file is created for each input proto, containing all of the services & request/response messages
for the messages defined & configured in the input.

This generated proto file takes the name `{input file}_service.proto`. For example,
using `todo.proto` as input, a `todo_service.proto` file is generated.

*Note: currently, files will be overwritten*

## Looking forward

In the future, `lazyrpc` could support:

* configuration for streaming method generation
* plugin option to signal merging of services/messages as an alternative to overwritting existing files
