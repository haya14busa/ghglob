package ghglob

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/haya14busa/ghglob/ghmatcher"
	"github.com/s12chung/fastwalk"
)

type Option struct {
	FollowSymbolicLinks bool
	Root                string
}

func GlobList(patterns []string, opt Option) ([]string, error) {
	ch := make(chan string)
	errCh := make(chan error)
	go func() {
		errCh <- Glob(ch, patterns, opt)
	}()
	var files []string
	for f := range ch {
		files = append(files, f)
	}
	if err := <-errCh; err != nil {
		return nil, err
	}
	return files, nil
}

func Glob(files chan<- string, patterns []string, opt Option) error {
	defer close(files)
	matcher, err := ghmatcher.New(patterns)
	if err != nil {
		return err
	}
	subms, err := buildSubMatchers(patterns)
	if err != nil {
		return fmt.Errorf("fail to build submatchers: %v", err)
	}
	root := opt.Root
	rootPrefix := ""
	switch root {
	case "":
		root = "."
		rootPrefix = "./"
	case "/":
		rootPrefix = "/"
	default:
		rootPrefix = root + "/"
	}

	return fastwalk.Walk(root, func(path string, typ os.FileMode) error {
		p := strings.TrimPrefix(path, rootPrefix)
		if p == "" {
			return nil
		}

		if opt.FollowSymbolicLinks && typ == os.ModeSymlink {
			followedPath, err := filepath.EvalSymlinks(path)
			if err == nil {
				fi, err := os.Lstat(followedPath)
				if err == nil && fi.IsDir() {
					return fastwalk.TraverseLink
				}
			}
		}

		if typ.IsDir() {
			if p != strings.TrimSuffix(rootPrefix, "/") && shouldSkipDir(subms, p) {
				return filepath.SkipDir
			}
			return nil
		}
		if !matcher.Match(p) {
			return nil
		}
		files <- p
		return nil
	})
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
	if len(ms) == 0 || path == "." || path == "/" {
		return false
	}
	for _, m := range ms {
		if m.Match(path) {
			return false
		}
	}
	return true
}
