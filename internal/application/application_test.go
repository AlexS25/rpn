package application

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCalcHandler(t *testing.T) {

	DisableLogging = true

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
			name:           "division",
			expression:     "1/2",
			expectedResult: 0.5,
		},
	}

	response := new(Response)
	request := new(Request)

	for _, testCase := range testCasesSuccess {
		t.Run(testCase.name, func(t *testing.T) {
			request.Expression = testCase.expression
			reqJson, _ := json.Marshal(request)
			reader := bytes.NewReader(reqJson)

			req, err := http.NewRequest(http.MethodPost, "/api/v1/calculate", reader)
			if err != nil {
				t.Fatal(err)
			}

			w := httptest.NewRecorder()
			handler := http.HandlerFunc(CalcHandler)
			handler.ServeHTTP(w, req)
			res := w.Result()

			defer res.Body.Close()

			body, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal("server work incorrect: ", err)
			}

			if res.StatusCode != http.StatusOK {
				t.Fatalf("server work incorrect: got %d; want %d", res.StatusCode, http.StatusOK)
			}

			err = json.Unmarshal(body, &response)
			if err != nil {
				t.Fatal("incorrect answer received: ", err)
			}

			if response.Result != testCase.expectedResult {
				t.Fatalf("%f should be equal %f", response.Result, testCase.expectedResult)
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
			name:       "priority",
			expression: "((2+2-*(2",
		},
		{
			name:       "empty",
			expression: "",
		},
	}

	for _, testCase := range testCasesFail {
		t.Run(testCase.name, func(t *testing.T) {
			request.Expression = testCase.expression
			reqJson, _ := json.Marshal(request)
			reader := bytes.NewReader(reqJson)

			req, err := http.NewRequest(http.MethodPost, "/api/v1/calculate", reader)
			if err != nil {
				t.Fatal(err)
			}

			w := httptest.NewRecorder()
			handler := http.HandlerFunc(CalcHandler)
			handler.ServeHTTP(w, req)
			res := w.Result()

			defer res.Body.Close()

			body, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal("server work incorrect: ", err)
			}

			if res.StatusCode == http.StatusOK {
				err = json.Unmarshal(body, &response)
				if err != nil {
					t.Fatal("incorrect answer received: ", err)
				}

				t.Fatalf("expression %s is invalid but result  %f was obtained",
					testCase.expression, response.Result)
			}

			if res.StatusCode != http.StatusUnprocessableEntity {
				t.Fatalf("server work incorrect: got %d; want %d",
					res.StatusCode, http.StatusUnprocessableEntity)
			}
		})
	}
}

func main() {
	req := httptest.NewRequest("POST", "/api/v1/calculate", nil)
	w := httptest.NewRecorder()
	CalcHandler(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {

	}

}
