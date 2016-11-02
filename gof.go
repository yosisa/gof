package gof

import (
	"sync"

	"github.com/hkwi/gopenflow/ofp4"
)

type Handler interface {
	HandleMessage(*Writer, *Message)
}

type FancyHandler interface {
	Features(*Writer, ofp4.SwitchFeatures)
	PacketIn(*Writer, ofp4.PacketIn)
	FlowRemoved(*Writer, ofp4.FlowRemoved)
	MultipartReply(*Writer, ofp4.MultipartReply)
	PortStatus(*Writer, ofp4.PortStatus)
}

func FancyHandle(h FancyHandler) Handler {
	return &fancyHandler{h}
}

type fancyHandler struct {
	h FancyHandler
}

func (h fancyHandler) HandleMessage(w *Writer, m *Message) {
	switch m.Type {
	case ofp4.OFPT_FEATURES_REPLY:
		h.h.Features(w, ofp4.SwitchFeatures(m.Payload))
	case ofp4.OFPT_PACKET_IN:
		h.h.PacketIn(w, ofp4.PacketIn(m.Payload))
	case ofp4.OFPT_FLOW_REMOVED:
		h.h.FlowRemoved(w, ofp4.FlowRemoved(m.Payload))
	case ofp4.OFPT_MULTIPART_REPLY:
		h.h.MultipartReply(w, ofp4.MultipartReply(m.Payload))
	case ofp4.OFPT_PORT_STATUS:
		h.h.PortStatus(w, ofp4.PortStatus(m.Payload))
	}
}

type NopFancyHandler struct{}

func (h NopFancyHandler) Features(w *Writer, d ofp4.SwitchFeatures)       {}
func (h NopFancyHandler) PacketIn(w *Writer, d ofp4.PacketIn)             {}
func (h NopFancyHandler) FlowRemoved(w *Writer, d ofp4.FlowRemoved)       {}
func (h NopFancyHandler) MultipartReply(w *Writer, d ofp4.MultipartReply) {}
func (h NopFancyHandler) PortStatus(W *Writer, d ofp4.PortStatus)         {}

func AutoInstantiate(f func(uint64) Handler) Handler {
	return &autoInstantiateHandler{
		h: make(map[uint64]Handler),
		f: f,
	}
}

type autoInstantiateHandler struct {
	h map[uint64]Handler
	f func(uint64) Handler
	m sync.Mutex
}

func (h *autoInstantiateHandler) HandleMessage(w *Writer, m *Message) {
	h.m.Lock()
	handler := h.h[m.DatapathID]
	if handler == nil {
		handler = h.f(m.DatapathID)
		h.h[m.DatapathID] = handler
	}
	h.m.Unlock()
	handler.HandleMessage(w, m)
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

type Message struct {
	DatapathID uint64
	Type       uint8
	Payload    []byte
}
