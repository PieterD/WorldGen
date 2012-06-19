package worldgen

import (
	"fmt"
	//"os"
	"image"
	"image/color"
	//"image/png"
	"math"
	"math/rand"
	//"time"
)

var LOGGING bool // = true

func log(format string, args ...interface{}) {
	if LOGGING {
		fmt.Printf(format, args...)
	}
}

type World struct {
	Size         int
	Map          []float64
	roughness    float64
	displacement float64
	wrap         bool
}

/*
func main() {
  rand.Seed(time.Nanoseconds())
  w := NewWorld(1024)
  i := math.NaN()
  w.Seed([]float64{0, i, 0,
                   i, 5, i,
                   0, i, 0})
  w.Generate(0.5, true)

  img := w.Image()
  wr, err := os.Create("test.png")
  if err != nil {
    panic(err)
  }
  if err = png.Encode(wr, img); err != nil {
    panic(err)
  }
}
*/

func NewWorld(size int) (w *World) {
	w = new(World)
	w.Size = size + 1
	w.Map = make([]float64, w.Size*w.Size)
	for i := range w.Map {
		w.Map[i] = math.NaN()
	}

	return
}

func (w *World) Seed(init []float64) bool {
	log("Seed\n")
	size := getinitsize(init)
	if size == -1 {
		return false
	}

	if size == 1 {
		w.set(w.Size/2, w.Size/2, init[0])
	} else {
		step := (w.Size - 1) / (size - 1)
		for y := 0; y < size; y++ {
			for x := 0; x < size; x++ {
				w.set(x*step, y*step, init[x+y*size])
			}
		}
	}
	return true
}

/* Generating wrapped terrain means that the opposing edges of the texture are the same thing.
 * Thus a wrapped texture is size /w.Size/ - 1. */
func (w *World) Generate(roughness float64, wrap bool) {
	w.roughness = math.Pow(2, -roughness)
	w.wrap = wrap
	log("Generate\n")
	w.displacement = 1
	for s := w.Size - 1; s > 1; s /= 2 {
		log("step %d (%f) (%f %f %f %f)\n", s, w.displacement, w.rnd(), w.rnd(), w.rnd(), w.rnd())
		/* Diamonds */
		log("Diamonds\n")
		for y := 0; y < w.Size-1; y += s {
			for x := 0; x < w.Size-1; x += s {
				tl := w.getcorner(x, y)
				tr := w.getcorner(x+s, y)
				br := w.getcorner(x+s, y+s)
				bl := w.getcorner(x, y+s)
				c := w.get(x+s/2, y+s/2)
				if math.IsNaN(c) {
					c = (tl+tr+br+bl)/4 + w.rnd()
					w.set(x+s/2, y+s/2, c)
				}
			}
		}
		/* Squares */
		log("Squares\n")
		for y := 0; y < w.Size-1; y += s {
			for x := 0; x < w.Size-1; x += s {
				tl := w.getcorner(x, y)
				tr := w.getcorner(x+s, y)
				br := w.getcorner(x+s, y+s)
				bl := w.getcorner(x, y+s)

				c := w.get(x+s/2, y+s/2)

				w.edge(tl+tr+c, x+s/2, y, 0, -s/2)
				w.edge(tr+br+c, x+s, y+s/2, s/2, 0)
				w.edge(bl+br+c, x+s/2, y+s, 0, s/2)
				w.edge(tl+bl+c, x, y+s/2, -s/2, 0)
			}
		}
		w.displacement *= w.roughness
	}
}

/* Don't render the right and bottom edges of the texture. This gives it a power of 2 size,
 * and if the generated terrain is wrapped those pixels are dead anyway. */
func (w *World) Image() image.Image {
	xtra := 0
	//if w.wrap {
	xtra = 1
	//}
	img := image.NewRGBA(image.Rect(0, 0, int(w.Size-xtra), int(w.Size-xtra)))
	var highest float64 = -900000000
	var lowest float64 = 9000000000
	var avg float64 = 0

	avgnum := 0
	for y := 0; y < w.Size-xtra; y++ {
		for x := 0; x < w.Size-xtra; x++ {
			v := w.Map[w.i(x, y)]
			avg += v
			avgnum++
			if v > highest {
				highest = v
			}
			if v < lowest {
				lowest = v
			}
		}
	}

	avg /= float64(avgnum)
	low := (avg*2 + lowest) / 3
	high := (avg*2 + highest) / 3

	log("stats: %f - %f - %f - %f - %f\n", lowest, low, avg, high, highest)
	for y := 0; y < w.Size-xtra; y++ {
		for x := 0; x < w.Size-xtra; x++ {
			val := w.get(x, y)
			below := false
			if val < low {
				//if val < 0 {
				below = true
			}
			above := false
			if val > high {
				//if val > 15 {
				above = true
			}
			val -= lowest
			val /= highest - lowest
			val *= 255
			log("val %d,%d = %f\n", x, y, val)
			var c color.Color
			if below {
				c = color.RGBA{0, 0, uint8(val), 255}
				//c = image.RGBAColor{0, 0, 120+uint8(val/2), 255}
				//c = image.RGBAColor{0, 0, 255, 255}
			} else if above {
				val /= 2
				c = color.RGBA{uint8(val), uint8(val), uint8(val), 255}
				//c = image.RGBAColor{120+uint8(val/2), 0, 0, 255}
				//c = image.RGBAColor{255, 0, 0, 255}
			} else {
				c = color.RGBA{0, uint8(val), 0, 255}
			}
			img.Set(int(x), int(y), c)
		}
	}
	return img
}

func (w *World) rnd() float64 {
	return rand.Float64()*w.displacement*2 - w.displacement
}

func (w *World) i(x, y int) int {
	xtra := 0
	if w.wrap {
		xtra++
	}
	if x < 0 {
		x += w.Size - 1
	}
	if y < 0 {
		y += w.Size - 1
	}
	if x >= w.Size-xtra {
		x -= w.Size - 1
		log("Wrapping\n")
	}
	if y >= w.Size-xtra {
		y -= w.Size - 1
		log("Wrapping\n")
	}
	return y*w.Size + x
}
func (w *World) get(x, y int) float64 {
	log("  %d,%d -> %f\n", x, y, w.Map[w.i(x, y)])
	return w.Map[w.i(x, y)]
}
func (w *World) set(x, y int, v float64) {
	log("  %d,%d <- %f\n", x, y, v)
	w.Map[w.i(x, y)] = v
}
func (w *World) getcorner(x, y int) float64 {
	log("  corner\n")
	v := w.get(x, y)
	if math.IsNaN(v) {
		log("   corner generate\n")
		v = w.rnd()
		w.set(x, y, v)
	}
	return v
}
func (w *World) edge(base float64, x, y, addx, addy int) float64 {
	if !w.wrap && (x+addx < 0 || y+addy < 0 || x+addx >= w.Size || y+addy >= w.Size) {
		base /= 3
	} else {
		base += w.get(x+addx, y+addy)
		base /= 4
	}
	e := w.get(x, y)
	if math.IsNaN(e) {
		base += w.rnd()
		w.set(x, y, base)
	} else {
		base = e
	}
	return base
}

func getinitsize(init []float64) int {
	var n uint32
	l := uint32(len(init))
	switch l {
	case 0:
		return -1
	case 1:
		return 1
	case 4:
		return 2
	case 9:
		return 3
	default:
		/* We only want (n+1)(n+1) where n is a power of 2.
		 * if this is the case and n>2, each of the terms in (n^2 + 2n + 1) takes up one bit. */
		n = 0
		for shift := uint32(2); shift > 0; shift <<= 1 {
			if l&shift != 0 {
				n = shift / 2
				break
			}
		}
		if n > 0 && n*n > 0 && l == 1+2*n+n*n {
			return int(n + 1)
		}
	}
	return -1
}
