package main

import "testing"

func Test_rng_next(t *testing.T) {
	seq := []int{
		123,
		15887950,
		16495136,
		527345,
		704524,
		1553684,
		12683156,
		11100544,
		12249484,
		7753432,
		5908254,
	}
	r := &rng{seed: seq[0]}
	for i := 1; i < len(seq); i++ {
		if n := r.next(); n != seq[i] {
			t.Errorf("r{%d}.next() = %d, want %d", seq[i-1], n, seq[i])
		}
	}

	want2 := map[int]int{
		1:    8685429,
		10:   4700978,
		100:  15273692,
		2024: 8667524,
	}

	for seed, want := range want2 {
		r := &rng{seed: seed}
		for i := 0; i < 2000; i++ {
			r.next()
		}
		if r.seed != want {
			t.Errorf("r{%d}.next()^2000 = %d, want %d", seed, r.seed, want)
		}
	}

	t.Run("part2-ex1", func(t *testing.T) {
		seq := [4]int{-1, -1, 0, 2}
		want := 6
		init := 123

		r := &rng{seed: init}
		ps, ds := r.prices()
		p := search(ps, ds, seq)

		if p != want {
			t.Errorf("search(r{%d}.prices()) = %d, want %d", init, p, want)

		}
	})

	part2 := map[int]int{
		1:    7,
		2:    7,
		3:    0,
		2024: 9,
	}

	seqPart2 := [4]int{-2, 1, -1, 3}

	for init, want := range part2 {
		r := &rng{seed: init}
		ps, ds := r.prices()
		p := search(ps, ds, seqPart2)
		if p != want {
			t.Errorf("search(r{%d}.prices()) = %d, want %d", init, p, want)
		}
	}
}
