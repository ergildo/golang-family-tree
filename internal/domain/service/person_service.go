package service

import (
	"context"
	"fmt"
	"golang-family-tree/internal/domain/model"
	"golang-family-tree/internal/repository"
	"golang-family-tree/internal/repository/entity"

	log "github.com/sirupsen/logrus"
)

func NewPersonService(repository repository.PersonRepository) PersonService {
	return PersonServiceImpl{
		repository: repository,
	}
}

type PersonServiceImpl struct {
	repository repository.PersonRepository
}

func (p PersonServiceImpl) Add(ctx context.Context, person model.Person) error {
	//TODO: transactional control not implemented yet
	//When add relationship operations fail, it should be able to rollback the entire transaction

	//Add new person
	newPerson := entity.Person{
		Name: person.Name,
	}

	savedPerson, err := p.repository.Save(ctx, newPerson)

	if err != nil {
		log.Error(err)
		return err
	}
	//Add parent relationship
	if person.Parent != 0 {
		relationship := entity.Relationship{
			ParentId: person.Parent,
			ChildId:  savedPerson.Id,
		}
		p.repository.AddRelationship(ctx, relationship)
		if err != nil {
			err := fmt.Errorf("could not add parent relationship. %v", err)
			log.Error(err)
			return err
		}
	}
	//Add child relationships
	for _, childId := range person.Children {
		parents, err := p.repository.FindParents(ctx, childId)
		//Validating number of parents, a child should not have more than 2 parents
		if err == nil && len(parents) == 2 {
			if err != nil {
				log.Error(err)
			}
			return fmt.Errorf("could not add child relationship. child %d already has 2 parents", childId)
		}
		relationship := entity.Relationship{
			ParentId: savedPerson.Id,
			ChildId:  childId,
		}
		err = p.repository.AddRelationship(ctx, relationship)
		if err != nil {
			err := fmt.Errorf("could not add child relationship. %v", err)
			log.Error(err)
			return err
		}
	}
	return nil
}

func (p PersonServiceImpl) FindAscendantsById(ctx context.Context, id int64) ([]*model.Ascendancy, error) {
	ascendants, err := p.repository.FindAscendantsById(ctx, id)
	log.Error(err)
	return ascendants, err
}

func (p PersonServiceImpl) FindAll(ctx context.Context) ([]*entity.Person, error) {
	people, err := p.repository.FindBAll(ctx)
	log.Error(err)
	return people, err
}
