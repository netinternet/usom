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
	"github.com/netinternet/usom"
)

func main() {
	list := usom.Scan([]string{"89.43.28.0/22"}, 1)
	for _, v := range list {
		fmt.Println(v.Hostname, v.IP)
	}
}



```
