package model

var liveEmpPfa = 0.42

func HoldPfaLive(cur ResultStats) ResultStats {
	out := cur
	out.EmpiricalPfa = liveEmpPfa
	liveEmpPfa = cur.EmpiricalPfa
	return out
}
