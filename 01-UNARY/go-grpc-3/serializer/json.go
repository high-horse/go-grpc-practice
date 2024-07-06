package serializer

import (
	"fmt"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func ProtobuffToJSON(message proto.Message) (string, error) {

	options := protojson.MarshalOptions{
		UseEnumNumbers: false, 
		EmitUnpopulated: false,
		UseProtoNames: true,
		Indent: "  ",
	}

	marshalled, err := options.Marshal(message)
	if err != nil {
		return "", fmt.Errorf("cannot marshal proto message to json: %w", err)
	}
	

	return string(marshalled), nil
}

