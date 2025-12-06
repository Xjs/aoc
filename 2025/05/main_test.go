package main

import (
	"github.com/Xjs/aoc/integer"
	"github.com/google/go-cmp/cmp"
	"testing"
)

func mustParseRanges(rawRanges []string) []integer.Range {
	var ranges []integer.Range
	for _, rawRange := range rawRanges {
		r, err := integer.ParseRange(rawRange)
		if err != nil {
			panic(err)
		}
		ranges = append(ranges, r)
	}
	return ranges
}

func TestMergeRanges(t *testing.T) {
	for _, tt := range []struct{ input, want []string }{
		{[]string{"3-5", "10-14", "16-20", "12-18"}, []string{"3-5", "10-20"}},
		{[]string{"15-19", "27-43", "82-84", "83-87", "47-86", "59-66", "20-22"}, []string{"15-19", "20-22", "27-43", "47-87"}},
	} {
		ranges := mustParseRanges(tt.input)

		wantRanges := mustParseRanges(tt.want)

		ranges = mergeRanges(ranges)

		if !cmp.Equal(ranges, wantRanges) {
			t.Errorf("ranges = %v, wantRanges = %v", ranges, wantRanges)
		}
	}
}
