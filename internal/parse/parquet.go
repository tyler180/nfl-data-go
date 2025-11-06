package parse

import (
	"fmt"
	"io"

	"github.com/tyler180/nfl-data-go/internal/schema"
)

// SnapCountsParquet parses a Parquet stream into []nflreadgo.SnapCount.
// TODO: implement with a parquet reader (e.g., github.com/apache/arrow/go/arrow/parquet)
func SnapCountsParquet(_ io.Reader) ([]schema.SnapCount, error) {
	return nil, fmt.Errorf("parquet parsing not implemented yet (use CSV)")
}
