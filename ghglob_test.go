package ghglob

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGlobList(t *testing.T) {
	tests := []struct {
		ps   []string
		want []string
	}{
		{
			ps: []string{`*`},
			want: []string{
				`README.doc`,
				`README.md`,
				`README.txt`,
				`app.js`,
				`hello.md`,
				`index.js`,
				`main.go`,
				`master`,
				`migrate-10909.sql`,
				`mona-feature`,
				`my-post.md`,
				`page.js`,
				`page.jsx`,
				`page.jsxx`,
				`test-feature`,
				`v1`,
				`v1.10.1`,
				`v2`,
				`v2.0`,
				`v2.0.0`,
				`v2.9`,
				`v3.1.0`,
				`ver-10-feature`,
			},
		},
		{
			ps: []string{`**`},
			want: []string{
				"README.doc",
				"README.md",
				"README.txt",
				"a/src/app.js",
				"all/the/branches",
				"all/the/files.md",
				"app.js",
				"db/migrate-v1.0.sql",
				"db/sept/migrate-v1.sql",
				"dir/docs/my-file.txt",
				"docs/README.md",
				"docs/a/markdown/file.md",
				"docs/file.txt",
				"docs/hello.md",
				"docs/mona/hello-world.md",
				"docs/mona/octocat.txt",
				"every/tag",
				"feature/beta-a/my-branch",
				"feature/mona/the/octoca",
				"feature/my-branch",
				"feature/your-branch",
				"hello.md",
				"hoge-docs/a",
				"index.js",
				"js/README.md",
				"js/index.js",
				"main.go",
				"master",
				"math/their-post.md",
				"migrate-10909.sql",
				"mona-feature",
				"my-post.md",
				"my-src/code/js/app.js",
				"page.js",
				"page.jsx",
				"page.jsxx",
				"releases/mona-the-octcat",
				"space/docs/plan/space.doc",
				"src/js/index.js",
				"test-feature",
				"v1",
				"v1.10.1",
				"v2",
				"v2.0",
				"v2.0.0",
				"v2.9",
				"v3.1.0",
				"ver-10-feature",
			},
		},
		{
			ps: []string{`*.jsx?`},
			want: []string{
				"app.js",
				"index.js",
				"page.js",
				"page.jsx",
			},
		},
		{
			ps: []string{`**.js`},
			want: []string{
				"a/src/app.js",
				"app.js",
				"index.js",
				"js/index.js",
				"my-src/code/js/app.js",
				"page.js",
				"src/js/index.js",
			},
		},
		{
			ps: []string{`**/docs/**`},
			want: []string{
				"dir/docs/my-file.txt",
				"docs/README.md",
				"docs/a/markdown/file.md",
				"docs/file.txt",
				"docs/hello.md",
				"docs/mona/hello-world.md",
				"docs/mona/octocat.txt",
				"space/docs/plan/space.doc",
			},
		},
		{
			ps: []string{`feature/*`},
			want: []string{
				"feature/my-branch",
				"feature/your-branch",
			},
		},
		{
			ps: []string{`master`, `releases/mona-the-octcat`},
			want: []string{
				"master",
				"releases/mona-the-octcat",
			},
		},
		{
			ps: []string{`*.md`, `!README.md`},
			want: []string{
				"hello.md",
				"my-post.md",
			},
		},
		{
			ps: []string{`*.md`, `!README.md`, `README*`},
			want: []string{
				"README.doc",
				"README.md",
				"README.txt",
				"hello.md",
				"my-post.md",
			},
		},
	}
	opt := Option{
		Sort: true,
		Root: "./_testdir/",
	}
	for _, tt := range tests {
		got, err := GlobList(tt.ps, opt)
		if err != nil {
			t.Errorf("GlobList(%v, ...) got error: %v", tt.ps, err)
			continue
		}
		if diff := cmp.Diff(got, tt.want); diff != "" {
			t.Errorf("GlobList(%v, ...) got diff:\n%s\n\ngot:\n%s", tt.ps, diff, strings.Join(got, "\n"))
		}
	}
}
