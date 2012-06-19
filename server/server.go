package main

import "os"
import "io"
import "image/png"
import "flag"
import "net/http"
import "log"
import "time"
import "math/rand"

import "github.com/PieterD/WorldGen/noise"
import "github.com/PieterD/WorldGen/worldgen"

var testrnd *bool = flag.Bool("testrnd", false, "Test the Random Number Generator")
var listen *string = flag.String("listen", "127.0.0.1:80", "Address to listen on")

func main() {
	flag.Parse()

	if *testrnd {
		rand.Seed(time.Now().UnixNano())
		rnd := noise.NewRnd()
		rnd.Seed6()
		rnd.AutoSeed(true)
		io.Copy(os.Stdout, rnd)
		return
	}

	http.HandleFunc("/dynamic/worldgen.png", handle)

	err := http.ListenAndServe(*listen, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
	}
}

func handle(w http.ResponseWriter, req *http.Request) {
	wrap := true
	log.Print("HTTP")
	req.ParseForm()
	if req.Form["nowrap"] != nil && len(req.Form["nowrap"]) > 0 && req.Form["nowrap"][0] == "true" {
		wrap = false
	}
	w.Header().Set("Content-Type", "image/png; charset=binary")
	world := worldgen.NewWorld(512)
	world.Generate(0.5, wrap)
	img := world.Image()
	err := png.Encode(w, img)
	if err != nil {
		log.Print(err.Error())
	}
}
