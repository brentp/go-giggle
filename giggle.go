package giggle

/*
#cgo CFLAGS: -g -O2 -fPIC -m64 -I${SRCDIR}/htslib -I${SRCDIR}/giggle/src/
#cgo LDFLAGS: -ldl -lz -lm -lpthread -lhts
#include "giggle_include.h"

typedef struct giggle_index giggle_index;
typedef struct giggle_query_result giggle_query_result;
typedef struct giggle_query_iter giggle_query_iter;

giggle_index *giggle_iload(char *data_dir) {

	giggle_index *gi = giggle_load(data_dir, uint64_t_ll_giggle_set_data_handler);
    giggle_data_handler.giggle_collect_intersection =
            giggle_collect_intersection_data_in_block;

    giggle_data_handler.map_intersection_to_offset_list =
            leaf_data_map_intersection_to_offset_list;

	return gi;
}

int giggle_hits(giggle_query_result *gqr, uint32_t *counts) {
	int i;
	for(i=0;i<gqr->num_files;i++) {
		counts[i] = giggle_get_query_len(gqr, i);
	}
	return 0;
}

char *index_file_name(giggle_index *gi, int i) {
	return file_index_get(gi->file_idx, i)->file_name;
}

giggle_index *giggle_init2(uint32_t num_chroms, char *data_dir, uint32_t force) {

	return giggle_init(num_chroms, data_dir, force, uint64_t_ll_giggle_set_data_handler);
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
	path  *C.char
	files []string
}

// Result is returned from a giggle query.
type Result struct {
	gqr *C.giggle_query_result
}

// Open an existing index at the given path.
func Open(path string) (*Index, error) {
	cs := C.CString(path)
	gi := C.giggle_iload(cs)
	idx := &Index{gi: gi, path: cs}
	runtime.SetFinalizer(idx, destroy)
	idx.setFiles()
	return idx, nil
}

func destroy(i *Index) {
	C.giggle_index_destroy(&i.gi)
	C.free(unsafe.Pointer(i.path))
}

// New creates a new index in the given directory and files it with files matching the glob.
func New(path string, glob string) (*Index, error) {
	cs := C.CString(path)
	defer C.free(unsafe.Pointer(cs))
	cg := C.CString(glob)
	defer C.free(unsafe.Pointer(cg))
	C.giggle_bulk_insert(cg, cs, C.uint32_t(1))
	return Open(path)
}

func (i *Index) setFiles() {

	var names **C.char
	var numIntervals *C.uint32_t
	var meanIntervalSice *C.double

	nFiles := C.giggle_get_indexed_files(i.path, &names, &numIntervals, &meanIntervalSice)
	n := int(nFiles)
	tmpslice := (*[1 << 30]*C.char)(unsafe.Pointer(names))[:n:n]

	i.files = make([]string, int(nFiles))
	for k, name := range tmpslice {
		i.files[k] = C.GoString(name)
	}
}

// Query gives the results for the given genomic location.
func (i *Index) Query(chrom string, start, end int) *Result {
	cs := C.CString(chrom)
	gqr := C.giggle_query(i.gi, cs, C.uint32_t(start), C.uint32_t(end), nil)
	C.free(unsafe.Pointer(cs))
	r := &Result{gqr: gqr}
	runtime.SetFinalizer(r, destroy_query_result)
	return r
}

func destroy_query_result(r *Result) {
	C.giggle_query_result_destroy(&r.gqr)
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
