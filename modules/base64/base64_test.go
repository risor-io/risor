package base64

import (
	"context"
	"encoding/base64"
	"testing"

	"github.com/risor-io/risor/object"
	"github.com/stretchr/testify/require"
)

func TestBase64Encode(t *testing.T) {
	type testCase struct {
		input string
		want  string
	}
	tests := []testCase{
		{"", ""},
		{"f", "Zg=="},
		{"fo", "Zm8="},
		{"foo", "Zm9v"},
		{"foob", "Zm9vYg=="},
		{"fooba", "Zm9vYmE="},
		{"foobar", "Zm9vYmFy"},
	}
	for _, test := range tests {
		got := Encode(context.Background(), object.NewString(test.input))
		ideal := base64.StdEncoding.EncodeToString([]byte(test.input))
		require.Equal(t, object.NewString(test.want), got)
		require.Equal(t, object.NewString(ideal), got)
	}
}

func TestBase64EncodeRaw(t *testing.T) {
	type testCase struct {
		input string
		want  string
	}
	tests := []testCase{
		{"", ""},
		{"f", "Zg"},
		{"fo", "Zm8"},
		{"foo", "Zm9v"},
		{"foob", "Zm9vYg"},
		{"fooba", "Zm9vYmE"},
		{"foobar", "Zm9vYmFy"},
	}
	for _, test := range tests {
		got := Encode(context.Background(), object.NewString(test.input), object.False)
		ideal := base64.RawStdEncoding.EncodeToString([]byte(test.input))
		require.Equal(t, object.NewString(test.want), got)
		require.Equal(t, object.NewString(ideal), got)
	}
}

func TestBase64Decode(t *testing.T) {
	type testCase struct {
		input string
		want  string
	}
	tests := []testCase{
		{"", ""},
		{"Zg==", "f"},
		{"Zm8=", "fo"},
		{"Zm9v", "foo"},
		{"Zm9vYg==", "foob"},
		{"Zm9vYmE=", "fooba"},
		{"Zm9vYmFy", "foobar"},
	}
	for _, test := range tests {
		got := Decode(context.Background(), object.NewString(test.input))
		ideal, _ := base64.StdEncoding.DecodeString(test.input)
		require.Equal(t, object.NewByteSlice([]byte(test.want)), got)
		require.Equal(t, object.NewByteSlice([]byte(ideal)), got)
	}
}

func TestBase64DecodeRaw(t *testing.T) {
	type testCase struct {
		input   string
		want    string
		wantErr string
	}
	tests := []testCase{
		{"", "", ""},
		{"Zg==", "", "illegal base64 data at input byte 2"},
		{"Zm8=", "fo", "illegal base64 data at input byte 3"},
		{"Zm9v", "foo", ""},
		{"Zm9vYg==", "foob", "illegal base64 data at input byte 6"},
		{"Zm9vYmE", "fooba", ""},
		{"Zm9vYmFy", "foobar", ""},
	}
	for _, test := range tests {
		got := Decode(context.Background(), object.NewString(test.input), object.False)
		dst := make([]byte, 32)
		count, err := base64.RawStdEncoding.Decode(dst, []byte(test.input))
		ideal := dst[:count]
		if test.wantErr != "" {
			gotErr, ok := got.(*object.Error)
			require.True(t, ok)
			require.Equal(t, test.wantErr, gotErr.Value().Error())
			require.Equal(t, test.wantErr, err.Error())
		} else {
			require.Equal(t, object.NewByteSlice([]byte(test.want)), got)
			require.Equal(t, object.NewByteSlice([]byte(ideal)), got)
		}
	}
}

func TestBase64URLEncoded(t *testing.T) {
	input := object.NewByteSlice([]byte{251})

	got := URLEncode(context.Background(), input)
	require.Equal(t, object.NewString("-w=="), got)

	got = URLEncode(context.Background(), input, object.False)
	require.Equal(t, object.NewString("-w"), got)

	got = URLDecode(context.Background(), object.NewString("-w=="))
	require.Equal(t, object.NewByteSlice([]byte{251}), got)

	got = URLDecode(context.Background(), object.NewString("-w"), object.False)
	require.Equal(t, object.NewByteSlice([]byte{251}), got)
}
