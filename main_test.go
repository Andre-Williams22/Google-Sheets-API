package main

import (
	"testing"
)

func BenchmarkwriteJSON(b *testing.B) {

	type Test struct {
		ID     int    `json:"id"`
		Name   string `json:"Name"`
		Number string `json:"Number"`
	}

	for i := 0; i < b.N; i++ {
		sendMessage("+2142846514", "2144989333", "Testing")
	}
}

//go test -bench=.

func TestTableMessage(t *testing.T) {
	var tests = []struct {
		input    int
		expected int
	}{
		{2, 3},
		{-1, 0},
		{0, 1},
		{100, 101},
	}

	for _, test := range tests {
		if output := numPeople(test.input); output != test.expected {
			t.Error("Test Failed: {} inputted, {} expected, recieved: {}", test.input, test.expected, output)
		}
	}

}

// go test
