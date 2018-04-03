package paginator

import (
	"testing"
)

func BenchmarkPaginator(b *testing.B) {

	b.ReportAllocs()

	for i := 0; i < b.N; i++ {

		NewPaginator(25, 500, 10, "/page/param/", nil)
	}
}
