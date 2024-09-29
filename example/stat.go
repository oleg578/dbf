package main

import (
	"fmt"
	"runtime"
)

func printMemUsage() string {
	var m runtime.MemStats
	out := ""
	runtime.ReadMemStats(&m)
	div := float64(1024 * 1024)
	out += fmt.Sprintf("HeapAlloc = %.3f MB, ", float64(m.HeapAlloc)/div)
	out += fmt.Sprintf("Sys = %.3f MB\n", float64(m.Sys)/div)
	return out
}
