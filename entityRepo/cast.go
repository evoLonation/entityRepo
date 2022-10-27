package entityRepo

type entityForCast interface {
	GetRealType() GoenInheritType
}

func (p *repo[T, PT]) CastFrom(e any) (PT, error) {
	//e.GetRealType()
	ei := (any(e)).(EntityForRepo)
	return p.Get(ei.GetGoenId())
}

func (p *repo[T, PT]) GetRealType(e PT) GoenInheritType {
	return (any(e)).(EntityForInheritRepo).GetRealType()
}
