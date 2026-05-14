package diff

import (
	"testing"
)

func TestOverlap_AllShared(t *testing.T) {
	a := map[string]string{"FOO": "1", "BAR": "2"}
	b := map[string]string{"FOO": "x", "BAR": "y"}

	r := Overlap("dev", a, "prod", b)

	if len(r.InBoth) != 2 {
		t.Errorf("expected 2 shared keys, got %d", len(r.InBoth))
	}
	if len(r.OnlyInA) != 0 {
		t.Errorf("expected 0 only-in-A keys, got %d", len(r.OnlyInA))
	}
	if r.JaccardIndex != 1.0 {
		t.Errorf("expected Jaccard=1.0, got %f", r.JaccardIndex)
	}
}

func TestOverlap_NoShared(t *testing.T) {
	a := map[string]string{"FOO": "1"}
	b := map[string]string{"BAR": "2"}

	r := Overlap("dev", a, "prod", b)

	if len(r.InBoth) != 0 {
		t.Errorf("expected 0 shared keys, got %d", len(r.InBoth))
	}
	if r.JaccardIndex != 0.0 {
		t.Errorf("expected Jaccard=0.0, got %f", r.JaccardIndex)
	}
}

func TestOverlap_PartialOverlap(t *testing.T) {
	a := map[string]string{"FOO": "1", "BAR": "2", "BAZ": "3"}
	b := map[string]string{"FOO": "x", "QUX": "y"}

	r := Overlap("dev", a, "prod", b)

	if len(r.InBoth) != 1 || r.InBoth[0] != "FOO" {
		t.Errorf("expected InBoth=[FOO], got %v", r.InBoth)
	}
	if len(r.OnlyInA) != 2 {
		t.Errorf("expected 2 only-in-A keys, got %d", len(r.OnlyInA))
	}
	if len(r.OnlyInB) != 1 || r.OnlyInB[0] != "QUX" {
		t.Errorf("expected OnlyInB=[QUX], got %v", r.OnlyInB)
	}
	// union=4, intersect=1 => 0.25
	expected := 0.25
	if r.JaccardIndex != expected {
		t.Errorf("expected Jaccard=%f, got %f", expected, r.JaccardIndex)
	}
}

func TestOverlap_EmptyInputs(t *testing.T) {
	r := Overlap("dev", map[string]string{}, "prod", map[string]string{})

	if r.JaccardIndex != 0.0 {
		t.Errorf("expected Jaccard=0.0 for empty inputs, got %f", r.JaccardIndex)
	}
	if len(r.InBoth) != 0 {
		t.Errorf("expected no shared keys")
	}
}

func TestOverlap_LabelsSet(t *testing.T) {
	r := Overlap("alpha", map[string]string{"X": "1"}, "beta", map[string]string{"X": "2"})

	if r.EnvA != "alpha" {
		t.Errorf("expected EnvA=alpha, got %s", r.EnvA)
	}
	if r.EnvB != "beta" {
		t.Errorf("expected EnvB=beta, got %s", r.EnvB)
	}
}

func TestOverlap_SortedOutput(t *testing.T) {
	a := map[string]string{"ZZZ": "1", "AAA": "2", "MMM": "3"}
	b := map[string]string{"QQQ": "x"}

	r := Overlap("dev", a, "prod", b)

	for i := 1; i < len(r.OnlyInA); i++ {
		if r.OnlyInA[i] < r.OnlyInA[i-1] {
			t.Errorf("OnlyInA not sorted: %v", r.OnlyInA)
		}
	}
}
