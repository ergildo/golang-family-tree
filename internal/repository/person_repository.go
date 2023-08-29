package repository

import (
	"context"
	"database/sql"
	"fmt"
	"golang-family-tree/internal/domain/model"
	"golang-family-tree/internal/repository/entity"

	log "github.com/sirupsen/logrus"
)

func NewPersonRepository(db *sql.DB) PersonRepository {
	return PersonRepositoryImpl{db: db}
}

func (p *PersonRepositoryImpl) BeginTx(ctx context.Context) (*PersonRepositoryImpl, error) {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	p.tx = tx
	return p, nil
}

func (p PersonRepositoryImpl) Commit(ctx context.Context) error {
	return p.tx.Commit()
}

type PersonRepositoryImpl struct {
	db *sql.DB
	tx *sql.Tx
}

func (p PersonRepositoryImpl) FindById(ctx context.Context, id int64) (*entity.Person, error) {

	query := "select id, name from person where id =$1"
	stm, err := p.db.PrepareContext(ctx, query)

	if err != nil {
		return nil, fmt.Errorf("could not prepare query: %w", err)
	}
	defer stm.Close()
	row := stm.QueryRowContext(ctx, id)
	if row.Err() != nil {
		return nil, fmt.Errorf("could not prepare query:%w", err)
	}
	var person entity.Person
	row.Scan(&person.Id, &person.Name)
	return &person, nil
}

func (p PersonRepositoryImpl) FindBAll(ctx context.Context) ([]*entity.Person, error) {
	query := "select id, name from person"
	stm, err := p.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error when preparing query: %w", err)
	}
	defer stm.Close()

	rows, err := stm.QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("error when executing query:%w", err)
	}

	var persons []*entity.Person
	for rows.Next() {
		var person entity.Person
		rows.Scan(&person.Id, &person.Name)
		persons = append(persons, &person)
	}

	return persons, nil
}

func (p PersonRepositoryImpl) FindAscendantsById(ctx context.Context, id int64) ([]*model.Ascendancy, error) {
	query := `	WITH RECURSIVE ascendancy (id, name, parentid, depth) AS (
		select p.id, p.name, r.parentid, 0 as depth from person p
			join relationship r 
			on r.childid  = p.id
			where p.id = $1
		
			UNION
			
		select p.id, p.name, r.parentid, depth+1  from ascendancy a
			join person p  
			on a.parentid  = p.id 
			left join relationship r 
			on r.childid = p.id 
		)
	SELECT id, name, depth FROM ascendancy
		group by id, name, depth order by depth;
`
	stm, err := p.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("could not prepare query: %w", err)
	}
	defer stm.Close()

	rows, err := stm.QueryContext(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("could not execute query: %w", err)
	}

	var ascendants []*model.Ascendancy
	for rows.Next() {
		var ascendant model.Ascendancy
		rows.Scan(&ascendant.Id, &ascendant.Name, &ascendant.Depth)
		parents, err := p.FindParents(ctx, ascendant.Id)
		if err == nil && parents != nil {
			ascendant.Parents = parents
		}
		ascendants = append(ascendants, &ascendant)
	}

	return ascendants, nil
}

func (p PersonRepositoryImpl) Save(ctx context.Context, person entity.Person) (*entity.Person, error) {
	sql := "insert into person (name) values ($1) RETURNING id"

	stm, err := p.db.PrepareContext(ctx, sql)

	if err != nil {
		return nil, fmt.Errorf("could not prepare save sql: %w", err)
	}
	var lastId int64
	rs := stm.QueryRowContext(ctx, person.Name)
	if rs.Err() != nil {
		return nil, fmt.Errorf("could not save person: %w", rs.Err())
	}
	rs.Scan(&lastId)
	return p.FindById(ctx, lastId)
}

func (p PersonRepositoryImpl) AddRelationship(ctx context.Context, r entity.Relationship) error {
	sql := "insert into relationship (parentId, childId) values ($1, $2)"

	stm, err := p.db.PrepareContext(ctx, sql)

	if err != nil {
		return fmt.Errorf("could not prepare update sql: %w", err)
	}
	_, err = stm.ExecContext(ctx, r.ParentId, r.ChildId)
	if err != nil {
		return fmt.Errorf("could not save relationship: %w", err)
	}
	return nil
}

func (p PersonRepositoryImpl) FindParents(ctx context.Context, id int64) ([]*model.Parent, error) {

	query := "select p.id, p.name from person p left join relationship r on p.id = r.parentid where r.childid = $1"
	stm, err := p.db.PrepareContext(ctx, query)
	if err != nil {
		err = fmt.Errorf("could not prepare query: %w", err)
		log.Error(err)
		return nil, err
	}
	defer stm.Close()

	rows, err := stm.QueryContext(ctx, id)
	if err != nil {
		err = fmt.Errorf("could not execute query: %w", err)
		log.Error(err)
		return nil, err
	}

	var parents []*model.Parent
	for rows.Next() {
		var parent model.Parent
		rows.Scan(&parent.Id, &parent.Name)
		parents = append(parents, &parent)
	}
	return parents, nil
}
