package tools_test

import (
	"image"
	"image/png"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strings"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/arsham/blush/internal/tools"
	"github.com/google/go-cmp/cmp"
)

func stringSliceEq(t *testing.T, a, b []string) {
	t.Helper()
	sort.Strings(a)
	sort.Strings(b)
	if diff := cmp.Diff(a, b); diff != "" {
		t.Errorf("(-want +got):\\n%s", diff)
	}
}

func inSlice(niddle string, haystack []string) bool {
	for _, s := range haystack {
		if s == niddle {
			return true
		}
	}
	return false
}

func setup(t *testing.T, count int) (dirs, expect []string) {
	t.Helper()
	ret := make(map[string]struct{})
	tmp := t.TempDir()

	files := []struct {
		dir   string
		count int
	}{
		{"a", count},     // keep this here.
		{"a/b/c", count}, // this one is in the above folder, keep!
		{"abc", count},   // this one is outside.
		{"f", 0},         // this should not be matched.
	}
	for _, f := range files {
		l := path.Join(tmp, f.dir)
		err := os.MkdirAll(l, os.ModePerm)
		assert.NoError(t, err)

		for i := 0; i < f.count; i++ {
			f, err := ioutil.TempFile(l, "file_")
			assert.NoError(t, err)
			ret[path.Dir(f.Name())] = struct{}{}
			expect = append(expect, f.Name())
			f.WriteString("test")
		}
	}
	for d := range ret {
		dirs = append(dirs, d)
	}
	sort.Strings(dirs)
	return dirs, expect
}

func TestFilesError(t *testing.T) {
	t.Parallel()
	f, err := tools.Files(false)
	assert.Nil(t, f)
	assert.Error(t, err)
	f, err = tools.Files(false, "/path to heaven")
	assert.Error(t, err)
	assert.Nil(t, f)
}

func TestFiles(t *testing.T) {
	t.Parallel()
	dirs, expect := setup(t, 10)

	f, err := tools.Files(false, dirs...)
	assert.NoError(t, err)
	assert.NotNil(t, f)
	stringSliceEq(t, expect, f)

	// the a and abc should match, a/b/c should not
	f, err = tools.Files(false, dirs[0], dirs[2])
	assert.NoError(t, err)
	if len(f) != 20 { // all files in `a` and `abc`
		t.Errorf("len(f) = %d, want %d: %v", len(f), 20, f)
	}
}

func TestFilesOnSingleFile(t *testing.T) {
	t.Parallel()
	file, err := ioutil.TempFile(t.TempDir(), "blush_tools")
	assert.NoError(t, err)
	name := file.Name()

	f, err := tools.Files(true, name)
	assert.NoError(t, err)
	if len(f) != 1 {
		t.Fatalf("len(f) = %d, want 1", len(f))
	}
	if f[0] != name {
		t.Errorf("f[0] = %s, want %s", f[0], name)
	}

	f, err = tools.Files(false, name)
	assert.NoError(t, err)
	if len(f) != 1 {
		t.Fatalf("len(f) = %d, want 1", len(f))
	}
	if f[0] != name {
		t.Errorf("f[0] = %s, want %s", f[0], name)
	}
}

func TestFilesRecursive(t *testing.T) {
	t.Parallel()
	f, err := tools.Files(true, "/path to heaven")
	assert.Error(t, err)
	assert.Nil(t, f)

	dirs, expect := setup(t, 10)

	f, err = tools.Files(true, dirs...)
	assert.NoError(t, err)
	assert.NotNil(t, f)
	stringSliceEq(t, expect, f)

	f, err = tools.Files(true, dirs[0]) // expecting `a`
	assert.NoError(t, err)
	if len(f) != 20 { // all files in `a`
		t.Errorf("len(f) = %d, want %d: %v", len(f), 20, f)
	}
}

func setupUnpermissioned(t *testing.T) (rootDir, fileB string) {
	t.Helper()
	// creates this structure:
	// /tmp
	//     /a <--- 0222 perm
	//       /aaa.txt
	//     /b <--- 0777 perm
	//       /bbb.txt
	rootDir = t.TempDir()
	dirA := path.Join(rootDir, "a")
	dirB := path.Join(rootDir, "b")
	dirs := []struct {
		dir, file string
	}{
		{dirA, "aaa.txt"},
		{dirB, "bbb.txt"},
	}
	for _, d := range dirs {
		err := os.Mkdir(d.dir, 0o777)
		assert.NoError(t, err)
		name := path.Join(d.dir, d.file)
		_, err = os.Create(name)
		assert.NoError(t, err)
	}
	err := os.Chmod(dirA, 0o222)
	assert.NoError(t, err)
	fileB = path.Join(dirB, "bbb.txt")
	t.Cleanup(func() {
		err := os.Chmod(dirA, 0o777)
		assert.NoError(t, err)
	})
	return rootDir, fileB
}

func TestIgnoreNontPermissionedFolders(t *testing.T) {
	t.Parallel()
	rootDir, fileB := setupUnpermissioned(t)
	f, err := tools.Files(true, rootDir)
	assert.NoError(t, err)
	assert.NotNil(t, f)
	expect := []string{
		fileB,
	}
	stringSliceEq(t, expect, f)
}

// first returned string is text format, second is binary.
func setupBinaryFile(t *testing.T) (str, name string) {
	t.Helper()
	dir := t.TempDir()

	txt, err := os.Create(path.Join(dir, "txt"))
	assert.NoError(t, err)

	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	binary, err := os.Create(path.Join(dir, "binary"))
	assert.NoError(t, err)

	txt.WriteString("aaaaa")
	png.Encode(binary, img)

	return txt.Name(), binary.Name()
}

func TestIgnoreNonTextFiles(t *testing.T) {
	t.Parallel()
	txt, binary := setupBinaryFile(t)
	paths := path.Dir(txt)
	got, err := tools.Files(false, paths)
	assert.NoError(t, err)
	assert.True(t, inSlice(txt, got))
	assert.False(t, inSlice(binary, got))
}

func TestUnPrintableButTextContents(t *testing.T) {
	t.Parallel()
	tcs := []struct {
		name    string
		input   string
		wantLen int
	}{
		{"null", string(rune(0)), 1},
		{"space", " ", 1},
		{"return", "\r", 1},
		{"line feed", "\n", 1},
		{"tab", "\t", 1},
		{"mix", "\na\tbbb\n\n\n\t\t\n \t \n \r\nsjdk", 1},
		{"one", string(rune(1)), 0},
		{"bell", "\b", 0},
		{"bell in middle", "a\bc", 0},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			file, err := ioutil.TempFile(t.TempDir(), "blush_text")
			assert.NoError(t, err)
			file.WriteString(tc.input)
			got, err := tools.Files(false, file.Name())
			assert.NoError(t, err)
			assert.Len(t, got, tc.wantLen, strings.Join(got, "\n"))
		})
	}
}

func TestFilesIgnoreDirs(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()

	p := path.Join(dir, "a")
	err := os.MkdirAll(p, 0o777)
	assert.NoError(t, err)

	file, err := ioutil.TempFile(dir, "b")
	assert.NoError(t, err)

	g, err := tools.Files(true, dir)
	assert.NoError(t, err)
	assert.NotNil(t, g)

	assert.False(t, inSlice(p, g))
	assert.True(t, inSlice(file.Name(), g))
}
