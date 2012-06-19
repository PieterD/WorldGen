package noise

import "testing"

func TestPeriod(t *testing.T) {
	/* Ordinary usage, with prime seeds */
	testPeriod(t, 0, 7, 11)
	testPeriod(t, 0, 7, 11, 29)
	testPeriod(t, 0, 11, 29, 37)
	testPeriod(t, 0, 2, 3, 5, 7, 11)
	/* Composites will repeat themselves */
	testPeriod(t, 12, 4, 6)
	testPeriod(t, 30, 6, 10)
	testPeriod(t, 120, 6, 8, 10)
	/* Check if autoseeding works */
	testAutoSeed(t, 7, 11)
	testAutoSeed(t, 7, 11, 17)
}
func testPeriod(t *testing.T, repeat uint32, primes ...int) {
	var period uint32
	rnd := NewRnd()
	for _, p := range primes {
		rnd.AddSeed(NewSeed(p).ReSeed())
		if period == 0 {
			period = uint32(p)
		} else {
			period *= uint32(p)
		}
	}

	if uint64(period) != rnd.GetPeriod() {
		t.Fatal("My calculated period (%d) and GetPeriod(%d) don't match!", period, rnd.GetPeriod())
	}

	/* See if Index() runs along with rnd.Uint32() when called sequentially */
	for count := uint32(0); count < period; count++ {
		index := rnd.Index(count)
		generate := rnd.Uint32()
		if index != rnd.Index(count+period) {
			t.Fatalf("Index(n) and Index(n+period) don't match!")
		}
		if index != generate {
			t.Fatalf("Index() and Uint32() don't line up! %08X != %08X", index, generate)
		}
	}
	/* Index() works, re-run period */
	for count := uint32(0); count < period; count++ {
		if rnd.Index(count) != rnd.Uint32() {
			t.Fatalf("Function should repeat, but it doesn't!")
		}
	}
	diff := testRepeat(rnd, period)
	if diff == 0 && repeat != 0 {
		t.Fatalf("Pattern repetition expected, none found")
	} else if diff != 0 && repeat == 0 {
		t.Fatalf("Pattern repetition detected, width %d (period %d)", diff, period)
	} else if diff != 0 && repeat != diff {
		t.Fatalf("Pattern repetition detected, width %d (period %d), but expected %d", diff, period, repeat)
	}
}
func testRepeat(rnd *Rnd, period uint32) (width uint32) {
	/* See if the sequence repeats itself before we expect it to */
	for count := uint32(0); count < period; count++ {
		for inner := count + 1; inner < period; inner++ {
			if rnd.Index(count) == rnd.Index(inner) {
				diff := inner - count
				repeat := uint32(0)
				for try := uint32(0); try < period; try++ {
					if rnd.Index(try) == rnd.Index(try+diff) {
						repeat++
					}
				}
				if repeat == period {
					width = diff
				}
			}
		}
	}
	return
}

func testAutoSeed(t *testing.T, primes ...int) {
	rnd := NewRnd()
	for _, p := range primes {
		rnd.AddSeed(NewSeed(p).ReSeed())
	}
	rnd.AutoSeed(true)
	period := uint32(rnd.GetPeriod())
	arr := make([]uint32, period)
	for i := range arr {
		generate := rnd.Uint32()
		arr[i] = rnd.Index(uint32(i))
		if arr[i] != generate {
			t.Fatalf("Before autoseeding, Index() and Uint32() don't match!")
		}
	}
	count := uint32(0)
	for i := range arr {
		generate := rnd.Uint32()
		if arr[i] == rnd.Index(uint32(i)) {
			count++
		}
		if generate != rnd.Index(uint32(i)) {
			t.Fatalf("After autoseeding, Index() and Uint32() don't match!")
		}
	}
	if count > period/2 {
		t.Fatalf("Autoseeding is broken!")
	}
}
