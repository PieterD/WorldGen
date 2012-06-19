package noise

// Implementation based on documentation by Matt Zucker <mazucker@vassar.edu>
// http://webstaff.itn.liu.se/~stegu/TNM022-2005/perlinnoiselinks/perlin-noise-math-faq.html

// Persistence device by Hugo Elias
// from http://freespace.virgin.net/hugo.elias/models/m_perlin.htm

// Implementation changed to a much simpler, but most of all functional
// interpolation method. I am not convinced it's Perlin noise anymore.

import "math"
import "fmt"

type NoiseMap struct {
	focus uint32
	wrap bool
	octaves uint32
	persistence float64
	rnd *Rnd
	Log bool
	cache []NoiseCache
	normalization float64
}
type NoiseCache struct {
	x, y int64
	grade [4]NoiseGrade
}
type NoiseGrade struct {
	x, y float64
}
func NewNoiseMap(rnd *Rnd, focus uint32, octaves int, persistence float64, wrap bool) (nm *NoiseMap) {
	nm = new(NoiseMap)
	nm.focus = focus
	nm.wrap = wrap
	nm.rnd = rnd
	nm.octaves = uint32(octaves)
	nm.persistence = persistence
	amplitude := float64(1)
	for o:=0; o<octaves; o++ {
		nm.normalization += amplitude
		amplitude *= nm.persistence
	}
	return
}
// Generate a random number (0 - 1) for the given integer coordinates
func (nm *NoiseMap) Noise2d(x, y uint32, octave uint32) float64 {
	size := nm.focus*(octave+1)
	if !nm.wrap {
		size ++
	}
	x %= size
	y %= size
	idx := y*size + x
	// FIXME: Doesn't really jump over the previous octaves
	idx += size*size
	rv := float64(nm.rnd.Index(idx)&0x7FFFFFFF)/(1<<31)
	if nm.Log {
		fmt.Printf("  noise2d(%d,%d) %d=%08X = %f\n", x, y, idx, nm.rnd.Index(idx), rv)
	}
	return rv
}

func (nm *NoiseMap) Perlin(x, y float64) (avg float64) {
	var frequency float64 = float64(nm.focus)
	var amplitude float64 = 1
	for octave:=uint32(0); octave<nm.octaves; octave++ {
		if nm.Log {
			fmt.Printf("  Octave %d freq=%f amp=%f\n", octave, frequency, amplitude)
		}

		posx := x*frequency
		posy := y*frequency
		// Integer grid coordinates
		gridx := uint32(posx)
		gridy := uint32(posy)
		// What's left of the fractional part of the coordinates
		fracx := posx - float64(gridx)
		fracy := posy - float64(gridy)
		if nm.Log {
			fmt.Printf(" Perlin(%d) %f,%f -> %d,%d (%f,%f)\n", octave, x, y, gridx, gridy, fracx, fracy)
		}

		/* Raw noise method */
		s := nm.Noise2d(gridx, gridy, octave)
		t := nm.Noise2d(gridx+1, gridy, octave)
		u := nm.Noise2d(gridx, gridy+1, octave)
		v := nm.Noise2d(gridx+1, gridy+1, octave)
		if nm.Log {
			fmt.Printf("  s=%f\n", s)
			fmt.Printf("  t=%f\n", t)
			fmt.Printf("  u=%f\n", u)
			fmt.Printf("  v=%f\n", v)
		}

		/* Cosine interpolation */
		a := interpolate_cosine(s, t, fracx)
		b := interpolate_cosine(u, v, fracx)
		z := interpolate_cosine(a, b, fracy)

		avg += z*amplitude

		frequency *= 2
		amplitude *= nm.persistence
	}
	avg /= nm.normalization
	if nm.Log {
		fmt.Printf("     Average %f\n", avg)
	}
	return
}

func interpolate_cosine(a, b, delta float64) float64 {
	theta := delta * math.Pi
	f := (1-math.Cos(theta)) / 2
	return a*(1-f) + b*f
}


