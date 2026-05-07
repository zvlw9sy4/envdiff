package filter

// Options holds all user-configurable filtering options.
type Options struct {
	// ShowAll includes matching keys in the output.
	ShowAll bool

	// OnlyMissing restricts output to missing keys only.
	OnlyMissing bool

	// OnlyMismatch restricts output to mismatched keys only.
	OnlyMismatch bool

	// Keys restricts output to a specific set of key names.
	Keys []string
}

// ToStatusFilter converts Options into a StatusFilter.
// OnlyMissing and OnlyMismatch take precedence over ShowAll.
func (o Options) ToStatusFilter() StatusFilter {
	if o.OnlyMissing {
		return StatusFilter{IncludeMissing: true}
	}
	if o.OnlyMismatch {
		return StatusFilter{IncludeMismatch: true}
	}
	if o.ShowAll {
		return AllFilter()
	}
	return DefaultFilter()
}
