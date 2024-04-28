package serializer

import (
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func ProtobufToJSON(message proto.Message) (string, error) {
	marshaler := protojson.MarshalOptions{
		Indent:          " ",
		UseProtoNames:   true,
		EmitUnpopulated: true,
		UseEnumNumbers:  false,
	}

	b, err := marshaler.Marshal(message)
	return string(b), err
}
