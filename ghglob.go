package ghglob

import (
	"fmt"
	"os"
	"strings"

	"github.com/haya14busa/ghglob/ghmatcher"
	"github.com/karrick/godirwalk"
)

type Option struct {
	// True for faster yet non-deterministic enumeration.
	Sort                bool
	FollowSymbolicLinks bool
	Root                string
}

func GlobList(patterns []string, opt Option) (files []string, err error) {
	ch := make(chan string)
	go func() {
		err = Glob(ch, patterns, opt)
	}()
	for f := range ch {
		files = append(files, f)
	}
	return files, nil
}

func Glob(files chan<- string, patterns []string, opt Option) error {
	defer close(files)
	matcher, err := ghmatcher.New(patterns)
	if err != nil {
		return err
	}
	isRoot := len(opt.Root) > 0 && opt.Root[0] == '/'
	subms, err := buildSubMatchers(patterns)
	if err != nil {
		return fmt.Errorf("fail to build submatchers: %v", err)
	}
	root := opt.Root
	if root == "" {
		root = "."
	}
	rootPrefix := strings.TrimPrefix(root, "./")
	if err := godirwalk.Walk(root, &godirwalk.Options{
		Callback: func(path string, de *godirwalk.Dirent) error {
			p := path
			if !isRoot {
				p = strings.TrimPrefix(path, rootPrefix)
			}
			if p == "" {
				return nil
			}
			if de.ModeType().IsDir() {
				if p != strings.TrimSuffix(rootPrefix, "/") && p != "/" && shouldSkipDir(subms, p) {
					return skipdir{}
				}
				return nil
			}
			if !matcher.Match(p) {
				return nil
			}
			files <- p
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

func buildSubMatchers(patterns []string) ([]*ghmatcher.Matcher, error) {
	var ms []*ghmatcher.Matcher
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
			m, err := ghmatcher.New([]string{ghpattern})
			if err != nil {
				return nil, err
			}
			ms = append(ms, m)
		}
	}
	return ms, nil
}

func shouldSkipDir(ms []*ghmatcher.Matcher, path string) bool {
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
