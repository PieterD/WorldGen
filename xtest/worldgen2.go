package main

import "image/color"
import "image/draw"
import "code.google.com/p/x-go-binding/ui/x11"
import "code.google.com/p/x-go-binding/ui"

import "fmt"

import "github.com/PieterD/WorldGen/noise"
import "github.com/PieterD/WorldGen/worldgen2"

func do_worldgen2() {
	win, err := x11.NewWindow()
	Panic(err)
	defer win.Close()
	eventchan := win.EventChan()
	scr := win.Screen()

	w := worldgen2.NewWorld(17, 5)
	translate := noise.NewTranslate(float64(scr.Bounds().Max.X), float64(scr.Bounds().Max.Y))
	translate.Zoom(0, 0, 1, 1)

	drawworld2(translate, w, scr)

	selecting := false
	selectx := uint32(0)
	selecty := uint32(0)

	stop := false
	for !stop {
		win.FlushImage()
		tlevent := <-eventchan
		switch event := tlevent.(type) {
		case nil:
			fmt.Printf("Window closed\n")
			stop = true
		//case draw.KeyEvent:
		case ui.MouseEvent:
			if event.Buttons&1 == 1 {
				if !selecting {
					selecting = true
					selectx = uint32(event.Loc.X)
					selecty = uint32(event.Loc.Y)
				} else {
					selecting = false
					translate.ZoomBox(float64(selectx), float64(selecty), float64(event.Loc.X), float64(event.Loc.Y))
					drawworld2(translate, w, scr)
				}
			} else {
				posx, posy := translate.Into(float64(event.Loc.X), float64(event.Loc.Y))
				per := w.GetHeight_Island(posx, posy)
				fmt.Printf("%f, %f = %f\n", posx, posy, per)
			}
		case ui.ConfigEvent:
			fmt.Printf("CONFIG\n")
			scr = win.Screen()
			translate.Resize(float64(scr.Bounds().Max.X), float64(scr.Bounds().Max.Y))
			drawworld2(translate, w, scr)
		case ui.ErrEvent:
			fmt.Printf("ERR\n")
			Panic(event.Err)
		}
	}
}

func drawworld2(translate *noise.Translate, w *worldgen2.World, img draw.Image) {
	w.Log = false
	bounds := img.Bounds()

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			var c color.RGBA
			c.A = 255
			posx, posy := translate.Into(float64(x), float64(y))
			per := w.GetHeight_Island(posx, posy)

			if per < 0 {
				v := uint8(255 - 127*(-per))
				c.R = 0
				c.G = 0
				c.B = v
			} else if per < 0.00002 {
				v := uint8(255 - 2550*per)
				c.R = v
				c.G = v
				c.B = 0
			} else if per > 0.34 {
				v := uint8(220 * (per))
				c.R = v
				c.G = v
				c.B = v
			} else {
				//v := uint8(255-per*255)
				v := uint8(255 - per*300)
				c.R = 0
				c.G = v
				c.B = 0
			}
			img.Set(x, y, c)
		}
	}
	w.Log = true
}
