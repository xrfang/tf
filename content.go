package main

import "time"

type Content struct {
	Caption string
	Type    string
	Data    []byte
	Creator int64
	Created time.Time
	Updated time.Time
}
