package noise

import "io"
import "crypto/rand"

/* The Stretto generator is extremely simple.
 * It uses a set of prime-length byte strings, initialized with random
 * data. For a given number i, it xors together the byte at i MOD l,
 * where l is the length of the string, for every string.
 * While each individual string repeats every l bytes, the combination
 * of all has a cycle length equal to the product of all l.
 *
 * The advantage is a direct mapping between sequence number and random
 * number, and it's cheap.
 *
 * I can't imagine this being very secure, but it holds up okay in
 * statistical tests if the initial random values are good, and if
 * the period is of decent size. Luckily, it's easy to get exceptionally
 * large periods.
 */

type Rnd struct {
	seed     []Seed
	autoseed bool
}

func NewRnd() (rnd *Rnd) {
	return new(Rnd)
}

func (rnd *Rnd) Seed2() {
	/* This has a period of 4.295.214.307, takes up 500k. */
	rnd.AddSeed(NewSeed(61751).ReSeed())
	rnd.AddSeed(NewSeed(69557).ReSeed())
}

func (rnd *Rnd) Seed4() {
	/* This has a period of 4.623.376.403, takes up 12k. */
	rnd.AddSeed(NewSeed(37).ReSeed())
	rnd.AddSeed(NewSeed(373).ReSeed())
	rnd.AddSeed(NewSeed(521).ReSeed())
	rnd.AddSeed(NewSeed(643).ReSeed())
}

func (rnd *Rnd) Seed6() {
	/* This has a period of 4.472.576.243, takes up 252 bytes. */
	rnd.AddSeed(NewSeed(29).ReSeed())
	rnd.AddSeed(NewSeed(31).ReSeed())
	rnd.AddSeed(NewSeed(37).ReSeed())
	rnd.AddSeed(NewSeed(43).ReSeed())
	rnd.AddSeed(NewSeed(53).ReSeed())
	rnd.AddSeed(NewSeed(59).ReSeed())
}

func (rnd *Rnd) GetPeriod() (period uint64) {
	for _, seed := range rnd.seed {
		if period == 0 {
			period = uint64(len(seed.data))
		} else {
			period *= uint64(len(seed.data))
		}
	}
	return
}

func (rnd *Rnd) ReSeed() {
	for _, seed := range rnd.seed {
		seed.ReSeed()
	}
}

func (rnd *Rnd) AddSeed(newseed *Seed) {
	rnd.seed = append(rnd.seed, *newseed)
}

func (rnd *Rnd) AutoSeed(yesno bool) {
	rnd.autoseed = yesno
}

func (rnd *Rnd) Index(index uint32) (rv uint32) {
	rv = 0
	for _, seed := range rnd.seed {
		idx := index % uint32(len(seed.data))
		rv ^= seed.data[idx]
	}
	return
}

func (rnd *Rnd) Uint32() (rv uint32) {
	rv = 0
	if rnd.autoseed {
		count := 0
		for _, seed := range rnd.seed {
			if seed.pos == 0 {
				count++
			}
		}
		if count == len(rnd.seed) {
			rnd.ReSeed()
		}
	}
	for i := range rnd.seed {
		rv = rv ^ rnd.seed[i].data[rnd.seed[i].pos]
		rnd.seed[i].pos++
		rnd.seed[i].pos %= len(rnd.seed[i].data)
	}
	return
}

func (rnd *Rnd) Uint64() (rv uint64) {
	rv = uint64(rnd.Uint32())
	rv <<= 32
	rv |= uint64(rnd.Uint64())
	return
}

func (rnd *Rnd) Int31() (rv int32) {
	ui := rnd.Uint32()
	rv = int32(ui & 0x7FFFFFFF)
	return
}

func (rnd *Rnd) Int63() (rv int64) {
	ui := rnd.Uint64()
	rv = int64(ui & 0x7FFFFFFFFFFFFFFF)
	return
}

func (rnd *Rnd) Float64() (rv float64) {
	i := rnd.Int63()
	rv = float64(i) / (1 << 63)
	return
}
func (rnd *Rnd) Float32() (rv float32) {
	return float32(rnd.Float64())
}

func (rnd *Rnd) Read(buf []byte) (n int, err error) {
	u := rnd.Uint32()
	pos := 0
	for i := range buf {
		buf[i] = byte(u & 0xFF)
		u >>= 8
		pos++
		if pos == 4 {
			u = rnd.Uint32()
			pos = 0
		}
	}
	return len(buf), nil
}

type Seed struct {
	data  []uint32
	pos   int
	input io.Reader
}

func NewSeed(size int) (s *Seed) {
	s = new(Seed)
	s.data = make([]uint32, size)
	s.input = rand.Reader
	return
}
func (s *Seed) SetReader(input io.Reader) {
	s.input = input
}
func (s *Seed) ReSeed() *Seed {
	buf := make([]byte, 4)
	for i := range s.data {
		_, err := io.ReadFull(s.input, buf)
		if err != nil {
			panic(err)
		}
		s.data[i] = uint32(buf[0])
		s.data[i] <<= 8
		s.data[i] |= uint32(buf[1])
		s.data[i] <<= 8
		s.data[i] |= uint32(buf[2])
		s.data[i] <<= 8
		s.data[i] |= uint32(buf[3])
	}
	return s
}
