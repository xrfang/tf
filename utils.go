package main

import (
	"bytes"
	"encoding/binary"
	"time"

	"go.etcd.io/bbolt"
	"go.xrfang.cn/act"
)

type (
	integers interface {
		int8 | int16 | int32 | int64 | int | uint8 | uint16 | uint32 | uint64 | uint
	}
	floats interface {
		float32 | float64
	}
	value interface {
		integers | floats | string | time.Time
	}
)

func encode[T value](v T) []byte {
	switch t := any(v).(type) {
	case string:
		return []byte(t)
	case time.Time:
		return encode(t.UnixNano())
	default:
		var bb bytes.Buffer
		act.Assert(binary.Write(&bb, binary.BigEndian, t))
		return bb.Bytes()
	}
}

func decode[T value](buf []byte, v *T) {
	switch t := any(v).(type) {
	case *string:
		*t = string(buf)
	case *time.Time:
		var ts int64
		decode(buf, &ts)
		*t = time.Unix(0, ts)
	default:
		br := bytes.NewReader(buf)
		act.Assert(binary.Read(br, binary.BigEndian, t))
	}
}

func bput[T integers | floats | string](b *bbolt.Bucket, key string, val T) {
	act.Assert(b.Put([]byte(key), encode(val)))
}
