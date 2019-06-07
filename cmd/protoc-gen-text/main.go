// Copyright 2019 Daniel Harrison. All Rights Reserved.

package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/danhhz/zeropb/generate"
	"github.com/golang/protobuf/proto"
)

func run() error {
	reqBytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return err
	}
	var req generate.CodeGeneratorRequest
	if err := proto.Unmarshal(reqBytes, &req); err != nil {
		return err
	}
	if err := proto.MarshalText(os.Stderr, &req); err != nil {
		return err
	}
	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
	}
}
