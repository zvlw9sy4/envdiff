package config

import (
	"flag"
	"strings"
)

// multiFlag is a flag.Value that collects repeated string flags.
type multiFlag []string

func (m *multiFlag) String() string  { return strings.Join(*m, ",") }
func (m *multiFlag) Set(v string) error { *m = append(*m, v); return nil }

// FromFlags builds a Config by parsing os.Args via the provided FlagSet.
// Call fs.Parse(args) before passing to this function, or use BuildFromArgs.
func FromFlags(fs *flag.FlagSet) *Config {
	c := &Config{}
	var labels multiFlag
	var filterKeys string

	fs.Var(&labels, "label", "label for each file (repeatable, must match file count)")
	fs.StringVar(&c.Output, "output", "text", "output format: text|json|csv|markdown|table|yaml")
	fs.BoolVar(&c.ShowAll, "all", false, "show all keys including matching ones")
	fs.BoolVar(&c.OnlyMissing, "only-missing", false, "show only missing keys")
	fs.BoolVar(&c.OnlyMismatch, "only-mismatch", false, "show only mismatched keys")
	fs.StringVar(&filterKeys, "keys", "", "comma-separated list of keys to include")
	fs.BoolVar(&c.NoColor, "no-color", false, "disable colored output")

	_ = fs.Parse(nil) // no-op if already parsed; caller should parse first

	c.Files = fs.Args()
	c.Labels = []string(labels)
	if filterKeys != "" {
		for _, k := range strings.Split(filterKeys, ",") {
			if t := strings.TrimSpace(k); t != "" {
				c.FilterKeys = append(c.FilterKeys, t)
			}
		}
	}
	return c
}

// BuildFromArgs parses args into a Config using a new FlagSet.
func BuildFromArgs(args []string) (*Config, error) {
	fs := flag.NewFlagSet("envdiff", flag.ContinueOnError)
	c := &Config{}
	var labels multiFlag
	var filterKeys string

	fs.Var(&labels, "label", "label for each file")
	fs.StringVar(&c.Output, "output", "text", "output format")
	fs.BoolVar(&c.ShowAll, "all", false, "show all keys")
	fs.BoolVar(&c.OnlyMissing, "only-missing", false, "only missing keys")
	fs.BoolVar(&c.OnlyMismatch, "only-mismatch", false, "only mismatched keys")
	fs.StringVar(&filterKeys, "keys", "", "comma-separated keys to include")
	fs.BoolVar(&c.NoColor, "no-color", false, "disable color")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}
	c.Files = fs.Args()
	c.Labels = []string(labels)
	if filterKeys != "" {
		for _, k := range strings.Split(filterKeys, ",") {
			if t := strings.TrimSpace(k); t != "" {
				c.FilterKeys = append(c.FilterKeys, t)
			}
		}
	}
	return c, nil
}
