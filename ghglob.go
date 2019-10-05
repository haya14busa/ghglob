package ghglob

import (
	"fmt"
	"os"
	"strings"

	"github.com/haya14busa/ghglob/pattern"
	"github.com/karrick/godirwalk"
)

type Option struct {
	// True for faster yet non-deterministic enumeration.
	Sort                bool
	FollowSymbolicLinks bool
}

func Glob(files chan<- string, patterns []string, opt Option) error {
	defer close(files)
	matcher, err := pattern.New(patterns)
	if err != nil {
		return err
	}
	subms, err := buildSubMatchers(patterns)
	if err != nil {
		return fmt.Errorf("fail to build submatchers: %v", err)
	}
	if err := godirwalk.Walk(".", &godirwalk.Options{
		Callback: func(path string, de *godirwalk.Dirent) error {
			if de.ModeType().IsDir() {
				if shouldSkipDir(subms, path) {
					return skipdir{}
				}
				return nil
			}
			if !matcher.Match(path) {
				return nil
			}
			files <- path
			return nil
		},
		ErrorCallback: func(path string, err error) godirwalk.ErrorAction {
			if _, ok := err.(skipdir); ok {
				return godirwalk.SkipNode
			}
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			return godirwalk.Halt
		},

		Unsorted:            !opt.Sort,
		FollowSymbolicLinks: opt.FollowSymbolicLinks,
	}); err != nil {
		return err
	}
	return nil
}

type skipdir struct{}

func (skipdir) Error() string {
	return "skipdir"
}

func buildSubMatchers(patterns []string) ([]*pattern.Matcher, error) {
	var ms []*pattern.Matcher
	for _, p := range patterns {
		if len(p) > 0 && p[0] == '!' {
			continue
		}
		seps := strings.Split(p, "/")
		for i := range seps {
			ghpattern := strings.Join(seps[:i+1], "/")
			if strings.Contains(seps[i], "**") {
				ghpattern = strings.Join(append(seps[:i], "**"), "/")
			}
			m, err := pattern.New([]string{ghpattern})
			if err != nil {
				return nil, err
			}
			ms = append(ms, m)
		}
	}
	return ms, nil
}

func shouldSkipDir(ms []*pattern.Matcher, path string) bool {
	if len(ms) == 0 || path == "." {
		return false
	}
	for _, m := range ms {
		if m.Match(path) {
			return false
		}
	}
	return true
}
