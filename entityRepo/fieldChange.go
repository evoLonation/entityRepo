package entityRepo

type MultiAssChangeType int

const (
	Include MultiAssChangeType = iota
	Exclude
)

type multiAssInfo struct {
	typ       MultiAssChangeType
	tableName string
	targetId  int
}

type FieldChange struct {
	entityStatus
	basicFieldChange []string
	assFieldChange   []string
	multiAssChange   []multiAssInfo
}

func (p *FieldChange) getEntityStatus() entityStatus {
	return p.entityStatus
}

func (p *FieldChange) getBasicFieldChange() []string {
	return p.basicFieldChange
}

func (p *FieldChange) getAssFieldChange() []string {
	return p.assFieldChange
}

func (p *FieldChange) getMultiAssChange() []multiAssInfo {
	return p.multiAssChange
}

func (p *FieldChange) setCreated() {
	p.entityStatus = Created
}
func (p *FieldChange) setExistent() {
	p.entityStatus = Existent
}
func (p *FieldChange) afterBasicSave() {
	p.entityStatus = Existent
	p.basicFieldChange = nil
}
func (p *FieldChange) afterAssSave() {
	p.entityStatus = Existent
	p.assFieldChange = nil
	p.multiAssChange = nil
}

func (p *FieldChange) AddBasicFieldChange(field string) {
	p.basicFieldChange = append(p.basicFieldChange, field)
}

func (p *FieldChange) AddAssFieldChange(field string) {
	p.assFieldChange = append(p.assFieldChange, field)
}

func (p *FieldChange) AddMultiAssChange(typ MultiAssChangeType, tableName string, targetId int) {
	p.multiAssChange = append(p.multiAssChange, multiAssInfo{
		typ:       typ,
		targetId:  targetId,
		tableName: tableName,
	})
}
