package service

import (
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

func (p PersonServiceImpl) Add(person model.Person) error {
	//TODO: transactional control not implemented yet
	//When add relationship operations fail, it should be able to rollback the entire transaction

	//Add new person
	savedPerson, err := p.repository.Save(entity.Person{
		Name: person.Name,
	})

	if err != nil {
		log.Error(err)
		return err
	}
	//Add parent relationship
	if person.Parent != 0 {
		p.repository.AddRelationship(entity.Relationship{
			ParentId: person.Parent,
			ChildId:  savedPerson.Id,
		})
		if err != nil {
			err := fmt.Errorf("could not add parent relationship. %v", err)
			log.Error(err)
			return err
		}
	}
	//Add child relationships
	for _, childId := range person.Children {
		parents, err := p.repository.FindParents(childId)
		//Validating number of parents, a child should not have more than 2 parents
		if err == nil && len(parents) == 2 {
			return fmt.Errorf("could not add child relationship. child %d already has 2 parents", childId)
		}
		relationship := entity.Relationship{
			ParentId: savedPerson.Id,
			ChildId:  childId,
		}
		err = p.repository.AddRelationship(relationship)
		if err != nil {
			err := fmt.Errorf("could not add child relationship. %v", err)
			log.Error(err)
			return err
		}
	}

	return nil
}

func (p PersonServiceImpl) FindAscendantsById(id int64) ([]*model.Ascendancy, error) {
	ascendants, err := p.repository.FindAscendantsById(id)
	log.Error(err)
	return ascendants, err
}

func (p PersonServiceImpl) FindAll() ([]*entity.Person, error) {
	people, err := p.repository.FindBAll()
	log.Error(err)
	return people, err
}
