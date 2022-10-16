package main

import (
	"fmt"
	"github.com/mattfenwick/collections/pkg/json"
	gunparse "github.com/mattfenwick/gunparse/pkg"
	"github.com/mattfenwick/gunparse/pkg/example"
	"github.com/mattfenwick/json-lens/pkg"
	"github.com/sirupsen/logrus"
)

func main() {
	result := JsonAST(`{"a": 1, "b": 2, "a": 3}`)
	if result.Success == nil {
		panic(result)
	}

	logrus.Infof(json.MustMarshalToString(result.Success.Value))

	obj := pkg.Object(map[string]pkg.JsonValue{
		"a": pkg.Number(1),
		"b": pkg.String("qrs"),
	})

	keyMatch := pkg.Traverse(obj, pkg.MatchKey("a"))
	for _, match := range keyMatch {
		fmt.Printf("match key: %s\n", json.MustMarshalToString(match))
	}
	fmt.Printf("\n\n")

	allKeyMatches := pkg.Traverse(obj, pkg.MatchAllKeys())
	for _, match := range allKeyMatches {
		fmt.Printf("match all keys: %s\n", json.MustMarshalToString(match))
	}
	fmt.Printf("\n\n")

	firstKeyMatches := pkg.Traverse(obj, pkg.MatchFirstKey())
	for _, match := range firstKeyMatches {
		fmt.Printf("match first key: %s\n", json.MustMarshalToString(match))
	}
	fmt.Printf("\n\n")
}

func JsonAST(input string) gunparse.Result[example.ParseError, *gunparse.Pair[int, int], rune, *example.Object] {
	return example.ObjectParser.Parse(example.StringToRunes(input), gunparse.NewPair[int, int](1, 1))
}
