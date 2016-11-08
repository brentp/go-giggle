go-giggle
---------

golang wrapper to [giggle](https://github.com/ryanlayer/giggle)

[![GoDoc] (https://godoc.org/github.com/brentp/go-giggle?status.png)](https://godoc.org/github.com/brentp/go-giggle)

```Go

import (
    giggle "github.com/brentp/go-giggle"
    "fmt"
) 

func main() {

    g := giggle.Open("/path/to/index")
    res := g.Query("1", 565657, 567999)

    var lines []string
    # access results by index of file.
    lines = res.Of(0)
    fmt.Println(strings.Join(lines, "\n"))
    lines = res.Of(1)
}
```
