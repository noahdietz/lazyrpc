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

### Configuration Options

The configuration of `lazyrpc` generation is defined using proto Message options.
The option available options are as follows:

* Methods: a comma-delimited list of methods to generate for this message's service
  * Supported methods are: `"create", "get", "list", "update", "delete"`
* Key: a single field name for use as a key in get & delete methods.
  * Omission of this phrase results in the entire message being used as input

Example:

```
import "lazy.proto"

...

message Todo {
  option (lazy.config) = {
    methods: ["create", "get", "list", "update", "delete"]
    key: "id"
  };
  
  int32 id = 1;
  string task = 2;
  bool done = 3;
}
```

A message without a `lazy.config.methods` field will be ignored entirely.

The extension defintion can be found in [proto/lazy.proto](./protos/lazy.proto).

The compiled Go definition can be imported as `github.com/noahdietz/lazyrpc/config` and found in [config/lazy.pb.go](./config/lazy.pb.go).

### Invocation

Invoke the plugin like any other protoc plugin:

```
export LAZY_RPC_PROTOS=$GOPATH/src/github.com/noahdietz/lazyrpc/protos
$ protoc -I $LAZY_RPC_PROTOS -I my/protos/ --lazy_rpc_out my/protos/ my/protos/a.proto
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
