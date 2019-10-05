package pattern

import "testing"

func TestMatch(t *testing.T) {
	tests := []struct {
		ps  []string
		ok  []string
		bad []string
	}{
		// https://help.github.com/en/articles/workflow-syntax-for-github-actions#patterns-to-match-file-paths
		{
			ps: []string{`*`},
			ok: []string{`README.md`, ``},
		},
		{
			ps:  []string{`*.jsx?`},
			ok:  []string{`page.js`, `page.jsx`},
			bad: []string{`page.jsxx`},
		},
		{
			ps: []string{`**`},
			ok: []string{`all/the/files.md`, ``},
		},
		{
			ps:  []string{`*.js`},
			ok:  []string{`app.js`, `index.js`},
			bad: []string{`js/index.js`, `src/js/index.js`, `main.go`},
		},
		{
			ps:  []string{`**.js`},
			ok:  []string{`index.js`, `js/index.js`, `src/js/index.js`},
			bad: []string{`main.go`},
		},
		{
			ps:  []string{`docs/*`},
			ok:  []string{`docs/README.md`, `docs/file.txt`},
			bad: []string{`README.md`, `docs/mona/octocat.txt`, `dir/docs/my-file.txt`},
		},
		{
			ps:  []string{`docs/**`},
			ok:  []string{`docs/README.md`, `docs/mona/octocat.txt`},
			bad: []string{`README.md`, `dir/docs/my-file.txt`},
		},
		{
			ps:  []string{`docs/**/*.md`},
			ok:  []string{`docs/README.md`, `docs/mona/hello-world.md`, `docs/a/markdown/file.md`},
			bad: []string{`README.md`, `docs/mona/octocat.txt`},
		},
		{
			ps:  []string{`**/docs/**`},
			ok:  []string{`/docs/hello.md`, `dir/docs/my-file.txt`, `space/docs/plan/space.doc`},
			bad: []string{`README.md`, `hoge-docs/a`},
		},
		{
			ps:  []string{`**/README.md`},
			ok:  []string{`README.md`, `js/README.md`},
			bad: []string{`README.txt`},
		},
		{
			ps: []string{`**/*src/**`},
			ok: []string{`a/src/app.js`, `my-src/code/js/app.js`},
		},
		{
			ps: []string{`**/*-post.md`},
			ok: []string{`my-post.md`, `math/their-post.md`},
		},
		{
			ps: []string{`**/migrate-*.sql`},
			ok: []string{`migrate-10909.sql`, `db/migrate-v1.0.sql`, `db/sept/migrate-v1.sql`},
		},
		{
			ps:  []string{`*.md`, `!README.md`},
			ok:  []string{`hello.md`},
			bad: []string{`README.md`, `docs/hello.md`},
		},
		{
			ps:  []string{`*.md`, `!README.md`, `README*`},
			ok:  []string{`hello.md`, `README.md`, `README.doc`},
			bad: []string{`docs/hello.md`},
		},
		// https://help.github.com/en/articles/workflow-syntax-for-github-actions#patterns-to-match-branches-and-tags
		{
			ps:  []string{`feature/*`},
			ok:  []string{`feature/my-branch`, `feature/your-branch`},
			bad: []string{`test-feature`, `feature/beta-a/my-branch`},
		},
		{
			ps:  []string{`feature/**`},
			ok:  []string{`feature/beta-a/my-branch`, `feature/your-branch`, `feature/mona/the/octoca`},
			bad: []string{`test-feature`},
		},
		{
			ps:  []string{`master`, `releases/mona-the-octcat`},
			ok:  []string{`master`, `releases/mona-the-octcat`},
			bad: []string{`test-feature`},
		},
		{
			ps:  []string{`*`},
			ok:  []string{`master`, `releases`},
			bad: []string{`releases/mona-the-octcat`},
		},
		{
			ps: []string{`**`},
			ok: []string{`all/the/branches`, `every/tag`},
		},
		{
			ps: []string{`*feature`},
			ok: []string{`mona-feature`, `feature`, `ver-10-feature`},
		},
		{
			ps: []string{`v2*`},
			ok: []string{`v2`, `v2.0`, `v2.9`,
				`v20`, // expected?
			},
			bad: []string{`v1`},
		},
		{
			ps:  []string{`v[12].[0-9]+.[0-9]+`},
			ok:  []string{`v1.10.1`, `v2.0.0`},
			bad: []string{`v3.1.0`},
		},
	}

	for _, tt := range tests {
		for _, s := range tt.ok {
			b, err := Match(s, tt.ps)
			if err != nil {
				t.Errorf("Match(%q, %v) returns unexpected error: %v", s, tt.ps, err)
				continue
			}
			if !b {
				t.Errorf("Match(%q, %v) = false, want true", s, tt.ps)
				dumpPattern(t, tt.ps)
			}
		}

		for _, s := range tt.bad {
			b, err := Match(s, tt.ps)
			if err != nil {
				t.Errorf("Match(%q, %v) returns unexpected error: %v", s, tt.ps, err)
				continue
			}
			if b {
				t.Errorf("Match(%q, %v) = true, want false", s, tt.ps)
				dumpPattern(t, tt.ps)
			}
		}
	}
}

func dumpPattern(t *testing.T, patterns []string) {
	t.Helper()
	for _, p := range patterns {
		negate := ""
		pstr := p
		if len(p) > 0 && p[0] == '!' {
			pstr = p[1:]
			negate = "!"
		}
		t.Log("REGEX PATTERNS:")
		t.Logf("\t%s %s", negate, replace(pstr))
	}
}
