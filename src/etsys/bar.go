package etsys

import "time"

type Bar struct {
	Started       time.Time
	Ended         time.Time
	O, H, L, C, V float64
}
