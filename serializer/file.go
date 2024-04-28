package serializer

import (
	"fmt"
	"os"

	"google.golang.org/protobuf/proto"
)

func WriteProtobufToJSONFile(message proto.Message, fileName string) error {
	data, err := ProtobufToJSON(message)
	if err != nil {
		return fmt.Errorf("Can't marshal proto message to json,%w", err)
	}

	err = os.WriteFile(fileName, []byte(data), 0644)
	if err != nil {
		return fmt.Errorf("Can't write to file,%w", err)
	}
	return nil
}

func WriteProtobufToBinaryFile(message proto.Message, fileName string) error {
	data, err := proto.Marshal(message)
	if err != nil {
		return fmt.Errorf("Can't marshal message,%w", err)
	}
	err = os.WriteFile(fileName, data, 0644)
	if err != nil {
		return fmt.Errorf("can't write binary to file,%w", err)
	}

	return nil
}

func ReadProtobufToBinaryFile(fileName string, message proto.Message) error {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("can't read file, %w", err)
	}

	err = proto.Unmarshal(data, message)
	if err != nil {
		return fmt.Errorf("Can't unmarshal data,%w", err)
	}

	return nil
}
