package bostrom

import "github.com/cosmos/gogoproto/proto"

type (
	protoType struct {
		t    proto.Message
		name string
	}
	protoFile struct {
		filename   string
		descriptor []byte
	}
)

func descriptor(t interface{ Descriptor() ([]byte, []int) }) []byte {
	desc, _ := t.Descriptor()
	return desc
}

func registerTypes(types []protoType) {
	for _, t := range types {
		proto.RegisterType(t.t, t.name)
	}
}

func registerFiles(files []protoFile) {
	for _, f := range files {
		proto.RegisterFile(f.filename, f.descriptor)
	}
}
