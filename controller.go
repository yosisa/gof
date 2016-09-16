package gof

import (
	"bufio"
	"encoding/binary"
	"io"
	"net"
	"sync"
	"time"

	"github.com/hkwi/gopenflow/ofp4"
)

func ListenAndServe(addr string, handler Handler) error {
	ctrl := &Controller{Addr: addr, Handler: handler}
	return ctrl.ListenAndServe()
}

type Controller struct {
	Addr         string
	Handler      Handler
	WriteTimeout time.Duration
	Concurrency  uint
}

func (c *Controller) ListenAndServe() error {
	addr := c.Addr
	if addr == "" {
		addr = ":6653"
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return c.Serve(ln)
}

func (c *Controller) Serve(l net.Listener) error {
	defer l.Close()
	for {
		nc, err := l.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				time.Sleep(5 * time.Millisecond)
				continue
			}
			return err
		}
		conn := c.newConn(nc)
		go conn.serve()
	}
}

func (c *Controller) newConn(nc net.Conn) *conn {
	return &conn{
		ctrl:     c,
		conn:     nc,
		inbound:  make(chan *Message, 100),
		outbound: make(chan []byte, 100),
	}
}

type conn struct {
	ctrl     *Controller
	conn     net.Conn
	bufr     *bufio.Reader
	inbound  chan *Message
	outbound chan []byte
	err      error
}

func (c *conn) serve() {
	w := &Writer{c: c}
	concurrency := int(c.ctrl.Concurrency)
	if concurrency == 0 {
		concurrency = 1
	}

	var wg sync.WaitGroup
	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go c.handleLoop(w)
	}

	go c.writeLoop()
	c.bufr = newBufioReader(c.conn)
	c.outbound <- ofp4.MakeHello(ofp4.MakeHelloElemVersionbitmap([]uint32{0x10}))
	c.outbound <- ofp4.MakeHeader(ofp4.OFPT_FEATURES_REQUEST)

	defer func() {
		putBufioReader(c.bufr)
		close(c.inbound)
		wg.Wait()
		close(c.outbound)
	}()

	var dpid uint64
	for c.err == nil {
		typ, payload, err := c.readData()
		if err != nil {
			if err == io.EOF {
				return
			}
			c.err = err
			return
		}

		switch typ {
		case ofp4.OFPT_ECHO_REQUEST:
			payload[1] = ofp4.OFPT_ECHO_REPLY
			c.outbound <- payload
		case ofp4.OFPT_FEATURES_REPLY:
			dpid = ofp4.SwitchFeatures(payload).DatapathId()
			fallthrough
		default:
			c.inbound <- &Message{dpid, typ, payload}
		}
	}
}

func (c *conn) readData() (typ uint8, payload []byte, err error) {
	var b []byte
	if b, err = c.bufr.Peek(4); err != nil {
		return
	}
	hdr := ofp4.Header(b)

	payload = make([]byte, hdr.Length())
	if _, err = io.ReadFull(c.bufr, payload); err != nil {
		return
	}

	typ = hdr.Type()
	if typ == ofp4.OFPT_ERROR {
		err = ofp4.ErrorMsg(payload)
	}
	return
}

func (c *conn) handleLoop(w *Writer) {
	for m := range c.inbound {
		c.ctrl.Handler.HandleMessage(w, m)
	}
}

func (c *conn) writeLoop() {
	var xid uint32 = 1
	defer c.conn.Close()
	for b := range c.outbound {
		binary.BigEndian.PutUint32(b[4:], xid)
		xid++

		if d := c.ctrl.WriteTimeout; d != 0 {
			c.conn.SetWriteDeadline(time.Now().Add(d))
		}
		_, c.err = c.conn.Write(b)
		if c.err != nil {
			return
		}
	}
}

var (
	bufioReaderPool sync.Pool
)

func newBufioReader(r io.Reader) *bufio.Reader {
	if v := bufioReaderPool.Get(); v != nil {
		br := v.(*bufio.Reader)
		br.Reset(r)
		return br
	}
	return bufio.NewReader(r)
}

func putBufioReader(br *bufio.Reader) {
	br.Reset(nil)
	bufioReaderPool.Put(br)
}
