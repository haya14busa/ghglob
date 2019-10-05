// Package pattern implements pattern matcher which imitates GitHub Actions
// filter patterns (UNOFFICIAL). https://help.github.com/en/articles/workflow-syntax-for-github-actions#filter-pattern-cheat-sheet.
package pattern

import (
	"regexp"
	"strings"
)

// Match returns true if given string (s) matches the given GitHub Actions
// filter patterns [1].
//
// - `*` matches zero or more characters, but does not match the / character
// - `**` matches zero or more of any character
// - `?` matches zero or one of the proceeding character
// - `+` matches one or more of the proceeding character
// - `[]` matches any character listed, or included in ranges. Ranges can only include a-zA-Z0-9. e.g [123abc] or [0-9a-f]
// - `!` at the start of a pattern makes it negate previous positive patterns. It has no special meaning if not the first character
//
// [1]: https://help.github.com/en/articles/workflow-syntax-for-github-actions#filter-pattern-cheat-sheet
func Match(s string, patterns []string) (bool, error) {
	ms, err := buildMatchers(patterns)
	if err != nil {
		return false, err
	}
	ok := false
	for _, m := range ms {
		if m.r.MatchString(s) {
			if m.negate {
				ok = false
			} else {
				ok = true
			}
		}
	}
	return ok, nil
}

type matcher struct {
	r      *regexp.Regexp
	negate bool
}

func buildMatchers(patterns []string) ([]matcher, error) {
	ms := make([]matcher, len(patterns))
	for i, p := range patterns {
		m := matcher{}
		pstr := p
		if len(p) > 0 && p[0] == '!' {
			pstr = p[1:]
			m.negate = true
		}
		r, err := buildRegex(pstr)
		if err != nil {
			return nil, err
		}
		m.r = r
		ms[i] = m
	}
	return ms, nil
}

func buildRegex(pattern string) (*regexp.Regexp, error) {
	return regexp.Compile(replace(pattern))
}

// replace replaces given GitHub Actions pattern to regex string.
func replace(pattern string) string {
	// https://help.github.com/en/articles/workflow-syntax-for-github-actions#filter-pattern-cheat-sheet
	r := strings.NewReplacer(
		// Special patterns.
		`**/*`, `.*`,
		`**/`, `.*(^|/)`,

		`**`, `.*`,
		`*`, `[^/]*`,
		`*`, `[^/]*`,
		`.`, `\.`,
		`(`, `\(`,
		`)`, `\)`,
		`|`, `\|`,
		`{`, `\{`,
		`}`, `\}`,
		`^`, `\^`,
		`$`, `\$`,
		`\`, `\\`,
	)
	return "^" + r.Replace(pattern) + "$"
}
