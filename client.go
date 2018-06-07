package memgo

import (
	"fmt"
	"strconv"
	"strings"
	"bufio"
)

// https://github.com/memcached/memcached/blob/master/doc/protocol.txt

type Client struct {
	Host      string
	Port      int
	Transport string
}

var DefaultClient = &Client{Host: "localhost", Port: 11211, Transport: "tcp"}

type Response struct {
	Status string
	Val    string // TODO: use generics?
}

//TODO: define as specific type
func Set(k string, v string, flags int, exptime int) (resp *Response, err error) {
	return DefaultClient.Set(k, v, flags, exptime)
}

func Get(k string) (resp *Response, err error) {
	return DefaultClient.Get(k)
}

const Newline = "\r\n"

func (c *Client) Set(k string, v string, flags int, exptime int) (resp *Response, err error) {
	conn := NewConnection(c)
	defer conn.Close()

	command := "set"
	byteSize := len(v)

	req := []string{command, k, strconv.Itoa(flags), strconv.Itoa(exptime), strconv.Itoa(byteSize)}
	conn.Write([]byte(strings.Join(req, " ") + Newline + v + Newline))

	reply := make([]byte, 1024)
	_, err = conn.Read(reply)
	if err != nil {
		panic(err)
	}
	fmt.Println("SET", string(reply))

	var r = &Response{}
	r.Status = k + ":" + v

	return r, nil
}

func (c *Client) Get(k string) (resp *Response, err error) {
	conn := NewConnection(c)
	defer conn.Close()

	command := "get"

	req := []string{command, k}
	conn.Write([]byte(strings.Join(req, " ") + Newline))

	var r = &Response{}
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "VALUE") {
			//TODO: parse key, flag, and exptime
			scanner.Scan()
			r.Val = scanner.Text()
		}

		if scanner.Text() == "END" {
			break
		}
	}

	return r, nil
}
