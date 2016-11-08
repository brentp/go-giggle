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

int giggle_hits(giggle_query_result *gqr, uint32_t *counts) {
	int i;
	for(i=0;i<gqr->num_files;i++) {
		counts[i] = giggle_get_query_len(gqr, i);
	}
	return 0;
}

char ** giggle_index_files(giggle_index *gi) {
	int i, n = gi->file_idx->index->num;
	char **names = (char **)malloc(sizeof(char *) * n);
	for(i=0;i<n;i++){
		names[i] = gi->file_idx[i].file_name;
	}
	return names;
}
*/
import "C"
import (
	"runtime"
	"unsafe"
)

// Index wraps the giggle index.
type Index struct {
	gi    *C.giggle_index
	files []string
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
	idx := &Index{gi: gi}
	idx.setFiles()
	return idx, nil
}

func (i *Index) setFiles() {
	files := C.giggle_index_files(i.gi)
	n := int(i.gi.file_idx.index.num)
	tmpslice := (*[1 << 30]*C.char)(unsafe.Pointer(files))[:n:n]
	i.files = make([]string, n)
	for k := 0; k < n; k++ {
		i.files[k] = C.GoString(tmpslice[k])
	}
	C.free(unsafe.Pointer(files))
}

// Query gives the results for the given genomic location.
func (i *Index) Query(chrom string, start, end int) *Result {
	cs := C.CString(chrom)
	gqr := C.giggle_query(i.gi, cs, C.uint32_t(start), C.uint32_t(end), nil)
	C.free(unsafe.Pointer(cs))
	runtime.SetFinalizer(gqr, C.giggle_query_result_destroy)
	return &Result{gqr: gqr}
}

// Files returns the files associated with the index.
func (i *Index) Files() []string {
	return i.files
}

// NFiles returns the number of files in the result-set.
func (r *Result) NFiles() int {
	return int(r.gqr.num_files)
}

// TotalHits returns the total number of overlaps in the result-set.
func (r *Result) TotalHits() int {
	return int(r.gqr.num_hits)
}

// Hits returns the number of overlaps for each file in the result-set.
func (r *Result) Hits() []uint32 {
	hits := make([]uint32, r.gqr.num_files)
	C.giggle_hits(r.gqr, (*C.uint32_t)(unsafe.Pointer(&hits[0])))
	return hits
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
