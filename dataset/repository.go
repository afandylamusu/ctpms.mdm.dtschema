package dataset

import (
	"context"

	"github.com/afandylamusu/ctpms.mdm.dtschema/models"
)

// Repository represent the DataSet's repository contract
type Repository interface {
	Fetch(ctx context.Context, cursor string, num int64) (res []*models.DataSet, nextCursor string, err error)
	GetByID(ctx context.Context, id int64) (*models.DataSet, error)
	GetByTitle(ctx context.Context, title string) (*models.DataSet, error)
	Update(ctx context.Context, ar *models.DataSet) error
	Store(ctx context.Context, a *models.DataSet) error
	Delete(ctx context.Context, id int64) error
}
