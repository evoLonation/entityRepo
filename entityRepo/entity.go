package entityRepo

type entityStatus int

const (
	Created entityStatus = iota
	Existent
)

type Entity struct {
	FieldChange

	GoenId            int  `db:"goen_id"`
	GoenInAllInstance bool `db:"goen_in_all_instance"`
}

//剩下的after都继承了FieldChange
func (p *Entity) afterNew(goenId int) {
	p.FieldChange.setCreated()
	p.GoenId = goenId
}
func (p *Entity) afterFind() {
	p.FieldChange.setExistent()
}

func (p *Entity) GetGoenId() int {
	return p.GoenId
}

func (p *Entity) setGoenInAllInstance(goenInAllInstance bool) {
	p.GoenInAllInstance = goenInAllInstance
	p.AddBasicFieldChange("goen_in_all_instance")
}
