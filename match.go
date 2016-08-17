package gof

import (
	"encoding/binary"
	"net"

	"github.com/hkwi/gopenflow/ofp4"
	"github.com/yosisa/gof/nxm"
)

type MatchMarshaler interface {
	MarshalMatch() []byte
}

type RawMatch []byte

func (v RawMatch) MarshalMatch() []byte {
	return v
}

func Matches(ms ...MatchMarshaler) []MatchMarshaler {
	return ms
}

func InPort(port uint32) *MatchInPort {
	return &MatchInPort{port}
}

type MatchInPort struct {
	Port uint32
}

func (v *MatchInPort) MarshalMatch() []byte {
	b := make([]byte, 4+4)
	binary.BigEndian.PutUint16(b, ofp4.OFPXMC_OPENFLOW_BASIC)
	b[2] = ofp4.OFPXMT_OFB_IN_PORT << 1
	b[3] = 4
	binary.BigEndian.PutUint32(b[4:], v.Port)
	return b
}

func EthDst(addr net.HardwareAddr) *MatchEthDst {
	return &MatchEthDst{Addr: addr}
}

type MatchEthDst struct {
	Addr net.HardwareAddr
}

func (v *MatchEthDst) MarshalMatch() []byte {
	b := make([]byte, 4+6)
	binary.BigEndian.PutUint16(b, ofp4.OFPXMC_OPENFLOW_BASIC)
	b[2] = ofp4.OFPXMT_OFB_ETH_DST << 1
	b[3] = 6
	copy(b[4:], v.Addr)
	return b
}

func EthSrc(addr net.HardwareAddr) *MatchEthSrc {
	return &MatchEthSrc{Addr: addr}
}

type MatchEthSrc struct {
	Addr net.HardwareAddr
}

func (v *MatchEthSrc) MarshalMatch() []byte {
	b := make([]byte, 4+6)
	binary.BigEndian.PutUint16(b, ofp4.OFPXMC_OPENFLOW_BASIC)
	b[2] = ofp4.OFPXMT_OFB_ETH_SRC << 1
	b[3] = 6
	copy(b[4:], v.Addr)
	return b
}

const (
	EthTypeIP  = 0x0800
	EthTypeARP = 0x0806
)

func EthType(n uint16) *MatchEthType {
	return &MatchEthType{Type: n}
}

type MatchEthType struct {
	Type uint16
}

func (v *MatchEthType) MarshalMatch() []byte {
	b := make([]byte, 4+2)
	binary.BigEndian.PutUint16(b, ofp4.OFPXMC_OPENFLOW_BASIC)
	b[2] = ofp4.OFPXMT_OFB_ETH_TYPE << 1
	b[3] = 2
	binary.BigEndian.PutUint16(b[4:], v.Type)
	return b
}

func IPv4Src(ip net.IP) MatchMarshaler {
	return ipv4Src(ip)
}

type ipv4Src net.IP

func (v ipv4Src) MarshalMatch() []byte {
	b := makeOXM(ofp4.OFPXMC_OPENFLOW_BASIC, ofp4.OFPXMT_OFB_IPV4_SRC, 4)
	copy(b[4:], net.IP(v).To4())
	return b
}

func IPv4Dst(ip net.IP) MatchMarshaler {
	return ipv4Dst(ip)
}

type ipv4Dst net.IP

func (v ipv4Dst) MarshalMatch() []byte {
	b := makeOXM(ofp4.OFPXMC_OPENFLOW_BASIC, ofp4.OFPXMT_OFB_IPV4_DST, 4)
	copy(b[4:], net.IP(v).To4())
	return b
}

const (
	ARPOpRequest = 1
	ARPOpReply   = 2
)

func ARPOp(op uint16) MatchMarshaler {
	return arpOp(op)
}

type arpOp uint16

func (v arpOp) MarshalMatch() []byte {
	b := makeOXM(ofp4.OFPXMC_OPENFLOW_BASIC, ofp4.OFPXMT_OFB_ARP_OP, 2)
	binary.BigEndian.PutUint16(b[4:], uint16(v))
	return b
}

func ARPSpa(ip net.IP) MatchMarshaler {
	return arpSpa(ip)
}

type arpSpa net.IP

func (v arpSpa) MarshalMatch() []byte {
	b := makeOXM(ofp4.OFPXMC_OPENFLOW_BASIC, ofp4.OFPXMT_OFB_ARP_SPA, 4)
	copy(b[4:], net.IP(v).To4())
	return b
}

func ARPTpa(ip net.IP) MatchMarshaler {
	return arpTpa(ip)
}

type arpTpa net.IP

func (v arpTpa) MarshalMatch() []byte {
	b := makeOXM(ofp4.OFPXMC_OPENFLOW_BASIC, ofp4.OFPXMT_OFB_ARP_TPA, 4)
	copy(b[4:], net.IP(v).To4())
	return b
}

func ARPSha(addr net.HardwareAddr) MatchMarshaler {
	return arpSha(addr)
}

type arpSha net.HardwareAddr

func (v arpSha) MarshalMatch() []byte {
	b := makeOXM(ofp4.OFPXMC_OPENFLOW_BASIC, ofp4.OFPXMT_OFB_ARP_SHA, 6)
	copy(b[4:], v)
	return b
}

func ARPTha(addr net.HardwareAddr) MatchMarshaler {
	return arpTha(addr)
}

type arpTha net.IP

func (v arpTha) MarshalMatch() []byte {
	b := makeOXM(ofp4.OFPXMC_OPENFLOW_BASIC, ofp4.OFPXMT_OFB_ARP_THA, 6)
	copy(b[4:], v)
	return b
}

func makeOXM(class uint16, field uint8, size uint8) []byte {
	b := make([]byte, 4+size)
	binary.BigEndian.PutUint16(b, class)
	b[2] = field << 1
	b[3] = size
	return b
}

func TunnelID(id uint64) *MatchTunnelID {
	return &MatchTunnelID{id}
}

type MatchTunnelID struct {
	ID uint64
}

func (v *MatchTunnelID) MarshalMatch() []byte {
	b := make([]byte, 4+8)
	binary.BigEndian.PutUint16(b, ofp4.OFPXMC_OPENFLOW_BASIC)
	b[2] = ofp4.OFPXMT_OFB_TUNNEL_ID << 1
	b[3] = 8
	binary.BigEndian.PutUint64(b[4:], v.ID)
	return b
}

func ParseOXMFields(b []byte) (*OXMFields, error) {
	oxm := &OXMFields{make(map[uint32]interface{})}
	for len(b) > 0 {
		class := binary.BigEndian.Uint16(b)
		field := b[2] >> 1
		size := b[3]
		payload := b[4 : 4+size]
		set := func(v interface{}) {
			oxm.set(class, field, v)
		}

		switch class {
		case ofp4.OFPXMC_OPENFLOW_BASIC:
			switch field {
			case ofp4.OFPXMT_OFB_IN_PORT:
				set(InPort(binary.BigEndian.Uint32(payload)))
			case ofp4.OFPXMT_OFB_ETH_DST:
				set(EthDst(net.HardwareAddr(payload)))
			case ofp4.OFPXMT_OFB_ETH_SRC:
				set(EthSrc(net.HardwareAddr(payload)))
			case ofp4.OFPXMT_OFB_TUNNEL_ID:
				set(TunnelID(binary.BigEndian.Uint64(payload)))
			}
		case nxm.OFPXMC_NXM_1:
			switch field {
			case nxm.OFPXMC_NXM1_TUN_IPV4_DST:
				set(nxm.TunnelIPv4Dst(payload))
			}
		}
		b = b[4+size:]
	}
	return oxm, nil
}

type OXMFields struct {
	m map[uint32]interface{}
}

func (o *OXMFields) set(class uint16, field uint8, v interface{}) {
	o.m[o.key(class, field)] = v
}

func (o *OXMFields) key(class uint16, field uint8) uint32 {
	return uint32(class)<<16 | uint32(field)
}

func (o *OXMFields) Lookup(class uint16, field uint8) interface{} {
	return o.m[o.key(class, field)]
}

func (o *OXMFields) InPort() *MatchInPort {
	v, _ := o.Lookup(ofp4.OFPXMC_OPENFLOW_BASIC, ofp4.OFPXMT_OFB_IN_PORT).(*MatchInPort)
	return v
}

func (o *OXMFields) EthDst() *MatchEthDst {
	v, _ := o.Lookup(ofp4.OFPXMC_OPENFLOW_BASIC, ofp4.OFPXMT_OFB_ETH_DST).(*MatchEthDst)
	return v
}

func (o *OXMFields) EthSrc() *MatchEthSrc {
	v, _ := o.Lookup(ofp4.OFPXMC_OPENFLOW_BASIC, ofp4.OFPXMT_OFB_ETH_SRC).(*MatchEthSrc)
	return v
}

func (o *OXMFields) TunnelID() *MatchTunnelID {
	v, _ := o.Lookup(ofp4.OFPXMC_OPENFLOW_BASIC, ofp4.OFPXMT_OFB_TUNNEL_ID).(*MatchTunnelID)
	return v
}
