package calc

import (
	"bytes"
	"strings"
	"testing"
)

func equal(a, b []Lexeme) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

type testLexer struct {
	input  string
	output []Lexeme
}

type testParser struct {
	input  []Lexeme
	output []Lexeme
}

type testCalc struct {
	input  string
	output int
}

var testsLexer = []testLexer{
	{"", nil},
	{"---5", []Lexeme{{3, "-"}, {3, "-"}, {3, "-"}, {1, "5"}}},
	{"(5 +   2)*3   -  (   -           2)", []Lexeme{{6, "("}, {1, "5"}, {2, "+"}, {1, "2"}, {7, ")"}, {4, "*"}, {1, "3"}, {3, "-"}, {6, "("}, {3, "-"}, {1, "2"}, {7, ")"}}},
	{"illegal", []Lexeme{{0, "i"}, {0, "l"}, {0, "l"}, {0, "e"}, {0, "g"}, {0, "a"}, {0, "l"}}},
}

var testsParser = []testParser{
	{[]Lexeme{{6, "("}, {1, "6"}, {2, "+"}, {1, "10"}, {3, "-"}, {1, "4"}, {7, ")"}, {5, "/"}, {6, "("}, {1, "1"}, {2, "+"}, {1, "1"}, {4, "*"}, {1, "2"}, {7, ")"}, {2, "+"}, {1, "1"}},
		[]Lexeme{{1, "6"}, {1, "10"}, {2, "+"}, {1, "4"}, {3, "-"}, {1, "1"}, {1, "1"}, {1, "2"}, {4, "*"}, {2, "+"}, {5, "/"}, {1, "1"}, {2, "+"}}},
	{[]Lexeme{{1, "333"}, {2, "+"}, {6, "("}, {3, "-"}, {1, "11"}, {7, ")"}, {3, "-"}, {6, "("}, {1, "22"}, {2, "+"}, {1, "22"}, {7, ")"}, {2, "+"}, {6, "("}, {6, "("}, {1, "5"}, {3, "-"}, {1, "40"}, {7, ")"}, {7, ")"}},
		[]Lexeme{{1, "333"}, {1, "0"}, {1, "11"}, {3, "-"}, {2, "+"}, {1, "22"}, {1, "22"}, {2, "+"}, {3, "-"}, {1, "5"}, {1, "40"}, {3, "-"}, {2, "+"}}},
	{[]Lexeme{{3, "-"}, {6, "("}, {3, "-"}, {6, "("}, {3, "-"}, {1, "5"}, {7, ")"}, {7, ")"}},
		[]Lexeme{{1, "0"}, {1, "0"}, {1, "0"}, {1, "5"}, {3, "-"}, {3, "-"}, {3, "-"}}},
}

var testsCalc = []testCalc{
	{"5 * (2 + 3)", 25},
	{"(6+10-4)/(1+1*2)+1", 5},
	{"333+(-11)-(22+22)+((5-40))", 243},
	{"-(-(-5))", -5},
}

func TestLexer(t *testing.T) {
	for _, test := range testsLexer {
		in := strings.NewReader(test.input)
		if res := Lexer(in); !equal(res, test.output) {
			t.Errorf("test for Lexer Failed - results not match\n %v %v", res, test.output)
		}
	}
}

func TestParser(t *testing.T) {
	for _, test := range testsParser {
		res, err := Parser(test.input)
		if err != nil {
			t.Errorf("test for Parser Failed - error")
		}
		if !equal(res, test.output) {
			t.Errorf("test for Parser Failed - results not match\n %v %v", res, test.output)
		}
	}
}

func TestCalc(t *testing.T) {
	for _, test := range testsCalc {
		in := strings.NewReader(test.input)
		out := new(bytes.Buffer)
		res, err := Calc(in, out)
		if err != nil {
			t.Errorf("test for Calc Failed - error")
		}
		if res != test.output {
			t.Errorf("test for Calc Failed - results not match\n %v %v", res, test.output)
		}
	}
}
