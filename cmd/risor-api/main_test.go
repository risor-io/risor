package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestExecuteHandler(t *testing.T) {
	payload := Request{
		Code: "['welcome', 'to', 'risor', '👋'] | strings.join(' ')",
	}
	payloadBytes, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", "/execute", bytes.NewBuffer(payloadBytes))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	executeHandler(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	var response Response
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}

	expectedResult := "\"welcome to risor 👋\""
	if string(response.Result) != expectedResult {
		t.Errorf("handler returned unexpected result: got %s want %s", response.Result, expectedResult)
	}

	if response.Time <= 0 {
		t.Errorf("handler returned invalid response time: %f", response.Time)
	}
}
