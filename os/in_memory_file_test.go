package os

import (
	"fmt"
	"testing"
)

func TestInMemoryFile(t *testing.T) {
	var res string
	buff := &InMemoryFile{}
	buff.Write([]byte("123"))
	buff.Rewind()
	res = fmt.Sprint(string(buff.Bytes()))
	if res != "123" {
		t.Errorf("unexpected result: %s", res)
	}
	buff.Seek(buff.Len())
	buff.Write([]byte("456"))
	buff.Rewind()
	res = fmt.Sprint(string(buff.Bytes()))
	if res != "123456" {
		t.Errorf("unexpected result: %s", res)
	}
	buff.Seek(3)
	buff.Write([]byte("321"))
	buff.Rewind()
	res = fmt.Sprint(string(buff.Bytes()))
	if res != "123321" {
		t.Errorf("unexpected result: %s", res)
	}
	buff.Seek(3)
	buff.Write([]byte("456789"))
	buff.Rewind()
	res = fmt.Sprint(string(buff.Bytes()))
	if res != "123456789" {
		t.Errorf("unexpected result: %s", res)
	}
}
