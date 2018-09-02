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

func generateRandomString(n int) string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}

func TestSet(t *testing.T) {
	t.Run("when the Key size is 250", func(t *testing.T) {
		key := generateRandomString(250)
		_, err := Set(Item{Key: key, Value: "123"})
		if err != nil {
			t.Errorf("actual %v, expected %v", err, "nil")
		}
	})

	t.Run("when the Key size is 251", func(t *testing.T) {
		key := generateRandomString(251)
		_, err := Set(Item{Key: key, Value: "123"})
		if err.Error() != "client error: CLIENT_ERROR bad command line format" {
			t.Errorf("actual %v, expected %v", err.Error() , "client error: CLIENT_ERROR")
		}
	})
}

func TestSetAndGet(t *testing.T) {
	t.Run("Test Key size", func(t *testing.T) {
		flushAll(t)

		// The size is 250
		key := generateRandomString(250)
		Set(Item{Key: key, Value: "123", Exptime: 0})
		actual, err := Get(key)
		if actual.Value != "123" {
			t.Errorf("actual %v, expected %v", actual, "123")
		}
		if actual.Flags != 0 {
			t.Errorf("actual %v, expected %v", actual, "0")
		}
		if actual.CasId != 0 {
			t.Errorf("actual %v, expected %v", actual.CasId, "0")
		}
		if err != nil {
			t.Errorf("return error %v", err)
		}

		// The size is 251
		key = generateRandomString(251)
		Set(Item{Key: key, Value: "123", Exptime: 0})
		actual, err = Get(key)
		if actual != nil {
			t.Errorf("actual %v, expected %v", actual, "nil")
		}
		if err.Error() != "client error: CLIENT_ERROR" {
			t.Errorf("actual %v, expected %v", err.Error() , "memcached returned CLIENT_ERROR: CLIENT_ERROR")
		}
	})

	t.Run("Test without correspondent item", func(t *testing.T) {
		flushAll(t)
		
		actual, err := Get("hoge")
		if actual != nil {
			t.Errorf("actual %v, expected %v", actual, "nil")
		}
		if err != nil {
			t.Errorf("return error %v", err)
		}
	})
}

func TestGets(t *testing.T) {
	// The size is 250
	key := generateRandomString(250)
	client := NewClient([]string{testServer}, Config{})
	client.Set(Item{Key: key, Value: "123", Flags: Flags{Value:1}, Exptime: 0})

	actual, err := client.Gets(key)
	if actual.Value != "123" {
		t.Errorf("actual %v, expected %v", actual, "123")
	}
	if actual.Flags != 1 {
		t.Errorf("actual %v, expected %v", actual, "1")
	}
	if actual.CasId == 0 {
		t.Errorf("actual %v, expected %v", actual.CasId, "not 0")
	}
	if err != nil {
		t.Errorf("return error %v", err)
	}
}

func TestAdd(t *testing.T) {
	flushAll(t)

	key := "test_key"
	value := "123"
	client := NewClient([]string{testServer}, Config{})

	_, addedErr := client.Add(Item{Key: key, Value: value})
	if addedErr != nil {
		t.Errorf("expected addedErr is nil")
	}

	actual, _ := Get(key)
	if actual.Value != "123" {
		t.Errorf("actual %v, expected %v", actual, "123")
	}

	_, addAgain := client.Add(Item{Key: key, Value: value})
	if addAgain != ErrorNotStored {
		t.Errorf("Add must return ErrorNotStored given the same Key")
	}
}

//TODO: Test values are scattered
func TestSharding(t *testing.T) {
	flushAll(t)

	key := "test_key"
	value := "123"

	client := NewClient([]string{"localhost:11211", "localhost:11212"}, Config{})

	client.Set(Item{Key: key, Value: value})
	actual, err := Get(key)

	if actual.Value != "123" {
		t.Errorf("actual %v, expected %v", actual, "123")
	}

	if err != nil {
		t.Errorf("return error %v", err)
	}
}

func TestCompress(t *testing.T) {
	key := "Key"
	val := generateRandomString(1024 * 1024)

	t.Run("with 1MB value and no CompressThresholdByte", func(t *testing.T) {
		flushAll(t)

		client := NewClient([]string{testServer}, Config{})

		_, err := client.Set(Item{Key: key, Value: val})
		if err != nil {
			t.Errorf("actual %v, expected %v", err, "nil")
		}

		actual, err := client.Get(key)
		if actual.Value != val {
			t.Errorf("actual %v, expected %v", "ac", 1)
		}
		if actual.Flags != CompressFlag {
			t.Errorf("actual %v, expected %v", actual.Flags, CompressFlag)
		}
		if actual.CasId != 0 {
			t.Errorf("actual %v, expected %v", actual.CasId, "0")
		}
		if err != nil {
			t.Errorf("return error %v", err)
		}
	})

	t.Run("with 1MB value and 2MB CompressThresholdByte", func(t *testing.T) {
		flushAll(t)

		client := NewClient([]string{testServer}, Config{CompressThresholdByte: 1024 * 1024 * 2})

		_, err := client.Set(Item{Key: key, Value: val})
		if err.Error() != "server error: SERVER_ERROR object too large for cache" {
			t.Errorf("actual %v, expected %v", err.Error(), "server error: SERVER_ERROR object too large for cache")
		}
	})

	t.Run("with 1MB value and 2MB CompressThresholdByte and CompressFlag", func(t *testing.T) {
		flushAll(t)

		client := NewClient([]string{testServer}, Config{CompressThresholdByte: 1024 * 1024 * 2})

		_, err := client.Set(Item{Key: key, Value: val, Flags: Flags{Value: CompressFlag}})
		if err != nil {
			t.Errorf("actual %v, expected %v", err, "nil")
		}

		actual, err := client.Get(key)
		if actual.Value != val {
			t.Errorf("actual %v, expected %v", "", val)
		}
		if actual.Flags != CompressFlag {
			t.Errorf("it should compress if given CompressFlag: actual %v, expected %v", actual.Flags, CompressFlag)
		}
		if err != nil {
			t.Errorf("return error %v", err)
		}
	})
}
