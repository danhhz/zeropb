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
	res, err := generate.Generate(&req)
	if err != nil {
		return err
	}
	resBytes, err := proto.Marshal(res)
	if err != nil {
		return err
	}
	if _, err := os.Stdout.Write(resBytes); err != nil {
		return err
	}
	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
	}
}
