package literaladapter

import (
	"context"

	"hufschlaeger.net/markscribe/literal"
)

// Adapter implements ports.LiteralPort using the local literal package.
type Adapter struct{}

func New() *Adapter { return &Adapter{} }

func (a *Adapter) CurrentlyReading(ctx context.Context, count int) ([]literal.Book, error) {
	books, err := literal.CurrentlyReading()
	if err != nil {
		return nil, err
	}
	if len(books) > count {
		return books[:count], nil
	}
	return books, nil
}
