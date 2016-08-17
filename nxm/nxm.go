package nxm

import (
	"encoding/binary"
	"net"
)

const (
	NX_EXPERIMENTER_ID = 0x2320
)

const (
	OFPXMC_NXM_0 uint16 = iota
	OFPXMC_NXM_1
)

const (
	NXM_OF_ETH_DST uint32 = uint32(OFPXMC_NXM_0)<<16 | 1<<9 | 6
	NXM_OF_ETH_SRC uint32 = uint32(OFPXMC_NXM_0)<<16 | 2<<9 | 6
	NXM_OF_ARP_SPA uint32 = uint32(OFPXMC_NXM_0)<<16 | 16<<9 | 4
	NXM_OF_ARP_TPA uint32 = uint32(OFPXMC_NXM_0)<<16 | 17<<9 | 4
)

const (
	NXM_NX_ARP_SHA uint32 = uint32(OFPXMC_NXM_1)<<16 | 17<<9 | 6
	NXM_NX_ARP_THA uint32 = uint32(OFPXMC_NXM_1)<<16 | 18<<9 | 6
)

const (
	OFPXMC_NXM1_TUN_IPV4_SRC uint8 = 31
	OFPXMC_NXM1_TUN_IPV4_DST uint8 = 32
)

const (
	NXAST_REG_MOVE = 6
)

func TunnelIPv4Dst(ip net.IP) *MatchTunnelIPv4Dst {
	return &MatchTunnelIPv4Dst{ip: ip}
}

type MatchTunnelIPv4Dst struct {
	ip net.IP
}

func (v *MatchTunnelIPv4Dst) MarshalMatch() []byte {
	b := make([]byte, 4+4)
	binary.BigEndian.PutUint16(b, OFPXMC_NXM_1)
	b[2] = OFPXMC_NXM1_TUN_IPV4_DST << 1
	b[3] = 4
	copy(b[4:], v.ip.To4())
	return b
}
