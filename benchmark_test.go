package my_protocol

import (
	"bytes"
	"crypto/cipher"
	"encoding/binary"
	"io"
	"testing"

	"golang.org/x/text/transform"
)

func BenchmarkWithTransformer(b *testing.B) {
	payload := []byte("hello,world")
	cip := newAES256WithRandomKey()
	wtr := transform.Chain(&EncryptedPacketPacker{cip: cip}, &PacketPacker{})
	tr := transform.Chain(&PacketUnpacker{}, &EncryptedPacketUnpacker{cip: cip})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		w := transform.NewWriter(&buf, wtr)
		if _, err := w.Write(payload); err != nil {
			b.Fatal(err)
		}
		pr := transform.NewReader(&buf, tr)
		rp := make([]byte, 20)
		if _, err := pr.Read(rp); err != nil {
			b.Fatal(err)
		}
	}
}

// without transformer impl
func readPacket(r io.Reader) ([]byte, error) {
	var size uint16
	if err := binary.Read(r, binary.LittleEndian, &size); err != nil {
		return nil, err
	}
	p := make([]byte, size)
	if _, err := io.ReadFull(r, p); err != nil {
		return nil, err
	}
	return p, nil
}

func writePacket(w io.Writer, p []byte) error {
	if err := binary.Write(w, binary.LittleEndian, uint16(len(p))); err != nil {
		return err
	}
	if _, err := w.Write(p); err != nil {
		return err
	}
	return nil
}

func encryptPacket(p []byte, cip cipher.Block) []byte {
	pp := padPKCS7(p)
	enc := make([]byte, len(pp))
	cip.Encrypt(enc, pp)
	return enc
}

func decryptPacket(p []byte, cip cipher.Block) []byte {
	dec := make([]byte, len(p))
	cip.Decrypt(dec, p)
	return unpadPKCS7(dec)
}

func BenchmarkWithoutTransformer(b *testing.B) {
	payload := []byte("hello,world")
	cip := newAES256WithRandomKey()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		p := encryptPacket(payload, cip)
		if err := writePacket(&buf, p); err != nil {
			b.Fatal(err)
		}
		rpenc, err := readPacket(&buf)
		if err != nil {
			b.Fatal(err)
		}
		_ = decryptPacket(rpenc, cip)
	}
}
