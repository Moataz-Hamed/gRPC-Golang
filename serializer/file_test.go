package serializer

import (
	"testing"

	"github.com/moataz-hamed/pb/pb"
	"github.com/moataz-hamed/sample"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func TestFileSerializer(t *testing.T) {
	t.Parallel()

	binaryFile := "../tmp/laptop.bin"
	jsonFile := "../tmp/laptop.json"

	laptop := sample.NewLaptop()
	err := WriteProtobufToBinaryFile(laptop, binaryFile)
	require.NoError(t, err)

	laptop2 := &pb.Laptop{}
	err = ReadProtobufToBinaryFile(binaryFile, laptop2)
	require.NoError(t, err)

	require.True(t, proto.Equal(laptop, laptop2)) //check if they are the same

	err = WriteProtobufToJSONFile(laptop, jsonFile)
	require.NoError(t, err)
	
}
