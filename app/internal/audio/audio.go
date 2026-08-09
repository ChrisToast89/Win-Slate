// Package audio — local DSP fingerprint after ffmpeg decode.
// Parity with slate-0.3.2/src/main/audio.ts
package audio

import (
	"bytes"
	"fmt"
	"math"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/wassermanproductions/slate-windows/internal/types"
)

const (
	sampleRate = 16000
	frame      = 512
	maxSeconds = 90
)

func decodePCM(path string) (pcm []int16, fullDurationSec float64, err error) {
	cmd := exec.Command("ffmpeg",
		"-i", path,
		"-map", "a:0",
		"-ac", "1",
		"-ar", strconv.Itoa(sampleRate),
		"-t", strconv.Itoa(maxSeconds),
		"-f", "s16le",
		"pipe:1",
	)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	runErr := cmd.Run()
	raw := stdout.Bytes()
	if runErr != nil && len(raw) == 0 {
		errStr := stderr.String()
		if strings.Contains(errStr, "does not contain any stream") || strings.Contains(errStr, "Stream map") {
			return nil, 0, fmt.Errorf("No audio track found in this file.")
		}
		return nil, 0, fmt.Errorf("Could not decode audio (ffmpeg): %v", runErr)
	}
	pcm = make([]int16, len(raw)/2)
	for i := 0; i+1 < len(raw); i += 2 {
		pcm[i/2] = int16(raw[i]) | int16(raw[i+1])<<8
	}
	re := regexp.MustCompile(`Duration:\s*(\d+):(\d+):(\d+\.?\d*)`)
	if m := re.FindStringSubmatch(stderr.String()); len(m) == 4 {
		h, _ := strconv.ParseFloat(m[1], 64)
		mi, _ := strconv.ParseFloat(m[2], 64)
		s, _ := strconv.ParseFloat(m[3], 64)
		fullDurationSec = h*3600 + mi*60 + s
	} else {
		fullDurationSec = float64(len(pcm)) / sampleRate
	}
	return pcm, fullDurationSec, nil
}

func frameStats(pcm []int16) (rms, zcr []float64) {
	n := len(pcm) / frame
	rms = make([]float64, n)
	zcr = make([]float64, n)
	for f := 0; f < n; f++ {
		var sum float64
		var crossings int
		base := f * frame
		for i := 0; i < frame; i++ {
			s := float64(pcm[base+i]) / 32768
			sum += s * s
			if i > 0 && (pcm[base+i-1] >= 0) != (pcm[base+i] >= 0) {
				crossings++
			}
		}
		rms[f] = math.Sqrt(sum / frame)
		zcr[f] = float64(crossings) / frame
	}
	return rms, zcr
}

func windowF0(pcm []int16, start int) *float64 {
	W := 1024
	if start+W*2 > len(pcm) {
		return nil
	}
	minLag := sampleRate / 400
	maxLag := sampleRate / 50
	var energy float64
	for i := 0; i < W; i++ {
		v := float64(pcm[start+i]) / 32768
		energy += v * v
	}
	if energy/float64(W) < 1e-5 {
		return nil
	}
	bestCorr := 0.0
	corrs := make([]float64, maxLag+1)
	for lag := minLag; lag <= maxLag; lag++ {
		var corr float64
		for i := 0; i < W; i++ {
			corr += (float64(pcm[start+i]) / 32768) * (float64(pcm[start+i+lag]) / 32768)
		}
		corr /= energy
		corrs[lag] = corr
		if corr > bestCorr {
			bestCorr = corr
		}
	}
	if bestCorr <= 0.45 {
		return nil
	}
	for lag := minLag; lag <= maxLag; lag++ {
		prev := 0.0
		if lag > 0 {
			prev = corrs[lag-1]
		}
		next := 0.0
		if lag+1 <= maxLag {
			next = corrs[lag+1]
		}
		if corrs[lag] >= bestCorr*0.9 && corrs[lag] >= prev && corrs[lag] >= next {
			v := float64(sampleRate) / float64(lag)
			return &v
		}
	}
	return nil
}

func estimateBPM(rms []float64) (bpm *float64, confidence string) {
	framesPerSec := float64(sampleRate) / frame
	onsets := make([]float64, len(rms))
	for i := 1; i < len(rms); i++ {
		onsets[i] = math.Max(0, rms[i]-rms[i-1])
	}
	var mean float64
	for _, v := range onsets {
		mean += v
	}
	mean /= float64(len(onsets))
	if mean < 1e-6 {
		return nil, "low"
	}
	for i := range onsets {
		onsets[i] -= mean
	}
	minLag := int(math.Round((60.0 / 190.0) * framesPerSec))
	maxLag := int(math.Round((60.0 / 60.0) * framesPerSec))
	bestLag := 0
	bestCorr := math.Inf(-1)
	var norm float64
	for _, v := range onsets {
		norm += v * v
	}
	if norm == 0 {
		return nil, "low"
	}
	for lag := minLag; lag <= maxLag; lag++ {
		var corr float64
		for i := 0; i+lag < len(onsets); i++ {
			corr += onsets[i] * onsets[i+lag]
		}
		corr /= norm
		if corr > bestCorr {
			bestCorr = corr
			bestLag = lag
		}
	}
	if bestLag == 0 {
		return nil, "low"
	}
	b := math.Round((60 * framesPerSec) / float64(bestLag))
	if bestCorr > 0.35 {
		confidence = "high"
	} else if bestCorr > 0.18 {
		confidence = "medium"
	} else {
		confidence = "low"
	}
	if confidence == "low" {
		return nil, confidence
	}
	return &b, confidence
}

func describeArc(rms []float64) string {
	thirds := 3
	seg := len(rms) / thirds
	if seg == 0 {
		return "too short to read"
	}
	levels := make([]float64, thirds)
	for t := 0; t < thirds; t++ {
		var s float64
		for i := t * seg; i < (t+1)*seg; i++ {
			s += rms[i]
		}
		levels[t] = s / float64(seg)
	}
	a, b, c := levels[0], levels[1], levels[2]
	rel := func(x, y float64) string {
		if x > y*1.35 {
			return ">"
		}
		if y > x*1.35 {
			return "<"
		}
		return "="
	}
	ab, bc := rel(a, b), rel(b, c)
	switch {
	case ab == "<" && bc == "<":
		return "steady build from quiet open to loud finish"
	case ab == ">" && bc == ">":
		return "loud open decaying to a quiet close"
	case ab == "<" && bc == ">":
		return "rises to a peak mid-way, then falls away"
	case ab == ">" && bc == "<":
		return "strong open, quiet middle, returns big at the end"
	default:
		return "relatively even energy throughout"
	}
}

func quantile(sorted []float64, p float64) float64 {
	if len(sorted) == 0 {
		return 1e-6
	}
	i := int(math.Floor(float64(len(sorted)) * p))
	if i >= len(sorted) {
		i = len(sorted) - 1
	}
	return sorted[i]
}

// ComputeFingerprint matches the TypeScript implementation.
func ComputeFingerprint(pcm []int16, fullDurationSec float64) types.AudioFingerprint {
	rms, zcr := frameStats(pcm)
	sorted := append([]float64{}, rms...)
	sort.Float64s(sorted)
	p95 := quantile(sorted, 0.95)
	if p95 == 0 {
		p95 = 1e-6
	}
	silenceThresh := math.Max(p95*0.06, 0.004)
	silent := 0
	longest := 0
	run := 0
	for _, v := range rms {
		if v < silenceThresh {
			silent++
			run++
			if run > longest {
				longest = run
			}
		} else {
			run = 0
		}
	}
	framesPerSec := float64(sampleRate) / frame

	var active []float64
	for _, v := range sorted {
		if v >= silenceThresh {
			active = append(active, v)
		}
	}
	dynamicRangeDb := 0.0
	if len(active) > 4 {
		dynamicRangeDb = 20 * math.Log10(quantile(active, 0.95)/math.Max(quantile(active, 0.2), 1e-6))
	}

	var f0s []float64
	hop := sampleRate / 4
	for s := 0; s+2048 < len(pcm); s += hop {
		frameIdx := s / frame
		if frameIdx >= len(rms) || rms[frameIdx] < silenceThresh {
			continue
		}
		if f0 := windowF0(pcm, s); f0 != nil {
			f0s = append(f0s, *f0)
		}
	}
	sort.Float64s(f0s)
	var pitchMedian *float64
	var pitchSpread *float64
	if len(f0s) > 6 {
		m := f0s[len(f0s)/2]
		rounded := math.Round(m)
		pitchMedian = &rounded
	}
	if pitchMedian != nil && len(f0s) > 10 {
		lo := f0s[int(float64(len(f0s))*0.1)]
		hi := f0s[int(float64(len(f0s))*0.9)]
		sp := math.Round(12 * math.Log2(hi/lo))
		pitchSpread = &sp
	}

	var zsum float64
	var zn int
	for i := range zcr {
		if rms[i] >= silenceThresh {
			zsum += zcr[i]
			zn++
		}
	}
	meanZcr := 0.0
	if zn > 0 {
		meanZcr = zsum / float64(zn)
	}
	brightness := "balanced"
	if meanZcr < 0.06 {
		brightness = "dark"
	} else if meanZcr > 0.14 {
		brightness = "bright"
	}

	bpm, conf := estimateBPM(rms)
	maxF0Slots := math.Max(1, float64(len(pcm)/hop))

	return types.AudioFingerprint{
		DurationSec:          math.Round(fullDurationSec*10) / 10,
		SampledSec:           math.Round((float64(len(pcm))/sampleRate)*10) / 10,
		BpmEstimate:          bpm,
		BpmConfidence:        conf,
		PitchMedianHz:        pitchMedian,
		PitchSpreadSemitones: pitchSpread,
		VoicedRatio:          math.Round((float64(len(f0s))/maxF0Slots)*100) / 100,
		Brightness:           brightness,
		DynamicRangeDb:       math.Round(dynamicRangeDb*10) / 10,
		EnergyArc:            describeArc(rms),
		SilenceRatio:         math.Round((float64(silent)/float64(len(rms)))*100) / 100,
		LongestSilenceSec:    math.Round((float64(longest)/framesPerSec)*10) / 10,
	}
}

// Analyze decodes and fingerprints an audio/video file.
func Analyze(path string) (types.AudioFingerprint, error) {
	pcm, fullDur, err := decodePCM(path)
	if err != nil {
		return types.AudioFingerprint{}, err
	}
	if len(pcm) < sampleRate {
		return types.AudioFingerprint{}, fmt.Errorf("Audio is too short to analyze (need at least 1 second).")
	}
	return ComputeFingerprint(pcm, fullDur), nil
}
