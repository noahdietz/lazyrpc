// Copyright 2018 Google
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package lazyrpc

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/golang/protobuf/proto"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/builder"
	"github.com/jhump/protoreflect/desc/protoprint"
	annotations "github.com/noahdietz/lazyrpc/config"
)

var keyRegex *regexp.Regexp
var methodRegex *regexp.Regexp

// Generate service(s) & method(s) for plain ol'Messages
func Generate(req *plugin.CodeGeneratorRequest) (*plugin.CodeGeneratorResponse, error) {
	var res plugin.CodeGeneratorResponse

	keyRegex = regexp.MustCompile(`Key:(\s*[a-z]+[^\s]*)`)
	methodRegex = regexp.MustCompile(`Methods:(\s*[a-z,]+[^\s]*)`)

	d, err := desc.CreateFileDescriptors(req.GetProtoFile())
	if err != nil {
		return nil, err
	}

	for _, name := range req.GetFileToGenerate() {
		rich, ok := d[name]
		if !ok {
			return nil, fmt.Errorf("FileToGenerate (%s) did not have a rich descriptor", name)
		}

		f, err := generate(rich)
		if err != nil {
			return nil, err
		}

		s, err := format(f)
		if err != nil {
			return nil, err
		}

		add(&res, f.GetName(), s)
	}

	return &res, nil
}

func generate(orig *desc.FileDescriptor) (*desc.FileDescriptor, error) {
	f, err := file(orig)
	if err != nil {
		return nil, err
	}

	for _, m := range orig.GetMessageTypes() {
		srv, msgs, err := build(m)
		if err != nil {
			return nil, err
		} else if srv == nil {
			continue
		}

		f.AddService(srv)

		for _, m := range msgs {
			f.AddMessage(m)
		}
	}

	return f.Build()
}

func file(orig *desc.FileDescriptor) (*builder.FileBuilder, error) {
	n := orig.GetName()
	f := builder.NewFile(n).SetProto3(true).SetPackageName(orig.GetPackage())

	if ndx := strings.LastIndex(n, "/"); ndx != -1 {
		n = n[ndx+1:]
	}

	if ndx := strings.LastIndex(n, "."); ndx != -1 {
		n = n[:ndx]
	}

	n += "_service.proto"

	return f.SetName(n), nil
}

func build(msg *desc.MessageDescriptor) (*builder.ServiceBuilder, []*builder.MessageBuilder, error) {
	srv := builder.NewService(msg.GetName() + "Service")
	msgs := []*builder.MessageBuilder{}
	opts := msg.GetOptions()

	if opts == nil {
		return nil, nil, nil
	}

	ext, err := proto.GetExtension(opts, annotations.E_Config)
	if err == proto.ErrMissingExtension {
		return nil, nil, nil
	} else if err != nil {
		return nil, nil, err
	}
	config := ext.(*annotations.Config)

	methods := config.GetMethods()
	if len(methods) == 0 {
		return nil, nil, nil
	}

	key := config.GetKey()

	og, err := builder.FromMessage(msg)
	if err != nil {
		return nil, nil, err
	}

	for _, m := range methods {
		var in, out *builder.RpcType

		verb := strings.Title(strings.ToLower(m.String()))
		name := verb + msg.GetName()

		switch m {
		case annotations.Method_UPDATE:
			fallthrough //update & create are almost the same
		case annotations.Method_CREATE:
			in = builder.RpcTypeMessage(og, false)
			out = in
		case annotations.Method_GET:
			out = builder.RpcTypeMessage(og, false)
			fallthrough // get & delete are almost the same
		case annotations.Method_DELETE:
			in = builder.RpcTypeMessage(og, false)
			if key != "" {
				ogF, err := builder.FromField(msg.FindFieldByName(key))
				if err != nil {
					return nil, nil, err
				}
				tmpF := builder.NewField(ogF.GetName(), ogF.GetType())

				tmpM := builder.NewMessage(name + "Request").AddField(tmpF)
				msgs = append(msgs, tmpM)
				in = builder.RpcTypeMessage(tmpM, false)
			}

			out = builder.RpcTypeMessage(og, false)
		case annotations.Method_LIST:
			tmpF := builder.NewField("content", builder.FieldTypeMessage(og)).SetRepeated()

			o := builder.
				NewMessage(name + "Response").
				AddField(builder.NewField("next_page_token", builder.FieldTypeString())).
				AddField(tmpF)
			msgs = append(msgs, o)
			out = builder.RpcTypeMessage(o, false)

			i := builder.
				NewMessage(name + "Request").
				AddField(builder.NewField("page_size", builder.FieldTypeInt32())).
				AddField(builder.NewField("page_token", builder.FieldTypeString()))
			msgs = append(msgs, i)
			in = builder.RpcTypeMessage(i, false)
		default:
			log.Printf("Found method value: %v\n", m)
			continue
		}

		srv.AddMethod(builder.NewMethod(name, in, out))
	}

	return srv, msgs, nil
}

func format(f *desc.FileDescriptor) (string, error) {
	p := protoprint.Printer{}

	return p.PrintProtoToString(f)
}

func add(res *plugin.CodeGeneratorResponse, name, data string) {
	file := &plugin.CodeGeneratorResponse_File{
		Name: proto.String(name),
	}

	file.Content = proto.String(data)

	res.File = append(res.File, file)
}
