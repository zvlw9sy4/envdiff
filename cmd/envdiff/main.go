package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/filter"
	"github.com/user/envdiff/internal/loader"
	"github.com/user/envdiff/internal/output"
)

func main() {
	var (
		labelsFlag   = flag.String("labels", "", "comma-separated labels for each file (e.g. dev,prod)")
		formatFlag   = flag.String("format", "text", "output format: text, json, csv")
		showAll      = flag.Bool("all", false, "show all keys including matches")
		onlyMissing  = flag.Bool("missing", false, "show only missing keys")
		onlyMismatch = flag.Bool("mismatch", false, "show only mismatched keys")
		keyFilter    = flag.String("keys", "", "comma-separated list of keys to include")
	)
	flag.Parse()

	paths := flag.Args()
	if len(paths) < 2 {
		fmt.Fprintln(os.Stderr, "usage: envdiff [flags] <file1> <file2> [file3...]")
		flag.PrintDefaults()
		os.Exit(1)
	}

	var labels []string
	if *labelsFlag != "" {
		labels = strings.Split(*labelsFlag, ",")
	}

	files, err := loader.LoadFiles(paths, labels)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading files: %v\n", err)
		os.Exit(1)
	}

	envLabels, envMaps := loader.ToEnvMaps(files)
	results := diff.Compare(envLabels, envMaps)

	opts := filter.Options{
		ShowAll:      *showAll,
		OnlyMissing:  *onlyMissing,
		OnlyMismatch: *onlyMismatch,
	}
	statusFilter := opts.ToStatusFilter()
	results = filter.Apply(results, statusFilter)

	if *keyFilter != "" {
		keys := strings.Split(*keyFilter, ",")
		results = filter.ByKey(results, keys)
	}

	fmtr, err := output.NewFormatter(*formatFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if err := fmtr.Write(os.Stdout, results); err != nil {
		fmt.Fprintf(os.Stderr, "error writing output: %v\n", err)
		os.Exit(1)
	}
}
