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

type shortReader struct {
	r    io.Reader
	size int
}

func (r *shortReader) Read(p []byte) (n int, err error) {
	return r.r.Read(p[:r.size])
}

func TestPacketUnpacker_Transform_WithShortBuffer(t *testing.T) {
	r := bytes.NewReader([]byte{0x05, 0x00, 0x68, 0x65, 0x6c, 0x6c, 0x6f})
	br := &shortReader{r: r, size: 1}
	pr := transform.NewReader(br, &PacketUnpacker{})
	var b bytes.Buffer
	if _, err := io.Copy(&b, pr); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "hello", string(b.Bytes()))
}

func TestPacketUnpacker_Transform_WithMultiPacket(t *testing.T) {
	r := bytes.NewReader([]byte{0x05, 0x00, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x05, 0x00, 0x77, 0x6f, 0x72, 0x6c, 0x64})
	br := &shortReader{r: r, size: 8}
	pr := transform.NewReader(br, &PacketUnpacker{})

	{
		p1 := make([]byte, 1024)
		n, err := pr.Read(p1)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, "hello", string(p1[:n]))
	}

	{
		p2 := make([]byte, 1024)
		n, err := pr.Read(p2)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, "world", string(p2[:n]))
	}
}
