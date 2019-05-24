// Copyright 2019 Daniel Harrison. All Rights Reserved.

package generate

import (
	"flag"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/golang/protobuf/proto"
)

var rewrite = flag.Bool("rewrite", false, "WIP")

func TestGolden(t *testing.T) {
	goldenDir := filepath.Join(`..`, `golden`)
	reqBytes, err := ioutil.ReadFile(filepath.Join(goldenDir, `raft.parsed`))
	if err != nil {
		t.Fatal(err)
	}
	goldenPath := filepath.Join(goldenDir, `raftzeropb`, `raft.zeropb.go`)
	golden, err := ioutil.ReadFile(goldenPath)
	if err != nil {
		t.Fatal(err)
	}

	var req CodeGeneratorRequest
	if err := proto.UnmarshalText(string(reqBytes), &req); err != nil {
		t.Fatal(err)
	}
	res, err := Generate(&req)
	if err != nil {
		t.Fatal(err)
	}
	if res.Error != nil {
		t.Fatal(err)
	}
	if len(res.File) != 1 || res.File[0] == nil {
		t.Fatal("expected exactly 1 output file")
	}
	file := *res.File[0]
	if name, expected := file.GetName(), `raft.zeropb.go`; name != expected {
		t.Fatalf("expected name to be %s got %s", expected, name)
	}
	if content := file.GetContent(); content != string(golden) {
		if *rewrite {
			_ = ioutil.WriteFile(goldenPath, []byte(content), 0644)
		}
		t.Fatalf("did not match golden:\n\n%s\n\ngot:\n\n%s\n", golden, content)
	}
}
