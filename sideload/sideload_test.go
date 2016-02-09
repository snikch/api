package sideload

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"testing"
)

type TestStruct struct {
	UserID    string `sideload:"users"`
	ProductID string `sideload:"products"`
	OwnerID   string
}

func resetRegistry() {
	typeRegistry = map[reflect.Type]map[int]string{}
	handlerRegistry = map[string]EntityHandler{}
}

func TestIDsFromDataInvalidTypes(t *testing.T) {
	var rdr io.Reader
	var ptr *struct{}
	rdr = &bytes.Buffer{}
	for _, typ := range []interface{}{
		"hello",             // string
		int(0),              // int
		map[string]string{}, // map
		ptr,                 // pointer
		rdr,                 // interface
	} {
		_, err := idsFromData(typ)
		if err == nil {
			t.Errorf("Expected type %T to fail", typ)
		}
	}
}

func TestIDsFromDataNoHandler(t *testing.T) {
	resetRegistry()
	data := TestStruct{"u1", "p1", "o1"}
	ids, err := idsFromData(data, "users", "products")
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
	if !reflect.DeepEqual(ids, map[string]map[string]bool{}) {
		t.Errorf("Expected empty map, got %s", ids)
	}
}

func TestIDsFromDataSingleStructAllRequired(t *testing.T) {
	resetRegistry()
	data := TestStruct{"u1", "p1", "o1"}
	RegisterType(data)
	ids, err := idsFromData(data)
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
	if !reflect.DeepEqual(ids, map[string]map[string]bool{
		"users":    {"u1": true},
		"products": {"p1": true},
	}) {
		t.Errorf("Unexpected map received: %s", ids)
	}
}

func TestIDsFromDataSingleStructSomeRequired(t *testing.T) {
	resetRegistry()
	data := TestStruct{"u1", "p1", "o1"}
	RegisterType(data)
	ids, err := idsFromData(data, "products")
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
	if !reflect.DeepEqual(ids, map[string]map[string]bool{
		"products": {"p1": true},
	}) {
		t.Errorf("Unexpected map received: %s", ids)
	}
}

func TestIDsFromDataStructSliceAllRequired(t *testing.T) {
	resetRegistry()
	data := []TestStruct{
		{"u1", "p1", "o1"},
		{"u2", "p1", "o2"},
		{"u3", "p2", "o3"},
		{"u2", "p1", "o2"},
	}
	RegisterType(data[0])
	ids, err := idsFromData(data)
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
	if !reflect.DeepEqual(ids, map[string]map[string]bool{
		"users":    {"u1": true, "u2": true, "u3": true},
		"products": {"p1": true, "p2": true},
	}) {
		t.Errorf("Unexpected map received: %s", ids)
	}
}

func TestIDsFromDataStructSliceSomeRequired(t *testing.T) {
	resetRegistry()
	data := []TestStruct{
		{"u1", "p1", "o1"},
		{"u2", "p1", "o2"},
		{"u3", "p2", "o3"},
		{"u2", "p1", "o2"},
	}
	RegisterType(data[0])
	ids, err := idsFromData(data, "users")
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
	if !reflect.DeepEqual(ids, map[string]map[string]bool{
		"users": {"u1": true, "u2": true, "u3": true},
	}) {
		t.Errorf("Unexpected map received: %s", ids)
	}
}

func TestHydrateEntitiesFromMapEmptyMap(t *testing.T) {
	resetRegistry()
	entities, err := hydrateEntitiesFromMap(map[string]map[string]bool{})
	if err != nil {
		t.Errorf("Expected no error: %s", err)
	}
	if entities != nil {
		t.Errorf("Unexpected map received: %s", entities)
	}
}

// No handler
func TestHydrateEntitiesFromMapNoHandler(t *testing.T) {
	resetRegistry()
	_, err := hydrateEntitiesFromMap(map[string]map[string]bool{
		"users": {"u1": true},
	})
	if err == nil {
		t.Errorf("Expected error, got none: %s", err)
	}
}

var (
	u1 = &TestStruct{UserID: "u1"}
	u2 = &TestStruct{UserID: "u2"}
	p1 = &TestStruct{ProductID: "p1"}
	p2 = &TestStruct{ProductID: "p2"}
)

func TestHydrateEntitiesFromMap(t *testing.T) {
	resetRegistry()
	RegisterEntityHandler("users", func(ids []string) (map[string]interface{}, error) {
		if !(ids[0] == "u1" && ids[1] == "u2") &&
			!(ids[1] == "u1" && ids[0] == "u2") {
			t.Errorf("Unexpected ids: %s", ids)
		}
		return map[string]interface{}{
			"u1": u1,
			"u2": u2,
		}, nil
	})
	RegisterEntityHandler("products", func(ids []string) (map[string]interface{}, error) {
		if !(ids[0] == "p1" && ids[1] == "p2") &&
			!(ids[1] == "p1" && ids[0] == "p2") {
			t.Errorf("Unexpected ids: %s", ids)
		}
		return map[string]interface{}{
			"p1": p1,
			"p2": p2,
		}, nil
	})
	entities, err := hydrateEntitiesFromMap(map[string]map[string]bool{
		"users":    {"u1": true, "u2": true},
		"products": {"p1": true, "p2": true},
	})
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if !reflect.DeepEqual(entities, map[string]map[string]interface{}{
		"users":    {"u1": u1, "u2": u2},
		"products": {"p1": p1, "p2": p2},
	}) {
		t.Errorf("Unexpected map received: %s", entities)
	}
}

func TestHydrateEntitiesFromMapOneFailure(t *testing.T) {
	resetRegistry()
	RegisterEntityHandler("users", func(ids []string) (map[string]interface{}, error) {
		return nil, fmt.Errorf("Failure")
	})
	RegisterEntityHandler("products", func(ids []string) (map[string]interface{}, error) {
		return map[string]interface{}{
			"p1": p1,
			"p2": p2,
		}, nil
	})
	_, err := hydrateEntitiesFromMap(map[string]map[string]bool{
		"users":    {"u1": true, "u2": true},
		"products": {"p1": true, "p2": true},
	})
	if err == nil {
		t.Errorf("Expected an error")
	}
}

func TestHydrateEntitiesFromMapMultipleFailures(t *testing.T) {
	resetRegistry()
	RegisterEntityHandler("users", func(ids []string) (map[string]interface{}, error) {
		return nil, fmt.Errorf("Failure")
	})
	RegisterEntityHandler("products", func(ids []string) (map[string]interface{}, error) {
		return nil, fmt.Errorf("Failure")
	})
	_, err := hydrateEntitiesFromMap(map[string]map[string]bool{
		"users":    {"u1": true, "u2": true},
		"products": {"p1": true, "p2": true},
	})
	if err == nil {
		t.Errorf("Expected an error")
	}
}
