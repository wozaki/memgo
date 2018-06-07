package memgo

import (
	"testing"
)

func TestSetAndGet(t *testing.T) {
	key := "test_key"
	value := "123"

	Set(key, value, 0, 0)
	actual, err := Get(key)

	if actual.Val != "123" {
		t.Errorf("actual %v, expected %v", actual, "123")
	}

	if err != nil {
		t.Errorf("return error %v", err)
	}
}
