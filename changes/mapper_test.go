package changes

import (
	"reflect"
	"testing"
	"time"

	"github.com/kylelemons/godebug/pretty"
)

type TestMappingStruct struct {
	NoTag     string
	JSONTag   string `json:"json_tag"`
	DBTag     string `db:"db_tag"`
	BothTag   string `json:"both_tag_json" db:"both_tag_db"`
	String    string
	StringPtr *string
	Time      time.Time
	TimePtr   *time.Time
	Int       int
	IntPtr    *int
	Float     float64
	FloatPtr  *float64
	Bool      bool
	BoolPtr   *bool
}

func TestTagMapping(t *testing.T) {
	mapper := NewTagMapper("db", "json")
	data := TestMappingStruct{}
	val := reflect.ValueOf(data)
	expected := map[string][]int{
		"NoTag":       {0},
		"json_tag":    {1},
		"db_tag":      {2},
		"both_tag_db": {3},
		"String":      {4},
		"StringPtr":   {5},
		"Time":        {6},
		"TimePtr":     {7},
		"Int":         {8},
		"IntPtr":      {9},
		"Float":       {10},
		"FloatPtr":    {11},
		"Bool":        {12},
		"BoolPtr":     {13},
	}
	result, err := mapper.KeyIndexes(val)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Unexpected Results\n%s", pretty.Compare(result, expected))

	}
}
