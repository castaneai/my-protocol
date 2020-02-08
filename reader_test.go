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

	p := make([]byte, 10)
	n, err := pr.Read(p)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "hello", string(p[:n]))
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

	p := make([]byte, 10)
	n, err := pr.Read(p)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "hello", string(p[:n]))
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

func TestPacketUnpacker_Transform_WithShortSource(t *testing.T) {
	r := bytes.NewReader([]byte{0x05, 0x00, 0x68, 0x65})
	br := &shortReader{r: r, size: 8}
	pr := transform.NewReader(br, &PacketUnpacker{})

	p := make([]byte, 1024)
	_, err := pr.Read(p)
	assert.Equal(t, transform.ErrShortSrc, err)
}

func TestEncryptedPacketUnpacker_Transform(t *testing.T) {
	plaintext := padPKCS7([]byte("hello"))
	ciphertext := make([]byte, len(plaintext))
	cip := newAES256WithRandomKey()
	cip.Encrypt(ciphertext, plaintext)

	r := bytes.NewReader(ciphertext)
	pr := transform.NewReader(r, &EncryptedPacketUnpacker{cip: cip})

	p := make([]byte, 10)
	n, err := pr.Read(p)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "hello", string(p[:n]))
}

func TestChainTransformers(t *testing.T) {
	plaintext := padPKCS7([]byte("hello"))
	ciphertext := make([]byte, len(plaintext))
	cip := newAES256WithRandomKey()
	cip.Encrypt(ciphertext, plaintext)

	var b bytes.Buffer
	w := transform.NewWriter(&b, &PacketPacker{})
	if _, err := w.Write(ciphertext); err != nil {
		t.Fatal(err)
	}

	tr := transform.Chain(&PacketUnpacker{}, &EncryptedPacketUnpacker{cip: cip})
	pr := transform.NewReader(&b, tr)

	p := make([]byte, 10)
	n, err := pr.Read(p)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "hello", string(p[:n]))
}
