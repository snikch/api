package changes

import (
	"reflect"
	"testing"
	"time"
)

type TestStruct struct {
	StringField    string     `db:"string_field"`
	StringPtrField *string    `json:"string_ptr_field"`
	TimeField      time.Time  `db:"time_field" json:"potato_field"`
	TimePtrField   *time.Time `json:"time_ptr_field"`
	IntField       int        `json:"int_field"`
	IntPtrField    *int       `json:"int_ptr_field"`
	FloatField     float64    `json:"float_field"`
	FloatPtrField  *float64   `json:"float_ptr_field"`
	BoolField      bool       `json:"bool_field"`
	BoolPtrField   *bool      `json:"bool_ptr_field"`
	ExcludedStruct struct {
		ExcludedField string
	} `diff:"exclude"`
	IncludedStruct struct {
		IncludedField string
		NoTag         string
	} `diff:"include"`
}

var testDiffer = Differ{
	KeyMapper: NewTagMapper("db", "json"),
}

func TestStructNilValues(t *testing.T) {
	s := struct{}{}
	_, err := testDiffer.Between(nil, s)
	if err == nil {
		t.Errorf("Expected error")
	}
	_, err = testDiffer.Between(s, nil)
	if err == nil {
		t.Errorf("Expected error")
	}
	_, err = testDiffer.Between(nil, nil)
	if err == nil {
		t.Errorf("Expected error")
	}
}

func TestStructMismatchingTypes(t *testing.T) {
	old := struct{}{}
	new := TestStruct{}
	_, err := testDiffer.Between(old, new)
	if err == nil {
		t.Errorf("Expected error")
	}
}

func TestStructPlainDiff(t *testing.T) {
	foo := "foo"
	bar := "bar"
	now := time.Now()
	old := TestStruct{
		StringField:    foo,
		StringPtrField: &foo,
		TimeField:      now,
		TimePtrField:   &now,
	}
	old.IncludedStruct.IncludedField = foo
	then := time.Now()
	new := TestStruct{
		StringField:    bar,
		StringPtrField: &bar,
		TimeField:      then,
		TimePtrField:   &then,
	}
	new.IncludedStruct.IncludedField = bar

	diffs, err := testDiffer.Between(old, new)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if len(diffs) != 5 {
		t.Errorf("Unexpected diff count: %d, expected 4", len(diffs))
	}
	for _, expectedDiff := range []Diff{
		{"string_field", "foo", "bar"},
		{"string_ptr_field", "foo", "bar"},
		{"time_field", now, then},
		{"time_ptr_field", now, then},
		{"IncludedStruct.IncludedField", "foo", "bar"},
	} {
		resultDiff, ok := diffs[expectedDiff.Key]
		if !ok {
			t.Errorf("No diff for %s", expectedDiff.Key)
			continue
		}
		if !reflect.DeepEqual(resultDiff.Old, expectedDiff.Old) {
			t.Errorf("Unexpected old value for key %s: %s was not %s", resultDiff.Key, resultDiff.Old, expectedDiff.Old)
			continue
		}
		if !reflect.DeepEqual(resultDiff.New, expectedDiff.New) {
			t.Errorf("Unexpected new value for key %s: %s was not %s", resultDiff.Key, resultDiff.New, expectedDiff.New)
			continue
		}
	}
}

func TestStructTimePointerVariations(t *testing.T) {
	foo := time.Now()
	old := TestStruct{}
	new := TestStruct{}

	// Compare nil with nil
	diffs, _ := testDiffer.Between(old, new)
	_, ok := diffs["time_ptr_field"]
	if ok {
		t.Errorf("Unexpected diff when passing both nil")
	}

	// Compare time with nil
	old.TimePtrField = &foo
	diffs, _ = testDiffer.Between(old, new)
	if !reflect.DeepEqual(diffs["time_ptr_field"].Old, foo) {
		t.Errorf("Unexpected old value %s was not %s", diffs["time_ptr_field"].Old, foo)
	}
	if !reflect.DeepEqual(diffs["time_ptr_field"].New, nil) {
		t.Errorf("Unexpected new value %s was not nil", diffs["time_ptr_field"].New)
	}

	// Compare time with same time
	new.TimePtrField = &foo
	diffs, _ = testDiffer.Between(old, new)
	_, ok = diffs["time_ptr_field"]
	if ok {
		t.Errorf("Unexpected diff when passing both nil")
	}

	// Compare time with same time in different timezone.
	bazNZ, err := time.Parse(time.RFC3339, "2016-01-01T13:00:00+13:00")
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	bazAU, err := time.Parse(time.RFC3339, "2016-01-01T11:00:00+11:00")
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	old.TimePtrField = &bazNZ
	new.TimePtrField = &bazAU
	diffs, _ = testDiffer.Between(old, new)
	_, ok = diffs["time_ptr_field"]
	if ok {
		t.Errorf("Unexpected diff when passing difference timezones")
	}

	// Compare nil with time
	old.TimePtrField = nil
	new.TimePtrField = &foo
	diffs, _ = testDiffer.Between(old, new)
	if !reflect.DeepEqual(diffs["time_ptr_field"].Old, nil) {
		t.Errorf("Unexpected old value %s was not nil", diffs["time_ptr_field"].Old)
	}
	if !reflect.DeepEqual(diffs["time_ptr_field"].New, foo) {
		t.Errorf("Unexpected new value %s was not %s", diffs["time_ptr_field"].New, foo)
	}

	// Compare different times
	bar := time.Now()
	old.TimePtrField = &bar
	diffs, _ = testDiffer.Between(old, new)
	if !reflect.DeepEqual(diffs["time_ptr_field"].Old, bar) {
		t.Errorf("Unexpected old value %s was not %s", diffs["time_ptr_field"].Old, bar)
	}
	if !reflect.DeepEqual(diffs["time_ptr_field"].New, foo) {
		t.Errorf("Unexpected new value %s was not %s", diffs["time_ptr_field"].New, foo)
	}
}

func TestStructStringPointerVariations(t *testing.T) {
	foo := "foo"
	old := TestStruct{}
	new := TestStruct{}

	// Compare nil with nil
	diffs, _ := testDiffer.Between(old, new)
	_, ok := diffs["string_ptr_field"]
	if ok {
		t.Errorf("Unexpected diff when passing both nil")
	}

	// Compare string with nil
	old.StringPtrField = &foo
	diffs, _ = testDiffer.Between(old, new)
	if !reflect.DeepEqual(diffs["string_ptr_field"].Old, foo) {
		t.Errorf("Unexpected old value %s was not %s", diffs["string_ptr_field"].Old, foo)
	}
	if !reflect.DeepEqual(diffs["string_ptr_field"].New, nil) {
		t.Errorf("Unexpected new value %s was not nil", diffs["string_ptr_field"].New)
	}

	// Compare string with same string, different pointer
	foo2 := "foo"
	new.StringPtrField = &foo2
	diffs, _ = testDiffer.Between(old, new)
	_, ok = diffs["string_ptr_field"]
	if ok {
		t.Errorf("Unexpected diff when passing both nil")
	}

	// Compare nil with string
	old.StringPtrField = nil
	new.StringPtrField = &foo
	diffs, _ = testDiffer.Between(old, new)
	if !reflect.DeepEqual(diffs["string_ptr_field"].Old, nil) {
		t.Errorf("Unexpected old value %s was not nil", diffs["string_ptr_field"].Old)
	}
	if !reflect.DeepEqual(diffs["string_ptr_field"].New, foo) {
		t.Errorf("Unexpected new value %s was not %s", diffs["string_ptr_field"].New, foo)
	}

	// Compare different strings
	bar := "bar"
	old.StringPtrField = &bar
	diffs, _ = testDiffer.Between(old, new)
	if !reflect.DeepEqual(diffs["string_ptr_field"].Old, bar) {
		t.Errorf("Unexpected old value %s was not %s", diffs["string_ptr_field"].Old, bar)
	}
	if !reflect.DeepEqual(diffs["string_ptr_field"].New, foo) {
		t.Errorf("Unexpected new value %s was not %s", diffs["string_ptr_field"].New, foo)
	}
}

func TestStructIntPointerVariations(t *testing.T) {
	foo := 1337
	old := TestStruct{}
	new := TestStruct{}

	// Compare nil with nil
	diffs, _ := testDiffer.Between(old, new)
	_, ok := diffs["int_ptr_field"]
	if ok {
		t.Errorf("Unexpected diff when passing both nil")
	}

	// Compare int with nil
	old.IntPtrField = &foo
	diffs, _ = testDiffer.Between(old, new)
	if !reflect.DeepEqual(diffs["int_ptr_field"].Old, foo) {
		t.Errorf("Unexpected old value %s was not %d", diffs["int_ptr_field"].Old, foo)
	}
	if !reflect.DeepEqual(diffs["int_ptr_field"].New, nil) {
		t.Errorf("Unexpected new value %s was not nil", diffs["int_ptr_field"].New)
	}

	// Compare int with same int, different pointer
	foo2 := 1337
	new.IntPtrField = &foo2
	diffs, _ = testDiffer.Between(old, new)
	_, ok = diffs["int_ptr_field"]
	if ok {
		t.Errorf("Unexpected diff when passing both nil")
	}

	// Compare nil with int
	old.IntPtrField = nil
	new.IntPtrField = &foo
	diffs, _ = testDiffer.Between(old, new)
	if !reflect.DeepEqual(diffs["int_ptr_field"].Old, nil) {
		t.Errorf("Unexpected old value %s was not nil", diffs["int_ptr_field"].Old)
	}
	if !reflect.DeepEqual(diffs["int_ptr_field"].New, foo) {
		t.Errorf("Unexpected new value %s was not %d", diffs["int_ptr_field"].New, foo)
	}

	// Compare different ints
	bar := 137
	old.IntPtrField = &bar
	diffs, _ = testDiffer.Between(old, new)
	if !reflect.DeepEqual(diffs["int_ptr_field"].Old, bar) {
		t.Errorf("Unexpected old value %s was not %d", diffs["int_ptr_field"].Old, bar)
	}
	if !reflect.DeepEqual(diffs["int_ptr_field"].New, foo) {
		t.Errorf("Unexpected new value %s was not %d", diffs["int_ptr_field"].New, foo)
	}
}

func TestStructFloatPointerVariations(t *testing.T) {
	foo := 133.7
	old := TestStruct{}
	new := TestStruct{}

	// Compare nil with nil
	diffs, _ := testDiffer.Between(old, new)
	_, ok := diffs["float_ptr_field"]
	if ok {
		t.Errorf("Unexpected diff when passing both nil")
	}

	// Compare float with nil
	old.FloatPtrField = &foo
	diffs, _ = testDiffer.Between(old, new)
	if !reflect.DeepEqual(diffs["float_ptr_field"].Old, foo) {
		t.Errorf("Unexpected old value %s was not %f", diffs["float_ptr_field"].Old, foo)
	}
	if !reflect.DeepEqual(diffs["float_ptr_field"].New, nil) {
		t.Errorf("Unexpected new value %s was not nil", diffs["float_ptr_field"].New)
	}

	// Compare float with same float, different pointer
	foo2 := 133.7
	new.FloatPtrField = &foo2
	diffs, _ = testDiffer.Between(old, new)
	_, ok = diffs["float_ptr_field"]
	if ok {
		t.Errorf("Unexpected diff when passing both nil")
	}

	// Compare nil with float
	old.FloatPtrField = nil
	new.FloatPtrField = &foo
	diffs, _ = testDiffer.Between(old, new)
	if !reflect.DeepEqual(diffs["float_ptr_field"].Old, nil) {
		t.Errorf("Unexpected old value %s was not nil", diffs["float_ptr_field"].Old)
	}
	if !reflect.DeepEqual(diffs["float_ptr_field"].New, foo) {
		t.Errorf("Unexpected new value %s was not %f", diffs["float_ptr_field"].New, foo)
	}

	// Compare different floats
	bar := 13.7
	old.FloatPtrField = &bar
	diffs, _ = testDiffer.Between(old, new)
	if !reflect.DeepEqual(diffs["float_ptr_field"].Old, bar) {
		t.Errorf("Unexpected old value %s was not %f", diffs["float_ptr_field"].Old, bar)
	}
	if !reflect.DeepEqual(diffs["float_ptr_field"].New, foo) {
		t.Errorf("Unexpected new value %s was not %f", diffs["float_ptr_field"].New, foo)
	}
}

func TestStructBoolPointerVariations(t *testing.T) {
	foo := false
	old := TestStruct{}
	new := TestStruct{}

	// Compare nil with nil
	diffs, _ := testDiffer.Between(old, new)
	_, ok := diffs["bool_ptr_field"]
	if ok {
		t.Errorf("Unexpected diff when passing both nil")
	}

	// Compare bool with nil
	old.BoolPtrField = &foo
	diffs, _ = testDiffer.Between(old, new)
	if !reflect.DeepEqual(diffs["bool_ptr_field"].Old, foo) {
		t.Errorf("Unexpected old value %s was not %t", diffs["bool_ptr_field"].Old, foo)
	}
	if !reflect.DeepEqual(diffs["bool_ptr_field"].New, nil) {
		t.Errorf("Unexpected new value %s was not nil", diffs["bool_ptr_field"].New)
	}

	// Compare bool with same bool, different pointer
	foo2 := false
	new.BoolPtrField = &foo2
	diffs, _ = testDiffer.Between(old, new)
	_, ok = diffs["bool_ptr_field"]
	if ok {
		t.Errorf("Unexpected diff when passing both nil")
	}

	// Compare nil with bool
	old.BoolPtrField = nil
	new.BoolPtrField = &foo
	diffs, _ = testDiffer.Between(old, new)
	if !reflect.DeepEqual(diffs["bool_ptr_field"].Old, nil) {
		t.Errorf("Unexpected old value %s was not nil", diffs["bool_ptr_field"].Old)
	}
	if !reflect.DeepEqual(diffs["bool_ptr_field"].New, foo) {
		t.Errorf("Unexpected new value %s was not %t", diffs["bool_ptr_field"].New, foo)
	}

	// Compare different bools
	bar := true
	old.BoolPtrField = &bar
	diffs, _ = testDiffer.Between(old, new)
	if !reflect.DeepEqual(diffs["bool_ptr_field"].Old, bar) {
		t.Errorf("Unexpected old value %s was not %t", diffs["bool_ptr_field"].Old, bar)
	}
	if !reflect.DeepEqual(diffs["bool_ptr_field"].New, foo) {
		t.Errorf("Unexpected new value %s was not %t", diffs["bool_ptr_field"].New, foo)
	}
}

func TestStructPtrStructDiff(t *testing.T) {
	old := TestStruct{
		StringField: "foo",
	}
	new := TestStruct{
		StringField: "bar",
	}

	// Compare string with nil
	for _, example := range []struct {
		Name     string
		Old, New interface{}
	}{
		{"No ptr", old, new},
		{"Old ptr", &old, new},
		{"New ptr", old, &new},
		{"Both ptr", &old, &new},
	} {
		diffs, _ := testDiffer.Between(example.Old, example.New)
		if !reflect.DeepEqual(diffs["string_field"].Old, "foo") {
			t.Errorf("%s: Unexpected old value %s was not %s", example.Name, diffs["string_field"].Old, "foo")
		}
		if !reflect.DeepEqual(diffs["string_field"].New, "bar") {
			t.Errorf("%s: Unexpected new value %s was not %s", example.Name, diffs["string_field"].New, "bar")
		}
	}
}

func BenchmarkPlainStruct(b *testing.B) {
	foo := "foo"
	bar := "bar"
	now := time.Now()
	old := TestStruct{
		StringField:    foo,
		StringPtrField: &foo,
		TimeField:      now,
		TimePtrField:   &now,
	}
	then := time.Now()
	new := TestStruct{
		StringField:    bar,
		StringPtrField: &bar,
		TimeField:      then,
		TimePtrField:   &then,
	}
	for i := 0; i < b.N; i++ {
		testDiffer.Between(old, new)
	}
}
