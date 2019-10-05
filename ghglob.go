package ghglob

import (
	"github.com/haya14busa/ghglob/pattern"
	"github.com/karrick/godirwalk"
)

type Option struct {
	// True for faster yet non-deterministic enumeration.
	Sort bool
}

func Glob(patterns []string, opt Option) ([]string, error) {
	matcher, err := pattern.New(patterns)
	if err != nil {
		return nil, err
	}

	var files []string
	if err := godirwalk.Walk(".", &godirwalk.Options{
		Callback: func(path string, de *godirwalk.Dirent) error {
			if de.ModeType().IsDir() {
				return nil
			}
			if !matcher.Match(path) {
				return nil
			}
			files = append(files, path)
			return nil
		},
		Unsorted: !opt.Sort,
	}); err != nil {
		return nil, err
	}
	return files, nil
}
