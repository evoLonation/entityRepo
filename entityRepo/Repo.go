package entityRepo

import (
	"database/sql"
	"errors"
	"fmt"
)

type EntityForRepo interface {
	afterNew(int)
	afterFind()
	afterBasicSave()
	afterAssSave()
	getEntityStatus() entityStatus
	setGoenInAllInstance(bool)
	getBasicFieldChange() []string
	getAssFieldChange() []string
	getMultiAssChange() []multiAssInfo
	GetGoenId() int
}

//type repoTypeParam[T any] interface {
//	*T
//	EntityForRepo
//}

type repo[T any, PT any] struct {
	tableName string
	maxGoenId int
}

func NewRepo[T any, PT any](tableName string) (*repo[T, PT], error) {
	_, ok := (any(new(T))).(PT)
	if !ok {
		return nil, errors.New("the type value T does not implement PT ")
	}
	_, ok = (any(new(T))).(EntityForRepo)
	if !ok {
		return nil, errors.New("the type value T does not implement EntityForRepo ")
	}
	repo := &repo[T, PT]{}
	repo.tableName = tableName
	query := fmt.Sprintf("select goen_id from %s order by goen_id DESC limit 1", repo.tableName)
	err := Db.Get(&repo.maxGoenId, query)
	if err != nil && err != sql.ErrNoRows {
		return repo, err
	}
	return repo, nil
}

func (p *repo[T, PT]) getInterface(e any) EntityForRepo {
	return e.(EntityForRepo)
}

func (p *repo[T, PT]) getPT(ei any) PT {
	return ei.(PT)
}

func (p *repo[T, PT]) GetGoenId(e PT) int {
	return (any(e)).(EntityForRepo).GetGoenId()
}

func (p *repo[T, PT]) generateGoenId() int {
	p.maxGoenId = p.maxGoenId + 1
	return p.maxGoenId
}

func (p *repo[T, PT]) New() PT {
	e := p.getInterface(new(T))
	e.afterNew(p.generateGoenId())
	p.addInQueue(e)
	return p.getPT(e)
}

// Get if no rows, return nil, nil
func (p *repo[T, PT]) Get(goenId int) (PT, error) {
	e := p.getInterface(new(T))
	//query := fmt.Sprintf("select * from %s where goen_id=? and goen_in_all_instance = true", p.tableName)
	query := fmt.Sprintf("select * from %s where goen_id=?", p.tableName)
	err := Db.Get(e, query, goenId)
	if err != nil {
		var nilPT PT
		if err == sql.ErrNoRows {
			return nilPT, nil
		}
		return nilPT, err
	}
	e.afterFind()
	p.addInQueue(e)
	return p.getPT(e), nil
}

func (p *repo[T, PT]) GetAll() []PT {
	var entityArr []*T
	var interfaceArr []PT
	query := fmt.Sprintf("select * from %s where goen_in_all_instance = true", p.tableName)
	err := Db.Select(&entityArr, query)
	if err != nil {
		panic(err)
		return nil
	}
	for _, e := range entityArr {
		ei := p.getInterface(e)
		ei.afterFind()
		p.addInQueue(ei)
		interfaceArr = append(interfaceArr, p.getPT(ei))
	}
	return interfaceArr
}

func (p *repo[T, PT]) IsInAllInstance(entity PT) bool {
	ret, err := p.Get(p.getInterface(entity).GetGoenId())
	if err != nil {
		return false
	}
	// must cast to ant,otherwise ret is not comparable
	if any(ret) == nil {
		return false
	}
	return true
}

func (p *repo[T, PT]) GetFromAllInstanceBy(member string, value any) PT {
	e := p.getInterface(new(T))
	query := fmt.Sprintf("select * from %s where %s=? and goen_in_all_instance = true", p.tableName, member)
	err := Db.Get(e, query, value)
	if err != nil {
		var nilPT PT
		if err == sql.ErrNoRows {
			return nilPT
		}
		panic(err)
		return nilPT
	}
	e.afterFind()
	p.addInQueue(e)
	return p.getPT(e)
}

func (p *repo[T, PT]) FindFromAllInstanceBy(member string, value any) []PT {
	var entityArr []*T
	var interfaceArr []PT
	query := fmt.Sprintf("select * from %s where %s=? and goen_in_all_instance = true", p.tableName, member)
	err := Db.Select(&entityArr, query, value)
	if err != nil {
		panic(err)
		return nil
	}
	for _, e := range entityArr {
		ei := p.getInterface(e)
		ei.afterFind()
		p.addInQueue(ei)
		interfaceArr = append(interfaceArr, p.getPT(ei))
	}
	return interfaceArr
}

func (p *repo[T, PT]) AddInAllInstance(e PT) {
	(any(e)).(EntityForRepo).setGoenInAllInstance(true)
}

func (p *repo[T, PT]) RemoveFromAllInstance(e PT) {
	(any(e)).(EntityForRepo).setGoenInAllInstance(false)
}

func (p *repo[T, PT]) FindFromMultiAssTable(tableName string, ownerId int) ([]PT, error) {
	var entityArr []*T
	var interfaceArr []PT
	query := fmt.Sprintf("select tmp.* from %s as ass, %s as tmp where ass.owner_goen_id = ? and ass.possession_goen_id = tmp.goen_id ",
		tableName, p.tableName)
	if err := Db.Select(&entityArr, query, ownerId); err != nil {
		return nil, err
	}
	for _, e := range entityArr {
		ei := p.getInterface(e)
		ei.afterFind()
		p.addInQueue(ei)
		interfaceArr = append(interfaceArr, p.getPT(ei))
	}
	return interfaceArr, nil
}

func (p *repo[T, PT]) addInQueue(e EntityForRepo) {
	Saver.addInBasicSaveQueue(func() error {
		return p.saveBasic(e)
	})
	Saver.addInAssSaveQueue(func() error {
		return p.saveAss(e)
	})
}

// the length of changedField must > 0
func (p *repo[T, PT]) getUpdateQuery(changedField []string) string {
	query := fmt.Sprintf("update %s set ", p.tableName)
	for _, field := range changedField[0 : len(changedField)-1] {
		query += fmt.Sprintf("%s= :%s,", field, field)
	}
	field := changedField[len(changedField)-1]
	query += fmt.Sprintf("%s = :%s", field, field)
	query += fmt.Sprintf(" where goen_id = :goen_id")
	print(query)
	return query
}

// the length of changedField must > 0
func (p *repo[T, PT]) getInsertQuery(changedField []string) string {
	lastField := changedField[len(changedField)-1]
	query := fmt.Sprintf("insert into %s(goen_id, ", p.tableName)
	for _, field := range changedField[0 : len(changedField)-1] {
		query += fmt.Sprintf("%s, ", field)
	}
	query += fmt.Sprintf("%s) values(:goen_id, ", lastField)
	for _, field := range changedField[0 : len(changedField)-1] {
		query += fmt.Sprintf(":%s ,", field)
	}
	query += fmt.Sprintf(":%s )", lastField)
	print(query)
	return query
}

func (p *repo[T, PT]) getMultiAssInsertQuery(tableName string) string {
	query := fmt.Sprintf("insert into %s (owner_goen_id, possession_goen_id) values (?, ?)", tableName)
	return query
}
func (p *repo[T, PT]) getMultiAssDeleteQuery(tableName string) string {
	query := fmt.Sprintf("delete from %s where owner_goen_id=? and possession_goen_id=?", tableName)
	return query
}

func (p *repo[T, PT]) saveBasic(e EntityForRepo) error {
	if len(e.getBasicFieldChange()) != 0 {
		if e.getEntityStatus() == Created {
			if _, err := Db.NamedExec(p.getInsertQuery(e.getBasicFieldChange()), e); err != nil {
				return err
			}
		} else {
			if _, err := Db.NamedExec(p.getUpdateQuery(e.getBasicFieldChange()), e); err != nil {
				return err
			}
		}
	}
	e.afterBasicSave()
	return nil
}

// saveAss e 's entityStatus must be Existence
func (p *repo[T, PT]) saveAss(e EntityForRepo) error {

	if e.getEntityStatus() == Created {
		return errors.New("entityStatus must be Existence")
	}
	if len(e.getAssFieldChange()) != 0 {
		if _, err := Db.NamedExec(p.getUpdateQuery(e.getAssFieldChange()), e); err != nil {
			return err
		}
	}
	for _, info := range e.getMultiAssChange() {
		var query string
		if info.typ == Include {
			query = p.getMultiAssInsertQuery(info.tableName)
		} else {
			query = p.getMultiAssDeleteQuery(info.tableName)
		}
		if _, err := Db.Exec(query, e.GetGoenId(), info.targetId); err != nil {
			return err
		}
	}
	e.afterAssSave()
	return nil
}
