package my_protocol

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/transform"
)

func TestPacketPacket_Transform(t *testing.T) {
	var b bytes.Buffer
	w := transform.NewWriter(&b, &PacketPacker{})
	if _, err := io.Copy(w, bytes.NewReader([]byte("hello"))); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, []byte{0x05, 0x00, 0x68, 0x65, 0x6c, 0x6c, 0x6f}, b.Bytes())
}
