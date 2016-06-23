package changes

import (
	"reflect"
	"testing"
)

type TestSliceStruct struct {
	Foo string
	Bar string
}

var testSliceDiffer = SliceDiffer{
	KeyMapper: NewTagMapper("db", "json"),
}

func TestSliceNilValues(t *testing.T) {
	s := []struct{}{}
	_, err := testSliceDiffer.Between(nil, s)
	if err == nil {
		t.Errorf("Expected error")
	}
	_, err = testSliceDiffer.Between(s, nil)
	if err == nil {
		t.Errorf("Expected error")
	}
	_, err = testSliceDiffer.Between(nil, nil)
	if err == nil {
		t.Errorf("Expected error")
	}
}

func TestSliceNoSlice(t *testing.T) {
	old := struct{}{}
	new := TestStruct{}
	_, err := testSliceDiffer.Between(old, new)
	if err == nil {
		t.Errorf("Expected error")
	}
}

func TestSliceMatching(t *testing.T) {
	old := []TestSliceStruct{
		{Foo: "foo", Bar: "bar"},
	}
	new := []TestSliceStruct{
		{Foo: "foo", Bar: "bar"},
	}
	diff, err := testSliceDiffer.Between(old, new)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}
	if len(diff.Added) != 0 {
		t.Errorf("Unexected added length %d: %+v", len(diff.Added), diff)
	}
	if len(diff.Removed) != 0 {
		t.Errorf("Unexected removed length %d: %+v", len(diff.Removed), diff)
	}
}

func TestSliceNotMatching(t *testing.T) {
	old := []TestSliceStruct{
		{Foo: "foo", Bar: "baz"},
		{Foo: "foo", Bar: "bar"},
	}
	new := []TestSliceStruct{
		{Foo: "fooz", Bar: "bar"},
		{Foo: "foo", Bar: "bar"},
	}
	diff, err := testSliceDiffer.Between(old, new)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}
	if len(diff.Added) != 1 {
		t.Errorf("Unexected added length %d: %+v", len(diff.Added), diff)
	}
	if len(diff.Removed) != 1 {
		t.Errorf("Unexected removed length %d: %+v", len(diff.Removed), diff)
	}
	if !reflect.DeepEqual(diff.Added[0], TestSliceStruct{Foo: "fooz", Bar: "bar"}) {
		t.Errorf("Unexpected added value: %+v", diff.Added[0])
	}
	if !reflect.DeepEqual(diff.Removed[0], TestSliceStruct{Foo: "foo", Bar: "baz"}) {
		t.Errorf("Unexpected removed value: %+v", diff.Removed[0])
	}
}
