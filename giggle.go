package giggle

/*
#cgo CFLAGS: -g -O2 -fPIC -m64 -I${SRCDIR}/htslib
#cgo LDFLAGS: -ldl -lz -lm -lpthread -lhts
#include "giggle_index.h"
#include "ll.h"

typedef struct giggle_index giggle_index;
typedef struct giggle_query_result giggle_query_result;
typedef struct giggle_query_iter giggle_query_iter;

giggle_index *giggle_iload(char *data_dir) {
	return giggle_load(data_dir, uint32_t_ll_giggle_set_data_handler);
}


*/
import "C"
import (
	"runtime"
	"unsafe"
)

// Index wraps the giggle index.
type Index struct {
	gi *C.giggle_index
}

// Result is returned from a giggle query.
type Result struct {
	gqr *C.giggle_query_result
}

// Open gets an existing index at the given path.
func Open(path string) (*Index, error) {
	cs := C.CString(path)
	defer C.free(unsafe.Pointer(cs))
	gi := C.giggle_iload(cs)
	runtime.SetFinalizer(gi, C.giggle_index_destroy)
	return &Index{gi: gi}, nil
}

// Query gives the results for the given genomic location.
func (i *Index) Query(chrom string, start, end int) *Result {
	cs := C.CString(chrom)
	gqr := C.giggle_query(i.gi, cs, C.uint32_t(start), C.uint32_t(end), nil)
	C.free(unsafe.Pointer(cs))
	runtime.SetFinalizer(gqr, C.giggle_query_result_destroy)
	return &Result{gqr: gqr}
}

// Files returns the number of files in the result-set
func (r *Result) Files() int {
	return int(r.gqr.num_files)
}

// Of returns a slice of strings for the given file index.
func (r *Result) Of(i int) []string {
	gqi := C.giggle_get_query_itr(r.gqr, C.uint32_t(i))
	n := make([]string, 0, 4)
	var result *C.char
	for {
		v := C.giggle_query_next(gqi, &result)
		if v == 0 {
			break
		}
		n = append(n, C.GoString(result))
	}
	C.giggle_iter_destroy(&gqi)
	return n
}
