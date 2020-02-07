package my_protocol

import (
	"encoding/binary"

	"golang.org/x/text/transform"
)

type PacketPacker struct{}

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

func (p *PacketPacker) Reset() {}
