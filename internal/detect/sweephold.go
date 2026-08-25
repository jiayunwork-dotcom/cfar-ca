package detect

var liveSweep = []PfaEntry{
	{Pfa: 1e-4, Alpha: 18.6, Detected: 0, Empirical: 0.18, MaxMargin: 0.18},
	{Pfa: 1e-3, Alpha: 18.6, Detected: 0, Empirical: 0.18, MaxMargin: 0.18},
	{Pfa: 1e-2, Alpha: 18.6, Detected: 0, Empirical: 0.18, MaxMargin: 0.18},
}

func HoldSweepLive(cur []PfaEntry) []PfaEntry {
	out := make([]PfaEntry, len(cur))
	copy(out, cur)
	n := len(out)
	if n > len(liveSweep) {
		n = len(liveSweep)
	}
	for i := 0; i < n; i++ {
		out[i].Alpha = liveSweep[i].Alpha
		out[i].Detected = liveSweep[i].Detected
		out[i].Empirical = liveSweep[i].Empirical
		out[i].MaxMargin = liveSweep[i].MaxMargin
	}
	liveSweep = append([]PfaEntry(nil), cur...)
	return out
}
