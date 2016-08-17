package gof

import (
	"encoding/binary"

	"github.com/hkwi/gopenflow/ofp4"
)

type InstructionMarshaler interface {
	MarshalInstruction() []byte
}

func Instructions(v ...InstructionMarshaler) []InstructionMarshaler {
	return v
}

func GotoTable(id uint8) *InstGotoTable {
	return &InstGotoTable{TableID: id}
}

type InstGotoTable struct {
	TableID uint8
}

func (v *InstGotoTable) MarshalInstruction() []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint16(b, ofp4.OFPIT_GOTO_TABLE)
	binary.BigEndian.PutUint16(b[2:], 8)
	b[4] = v.TableID
	return b
}

func ApplyActions(actions ...ActionMarshaler) *InstApplyActions {
	return &InstApplyActions{actions}
}

type InstApplyActions struct {
	Actions []ActionMarshaler
}

func (v *InstApplyActions) MarshalInstruction() []byte {
	var actions []byte
	for _, a := range v.Actions {
		actions = append(actions, a.MarshalAction()...)
	}

	b := make([]byte, 8+len(actions))
	binary.BigEndian.PutUint16(b, ofp4.OFPIT_APPLY_ACTIONS)
	binary.BigEndian.PutUint16(b[2:], uint16(len(b)))
	copy(b[8:], actions)
	return b
}
