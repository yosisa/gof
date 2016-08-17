package gof

import "github.com/hkwi/gopenflow/ofp4"

type Handler interface {
	Features(*Writer, ofp4.SwitchFeatures)
	PacketIn(*Writer, ofp4.PacketIn)
	FlowRemoved(*Writer, ofp4.FlowRemoved)
	MultipartReply(*Writer, ofp4.MultipartReply)
}

type Marshaler interface {
	Marshal() []byte
}

type Writer struct {
	c *conn
}

func (w *Writer) Write(m Marshaler) {
	w.c.outbound <- m.Marshal()
}

func (w *Writer) WriteBytes(b []byte) {
	w.c.outbound <- b
}
