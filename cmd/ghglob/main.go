package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/haya14busa/ghglob"
	"github.com/haya14busa/ghglob/ghmatcher"
	"github.com/mattn/go-isatty"
)

var (
	all       = flag.Bool("all", false, "do not ignore entries starting with .")
	sort      = flag.Bool("sort", false, "sort results.")
	followSym = flag.Bool("symlink", true, "follow symlink if true")
)

func usage() {
	fmt.Fprintln(os.Stderr, "Usage: ghglob [FLAGS] [PATTERN]...")
	fmt.Fprintln(os.Stderr, "\tIf STDIN is provided, read it and return matched results.")
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

	if isatty.IsTerminal(os.Stdin.Fd()) {
		glob()
	} else {
		filter()
	}
}

func glob() {
	ps := flag.Args()
	if !*all && shouldIgnoreDot(ps) {
		ps = append(ps, "!**/.*")
		ps = append(ps, "!**/.*/**")
	}
	opt := ghglob.Option{Sort: *sort, FollowSymbolicLinks: *followSym}
	if shouldRoot(flag.Args()) {
		opt.Root = "/"
	}
	files := make(chan string, 100)
	go func() {
		if err := ghglob.Glob(files, ps, opt); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}()
	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()
	for f := range files {
		fmt.Fprintln(w, f)
	}
}

func filter() {
	s := bufio.NewScanner(os.Stdin)
	m, err := ghmatcher.New(flag.Args())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	for s.Scan() {
		t := s.Text()
		if m.Match(t) {
			fmt.Fprintln(os.Stdout, t)
		}
	}
}

// Do not ignore dot files when -all=false for special cases that given
// patterns contains .
func shouldIgnoreDot(patterns []string) bool {
	for _, p := range patterns {
		if len(p) > 0 && p[0] == '.' || strings.Contains(p, "/.") {
			return false
		}
	}
	return true
}

func shouldRoot(patterns []string) bool {
	for _, p := range patterns {
		if len(p) > 0 && p[0] != '/' {
			fmt.Println(p, p[0])
			return false
		} else if len(p) > 1 && p[0] == '!' && p[1] != '/' {
			return false
		}
	}
	return true
}
