package test

import (
	"go/types"
	"testing"
)

func TestIsEqual(t *testing.T) {
	type Some struct {
		Val float64
	}
	type OtherStruct struct {
		Val          int
		Name         string
		Map          map[string]int
		IntArray     []int
		SomeSome     Some
		FloatPointer *float64
		Fifth        types.Nil
		Sixth        int
	}

	firstSome := Some{Val: 2.3}
	equalToFirstSome := Some{Val: 2.3}

	exampleFloat := 3.233
	matchingExampleFloat := 3.233

	firstStruct := OtherStruct{Val: 30, Name: "Jim", Map: map[string]int{"hello": 1},
		IntArray: []int{1, 2, 3, 4, 5}, SomeSome: firstSome, FloatPointer: &exampleFloat}
	equalToFirstStruct := OtherStruct{Val: 30, Name: "Jim", Map: map[string]int{"hello": 1},
		IntArray: []int{1, 2, 3, 4, 5}, SomeSome: equalToFirstSome, FloatPointer: &matchingExampleFloat}

	//Equal (no error)
	e := IsEqual(firstStruct, equalToFirstStruct)
	if e != nil {
		t.Fatal("There was an error:\n ", e)
	}

	mapIsEmpty := OtherStruct{Val: 30, Name: "Jim", Map: map[string]int{},
		IntArray: []int{1, 2, 3, 4, 5}, SomeSome: Some{Val: 2.32}}
	e = IsEqual(mapIsEmpty, mapIsEmpty)
	if e != nil {
		t.Fatal("There was an error:\n mapIsEmpty Should have succeeded", e)
	}

	allZeroNil := OtherStruct{Val: 0, Name: "", Map: nil,
		IntArray: nil, SomeSome: Some{}, FloatPointer: nil}
	e = IsEqual(allZeroNil, allZeroNil)
	if e != nil {
		t.Fatal("There was an error:\n allZeroNil Should have succeeded", e)
	}

	//Different (should be error)
	differentValThanFirstStruct := OtherStruct{Val: 24, Name: "Jim", Map: map[string]int{"hello": 1},
		IntArray: []int{1, 2, 3, 4, 5}, SomeSome: Some{Val: 2.3}, FloatPointer: &exampleFloat}
	e = IsEqual(firstStruct, differentValThanFirstStruct)
	if e == nil {
		t.Fatal("There was an error:\n differentValThanFirstStruct should have failed equality")
	}

	differentNameThanFirstStruct := OtherStruct{Val: 30, Name: "Jimbo", Map: map[string]int{"hello": 1},
		IntArray: []int{1, 2, 3, 4, 5}, SomeSome: Some{Val: 2.3}, FloatPointer: &exampleFloat}
	e = IsEqual(firstStruct, differentNameThanFirstStruct)
	if e == nil {
		t.Fatal("There was an error:\n differentNameThanFirstStruct should have failed equality")
	}

	differentMapThanFirstStruct := OtherStruct{Val: 24, Name: "Jim", Map: map[string]int{"hello": 1, "world": 2},
		IntArray: []int{1, 2, 3, 4, 5}, SomeSome: Some{Val: 2.3}, FloatPointer: &exampleFloat}
	e = IsEqual(firstStruct, differentMapThanFirstStruct)
	if e == nil {
		t.Fatal("There was an error:\n differentMapThanFirstStruct should have failed equality")
	}

	differentIntArrayThanFirstStruct := OtherStruct{Val: 30, Name: "Jim", Map: map[string]int{"hello": 1},
		IntArray: []int{1, 2, 4, 4, 5}, SomeSome: Some{Val: 2.3}, FloatPointer: &exampleFloat}
	e = IsEqual(firstStruct, differentIntArrayThanFirstStruct)
	if e == nil {
		t.Fatal("There was an error:\n differentIntArrayThanFirstStruct should have failed equality")
	}

	differentSomeSomeThanFirstStruct := OtherStruct{Val: 30, Name: "Jim", Map: map[string]int{"hello": 1},
		IntArray: []int{1, 2, 3, 4, 5}, SomeSome: Some{Val: 1820.42129}, FloatPointer: &exampleFloat}
	e = IsEqual(firstStruct, differentSomeSomeThanFirstStruct)
	if e == nil {
		t.Fatal("There was an error:\n differentSomeSomeThanFirstStruct should have failed equality")
	}

	notMatchingExampleFloat := 192.168
	differentFloatPointerThanFirstStruct := OtherStruct{Val: 30, Name: "Jim", Map: map[string]int{"hello": 1},
		IntArray: []int{1, 2, 3, 4, 5}, SomeSome: Some{Val: 2.32}, FloatPointer: &notMatchingExampleFloat}
	e = IsEqual(firstStruct, differentFloatPointerThanFirstStruct)
	if e == nil {
		t.Fatal("There was an error:\n differentFloatPointerThanFirstStruct should have failed equality")
	}

	missingFieldFromExampleStruct := OtherStruct{Val: 30, Name: "Jim", Map: map[string]int{"hello": 1},
		IntArray: []int{1, 2, 3, 4, 5}, SomeSome: Some{Val: 2.32}}
	e = IsEqual(firstStruct, missingFieldFromExampleStruct)
	if e == nil {
		t.Fatal("There was an error:\n missingFieldFromExampleStruct should have failed equality")
	}

}
