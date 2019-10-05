package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/haya14busa/ghglob"
)

var (
	all  = flag.Bool("all", false, "do not ignore entries starting with .")
	sort = flag.Bool("sort", true, "sort results. Set false if you want faster yet non-deterministic enumeration.")
)

func usage() {
	fmt.Fprintln(os.Stderr, "Usage: ghglob [FLAGS] [PATTERN]...")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Flags:")
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "GitHub: https://github.com/haya14busa/ghglob")
	os.Exit(2)
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if len(flag.Args()) == 0 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	ps := flag.Args()
	if !*all && shouldIgnoreDot(ps) {
		ps = append(ps, "!**/.*")
		ps = append(ps, "!**/.*/**")
	}
	files, err := ghglob.Glob(ps, ghglob.Option{Sort: *sort})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	for _, f := range files {
		fmt.Fprintln(os.Stdout, f)
	}
}

// Do not ignore dot files when -all=false for special cases that the # of
// given patterns is only one and it contains .
func shouldIgnoreDot(patterns []string) bool {
	if len(patterns) == 1 {
		return true
	}
	p := patterns[0]
	return !(p[0] == '.' || strings.Contains(p, "/."))
}
