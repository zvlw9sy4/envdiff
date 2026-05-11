package output

import (
	"encoding/xml"
	"fmt"
	"io"

	"github.com/user/envdiff/internal/diff"
)

type xmlFormatter struct{}

type xmlReport struct {
	XMLName xml.Name   `xml:"envdiff"`
	Results []xmlEntry `xml:"result"`
}

type xmlEntry struct {
	Key    string       `xml:"key,attr"`
	Status string       `xml:"status,attr"`
	Values []xmlEnvVal  `xml:"env"`
}

type xmlEnvVal struct {
	Name  string `xml:"name,attr"`
	Value string `xml:",chardata"`
}

func (f *xmlFormatter) Write(w io.Writer, results []diff.Result) error {
	if len(results) == 0 {
		_, err := fmt.Fprintln(w, `<?xml version="1.0" encoding="UTF-8"?><envdiff></envdiff>`)
		return err
	}

	report := xmlReport{}
	for _, r := range results {
		entry := xmlEntry{
			Key:    r.Key,
			Status: string(r.Status),
		}
		for env, val := range r.Values {
			entry.Values = append(entry.Values, xmlEnvVal{
				Name:  env,
				Value: val,
			})
		}
		report.Results = append(report.Results, entry)
	}

	_, err := fmt.Fprintln(w, `<?xml version="1.0" encoding="UTF-8"?>`)
	if err != nil {
		return err
	}

	enc := xml.NewEncoder(w)
	enc.Indent("", "  ")
	if err := enc.Encode(report); err != nil {
		return err
	}
	return enc.Flush()
}
