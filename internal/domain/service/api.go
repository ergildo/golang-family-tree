package service

import (
	"golang-family-tree/internal/domain/model"
	"golang-family-tree/internal/repository/entity"
)

type PersonService interface {
	Add(person model.Person) error
	FindAll() ([]*entity.Person, error)
	FindAscendantsById(id int64) ([]*model.Ascendancy, error)
}
