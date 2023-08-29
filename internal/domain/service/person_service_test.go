package service

import (
	"context"
	"errors"
	"golang-family-tree/internal/domain/model"
	"golang-family-tree/internal/repository/entity"
	"golang-family-tree/mocks"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repository := mocks.NewMockPersonRepository(ctrl)
	ctx := context.Background()

	addPerson := model.Person{
		Name:     "test",
		Parent:   1,
		Children: []int64{2, 3, 4},
	}

	newPerson := entity.Person{
		Name: addPerson.Name,
	}

	savedPerson := entity.Person{
		Id:   5,
		Name: addPerson.Name,
	}

	parentRelationship := entity.Relationship{
		ParentId: addPerson.Parent,
		ChildId:  savedPerson.Id,
	}

	repository.EXPECT().Save(ctx, newPerson).Return(&savedPerson, nil)

	repository.EXPECT().AddRelationship(ctx, parentRelationship)

	for _, childId := range addPerson.Children {
		repository.EXPECT().FindParents(ctx, childId)
		childRelationship := entity.Relationship{
			ParentId: savedPerson.Id,
			ChildId:  childId,
		}
		repository.EXPECT().AddRelationship(ctx, childRelationship)
	}

	service := NewPersonService(repository)
	err := service.Add(ctx, addPerson)
	assert.Nil(t, err)
}

func TestAddWhenSavePersonFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repository := mocks.NewMockPersonRepository(ctrl)
	ctx := context.Background()

	addPerson := model.Person{
		Name:     "test",
		Parent:   1,
		Children: []int64{2, 3, 4},
	}

	newPerson := entity.Person{
		Name: addPerson.Name,
	}

	repository.EXPECT().Save(ctx, newPerson).Return(nil, errors.New("unable to save person"))

	repository.EXPECT().AddRelationship(ctx, gomock.Any()).Times(0)
	repository.EXPECT().FindParents(ctx, gomock.Any()).Times(0)
	repository.EXPECT().AddRelationship(ctx, gomock.Any()).Times(0)

	service := NewPersonService(repository)
	err := service.Add(ctx, addPerson)
	assert.Error(t, err)
}

func TestAddWhenAddParentRelationshipFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repository := mocks.NewMockPersonRepository(ctrl)
	ctx := context.Background()

	addPerson := model.Person{
		Name:     "test",
		Parent:   1,
		Children: []int64{2, 3, 4},
	}

	newPerson := entity.Person{
		Name: addPerson.Name,
	}

	savedPerson := entity.Person{
		Id:   5,
		Name: addPerson.Name,
	}

	parentRelationship := entity.Relationship{
		ParentId: addPerson.Parent,
		ChildId:  savedPerson.Id,
	}

	repository.EXPECT().Save(ctx, newPerson).Return(&savedPerson, nil)

	repository.EXPECT().AddRelationship(ctx, parentRelationship).Return(errors.New("unable to add parent relationship"))

	for _, childId := range addPerson.Children {
		repository.EXPECT().FindParents(ctx, childId).Times(0)
		childRelationship := entity.Relationship{
			ParentId: savedPerson.Id,
			ChildId:  childId,
		}
		repository.EXPECT().AddRelationship(ctx, childRelationship).Times(0)
	}

	service := NewPersonService(repository)
	err := service.Add(ctx, addPerson)
	assert.Error(t, err)
}

func TestAddWhenAddChildRelationshipFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repository := mocks.NewMockPersonRepository(ctrl)
	ctx := context.Background()

	childId := int64(2)
	addPerson := model.Person{
		Name:     "test",
		Parent:   1,
		Children: []int64{childId},
	}

	newPerson := entity.Person{
		Name: addPerson.Name,
	}

	savedPerson := entity.Person{
		Id:   5,
		Name: addPerson.Name,
	}

	parentRelationship := entity.Relationship{
		ParentId: addPerson.Parent,
		ChildId:  savedPerson.Id,
	}

	repository.EXPECT().Save(ctx, newPerson).Return(&savedPerson, nil)

	repository.EXPECT().AddRelationship(ctx, parentRelationship)

	childRelationship := entity.Relationship{
		ParentId: savedPerson.Id,
		ChildId:  childId,
	}
	repository.EXPECT().AddRelationship(ctx, childRelationship).Return(errors.New("unable to add child relationship"))

	repository.EXPECT().FindParents(ctx, childId)

	service := NewPersonService(repository)
	err := service.Add(ctx, addPerson)
	assert.Error(t, err)
}

func TestAddWhenAddChildHasTwoParentsAlready(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repository := mocks.NewMockPersonRepository(ctrl)
	ctx := context.Background()

	childId := int64(2)
	addPerson := model.Person{
		Name:     "test",
		Parent:   1,
		Children: []int64{childId},
	}

	newPerson := entity.Person{
		Name: addPerson.Name,
	}

	savedPerson := entity.Person{
		Id:   5,
		Name: addPerson.Name,
	}

	parentRelationship := entity.Relationship{
		ParentId: addPerson.Parent,
		ChildId:  savedPerson.Id,
	}

	repository.EXPECT().Save(ctx, newPerson).Return(&savedPerson, nil)

	repository.EXPECT().AddRelationship(ctx, parentRelationship)

	childRelationship := entity.Relationship{
		ParentId: savedPerson.Id,
		ChildId:  childId,
	}
	repository.EXPECT().AddRelationship(ctx, childRelationship).Times(0)

	parents := []*model.Parent{
		{Id: savedPerson.Id,
			Name: "parent"},
		{Id: savedPerson.Id,
			Name: "parent"},
	}
	repository.EXPECT().FindParents(ctx, childId).Return(parents, nil)

	service := NewPersonService(repository)
	err := service.Add(ctx, addPerson)
	assert.Error(t, err)
}
