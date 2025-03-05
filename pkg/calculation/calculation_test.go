package calculation_test

import (
	"reflect"
	"strings"
	"testing"

	calc "github.com/AlexS25/rpn/pkg/calculation"
)

func TestIsNumber(t *testing.T) {
	val := "12.34"
	got := calc.IsNumber(val)
	want := true
	if got != want {
		t.Errorf("IsNumber(%q) = %v, want %v", val, got, want)
	}

	val = "1234"
	got = calc.IsNumber(val)
	want = true
	if got != want {
		t.Errorf("IsNumber(%q) = %v, want %v", val, got, want)
	}

	val = "a"
	got = calc.IsNumber(val)
	want = false
	if got != want {
		t.Errorf("IsNumber(%q) = %v, want %v", val, got, want)
	}
}

func TestParseExpr(t *testing.T) {

	val := "2 + 1(3 - 8/4 * 1) + 5"
	got, _ := calc.ParseExpr(val)
	want := []string{
		"2",
		"+",
		"1",
		"(",
		"3",
		"-",
		"8",
		"/",
		"4",
		"*",
		"1",
		")",
		"+",
		"5",
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ParseExpr(%q) = %q, want %q", val, strings.Join(got, ""), strings.Join(want, ""))
	}

	val = val + "aa"
	_, err := calc.ParseExpr(val)
	if err == nil {
		t.Fatalf("successfull case %q returns errors", val)
	}
}

func TestCheckBrackets(t *testing.T) {
	var cases = []struct {
		name   string
		values []string
		want   bool
	}{
		{
			name:   "correct values",
			values: []string{"2", "+", "1", "*", "(", "3", "-", "8", "/", "4", "*", "1", ")", "+", "5"},
			want:   true,
		},
		{
			name:   "no opening bracket",
			values: []string{"2", "+", "1", "*", "3", "-", "8", "/", "4", "*", "1", ")", "+", "5"},
			want:   false,
		},
		{
			name:   "no closing bracket",
			values: []string{"2", "+", "1", "*", "(", "3", "-", "8", "/", "4", "*", "1", "+", "5"},
			want:   false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, _ := calc.CheckBrackets(tc.values)
			if got != tc.want {
				t.Errorf("CheckBrackets(%q) = %t, but want %t", strings.Join(tc.values, ""), got, tc.want)
			}
		})
	}
}

func TestCheckSyntax(t *testing.T) {
	var cases = []struct {
		name   string
		values []string
		want   bool
	}{
		{
			name:   "correct values",
			values: []string{"2", "+", "1", "*", "(", "3", "-", "8", "/", "4", "*", "1", ")", "+", "5"},
			want:   true,
		},
		{
			name:   "extra operator",
			values: []string{"2", "+", "-", "1", "*", "(", "3", "-", "8", "/", "4", "*", "1", ")", "+", "5"},
			want:   false,
		},
		{
			name:   "letter symbol",
			values: []string{"2a", "+", "1", "*", "(", "3", "-", "8", "/", "4", "*", "1", "+", "5"},
			want:   false,
		},
		{
			name:   "empty",
			values: []string{},
			want:   false,
		},
		{
			name:   "nil value",
			values: nil,
			want:   false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, _ := calc.CheckSyntax(tc.values)
			if got != tc.want {
				t.Errorf("CheckSyntax(%q) = %t, but want %t", strings.Join(tc.values, ""), got, tc.want)
			}
		})
	}
}

func TestCalc(t *testing.T) {
	testCasesSuccess := []struct {
		name           string
		expression     string
		expectedResult float64
	}{
		{
			name:           "simple",
			expression:     "1+1",
			expectedResult: 2,
		},
		{
			name:           "priority",
			expression:     "(2+2)*2",
			expectedResult: 8,
		},
		{
			name:           "priority",
			expression:     "2+2*2",
			expectedResult: 6,
		},
		{
			name:           "/",
			expression:     "1/2",
			expectedResult: 0.5,
		},
	}

	for _, testCase := range testCasesSuccess {
		t.Run(testCase.name, func(t *testing.T) {
			val, err := calc.Calc(testCase.expression)
			if err != nil {
				t.Fatalf("successful case %s returns error", testCase.expression)
			}
			if val != testCase.expectedResult {
				t.Fatalf("%f should be equal %f", val, testCase.expectedResult)
			}
		})
	}

	testCasesFail := []struct {
		name        string
		expression  string
		expectedErr error
	}{
        {
            name:       "simple",
            expression: "1+1*",
        },
        {
            name:       "priority",
            expression: "2+2**2",
        },
        {
            name:       "extra bracket 1",
            expression: "(2+2-*2",
        },
        {
            name:       "extra bracket 2",
            expression: "2+2)-*2",
        },
        {
            name:       "extra bracket 3",
            expression: "2(+2-*2",
        },
        {
            name:       "extra bracket 4",
            expression: "2)+2-*2",
        },
        {
            name:       "empty",
            expression: "",
        },
        {
          name:         "division by zero",
          expression:   "1/0",
        },
	}

	for _, testCase := range testCasesFail {
		t.Run(testCase.name, func(t *testing.T) {
			val, err := calc.Calc(testCase.expression)
			if err == nil {
				t.Fatalf("expression %s is invalid but result  %f was obtained", testCase.expression, val)
			}
		})
	}
}
