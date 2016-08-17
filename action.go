package gof

import (
	"encoding/binary"

	"github.com/hkwi/gopenflow/ofp4"
)

type ActionMarshaler interface {
	MarshalAction() []byte
}

func Actions(v ...ActionMarshaler) []ActionMarshaler {
	return v
}

func Output(port uint32) *ActionOutput {
	return &ActionOutput{Port: port}
}

type ActionOutput struct {
	Port   uint32
	MaxLen uint16
}

func (v *ActionOutput) MarshalAction() []byte {
	maxlen := v.MaxLen
	if maxlen == 0 {
		maxlen = ofp4.OFPCML_MAX
	}

	b := make([]byte, 16)
	binary.BigEndian.PutUint16(b, ofp4.OFPAT_OUTPUT)
	binary.BigEndian.PutUint16(b[2:], uint16(len(b)))
	binary.BigEndian.PutUint32(b[4:], v.Port)
	binary.BigEndian.PutUint16(b[8:], maxlen)
	return b
}

func SetField(oxm MatchMarshaler) *ActionSetField {
	return &ActionSetField{OXM: oxm}
}

type ActionSetField struct {
	OXM MatchMarshaler
}

func (v *ActionSetField) MarshalAction() []byte {
	oxm := v.OXM.MarshalMatch()
	b := make([]byte, align8(4+len(oxm)))
	binary.BigEndian.PutUint16(b, ofp4.OFPAT_SET_FIELD)
	binary.BigEndian.PutUint16(b[2:], uint16(len(b)))
	copy(b[4:], oxm)
	return b
}

func align8(n int) int {
	return (n + 7) / 8 * 8
}
