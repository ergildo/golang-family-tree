package repository

import (
	"context"
	"golang-family-tree/internal/domain/model"
	"golang-family-tree/internal/repository/entity"
)

type PersonRepository interface {
	Save(ctx context.Context, person entity.Person) (*entity.Person, error)
	AddRelationship(ctx context.Context, person entity.Relationship) error
	FindById(ctx context.Context, id int64) (*entity.Person, error)
	FindBAll(ctx context.Context) ([]*entity.Person, error)
	FindAscendantsById(ctx context.Context, id int64) ([]*model.Ascendancy, error)
	FindParents(ctx context.Context, id int64) ([]*model.Parent, error)
}
