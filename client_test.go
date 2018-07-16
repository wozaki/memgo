package memgo

import (
	"testing"
	"net"
	"math/rand"
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

func generateRndomString(n int) string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}

func TestSet(t *testing.T) {
	t.Run("when the key size is 250", func(t *testing.T) {
		key := generateRndomString(250)
		err := Set(Item{Key: key, Value: "123"})
		if err != nil {
			t.Errorf("actual %v, expected %v", err, "nil")
		}
	})

	t.Run("when the key size is 251", func(t *testing.T) {
		key := generateRndomString(251)
		err := Set(Item{Key: key, Value: "123"})
		if err.Error() != "memcached returned CLIENT_ERROR: CLIENT_ERROR bad command line format" {
			t.Errorf("actual %v, expected %v", err.Error() , "memcached returned CLIENT_ERROR: CLIENT_ERROR")
		}
	})
}

func TestSetAndGet(t *testing.T) {
	flushAll(t)

	t.Run("when the key size is 250", func(t *testing.T) {
		key := generateRndomString(250)
		value := "123"

		Set(Item{Key: key, Value: value, Flags: 1, Exptime: 0})
		actual, err := Get(key)

		if actual.Value != "123" {
			t.Errorf("actual %v, expected %v", actual, "123")
		}
		if actual.Flags != 1 {
			t.Errorf("actual %v, expected %v", actual, "1")
		}

		if err != nil {
			t.Errorf("return error %v", err)
		}
	})

	t.Run("when the key size is 251", func(t *testing.T) {
		key := generateRndomString(251)
		actual, err := Get(key)
		if actual != nil {
			t.Errorf("actual %v, expected %v", actual, "nil")
		}
		if err.Error() != "memcached returned CLIENT_ERROR: CLIENT_ERROR" {
			t.Errorf("actual %v, expected %v", err.Error() , "memcached returned CLIENT_ERROR: CLIENT_ERROR")
		}
	})
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

	addedErr := Add(Item{Key: key, Value: value})
	if addedErr != nil {
		t.Errorf("expected addedErr is nil")
	}

	actual, _ := Get(key)
	if actual.Value != "123" {
		t.Errorf("actual %v, expected %v", actual, "123")
	}

	addAgain := Add(Item{Key: key, Value: value})
	if addAgain != ErrorNotStored {
		t.Errorf("Add must return ErrorNotStored given the same key")
	}
}

//TODO: Test values are scattered
func TestSharding(t *testing.T) {
	flushAll(t)

	key := "test_key"
	value := "123"

	client := NewClient([]string{"localhost:11211", "localhost:11212"})

	client.Set(Item{Key: key, Value: value})
	actual, err := Get(key)

	if actual.Value != "123" {
		t.Errorf("actual %v, expected %v", actual, "123")
	}

	if err != nil {
		t.Errorf("return error %v", err)
	}
}
