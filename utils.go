package main

import (
	"bytes"
	"encoding/binary"
	"io"
	"strings"

	"go.etcd.io/bbolt"
	"go.xrfang.cn/act"
)

type (
	integers interface {
		int8 | int16 | int32 | int64 | uint8 | uint16 | uint32 | uint64
	}
	floats interface {
		float32 | float64
	}
	value interface {
		integers | floats | string
	}
)

func encode[T value](vs ...T) []byte {
	act.Assert(len(vs) > 0, "encode: no value")
	var bb bytes.Buffer
	switch ts := any(vs).(type) {
	case []string:
		for _, t := range ts {
			bb.WriteByte(0)
			bb.WriteString(t)
		}
	case []int8:
		bb.WriteByte(0x41)
		for _, t := range ts {
			bb.WriteByte(byte(t))
		}
	case []int16:
		bb.WriteByte(0x42)
		for _, t := range ts {
			binary.Write(&bb, binary.BigEndian, t)
		}
	case []int32:
		bb.WriteByte(0x44)
		for _, t := range ts {
			binary.Write(&bb, binary.BigEndian, t)
		}
	case []int64:
		bb.WriteByte(0x48)
		for _, t := range ts {
			binary.Write(&bb, binary.BigEndian, t)
		}
	case []uint8:
		bb.WriteByte(0x01)
		for _, t := range ts {
			bb.WriteByte(byte(t))
		}
	case []uint16:
		bb.WriteByte(0x02)
		for _, t := range ts {
			binary.Write(&bb, binary.BigEndian, t)
		}
	case []uint32:
		bb.WriteByte(0x04)
		for _, t := range ts {
			binary.Write(&bb, binary.BigEndian, t)
		}
	case []uint64:
		bb.WriteByte(0x08)
		for _, t := range ts {
			binary.Write(&bb, binary.BigEndian, t)
		}
	case []float32:
		bb.WriteByte(0x84)
		for _, t := range ts {
			binary.Write(&bb, binary.BigEndian, t)
		}
	case []float64:
		bb.WriteByte(0x88)
		for _, t := range ts {
			binary.Write(&bb, binary.BigEndian, t)
		}
	}
	return bb.Bytes()
}

func decode[T value](buf []byte) (vs []T, err error) {
	defer act.Catch(&err)
	if len(buf) == 0 {
		return vs, nil
	}
	parse := func(buf []byte) {
		br := bytes.NewReader(buf)
		for {
			var v T
			if err := binary.Read(br, binary.BigEndian, &v); err != nil {
				if err == io.EOF {
					break
				}
				panic(err)
			}
			vs = append(vs, v)
		}
	}
	var z T
	switch any(z).(type) {
	case string:
		act.Assert(buf[0] == 0, "decode: data is not string (expected mark 0x00)")
		for _, s := range strings.Split(string(buf[1:]), string(0)) {
			vs = append(vs, any(s).(T))
		}
	case int8:
		act.Assert(buf[0] == 0x41, "decode: data is not int8 (expected mark 0x41)")
		for _, i := range buf[1:] {
			vs = append(vs, any(i).(T))
		}
	case int16:
		act.Assert(buf[0] == 0x42, "decode: data is not int16 (expected mark 0x42)")
		parse(buf[1:])
	case int32:
		act.Assert(buf[0] == 0x44, "decode: data is not int32 (expected mark 0x44)")
		parse(buf[1:])
	case int64:
		act.Assert(buf[0] == 0x48, "decode: data is not int64 (expected mark 0x48)")
		parse(buf[1:])
	case uint8:
		act.Assert(buf[0] == 0x01, "decode: data is not uint8 (expected mark 0x01)")
		for _, i := range buf[1:] {
			vs = append(vs, any(i).(T))
		}
	case uint16:
		act.Assert(buf[0] == 0x02, "decode: data is not uint16 (expected mark 0x02)")
		parse(buf[1:])
	case uint32:
		act.Assert(buf[0] == 0x04, "decode: data is not uint32 (expected mark 0x04)")
		parse(buf[1:])
	case uint64:
		act.Assert(buf[0] == 0x08, "decode: data is not uint64 (expected mark 0x08)")
		parse(buf[1:])
	case float32:
		act.Assert(buf[0] == 0x84, "decode: data is not float32 (expected mark 0x84)")
		parse(buf[1:])
	case float64:
		act.Assert(buf[0] == 0x88, "decode: data is not float64 (expected mark 0x88)")
		parse(buf[1:])
	}
	return vs, nil
}

func bput[T value](b *bbolt.Bucket, key string, val ...T) {
	act.Assert(b.Put([]byte(key), encode(val...)))
}
