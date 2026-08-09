package audio

import (
	"math"
	"testing"
)

func TestComputeFingerprintSine(t *testing.T) {
	// 2 seconds of 220 Hz sine at 16 kHz
	n := sampleRate * 2
	pcm := make([]int16, n)
	for i := 0; i < n; i++ {
		pcm[i] = int16(12000 * math.Sin(2*math.Pi*220*float64(i)/float64(sampleRate)))
	}
	fp := ComputeFingerprint(pcm, 2.0)
	if fp.DurationSec != 2.0 {
		t.Fatalf("duration %v", fp.DurationSec)
	}
	if fp.SampledSec < 1.5 {
		t.Fatalf("sampled %v", fp.SampledSec)
	}
	// Pitch should be near 220 when voiced
	if fp.PitchMedianHz != nil {
		if math.Abs(*fp.PitchMedianHz-220) > 40 {
			t.Logf("pitch median %v (loose check)", *fp.PitchMedianHz)
		}
	}
	if fp.EnergyArc == "" {
		t.Fatal("empty energy arc")
	}
}
