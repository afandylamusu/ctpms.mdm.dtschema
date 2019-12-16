package dataset

import (
	"context"

	"github.com/afandylamusu/ctpms.mdm.dtschema/models"
)

// Usecase represent the DataSet's usecases
type Usecase interface {
	Fetch(ctx context.Context, cursor string, num int64) ([]*models.DataSet, string, error)
	GetByID(ctx context.Context, id int64) (*models.DataSet, error)
	Update(ctx context.Context, ar *models.DataSet) error
	GetByTitle(ctx context.Context, title string) (*models.DataSet, error)
	Store(context.Context, *models.DataSet) error
	Delete(ctx context.Context, id int64) error
}
