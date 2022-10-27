package entityRepo

type GoenInheritType int

type BasicEntity struct {
	Entity
	GoenInheritType GoenInheritType `db:"goen_inherit_type"`
}

func (p *BasicEntity) inheritAfterNew(goenId int, inheritType GoenInheritType) {
	p.afterNew(goenId)
	p.GoenInheritType = inheritType
	p.AddBasicFieldChange("goen_inherit_type")
}

func (p *BasicEntity) GetParentEntity() EntityForInheritRepo {
	return nil
}

func (p *BasicEntity) GetRealType() GoenInheritType {
	return p.GoenInheritType
}
