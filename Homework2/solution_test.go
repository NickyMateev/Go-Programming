package main

import (
	"testing"
)

func TestSample(t *testing.T) {
	t.Run("origin", func(t *testing.T) {
		f := NewEditor("foobar")
		compare(t, "foobar", f.String())
	})

	t.Run("insert", func(t *testing.T) {
		f := NewEditor("foobar")
		compare(t, "fobazobar", f.Insert(2, "baz").String())
	})

	t.Run("append", func(t *testing.T) {
		f := NewEditor("foobar")
		compare(t, "foobarbaz", f.Insert(6, "baz").String())
	})

	t.Run("insert_append_front_append_back", func(t *testing.T) {
		f := NewEditor("A large span of text")
		f.Insert(16, "an English ")
		f.Insert(2, "very ")
		f.Insert(36, " message.")
		f.Insert(0, "This is ")
		compare(t, "This is A very large span of an English text message.", f.String())
	})

	t.Run("delete", func(t *testing.T) {
		f := NewEditor("foobar")
		compare(t, "far", f.Delete(1, 3).String())
	})
}

func compare(t *testing.T, exp, got string) {
	if got != exp {
		t.Errorf("Expect: %q; got %q", exp, got)
	}
}
