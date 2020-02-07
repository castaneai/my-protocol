package my_protocol

import (
	"encoding/binary"

	"golang.org/x/text/transform"
)

type PacketUnpacker struct{ transform.NopResetter }

func (p *PacketUnpacker) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) {
	if atEOF && len(src) == 0 {
		return
	}
	if len(src) < 2 {
		err = transform.ErrShortSrc
		return
	}
	size := int(binary.LittleEndian.Uint16(src[:2]))
	if len(src[2:]) < size {
		err = transform.ErrShortSrc
		return
	}
	nDst = copy(dst, src[2:2+size])
	nSrc = 2 + nDst
	if nDst < size {
		err = transform.ErrShortDst
	}
	return
}
