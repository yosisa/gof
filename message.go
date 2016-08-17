package gof

import (
	"encoding/binary"

	"github.com/hkwi/gopenflow/ofp4"
)

type FlowMod struct {
	Cookie       uint64
	CookieMask   uint64
	TableID      uint8
	Command      uint8
	IdleTimeout  uint16
	HardTimeout  uint16
	Priority     uint16
	BufferID     uint32
	OutPort      uint32
	OutGroup     uint32
	Flags        uint16
	Matches      []MatchMarshaler
	Instructions []InstructionMarshaler
}

func (v *FlowMod) Marshal() []byte {
	var match []byte
	for _, x := range v.Matches {
		match = append(match, x.MarshalMatch()...)
	}
	match = ofp4.MakeMatch(match)

	var inst []byte
	for _, x := range v.Instructions {
		inst = append(inst, x.MarshalInstruction()...)
	}
	bufferID := v.BufferID
	if bufferID == 0 {
		bufferID = ofp4.OFP_NO_BUFFER
	}

	b := make([]byte, 8+40+len(match)+len(inst))
	hdr, msg := b[:8], b[8:]
	hdr[0] = 4
	hdr[1] = ofp4.OFPT_FLOW_MOD
	binary.BigEndian.PutUint16(hdr[2:], uint16(len(b)))

	binary.BigEndian.PutUint64(msg, v.Cookie)
	binary.BigEndian.PutUint64(msg[8:], v.CookieMask)
	msg[16] = v.TableID
	msg[17] = v.Command
	binary.BigEndian.PutUint16(msg[18:], v.IdleTimeout)
	binary.BigEndian.PutUint16(msg[20:], v.HardTimeout)
	binary.BigEndian.PutUint16(msg[22:], v.Priority)
	binary.BigEndian.PutUint32(msg[24:], bufferID)
	binary.BigEndian.PutUint32(msg[28:], v.OutPort)
	binary.BigEndian.PutUint32(msg[32:], v.OutGroup)
	binary.BigEndian.PutUint16(msg[36:], v.Flags)

	n := copy(msg[40:], match)
	copy(msg[40+n:], inst)
	return b
}

type PacketOut struct {
	BufferID uint32
	InPort   uint32
	Actions  []ActionMarshaler
	Data     []byte
}

func (v *PacketOut) Marshal() []byte {
	var actions []byte
	for _, x := range v.Actions {
		actions = append(actions, x.MarshalAction()...)
	}

	b := make([]byte, 8+16+len(actions)+len(v.Data))
	hdr, msg := b[:8], b[8:]
	hdr[0] = 4
	hdr[1] = ofp4.OFPT_PACKET_OUT
	binary.BigEndian.PutUint16(hdr[2:], uint16(len(b)))

	binary.BigEndian.PutUint32(msg, v.BufferID)
	binary.BigEndian.PutUint32(msg[4:], v.InPort)
	binary.BigEndian.PutUint16(msg[8:], uint16(len(actions)))
	n := copy(msg[16:], actions)
	copy(msg[16+n:], v.Data)
	return b
}

type MultipartRequest struct {
	Type  uint16
	Flags uint16
	Body  []byte
}

func (v *MultipartRequest) Marshal() []byte {
	b := make([]byte, 8+8+len(v.Body))
	hdr, msg := b[:8], b[8:]
	hdr[0] = 4
	hdr[1] = ofp4.OFPT_MULTIPART_REQUEST
	binary.BigEndian.PutUint16(hdr[2:], uint16(len(b)))

	binary.BigEndian.PutUint16(msg, v.Type)
	binary.BigEndian.PutUint16(msg[2:], v.Flags)
	copy(msg[8:], v.Body)
	return b
}
