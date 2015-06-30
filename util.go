package main

import (
	"fmt"
	"strings"
)

func printdebug(format string, args ...interface{}) {
	if *debug {
		fmt.Printf(format, args...)
	}
}

func asciifold(s string) string {
	s = strings.ToLower(s)
	r := strings.NewReplacer("ç", "c", "ğ", "g", "ı", "i", "ö", "o", "ş", "s", "ü", "u")

	return r.Replace(s)
}
