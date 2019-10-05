// Package pattern implements pattern matcher which imitates GitHub Actions
// filter patterns (UNOFFICIAL). https://help.github.com/en/articles/workflow-syntax-for-github-actions#filter-pattern-cheat-sheet.
package pattern

import (
	"regexp"
	"strings"
)

type Matcher struct {
	ms []matcher
}

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
func (matcher *Matcher) Match(s string) bool {
	ok := false
	for _, m := range matcher.ms {
		if m.r.MatchString(s) {
			ok = !m.negate
		}
	}
	return ok
}

type matcher struct {
	r      *regexp.Regexp
	negate bool
}

// New creates new Matcher.
func New(patterns []string) (*Matcher, error) {
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
	return &Matcher{ms: ms}, nil
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
		`**/`, `(|.*(^|/))`,

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
