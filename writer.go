package my_protocol

import (
	"crypto/cipher"
	"encoding/binary"

	"golang.org/x/text/transform"
)

type PacketPacker struct{ transform.NopResetter }

func (p *PacketPacker) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) {
	nSrc = len(src)
	binary.LittleEndian.PutUint16(dst, uint16(nSrc))
	nDst = copy(dst[2:], src)
	if nDst < nSrc {
		err = transform.ErrShortDst
	}
	nDst += 2
	return
}

type EncryptedPacketPacker struct {
	transform.NopResetter
	cip cipher.Block
}

func (p *EncryptedPacketPacker) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) {
	nSrc = len(src)
	psrc := padPKCS7(src)
	p.cip.Encrypt(dst, psrc)
	nDst = len(psrc)
	return
}
