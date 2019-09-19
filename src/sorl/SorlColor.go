package main

import (
	"math/rand"
	"time"
)

func SorlGetColor() string {

	SorlColors := []string{
		"",
		"\x1b[31;1m",
		"\x1b[37;1m",
		"\x1b[32;1m",
		"\x1b[33;1m",
		"\x1b[34;1m",
		"\x1b[35;1m",
		"\x1b[36;1m",
		"\x1b[38;1m",
		"\x1b[39;1m",
	}
	max := len(SorlColors)
	min := 0
	rand.Seed(time.Now().UnixNano())
	trand := rand.Intn(max-min) + min

	return SorlColors[trand]

}
