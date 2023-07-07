package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExecuteHandler(t *testing.T) {
	code := "['welcome', 'to', 'risor', '👋'] | strings.join(' ')"
	req, err := http.NewRequest("POST", "/execute", bytes.NewBuffer([]byte(code)))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	executeHandler(rr, req)

	require.Equal(t, 200, rr.Code)

	var response interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}

	responseText, ok := response.(string)
	require.True(t, ok)
	require.Equal(t, "welcome to risor 👋", responseText)
}
