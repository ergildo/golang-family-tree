package repository

import (
	"golang-family-tree/internal/domain/model"
	"golang-family-tree/internal/repository/entity"
)

type PersonRepository interface {
	Save(person entity.Person) (*entity.Person, error)
	AddRelationship(person entity.Relationship) error
	FindById(id int64) (*entity.Person, error)
	FindBAll() ([]*entity.Person, error)
	FindAscendantsById(id int64) ([]*model.Ascendancy, error)
	FindParents(id int64) ([]*model.Parent, error)
}
