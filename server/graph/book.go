package graph

import (
	"context"

	"github.com/c5rogers/one-tap/salary-advance-loan-system/models"
)

type uuid string

func (gc *LMSGraphClient) GetBookById(id string) (u *models.Book, err error) {
	var query struct {
		Books []models.Book `graphql:"book(where:{id:{_eq:$id}})"` // get the book by id
	}

	variables := map[string]interface{}{
		"id": uuid(id),
	}

	if err := gc.Query(context.Background(), &query, variables); err != nil {
		return nil, err
	}
	if len(query.Books) > 0 {
		u = &query.Books[0]
	}
	return u, err
}
