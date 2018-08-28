package calc

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	sc "text/scanner"
)

type Token uint

const (
	ILLEGAL Token = iota

	NUMBER
	ADDITION
	SUBTRACTION
	MULTIPLICATION
	DIVISION

	LBRACKET
	RBRACKET
)

type Lexeme struct {
	Token Token
	Value string
}

// parsing the input string to tokens
func Lexer(r io.Reader) (out []Lexeme) {
	var s sc.Scanner
	s.Init(r)
	s.Mode = sc.ScanInts

	var tok rune
	for tok != sc.EOF {
		tok = s.Scan()
		mask := 1 << -uint(tok)

		var lexeme Lexeme
		if mask == sc.ScanInts {
			lexeme = Lexeme{NUMBER, s.TokenText()}
			out = append(out, lexeme)
			continue
		} else if tok == sc.EOF {
			break
		}

		switch string(tok) {
		case "+":
			lexeme = Lexeme{ADDITION, s.TokenText()}
		case "-":
			lexeme = Lexeme{SUBTRACTION, s.TokenText()}
		case "*":
			lexeme = Lexeme{MULTIPLICATION, s.TokenText()}
		case "/":
			lexeme = Lexeme{DIVISION, s.TokenText()}
		case "(":
			lexeme = Lexeme{LBRACKET, s.TokenText()}
		case ")":
			lexeme = Lexeme{RBRACKET, s.TokenText()}
		default:
			lexeme = Lexeme{ILLEGAL, s.TokenText()}
		}
		out = append(out, lexeme)
	}

	return
}

// converting a sequence of tokens into a RPN sequence
func Parser(in []Lexeme) (rpn []Lexeme, err error) {
	stack := make([]Lexeme, 0)
	var top Lexeme
	for num, lex := range in {
		if len(stack) > 0 {
			top = stack[len(stack)-1]
		}

		switch {
		case lex.Token == ILLEGAL:
			return nil, fmt.Errorf("illegal token detected")
		case lex.Token == NUMBER:
			rpn = append(rpn, lex)
		case len(stack) == 0:
			if len(rpn) == 0 && lex.Token == SUBTRACTION {
				rpn = append(rpn, Lexeme{NUMBER, "0"})
			}
			stack = append(stack, lex)
		case lex.Token == LBRACKET:
			stack = append(stack, lex)
		case lex.Token == RBRACKET:
			for i := len(stack) - 1; i >= 0; i-- {
				top := stack[i]
				stack = stack[:i]
				if top.Token == RBRACKET {
					return nil, fmt.Errorf("unpaired bracket")
				} else if top.Token == LBRACKET {
					break
				} else if i == 0 {
					return nil, fmt.Errorf("brackets not matched")
				}
				rpn = append(rpn, top)
			}
		case lex.Token == ADDITION || lex.Token == SUBTRACTION:
			prevLex := in[num-1]
			if top.Token != LBRACKET && top.Token != RBRACKET {
				rpn = append(rpn, top)
				stack = stack[:len(stack)-1]
			} else if lex.Token == SUBTRACTION && prevLex.Token == LBRACKET {
				rpn = append(rpn, Lexeme{NUMBER, "0"})
			}
			stack = append(stack, lex)
		case (lex.Token == MULTIPLICATION || lex.Token == DIVISION) &&
			(top.Token == MULTIPLICATION || top.Token == DIVISION):
			rpn = append(rpn, top)
			stack = stack[:len(stack)-1]
			stack = append(stack, lex)
		default:
			stack = append(stack, lex)
		}
	}

	for i := len(stack) - 1; i >= 0; i-- {
		top := stack[i]
		stack = stack[:i]
		if top.Token != LBRACKET && top.Token != RBRACKET {
			rpn = append(rpn, top)
		} else {
			return nil, fmt.Errorf("brackets not matched")
		}
	}

	return
}

// evaluation of expression
func Calc(input io.Reader, output io.Writer) (res int, err error) {
	in := bufio.NewReader(input)

	lexems := Lexer(in)
	if lexems == nil {
		return 0, nil
	}
	rpn, err := Parser(lexems)
	if err != nil {
		return 0, err
	}
	//fmt.Println("rpn:", rpn) // printing RPN sequence

	stack := make([]int, 0)
	for _, lex := range rpn {
		switch {
		case lex.Token == NUMBER:
			intTok, _ := strconv.Atoi(lex.Value)
			stack = append(stack, intTok)
		case len(stack) < 2:
			return 0, fmt.Errorf("invalid expression")
		case lex.Token == ADDITION:
			res := stack[len(stack)-2] + stack[len(stack)-1]
			stack = stack[:len(stack)-2]
			stack = append(stack, res)
		case lex.Token == SUBTRACTION:
			res := stack[len(stack)-2] - stack[len(stack)-1]
			stack = stack[:len(stack)-2]
			stack = append(stack, res)
		case lex.Token == MULTIPLICATION:
			res := stack[len(stack)-2] * stack[len(stack)-1]
			stack = stack[:len(stack)-2]
			stack = append(stack, res)
		case lex.Token == DIVISION:
			res := stack[len(stack)-2] / stack[len(stack)-1]
			stack = stack[:len(stack)-2]
			stack = append(stack, res)
		}
	}

	res = stack[0]
	fmt.Fprint(output, res)
	return
}
