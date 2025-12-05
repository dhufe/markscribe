package literal

import (
	"context"

	"hufschlaeger.net/markscribe/internal/usecase/ports"
	"hufschlaeger.net/markscribe/literal"
)

// Service wraps the LiteralPort and contains app-level logic for literal.club features.
type Service struct {
	lit ports.LiteralPort
}

func New(lit ports.LiteralPort) *Service { return &Service{lit: lit} }

func (s *Service) CurrentlyReading(count int) []literal.Book {
	books, err := s.lit.CurrentlyReading(context.Background(), count)
	if err != nil {
		panic(err)
	}
	return books
}
