package my_protocol

import (
	"encoding/binary"

	"golang.org/x/text/transform"
)

type PacketUnpacker struct {
	rest []byte
}

func (p *PacketUnpacker) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) {
	npRest := len(p.rest)
	if npRest > 0 {
		src = append(p.rest, src...)
		p.rest = nil
	}
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
	nSrc = 2 + nDst - npRest
	if nDst < size {
		err = transform.ErrShortDst
		return
	}
	nRest := len(src[2+size:])
	if nRest > 0 {
		p.rest = make([]byte, nRest)
		nSrc += copy(p.rest, src[2+size:])
	}
	return
}

func (p *PacketUnpacker) Reset() {
	p.rest = nil
}
