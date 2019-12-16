package organization

import (
	"context"

	"github.com/afandylamusu/ctpms.mdm.dtschema/models"
)

// Repository represent the article's repository contract
type Repository interface {
	Fetch(ctx context.Context, cursor string, num int64) (res []*models.Organization, nextCursor string, err error)
	GetByID(ctx context.Context, id int64) (*models.Organization, error)
	GetByTitle(ctx context.Context, title string) (*models.Organization, error)
	Update(ctx context.Context, ar *models.Organization) error
	Store(ctx context.Context, a *models.Organization) error
	Delete(ctx context.Context, id int64) error
}
