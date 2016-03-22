/*
Package changes implements a way to create changesets between two different
instances of the same object. This provides an easy mechanism for both audit
trails of object changes, and implementing PATCH style mechanics in APIs.

  type MyStruct struct {
   	FieldA string
   	FieldB string
  }

  original := MyStruct{"Foo", "Bar"}
  updated := MyStruct{"Foo", "Baz"}

  diff, err := changes.Diff(original, updated)

This will output a plain changeset that contains the field name by key, and a
struct with the Old and New values as the value.

  // map[string]struct{Old, New interface{}}{
  //	"FieldB": {"Bar", "Baz"},
  // }

*/
package changes
