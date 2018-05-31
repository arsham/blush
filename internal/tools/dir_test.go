package tools_test

import (
	"image"
	"image/png"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"testing"

	"github.com/arsham/blush/internal/tools"
)

func stringSliceEq(a, b []string) bool {
	sort.Strings(a)
	sort.Strings(b)
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func inStringSlice(a string, b []string) bool {
	for _, s := range b {
		if a == s {
			return true
		}
	}
	return false
}

func setup(count int) (dirs, expect []string, cleanup func(), err error) {
	ret := make(map[string]struct{})
	tmp, err := ioutil.TempDir("", "blush")
	if err != nil {
		return nil, nil, func() {}, err
	}
	cleanup = func() {
		os.RemoveAll(tmp)
	}
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
		if err != nil {
			return nil, nil, cleanup, err
		}

		for i := 0; i < f.count; i++ {
			f, err := ioutil.TempFile(l, "file_")
			if err != nil {
				return nil, nil, cleanup, err
			}
			ret[path.Dir(f.Name())] = struct{}{}
			expect = append(expect, f.Name())
			f.WriteString("test")
		}
	}
	for d := range ret {
		dirs = append(dirs, d)
	}
	sort.Strings(dirs)
	return
}

func TestFilesError(t *testing.T) {
	f, err := tools.Files(false)
	if f != nil {
		t.Errorf("f = %v, want nil", f)
	}
	if err == nil {
		t.Error("err = nil, want error")
	}
	f, err = tools.Files(false, "/path to heaven")
	if err == nil {
		t.Error("err = nil, want error")
	}
	if f != nil {
		t.Errorf("f = %v, want nil", f)
	}
}

func TestFiles(t *testing.T) {
	dirs, expect, cleanup, err := setup(10)
	defer cleanup()
	if err != nil {
		t.Fatal(err)
	}

	f, err := tools.Files(false, dirs...)
	if err != nil {
		t.Errorf("err = %v, want nil", err)
	}
	if f == nil {
		t.Error("f = nil, want []string")
	}
	if !stringSliceEq(expect, f) {
		t.Errorf("f = %v, \nwant %v", f, expect)
	}

	// the a and abc should match, a/b/c should not
	f, err = tools.Files(false, dirs[0], dirs[2])
	if err != nil {
		t.Errorf("err = %v, want nil", err)
	}
	if len(f) != 20 { // all files in `a` and `abc`
		t.Errorf("len(f) = %d, want %d: %v", len(f), 20, f)
	}
}

func TestFilesOnSingleFile(t *testing.T) {
	file, err := ioutil.TempFile("", "blush_tools")
	if err != nil {
		t.Fatal(err)
	}
	name := file.Name()
	defer func() {
		if err = os.Remove(name); err != nil {
			t.Error(err)
		}
	}()

	f, err := tools.Files(true, name)
	if err != nil {
		t.Errorf("err = %v, want nil", err)
	}
	if len(f) != 1 {
		t.Fatalf("len(f) = %d, want 1", len(f))
	}
	if f[0] != name {
		t.Errorf("f[0] = %s, want %s", f[0], name)
	}

	f, err = tools.Files(false, name)
	if err != nil {
		t.Errorf("err = %v, want nil", err)
	}
	if len(f) != 1 {
		t.Fatalf("len(f) = %d, want 1", len(f))
	}
	if f[0] != name {
		t.Errorf("f[0] = %s, want %s", f[0], name)
	}
}

func TestFilesRecursive(t *testing.T) {
	f, err := tools.Files(true, "/path to heaven")
	if err == nil {
		t.Error("err = nil, want error")
	}
	if f != nil {
		t.Errorf("f = %v, want nil", f)
	}

	dirs, expect, cleanup, err := setup(10)
	defer cleanup()
	if err != nil {
		t.Fatal(err)
	}

	f, err = tools.Files(true, dirs...)
	if err != nil {
		t.Errorf("err = %v, want nil", err)
	}
	if f == nil {
		t.Error("f = nil, want []string")
	}
	if !stringSliceEq(expect, f) {
		t.Errorf("f = %v, want %v", f, expect)
	}

	f, err = tools.Files(true, dirs[0]) // expecting `a`
	if err != nil {
		t.Errorf("err = %v, want nil", err)
	}
	if len(f) != 20 { // all files in `a`
		t.Errorf("len(f) = %d, want %d: %v", len(f), 20, f)
	}
}

func setupUnpermissioned(t *testing.T) (string, string, func()) {
	// creates this structure:
	// /tmp
	//     /a <--- 0222 perm
	//       /aaa.txt
	//     /b <--- 0777 perm
	//       /bbb.txt
	rootDir, err := ioutil.TempDir("", "blush_dir")
	if err != nil {
		t.Fatal(err)
	}
	dirA := path.Join(rootDir, "a")
	dirB := path.Join(rootDir, "b")
	dirs := []struct {
		dir, file string
	}{
		{dirA, "aaa.txt"},
		{dirB, "bbb.txt"},
	}
	for _, d := range dirs {
		err = os.Mkdir(d.dir, 0777)
		if err != nil {
			t.Fatal(err)
		}
		name := path.Join(d.dir, d.file)
		_, err = os.Create(name)
		if err != nil {
			t.Fatal(err)
		}
	}
	err = os.Chmod(dirA, 0222)
	if err != nil {
		t.Fatal(err)
	}
	fileB := path.Join(dirB, "bbb.txt")
	return rootDir, fileB, func() {
		err = os.Chmod(dirA, 0777)
		if err != nil {
			t.Error(err)
		}
		if err = os.RemoveAll(rootDir); err != nil {
			t.Error(err)
		}
	}
}

func TestIgnoreNontPermissionedFolders(t *testing.T) {
	rootDir, fileB, cleanup := setupUnpermissioned(t)
	defer cleanup()
	f, err := tools.Files(true, rootDir)
	if err != nil {
		t.Errorf("err = %v, want nil", err)
	}
	if f == nil {
		t.Error("f = nil, want []string")
	}
	expect := []string{
		fileB,
	}
	if !stringSliceEq(expect, f) {
		t.Errorf("f = %v, want %v", f, expect)
	}
}

// first returned string is text format, second is binary
func setupBinaryFile(t *testing.T) (string, string, func()) {
	dir, err := ioutil.TempDir("", "blush_dir")
	if err != nil {
		t.Fatal(err)
	}
	txt, err := os.Create(path.Join(dir, "txt"))
	if err != nil {
		t.Fatal(err)
	}
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	binary, err := os.Create(path.Join(dir, "binary"))
	if err != nil {
		t.Fatal(err)
	}
	txt.WriteString("aaaaa")
	png.Encode(binary, img)
	return txt.Name(), binary.Name(), func() {
		if err = os.RemoveAll(dir); err != nil {
			t.Error(err)
		}
	}
}

func TestIgnoreNonTextFiles(t *testing.T) {
	txt, binary, cleanup := setupBinaryFile(t)
	defer cleanup()
	paths := path.Dir(txt)
	got, err := tools.Files(false, paths)
	if err != nil {
		t.Error(err)
	}
	if !inStringSlice(txt, got) {
		t.Errorf("%s not found in %v", txt, got)
	}
	if inStringSlice(binary, got) {
		t.Errorf("%s was found in %v", binary, got)
	}
}

func TestUnPrintableButTextContents(t *testing.T) {
	tcs := []struct {
		name  string
		input string
	}{
		{"null", string(0)},
		{"space", " "},
		{"R", "\r"},
		{"Line feed", "\n"},
		{"Tab", "\t"},
		{"mix", "\na\tbbb\n\n\n\t\t\n \t \n \r\nsjdk"},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			file, err := ioutil.TempFile("", "blush_text")
			if err != nil {
				t.Fatal(err)
			}
			defer func() {
				if err = os.Remove(file.Name()); err != nil {
					t.Error(err)
				}
			}()
			file.WriteString(tc.input)
			got, err := tools.Files(false, file.Name())
			if err != nil {
				t.Error(err)
			}
			if len(got) != 1 {
				t.Errorf("len(%v) = %d, want 1", got, len(got))
			}
		})
	}
}
