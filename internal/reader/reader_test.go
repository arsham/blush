package reader_test

import (
	"bytes"
	"io"
	"path"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/arsham/blush/internal/reader"
)

func TestWithReader(t *testing.T) {
	t.Parallel()
	m, err := reader.NewMultiReader(reader.WithReader("name", nil))
	assert.Error(t, err)
	assert.Nil(t, m)

	r := io.NopCloser(&bytes.Buffer{})
	m, err = reader.NewMultiReader(reader.WithReader("name", r))
	assert.NoError(t, err)
	assert.NotNil(t, m)

	m, err = reader.NewMultiReader(reader.WithReader("", r))
	assert.NoError(t, err)
	assert.NotNil(t, m)
}

func TestWithReaderMultipleReadersClose(t *testing.T) {
	t.Parallel()
	var called []string
	input1 := "afmBEswIRYosG7"
	input2 := "UbMFeIFjvAhdA3sdT"
	r1 := nopCloser{
		Reader: bytes.NewBufferString(input1),
		closeFunc: func() error {
			called = append(called, "r1")
			return nil
		},
	}
	r2 := nopCloser{
		Reader: bytes.NewBufferString(input2),
		closeFunc: func() error {
			called = append(called, "r2")
			return nil
		},
	}
	m, err := reader.NewMultiReader(reader.WithReader("r1", r1), reader.WithReader("r2", r2))
	assert.NoError(t, err)

	b := make([]byte, 100)
	_, err = m.Read(b)
	assert.NoError(t, err)
	assert.EqualValues(t, input1, bytes.Trim(b, "\x00"))

	_, err = m.Read(b)
	assert.NoError(t, err)
	assert.True(t, inSlice("r1", called))
	assert.EqualValues(t, input2, bytes.Trim(b, "\x00"))

	_, err = m.Read(b)
	assert.EqualError(t, io.EOF, err.Error())
	assert.True(t, inSlice("r2", called))
}

func TestWithReaderMultipleReadersError(t *testing.T) {
	t.Parallel()
	r := nopCloser{
		Reader: &bytes.Buffer{},
		closeFunc: func() error {
			return nil
		},
	}
	m, err := reader.NewMultiReader(reader.WithReader("r", r), nil)
	assert.Error(t, err)
	assert.Nil(t, m)
}

func TestWithPathsError(t *testing.T) {
	t.Parallel()
	tcs := []struct {
		name  string
		input []string
	}{
		{"nil", nil},
		{"empty", []string{}},
		{"empty string", []string{""}},
		{"not found", []string{"nomansland2987349237"}},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			input := reader.WithPaths(tc.input, true)
			m, err := reader.NewMultiReader(input)
			assert.Error(t, err)
			assert.Nil(t, m)
		})
	}
}

func TestNewMultiReaderWithPaths(t *testing.T) {
	t.Parallel()
	var (
		c1 = "VJSNS5IeLCtEB"
		c2 = "kkNL8vGNJn"
		c3 = "o6Ubb5Taj"
	)
	input := []testCase{
		{"a/a.txt", c1},
		{"a/b.txt", c2},
		{"ab.txt", c3},
	}

	dirs := setup(t, input)
	m, err := reader.NewMultiReader(reader.WithPaths(dirs, false))
	assert.NoError(t, err)
	assert.NotNil(t, m)
	err = m.Close()
	assert.NoError(t, err)
}

func TestMultiReaderReadOneReader(t *testing.T) {
	t.Parallel()
	input := "sdlksjdljfQYawl5OEEg"
	r := io.NopCloser(bytes.NewBufferString(input))
	m, err := reader.NewMultiReader(reader.WithReader("r", r))
	assert.NoError(t, err)
	assert.NotNil(t, m)

	b := make([]byte, len(input))
	n, err := m.Read(b)
	assert.NoError(t, err)
	assert.Equal(t, len(input), n)
	assert.EqualValues(t, input, b)

	n, err = m.Read(b)
	assert.EqualError(t, err, io.EOF.Error())
	assert.Zero(t, n)
}

func TestMultiReaderReadZeroBytes(t *testing.T) {
	t.Parallel()
	input := "3wAgvZ4bSfQYawl5OEEg"
	r := io.NopCloser(bytes.NewBufferString(input))
	m, err := reader.NewMultiReader(reader.WithReader("r", r))
	assert.NoError(t, err)
	assert.NotNil(t, m)

	b := make([]byte, 0)
	n, err := m.Read(b)
	assert.NoError(t, err)
	assert.Zero(t, n)
	assert.Empty(t, b)
}

func TestMultiReaderReadOneReaderMoreSpace(t *testing.T) {
	t.Parallel()
	input := "3wAgvZ4bSfQYawl5OEEg"
	r := io.NopCloser(bytes.NewBufferString(input))
	m, err := reader.NewMultiReader(reader.WithReader("r", r))
	assert.NoError(t, err)
	assert.NotNil(t, m)
	b := make([]byte, len(input)+1)
	n, err := m.Read(b)
	assert.NoError(t, err)
	assert.EqualValues(t, len(input), n)
	assert.EqualValues(t, input, bytes.Trim(b, "\x00"))
}

func TestMultiReaderReadMultipleReaders(t *testing.T) {
	t.Parallel()
	input := []string{"P5tyugWXFn", "b8YbUO7pMX3G8j4Bi"}
	r1 := io.NopCloser(bytes.NewBufferString(input[0]))
	r2 := io.NopCloser(bytes.NewBufferString(input[1]))
	m, err := reader.NewMultiReader(
		reader.WithReader("r1", r1),
		reader.WithReader("r2", r2),
	)
	assert.NoError(t, err)
	assert.NotNil(t, m)

	tcs := []struct {
		name    string
		b       []byte
		wantErr error
		wantLen int
		wantOut string
	}{
		{"r1", make([]byte, len(input[0])), nil, len(input[0]), input[0]},
		{"r2", make([]byte, len(input[1])), nil, len(input[1]), input[1]},
		{"nothing left", make([]byte, 10), io.EOF, 0, ""},
	}
	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			n, err := m.Read(tc.b)
			assert.Equal(t, err, tc.wantErr, "error")
			assert.EqualValues(t, tc.wantLen, n)
			assert.Equal(t, tc.wantOut, string(bytes.Trim(tc.b, "\x00")), "output")
		})
	}
}

func TestMultiReaderNames(t *testing.T) {
	t.Parallel()
	input := []string{"Mw0mxekLYOpXaKl8PVT", "1V5MjHUXYTPChW"}
	r1 := io.NopCloser(bytes.NewBufferString(input[0]))
	r2 := io.NopCloser(bytes.NewBufferString(input[1]))
	m, err := reader.NewMultiReader(
		reader.WithReader("r1", r1),
		reader.WithReader("r2", r2),
	)
	assert.NoError(t, err)
	assert.NotNil(t, m)
	b := make([]byte, 100)
	tcs := []struct {
		name    string
		wantErr error
	}{
		{"r1", nil},
		{"r2", nil},
		{"", io.EOF},
	}
	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			_, err := m.Read(b)
			assert.Equal(t, err, tc.wantErr, "error")
			assert.Equal(t, tc.name, m.FileName())
		})
	}
}

func TestNewMultiReaderWithPathsRead(t *testing.T) {
	t.Parallel()
	var (
		c1 = "VJSNS5IeLCtEB"
		c2 = "kkNL8vGNJn"
		c3 = "o6Ubb5Taj"
	)
	input := []testCase{
		{"a/a.txt", c1},
		{"a/b.txt", c2},
		{"ab.txt", c3},
	}

	dirs := setup(t, input)
	w, err := reader.NewMultiReader(reader.WithPaths(dirs, false))
	assert.NoError(t, err)
	assert.NotNil(t, w)
	t.Cleanup(func() {
		err = w.Close()
		assert.NoError(t, err)
	})

	buf := &bytes.Buffer{}
	_, err = buf.ReadFrom(w)
	assert.NoError(t, err)
	for _, s := range []string{c1, c2, c3} {
		assert.Contains(t, buf.String(), s)
	}
}

func TestNewMultiReaderRecursive(t *testing.T) {
	t.Parallel()
	var (
		c1 = "1JQey4agQ3w9pqg3"
		c2 = "7ToNRMgsOAR6A"
		c3 = "EtOkn9C5zoH0Dla2rF9"
	)
	input := []testCase{
		{"a/a.txt", c1},
		{"a/b.txt", c2},
		{"a/b/c.txt", c3},
	}

	dirs := setup(t, input)
	base := path.Join(path.Dir(dirs[0]), "a")
	w, err := reader.NewMultiReader(reader.WithPaths([]string{base}, true))
	assert.NoError(t, err)
	assert.NotNil(t, w)

	t.Cleanup(func() {
		err = w.Close()
		assert.NoError(t, err)
	})

	buf := &bytes.Buffer{}
	_, err = buf.ReadFrom(w)
	assert.NoError(t, err)

	for _, s := range []string{c1, c2, c3} {
		assert.Contains(t, buf.String(), s)
	}
}

func TestNewMultiReaderNonRecursive(t *testing.T) {
	t.Parallel()
	var (
		c1 = "DRAjfSq2y"
		c2 = "ht3xCIQ"
		c3 = "jPqPoAbMNb"
	)
	input := []testCase{
		{"a/a.txt", c1},
		{"a/b.txt", c2},
		{"a/b/c.txt", c3},
	}

	dirs := setup(t, input)
	base := path.Join(path.Dir(dirs[0]), "a")
	w, err := reader.NewMultiReader(reader.WithPaths([]string{base}, false))
	assert.NoError(t, err)
	assert.NotNil(t, w)

	t.Cleanup(func() {
		err = w.Close()
		assert.NoError(t, err)
	})

	buf := &bytes.Buffer{}
	_, err = buf.ReadFrom(w)
	assert.NoError(t, err)
	for _, s := range []string{c1, c2} {
		assert.Contains(t, buf.String(), s)
	}
	assert.NotContains(t, buf.String(), c3)
}
