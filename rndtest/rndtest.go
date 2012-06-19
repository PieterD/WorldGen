package main

import "os"
import "io"

import "github.com/PieterD/WorldGen/noise"

func main() {
	rnd := noise.NewRnd()
	rnd.Seed6()
	rnd.AutoSeed(true)
	io.Copy(os.Stdout, rnd)
	return
}
