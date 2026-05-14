package diff

// OverlapResult holds the analysis of key overlap between two environments.
type OverlapResult struct {
	EnvA        string
	EnvB        string
	OnlyInA     []string
	OnlyInB     []string
	InBoth      []string
	JaccardIndex float64
}

// Overlap computes the key overlap between two named env maps.
// It returns which keys are exclusive to each side, which are shared,
// and a Jaccard similarity index (0.0–1.0).
func Overlap(labelA string, envA map[string]string, labelB string, envB map[string]string) OverlapResult {
	setA := make(map[string]struct{}, len(envA))
	for k := range envA {
		setA[k] = struct{}{}
	}

	setB := make(map[string]struct{}, len(envB))
	for k := range envB {
		setB[k] = struct{}{}
	}

	var onlyA, onlyB, both []string

	for k := range setA {
		if _, ok := setB[k]; ok {
			both = append(both, k)
		} else {
			onlyA = append(onlyA, k)
		}
	}

	for k := range setB {
		if _, ok := setA[k]; !ok {
			onlyB = append(onlyB, k)
		}
	}

	sortStrings(onlyA)
	sortStrings(onlyB)
	sortStrings(both)

	unionSize := len(onlyA) + len(onlyB) + len(both)
	var jaccard float64
	if unionSize > 0 {
		jaccard = float64(len(both)) / float64(unionSize)
	}

	return OverlapResult{
		EnvA:        labelA,
		EnvB:        labelB,
		OnlyInA:     onlyA,
		OnlyInB:     onlyB,
		InBoth:      both,
		JaccardIndex: jaccard,
	}
}

func sortStrings(s []string) {
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && s[j] < s[j-1]; j-- {
			s[j], s[j-1] = s[j-1], s[j]
		}
	}
}
