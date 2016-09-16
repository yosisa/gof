package gof

type DatapathMux struct {
	m  map[uint64]Handler
	dh Handler
}

func (m *DatapathMux) Handle(dpid uint64, h Handler) {
	if m.m == nil {
		m.m = make(map[uint64]Handler)
	}
	m.m[dpid] = h
}

func (m *DatapathMux) SetDefault(h Handler) {
	m.dh = h
}

func (m *DatapathMux) HandleMessage(w *Writer, msg *Message) {
	if h, ok := m.m[msg.DatapathID]; ok {
		h.HandleMessage(w, msg)
	} else if m.dh != nil {
		m.dh.HandleMessage(w, msg)
	}
}
