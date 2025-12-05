package ports

import (
	"context"

	"hufschlaeger.net/markscribe/literal"
)

// LiteralPort defines operations we use from literal.club integration.
type LiteralPort interface {
	CurrentlyReading(ctx context.Context, count int) ([]literal.Book, error)
}
