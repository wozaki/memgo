# memgo
memgo is a Memcached client for Go.

[![CircleCI](https://circleci.com/gh/wozaki/memgo/tree/master.svg?style=svg)](https://circleci.com/gh/wozaki/memgo/tree/master)

## Example
```go
import (
	"github.com/wozaki/memgo"
)

func main() {
	mc := memgo.NewClient([]string{"cache1.example.com:11211", "cache2.example.com:11211"}, memgo.Config{ConnectTimeout: 100 * time.Millisecond})
	mc.Set(mc.Item{Key: "foo", Value: "my value"})

	res, err := mc.Get("foo")
	...
}
```

## Installation
```
go get github.com/wozaki/memgo
```

## Features

Feature | Description
 --- |---
Compression | The format follows [compress/zlib](https://golang.org/pkg/compress/zlib/).<br><br> There are two kinds of compression:<br> 1: compress automatically if an item size over the threshold. The default threshold size is 1MB. You can configure the threshold size with `Config.CompressThresholdByte`<br> 2: compress manually if given CompressFlag with storage commands like `Set(Item{Key: key, Value: val, Flags: Flags{Value: CompressFlag}})`. It will compress the item even if the size of item under the threshold.  

## License
MIT Licensed. See the LICENSE file for details.
