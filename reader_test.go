package my_protocol

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"

	"golang.org/x/text/transform"
)

func TestPacketUnpacker_Transform(t *testing.T) {
	r := bytes.NewReader([]byte{0x05, 0x00, 0x68, 0x65, 0x6c, 0x6c, 0x6f})
	pr := transform.NewReader(r, &PacketUnpacker{})
	var b bytes.Buffer
	if _, err := io.Copy(&b, pr); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "hello", string(b.Bytes()))
}
