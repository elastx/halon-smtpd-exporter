package halon_smtpd_ctl

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

type HalonSMTPDCtl struct {
	socket string
}

func New(socket string) *HalonSMTPDCtl {
	return &HalonSMTPDCtl{socket: socket}
}

func (h *HalonSMTPDCtl) connectHandshake() (net.Conn, error) {
	c, err := net.Dial("unix", h.socket)
	if err != nil {
		return c, err
	}

	// Handshake
	var b = []byte{5, 15, 113}
	_, err = c.Write(b)
	if err != nil {
		c.Close()
		return c, err
	}

	return c, nil
}

type response struct {
	Error        error
	Status       string
	Response     []byte
	ResponseSize uint64
}

func (h *HalonSMTPDCtl) readResponse(c net.Conn) response {
	var r response

	status := make([]byte, 1)
	_, err := c.Read(status)
	if err != nil {
		r.Error = err
		return r
	}
	r.Status = string(status)

	sz := make([]byte, 8)
	_, err = c.Read(sz)
	if err != nil {
		r.Error = err
		return r
	}
	r.ResponseSize = binary.LittleEndian.Uint64(sz)

	data := make([]byte, r.ResponseSize)
	_, err = io.ReadFull(c, data)
	if err != nil {
		r.Error = err
		return r
	}

	r.Response = data
	return r
}

func (h *HalonSMTPDCtl) readResponseChan(c net.Conn, ch chan response) {
	ch <- h.readResponse(c)
}

func (h *HalonSMTPDCtl) Query(q string) ([]byte, error) {
	var r []byte

	c, err := h.connectHandshake()
	if err != nil {
		return r, err
	}
	defer c.Close()

	ch := make(chan response)
	go h.readResponseChan(c, ch)

	_, err = c.Write([]byte(q))
	if err != nil {
		return r, err
	}

	response := <-ch
	if response.Error != nil {
		return r, response.Error
	}
	if response.Status != "+" {
		return r, fmt.Errorf("Status response indicates error, body is: %s\n", response.Response)
	}

	return response.Response, nil
}
