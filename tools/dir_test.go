package tools_test

import (
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"sort"
	"testing"

	"github.com/arsham/blush/tools"
)

func setup(t *testing.T, count int) (dirs, expect []string, cleanup func()) {
	ret := make(map[string]struct{}, 0)
	tmp, err := ioutil.TempDir("", "blush")
	if err != nil {
		t.Fatal(err)
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
			t.Fatal(err)
		}

		for i := 0; i < f.count; i++ {
			f, err := ioutil.TempFile(l, "file_")
			if err != nil {
				t.Fatal(err)
			}
			ret[path.Dir(f.Name())] = struct{}{}
			expect = append(expect, f.Name())
		}
	}
	for d := range ret {
		dirs = append(dirs, d)
	}
	sort.Strings(dirs)
	return
}

func TestFiles(t *testing.T) {
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

	dirs, expect, cleanup := setup(t, 10)
	defer cleanup()

	f, err = tools.Files(false, dirs...)
	if err != nil {
		t.Errorf("err = %v, want nil", err)
	}
	if f == nil {
		t.Error("f = nil, want []string")
	}
	sort.Strings(expect)
	sort.Strings(f)
	if !reflect.DeepEqual(expect, f) {
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

func TestFilesRecursive(t *testing.T) {
	f, err := tools.Files(true, "/path to heaven")
	if err == nil {
		t.Error("err = nil, want error")
	}
	if f != nil {
		t.Errorf("f = %v, want nil", f)
	}

	dirs, expect, cleanup := setup(t, 10)
	defer cleanup()
	f, err = tools.Files(true, dirs...)
	if err != nil {
		t.Errorf("err = %v, want nil", err)
	}
	if f == nil {
		t.Error("f = nil, want []string")
	}
	sort.Strings(expect)
	sort.Strings(f)
	if !reflect.DeepEqual(expect, f) {
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
