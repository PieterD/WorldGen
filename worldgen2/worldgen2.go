package worldgen2

import "math"
import "fmt"

import "github.com/PieterD/WorldGen/noise"

type World struct {
	Focus float64
	rnd *noise.Rnd
	nm *noise.NoiseMap

	t *noise.Translate
	tm *noise.NoiseMap
	Log bool
}
func NewWorld(focus uint32, octaves int) (w *World) {
	w = new(World)
	w.Focus = float64(focus)
	w.rnd = noise.NewRnd()
	w.rnd.Seed6()
	w.nm = noise.NewNoiseMap(w.rnd, focus, octaves, 0.90, true)
	/*
	w.t = noise.NewTranslate(float64(size), float64(size))
	w.t.Zoom(0,0, 1,1)
	*/

	w.tm = noise.NewNoiseMap(w.rnd, 17, 3, 0.5, false)
	return
}

func (w *World) GetHeight_Island(x, y float64) float64 {
	posx := x
	posy := y
	per := w.nm.Perlin(posx, posy)
	cenx := 0.5 - posx
	ceny := 0.5 - posy
	distance := math.Sqrt(cenx*cenx+ceny*ceny)
	//perturb := (w.nm.Perlin(posx/3, posy/3)-0.5) * w.Focus
	perturb := (w.tm.Perlin(posx, posy)-0.5)/3
	distance += perturb
	b := bulge(distance, 0.25, 0.4)
	per += b
	per -= 1
	if w.Log {
		fmt.Printf("GetHeight(%f,%f)\n", x, y)
		fmt.Printf(" posx,posy: %f,%f\n", posx, posy)
		fmt.Printf(" per: %f\n", per)
		fmt.Printf(" distance: %f\n", distance)
		fmt.Printf(" perturb: %f\n", perturb)
		fmt.Printf(" bulge: %f\n", b)
	}
	// Smooth it out
	return per*per*per*2
}

func bulge (x, d, p float64) float64 {
  x = math.Abs(x)
  return (1+math.Tanh((d-x)*5/p))/2
}

