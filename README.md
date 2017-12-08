# Usom Url Scanner
Scan usom blocked url list inside your ip masks.

## Usage

```
go get -u github.com/netinternet/usom
```

## Example

```go


package main

import (
	"fmt"
	"time"
	s "github.com/kevsersrca/usom"
)

func main() {
	list := s.Usom([]string{"89.43.28.0/22"}, time.Millisecond*100)
	for _, v := range list {
		fmt.Println(v.Hostname, v.IP)
	}
}


```
