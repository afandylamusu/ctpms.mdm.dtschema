package author

import (
	"context"

	"github.com/afandylamusu/ctpms.mdm.dtschema/models"
)

// Repository represent the author's repository contract
type Repository interface {
	GetByID(ctx context.Context, id int64) (*models.Author, error)
}
