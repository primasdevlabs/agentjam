package runtime

// PrecedenceCandidate represents a rule or setting candidate.
type PrecedenceCandidate struct {
	Level      string `json:"level"`
	SourceName string `json:"sourceName"`
	Value      string `json:"value"`
	Priority   int    `json:"priority"`
}

// PrecedenceOrder defines precedence weights (1 = highest).
var PrecedenceOrder = map[string]int{
	"Current Codebase State":                       1,
	"Project Configuration (.agentjam/config.yaml)": 2,
	"User explicit instruction":                    3,
	"Pinned Rules":                                 4,
	"Authoritative Documentation":                  5,
	"Environment defaults":                         6,
}

// ResolveHighest returns candidate with highest precedence weight.
func ResolveHighest(candidates []PrecedenceCandidate) PrecedenceCandidate {
	if len(candidates) == 0 {
		return PrecedenceCandidate{}
	}

	best := candidates[0]
	bestWeight := 999
	if w, ok := PrecedenceOrder[best.Level]; ok {
		bestWeight = w
	}

	for _, c := range candidates[1:] {
		w := 999
		if val, ok := PrecedenceOrder[c.Level]; ok {
			w = val
		}
		if w < bestWeight {
			best = c
			bestWeight = w
		}
	}

	return best
}
