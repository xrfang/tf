package main

import (
	"bytes"
	"encoding/binary"
	"strings"
	"time"

	"go.xrfang.cn/act"
)

type Content struct {
	id      uint64
	Caption string
	Type    string
	Data    []byte
	Creator int64
	Created time.Time
	Updated time.Time
}

func (c Content) encode() []byte {
	var bb bytes.Buffer
	binary.Write(&bb, binary.BigEndian, c.Creator)
	binary.Write(&bb, binary.BigEndian, c.Created.UnixNano())
	binary.Write(&bb, binary.BigEndian, c.Updated.UnixNano())
	if cap := strings.TrimSpace(c.Caption); len(cap) > 0 {
		bb.WriteString(cap)
	}
	bb.WriteByte(0)
	if typ := strings.TrimSpace(c.Type); len(typ) > 0 {
		bb.WriteString(typ)
	}
	bb.WriteByte(0)
	if len(c.Data) > 0 {
		bb.Write(c.Data)
	}
	return bb.Bytes()
}

func (c *Content) decode(b []byte) *Content {
	r := bytes.NewReader(b)
	act.Assert(binary.Read(r, binary.BigEndian, &c.Creator))
	var ts int64
	act.Assert(binary.Read(r, binary.BigEndian, &ts))
	c.Created = time.Unix(0, ts)
	act.Assert(binary.Read(r, binary.BigEndian, &ts))
	c.Updated = time.Unix(0, ts)
	var bs [][]byte
	buf := make([]byte, len(b))
	if n, _ := r.Read(buf); n > 0 {
		bs = bytes.SplitN(buf[:n], []byte{0}, 3)
	}
	c.Caption = ""
	if len(bs) > 0 {
		c.Caption = string(bs[0])
	}
	c.Type = ""
	if len(bs) > 1 {
		c.Type = string(bs[1])
	}
	c.Data = nil
	if len(bs) > 2 {
		c.Data = bs[2]
	}
	return c
}
