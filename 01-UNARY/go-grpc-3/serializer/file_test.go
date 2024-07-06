package serializer_test

import (
	"grpc-3/sample"
	"grpc-3/serializer"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFileSerializer(t *testing.T) {
	t.Parallel()

	binaryFile := "../tmp/laptop.bin"

	laptop1 := sample.NewLaptop()
	err := serializer.WriteProtobuffToBinaryFile(laptop1, binaryFile)

	require.NoError(t, err)
}