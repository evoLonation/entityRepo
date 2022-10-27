package entityRepo

type saver struct {
	waitBasicSave       []func() error
	waitAssociationSave []func() error
}

var Saver saver

func (p *saver) Save() error {
	for _, foo := range p.waitBasicSave {
		if err := foo(); err != nil {
			return err
		}
	}
	p.waitBasicSave = nil
	for _, foo := range p.waitAssociationSave {
		if err := foo(); err != nil {
			return err
		}
	}
	p.waitAssociationSave = nil
	return nil
}
func (p *saver) addInBasicSaveQueue(foo func() error) {
	p.waitBasicSave = append(p.waitBasicSave, foo)
}
func (p *saver) addInAssSaveQueue(foo func() error) {
	p.waitAssociationSave = append(p.waitAssociationSave, foo)
}
