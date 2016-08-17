package nxm

import (
	"encoding/binary"

	"github.com/hkwi/gopenflow/ofp4"
)

func NXRegMove(src, dst uint32, nbits uint16) *NXActionRegMove {
	return &NXActionRegMove{
		Nbits:    nbits,
		SrcField: src,
		DstField: dst,
	}
}

type NXActionRegMove struct {
	Nbits    uint16
	SrcOfs   uint16
	DstOfs   uint16
	SrcField uint32
	DstField uint32
}

func (v *NXActionRegMove) MarshalAction() []byte {
	b := make([]byte, 24)
	binary.BigEndian.PutUint16(b, ofp4.OFPAT_EXPERIMENTER)
	binary.BigEndian.PutUint16(b[2:], uint16(len(b)))
	binary.BigEndian.PutUint32(b[4:], NX_EXPERIMENTER_ID)
	binary.BigEndian.PutUint16(b[8:], NXAST_REG_MOVE)
	binary.BigEndian.PutUint16(b[10:], v.Nbits)
	binary.BigEndian.PutUint16(b[12:], v.SrcOfs)
	binary.BigEndian.PutUint16(b[14:], v.DstOfs)
	binary.BigEndian.PutUint32(b[16:], v.SrcField)
	binary.BigEndian.PutUint32(b[20:], v.DstField)
	return b
}
