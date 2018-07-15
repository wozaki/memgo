# memgo
memgo is a Memcached client for Go.

## Example
```go
import (
	"github.com/wozaki/memgo"
)

func main() {
	mc := memgo.NewClient([]string{"cache1.example.com:11211", "cache2.example.com:11211"})
	mc.Set(mc.Item{Key: "foo", Value: "my value"})

	res, err := mc.Get("foo")
	...
}
```

## Installation
```
go get github.com/wozaki/memgo
```

## License
MIT Licensed. See the LICENSE file for details.
