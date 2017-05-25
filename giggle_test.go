package giggle_test

import (
	"log"
	"testing"

	giggle "github.com/brentp/go-giggle"
)

/*
func TestInit(t *testing.T) {

	g, err := giggle.New("test.idx", "giggle/test/data/many/0.1.bed.gz")
	if err != nil {
		t.Fatal(err)
	}
	if len(g.Files()) != 1 {
		t.Fatal("expected: 1 got: %d", len(g.Files()))
	}

}
*/

func TestMany(t *testing.T) {
	g, err := giggle.New("test.idx", "giggle/test/data/many/*.bed.gz")
	if err != nil {
		t.Fatal(err)
	}

	if len(g.Files()) != 22 {
		t.Fatalf("expected: 22 got: %d", len(g.Files()))
	}

	log.Println("OK")
	res := g.Query("chr1", 12676090, 12676629)
	log.Println(res.NFiles())
	log.Println(res.Hits())
	log.Println(res.TotalHits())
}
