go-giggle
---------

golang wrapper to [giggle](https://github.com/ryanlayer/giggle)

[![GoDoc](https://godoc.org/github.com/brentp/go-giggle?status.png)](https://godoc.org/github.com/brentp/go-giggle)

```Go

import (
    giggle "github.com/brentp/go-giggle"
    "fmt"
) 

func main() {

    index := giggle.Open("/path/to/index")
    res := index.Query("1", 565657, 567999)

    // all files in the index
    index.Files()

    // int showing total count
    res.TotalHits()

    // []uint32 giving number of hits for each file
    res.Hits()

    var lines []string
    # access results by index of file.
    lines = res.Of(0)
    fmt.Println(strings.Join(lines, "\n"))
    lines = res.Of(1)
}
```

setup
=====

```
grep -l "int main" giggle/src/*.c | xargs -I{} rm {} giggle/src/$(basename {} .c).h
# timer start() conflicts with jsw start()
rm giggle/src/timer.h
grep -v "timer.h" giggle/src/giggle_index.c > $$.tmp; mv $$.tmp giggle/src/giggle_index.c
```
