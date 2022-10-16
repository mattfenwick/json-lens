package pkg

import (
	"github.com/mattfenwick/collections/pkg/function"
	"github.com/mattfenwick/collections/pkg/slice"
	"github.com/sirupsen/logrus"
)

// nodes to visit:
// object
//   key/val
// array
//   index/val
// int/float
// null
// bool
// string

// examples:
//  using {"a": 1, "b": 2, "c": [14, true, null, "qrs", {"b": []}]}
//  string("a") -> nil
//  search(string("a")) -> [[key("a")]]
//  string("b") ->

type JsonValue struct {
	String *string
	Number *int // who cares about floats
	Bool   *bool
	Null   bool
	Object *map[string]JsonValue
	Array  *[]JsonValue
}

func String(s string) JsonValue {
	return JsonValue{String: &s}
}

func Number(i int) JsonValue {
	return JsonValue{Number: &i}
}

func Bool(b bool) JsonValue {
	return JsonValue{Bool: &b}
}

func Null() JsonValue {
	return JsonValue{Null: true}
}

func Object(o map[string]JsonValue) JsonValue {
	return JsonValue{Object: &o}
}

func Array(a []JsonValue) JsonValue {
	return JsonValue{Array: &a}
}

func (j *JsonValue) Get(path []*PathStep) *JsonValue {
	if path == nil {
		return j
	}
	first, rest := path[0], path[1:]
	if first.MapKey != nil && j.Object != nil {
		obj := *j.Object
		key := *first.MapKey
		if val, ok := obj[key]; ok {
			return val.Get(rest)
		} else {
			return nil
		}
	}
	if first.ArrayIndex != nil && j.Array != nil {
		arr := *j.Array
		key := *first.ArrayIndex
		if key >= 0 && key < len(arr) {
			return arr[key].Get(rest)
		} else {
			return nil
		}
	}
	logrus.Debugf("value not found: type mismatch (wanted %+v, got %+v)", first, j)
	return nil
}

type PathStep struct {
	MapKey     *string
	ArrayIndex *int
}

type Match struct {
	Path  []*PathStep
	Value JsonValue
}

type Matcher func(path []*PathStep, value JsonValue, step *PathStep) ([]*Match, bool)

func Traverse(value JsonValue, matcher Matcher) []*Match {
	matches, _ := TraverseHelp(nil, value, matcher)
	return matches
}

func TraverseHelp(path []*PathStep, value JsonValue, matcher Matcher) ([]*Match, bool) {
	copiedPath := slice.Map(function.Id[*PathStep], path)
	matches, done := matcher(copiedPath, value, nil)
	if done {
		return matches, true
	}

	// TODO need a way to control the traversal:
	//   how deep to go
	//   when to stop
	//   early exit
	//   whether to continue

	if value.Object != nil {
		for key, val := range *value.Object {
			// copy key and path
			keyCopy := key
			step := &PathStep{MapKey: &keyCopy}
			objPath := slice.Append(copiedPath, []*PathStep{step})

			// visit the key/val pair
			objMatches, done := matcher(objPath, val, step)
			matches = append(matches, objMatches...)
			if done {
				return matches, true
			}

			// recur
			recurMatches, done := TraverseHelp(objPath, val, matcher)
			matches = append(matches, recurMatches...)
			if done {
				return matches, true
			}
		}
	} else if value.Array != nil {
		for index, val := range *value.Array {
			// copy index and path
			indexCopy := index
			step := &PathStep{ArrayIndex: &indexCopy}
			arrayPath := slice.Append(copiedPath, []*PathStep{step})

			// visit the index/val pair
			arrayMatches, done := matcher(arrayPath, val, step)
			matches = append(matches, arrayMatches...)
			if done {
				return matches, true
			}

			// recur
			recurMatches, done := TraverseHelp(arrayPath, val, matcher)
			matches = append(matches, recurMatches...)
			if done {
				return matches, true
			}
		}
	}

	return matches, false
}

func MatchKey(key string) Matcher {
	return func(path []*PathStep, value JsonValue, step *PathStep) ([]*Match, bool) {
		if step != nil && step.MapKey != nil && *step.MapKey == key {
			return []*Match{{
				Path:  path,
				Value: value,
			}}, false
		}
		return nil, false
	}
}

func MatchAllKeys() Matcher {
	return func(path []*PathStep, value JsonValue, step *PathStep) ([]*Match, bool) {
		if step != nil && step.MapKey != nil {
			return []*Match{{
				Path:  path,
				Value: value,
			}}, false
		}
		return nil, false
	}
}

func MatchFirstKey() Matcher {
	return func(path []*PathStep, value JsonValue, step *PathStep) ([]*Match, bool) {
		if step != nil && step.MapKey != nil {
			return []*Match{{
				Path:  path,
				Value: value,
			}}, true
		}
		return nil, false
	}
}
