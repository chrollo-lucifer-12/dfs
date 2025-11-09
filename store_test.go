package main

import (
	"bytes"
	"io"
	"testing"
)

func TestPathTransformFunc(t *testing.T) {
	key := "gojo"
	pathkey := CASPathTransformFunc(key)
	expectedOriginalKey := "5ed64c90f26c784f26f571a1d2abc883d6eece7f"
	expectedPathName := "5ed64/c90f2/6c784/f26f5/71a1d/2abc8/83d6e/ece7f"

	if pathkey.Pathname != expectedPathName {
		t.Errorf("have %s want %s", pathkey.Pathname, expectedPathName)
	}

	if pathkey.Filename != expectedOriginalKey {
		t.Errorf("have %s want %s", pathkey.Filename, expectedOriginalKey)
	}
}

func TestStore(t *testing.T) {
	s := NewStore(StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	})
	data := []byte("hi")
	s.writeStream("sahil", bytes.NewReader(data))

	r, _ := s.Read("sahil")

	b, _ := io.ReadAll(r)

	if string(b) != string(data) {
		t.Errorf("want %s have %s", string(data), string(b))
	}
}
