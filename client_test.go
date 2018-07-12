package memgo

import (
	"testing"
	"net"
)

const testServer = "localhost:11211"

func flushAll(t *testing.T) bool {
	c, err := net.Dial("tcp", testServer)
	if err != nil {
		panic("can't connect" + testServer)
	}
	c.Write([]byte("flush_all\r\n"))
	c.Close()
	return true
}

func TestSetAndGet(t *testing.T) {
	flushAll(t)

	key := "test_key"
	value := "123"

	Set(key, value, 1, 0)
	actual, err := Get(key)

	if actual.Val != "123" {
		t.Errorf("actual %v, expected %v", actual, "123")
	}
	if actual.Flags != 1 {
		t.Errorf("actual %v, expected %v", actual, "1")
	}
	if actual.ByteSize != 3 {
		t.Errorf("actual %v, expected %v", actual, "3")
	}

	if err != nil {
		t.Errorf("return error %v", err)
	}
}

func TestGetNothing(t *testing.T) {
	flushAll(t)

	actual, err := Get("hoge")
	if actual != nil {
		t.Errorf("actual %v, expected %v", actual, "nil")
	}

	if err != nil {
		t.Errorf("return error %v", err)
	}
}

func TestAdd(t *testing.T) {
	flushAll(t)

	key := "test_key"
	value := "123"

	addedErr := Add(key, value, 0, 0)
	if addedErr != nil {
		t.Errorf("expected addedErr is nil")
	}

	actual, _ := Get(key)
	if actual.Val != "123" {
		t.Errorf("actual %v, expected %v", actual, "123")
	}

	addAgain := Add(key, value, 0, 0)
	if addAgain != ErrorNotStored {
		t.Errorf("Add must return ErrorNotStored given the same key")
	}
}

//TODO: Test values are scattered
func TestSharding(t *testing.T) {
	flushAll(t)

	key := "test_key"
	value := "123"

	client := NewClient([]string{"localhost:11211", "localhost:11212"}, "tcp")

	client.Set(key, value, 0, 0)
	actual, err := Get(key)

	if actual.Val != "123" {
		t.Errorf("actual %v, expected %v", actual, "123")
	}

	if err != nil {
		t.Errorf("return error %v", err)
	}
}
