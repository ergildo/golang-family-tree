package service

import (
	"context"
	"golang-family-tree/internal/domain/model"
	"golang-family-tree/internal/repository/entity"
)

type PersonService interface {
	Add(ctx context.Context, person model.Person) error
	FindAll(ctx context.Context) ([]*entity.Person, error)
	FindAscendantsById(ctx context.Context, id int64) ([]*model.Ascendancy, error)
}
