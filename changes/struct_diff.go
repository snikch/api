package changes

import (
	"errors"
	"reflect"
	"time"
)

// Diff represents a single field's change in values.
type Diff struct {
	Key      string
	Old, New interface{}
}

func (diff Diff) Match() bool {
	return diff.Old == diff.New
}

// DiffSet represents a set of multiple diffs, with the field name as a key.
type DiffSet map[string]Diff

var (
	ErrNil         = errors.New("nil interface supplied")
	ErrNotStruct   = errors.New("structs must be supplied")
	ErrNotSameType = errors.New("structs must be of the same type")
)

type Differ struct {
	KeyMapper KeyMapper
}

func (differ *Differ) Between(old, new interface{}) (DiffSet, error) {
	// No nils thanks.
	if old == nil || new == nil {
		return nil, ErrNil
	}

	oldVal := reflect.Indirect(reflect.ValueOf(old))
	newVal := reflect.Indirect(reflect.ValueOf(new))
	oldType := oldVal.Type()

	// Ensure both types are a match.
	if oldType != newVal.Type() {
		return nil, ErrNotSameType
	}

	// Ensure we're dealing with structs.
	if oldVal.Kind() != reflect.Struct {
		return nil, ErrNotStruct
	}

	// Get the key indexes and names for changes we care about.
	keyIndexes, err := differ.KeyMapper.KeyIndexes(newVal)
	if err != nil {
		return nil, err
	}

	// Loop over the keyIndexes and retrieve the values from both old and new.
	diffs := DiffSet{}
	for _, name := range keyIndexes.Keys {
		index := keyIndexes.Indexes[name]
		// Retrieve old and new values.
		oldField := oldVal.FieldByIndex(index)
		newField := newVal.FieldByIndex(index)
		switch oldFieldVal := oldField.Interface().(type) {
		case int:
			newFieldVal := newField.Interface().(int)
			if oldFieldVal != newFieldVal {
				diffs[name] = Diff{name, oldFieldVal, newFieldVal}
			}
			break
		case *int:
			newFieldVal := newField.Interface().(*int)
			diff := Diff{Key: name}
			if oldFieldVal == nil {
				diff.Old = nil
			} else {
				diff.Old = *oldFieldVal
			}
			if newFieldVal == nil {
				diff.New = nil
			} else {
				diff.New = *newFieldVal
			}
			if !diff.Match() {
				diffs[name] = diff
			}
			break
		case bool:
			newFieldVal := newField.Interface().(bool)
			if oldFieldVal != newFieldVal {
				diffs[name] = Diff{name, oldFieldVal, newFieldVal}
			}
			break
		case *bool:
			newFieldVal := newField.Interface().(*bool)
			diff := Diff{Key: name}
			if oldFieldVal == nil {
				diff.Old = nil
			} else {
				diff.Old = *oldFieldVal
			}
			if newFieldVal == nil {
				diff.New = nil
			} else {
				diff.New = *newFieldVal
			}
			if !diff.Match() {
				diffs[name] = diff
			}
			break
		case float64:
			newFieldVal := newField.Interface().(float64)
			if oldFieldVal != newFieldVal {
				diffs[name] = Diff{name, oldFieldVal, newFieldVal}
			}
			break
		case *float64:
			newFieldVal := newField.Interface().(*float64)
			diff := Diff{Key: name}
			if oldFieldVal == nil {
				diff.Old = nil
			} else {
				diff.Old = *oldFieldVal
			}
			if newFieldVal == nil {
				diff.New = nil
			} else {
				diff.New = *newFieldVal
			}
			if !diff.Match() {
				diffs[name] = diff
			}
			break
		case string:
			newFieldVal := newField.Interface().(string)
			if oldFieldVal != newFieldVal {
				diffs[name] = Diff{name, oldFieldVal, newFieldVal}
			}
			break
		case *string:
			newFieldVal := newField.Interface().(*string)
			diff := Diff{Key: name}
			if oldFieldVal == nil {
				diff.Old = nil
			} else {
				diff.Old = *oldFieldVal
			}
			if newFieldVal == nil {
				diff.New = nil
			} else {
				diff.New = *newFieldVal
			}
			if !diff.Match() {
				diffs[name] = diff
			}
			break
		case time.Time:
			newFieldVal := newField.Interface().(time.Time)
			if !oldFieldVal.Equal(newFieldVal) {
				diffs[name] = Diff{name, oldFieldVal, newFieldVal}
			}
			break

		case *time.Time:
			newFieldVal := newField.Interface().(*time.Time)
			diff := Diff{Key: name}
			if oldFieldVal == nil {
				diff.Old = nil
			} else {
				diff.Old = *oldFieldVal
			}
			if newFieldVal == nil {
				diff.New = nil
			} else {
				diff.New = *newFieldVal
			}
			// If both are nil, there's no change.
			if oldFieldVal == nil && newFieldVal == nil {
				break
			}
			// If either value is nil, i.e. one was nil and the other wasn't
			// OR if the values are not equal, then add the diff.
			if (oldFieldVal == nil || newFieldVal == nil) || !(*oldFieldVal).Equal(*newFieldVal) {
				diffs[name] = diff
			}
			break
		}
	}

	// Retrieve the fields from the mapper.
	return diffs, nil
}

// type Diff struct {
// 	FieldName string
// 	Value1    string
// 	Value2    string
// }

// func Compare(struct1 interface{}, struct2 interface{}) ([]*Diff, error) {
//
// 	numFields := structVal1.NumField()
// 	results := make([]*Diff, 0, numFields)
//
// 	for i := 0; i < numFields; i++ {
// 		//Get values of the structure's fields
// 		field1 := structVal1.Field(i)
// 		field2 := structVal2.Field(i)
//
// 		//If the field name is unexported, skip
// 		if structType.Field(i).PkgPath != "" {
// 			continue
// 		}
//
// 		//Handle nil pointers
// 		isPtr := field1.Kind() == reflect.Ptr
// 		isField1Nil := false
// 		isField2Nil := false
// 		if isPtr {
// 			isField1Nil = field1.IsNil()
// 			isField2Nil = field2.IsNil()
// 		}
//
// 		//If both fields are nil, continue the loop
// 		if isPtr && isField1Nil && isField2Nil {
// 			continue
// 		}
//
// 		switch val1 := field1.Interface().(type) {
// 		case int:
// 			val2 := field2.Interface().(int)
// 			if val1 != val2 {
// 				result := &Diff{structType.Field(i).Name, strconv.Itoa(val1), strconv.Itoa(val2)}
// 				results = append(results, result)
// 			}
// 			break
// 		case *int:
// 			val2 := field2.Interface().(*int)
// 			var int1 int
// 			var int2 int
// 			if isField1Nil {
// 				int1 = 0
// 			} else {
// 				int1 = *val1
// 			}
// 			if isField2Nil {
// 				int2 = 0
// 			} else {
// 				int2 = *val2
// 			}
// 			if int1 != int2 {
// 				result := &CompareResult{structType.Field(i).Name, strconv.Itoa(int1), strconv.Itoa(int2)}
// 				results = append(results, result)
// 			}
// 			break
// 		case bool:
// 			val2 := field2.Interface().(bool)
// 			if val1 != val2 {
// 				result := &CompareResult{structType.Field(i).Name, strconv.FormatBool(val1), strconv.FormatBool(val2)}
// 				results = append(results, result)
// 			}
// 			break
// 		case *bool:
// 			val2 := field2.Interface().(*bool)
// 			var bool1 bool
// 			var bool2 bool
// 			if isField1Nil {
// 				bool1 = false
// 			} else {
// 				bool1 = *val1
// 			}
// 			if isField2Nil {
// 				bool2 = false
// 			} else {
// 				bool2 = *val2
// 			}
// 			if bool1 != bool2 {
// 				result := &CompareResult{structType.Field(i).Name, strconv.FormatBool(bool1), strconv.FormatBool(bool2)}
// 				results = append(results, result)
// 			}
// 			break
// 		case float64:
// 			val2 := field2.Interface().(float64)
// 			if val1 != val2 {
// 				result := &CompareResult{structType.Field(i).Name, strconv.FormatFloat(val1, 'f', 2, 64), strconv.FormatFloat(val2, 'f', 2, 64)}
// 				results = append(results, result)
// 			}
// 			break
// 		case *float64:
// 			val2 := field2.Interface().(*float64)
// 			var float1 float64
// 			var float2 float64
// 			if isField1Nil {
// 				float1 = 0
// 			} else {
// 				float1 = *val1
// 			}
// 			if isField2Nil {
// 				float2 = 0
// 			} else {
// 				float2 = *val2
// 			}
// 			if float1 != float2 {
// 				result := &CompareResult{structType.Field(i).Name, strconv.FormatFloat(float1, 'f', 2, 64), strconv.FormatFloat(float2, 'f', 2, 64)}
// 				results = append(results, result)
// 			}
// 			break
// 		case string:
// 			val2 := field2.Interface().(string)
// 			if val1 != val2 {
// 				result := &CompareResult{structType.Field(i).Name, val1, val2}
// 				results = append(results, result)
// 			}
// 			break
// 		case *string:
// 			val2 := field2.Interface().(*string)
// 			var string1 string
// 			var string2 string
// 			if isField1Nil {
// 				string1 = ""
// 			} else {
// 				string1 = *val1
// 			}
// 			if isField2Nil {
// 				string2 = ""
// 			} else {
// 				string2 = *val2
// 			}
// 			if string1 != string2 {
// 				result := &CompareResult{structType.Field(i).Name, string1, string2}
// 				results = append(results, result)
// 			}
// 			break
// 		case time.Time:
// 			val2 := field2.Interface().(time.Time)
// 			if val1 != val2 {
// 				result := &CompareResult{structType.Field(i).Name, val1.Format(config.DateTimeFormat), val2.Format(config.DateTimeFormat)}
// 				results = append(results, result)
// 			}
// 			break
//
// 		case *time.Time:
// 			val2 := field2.Interface().(*time.Time)
// 			var time1 string
// 			var time2 string
// 			if isField1Nil {
// 				time1 = ""
// 			} else {
// 				time1 = val1.Format(config.DateTimeFormat)
// 			}
// 			if isField2Nil {
// 				time2 = ""
// 			} else {
// 				time2 = val2.Format(config.DateTimeFormat)
// 			}
// 			if time1 != time2 {
// 				result := &CompareResult{structType.Field(i).Name, time1, time2}
// 				results = append(results, result)
// 			}
// 			break
// 		default:
// 			return nil, errors.New(fmt.Sprintf("Unsupported type: %v", val1))
// 		}
//
// 	}
//
// 	return results, nil
// }
