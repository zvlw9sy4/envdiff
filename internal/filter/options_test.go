package filter_test

import (
	"testing"

	"github.com/user/envdiff/internal/filter"
)

func TestOptions_ToStatusFilter_Default(t *testing.T) {
	o := filter.Options{}
	f := o.ToStatusFilter()
	if !f.IncludeMissing || !f.IncludeMismatch || f.IncludeMatch {
		t.Errorf("default options should return DefaultFilter, got %+v", f)
	}
}

func TestOptions_ToStatusFilter_ShowAll(t *testing.T) {
	o := filter.Options{ShowAll: true}
	f := o.ToStatusFilter()
	if !f.IncludeMissing || !f.IncludeMismatch || !f.IncludeMatch {
		t.Errorf("ShowAll should return AllFilter, got %+v", f)
	}
}

func TestOptions_ToStatusFilter_OnlyMissing(t *testing.T) {
	o := filter.Options{OnlyMissing: true, ShowAll: true}
	f := o.ToStatusFilter()
	if !f.IncludeMissing || f.IncludeMismatch || f.IncludeMatch {
		t.Errorf("OnlyMissing should override ShowAll, got %+v", f)
	}
}

func TestOptions_ToStatusFilter_OnlyMismatch(t *testing.T) {
	o := filter.Options{OnlyMismatch: true, ShowAll: true}
	f := o.ToStatusFilter()
	if f.IncludeMissing || !f.IncludeMismatch || f.IncludeMatch {
		t.Errorf("OnlyMismatch should override ShowAll, got %+v", f)
	}
}

func TestOptions_ToStatusFilter_BothOnlyFlags(t *testing.T) {
	// When both OnlyMissing and OnlyMismatch are set, each respective field
	// should be included and IncludeMatch should remain false.
	o := filter.Options{OnlyMissing: true, OnlyMismatch: true}
	f := o.ToStatusFilter()
	if !f.IncludeMissing || !f.IncludeMismatch || f.IncludeMatch {
		t.Errorf("OnlyMissing+OnlyMismatch should include missing and mismatch only, got %+v", f)
	}
}
