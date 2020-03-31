package postgres

import (
	"context"

	"github.com/d7561985/questions/model"
)

//nolint:unused
// postgress is just example how we can handle clear architecture approach
type postgres struct {
}

func (p postgres) TakerList(ctx context.Context) (model.TakerList, error) {
	panic("implement me")
}
