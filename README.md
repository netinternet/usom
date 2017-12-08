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
	masks := []string{"89.43.28.0/22", "89.43.26.0/22"}
	list := usom.Scan(masks, 10)
	for _, v := range list {
		fmt.Println(v.Hostname, v.IP)
	}
}



```
