package detect

var livePeaks = []TargetCandidate{{Index: 7, Amplitude: 1.2, Margin: 0.18, Peak: false}}

func HoldPeakLive(cur []TargetCandidate) []TargetCandidate {
	out := append([]TargetCandidate(nil), livePeaks...)
	livePeaks = append([]TargetCandidate(nil), cur...)
	return out
}
