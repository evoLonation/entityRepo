package entityRepo

type InheritRepoForOther[T any] interface {
	RepoForOther[T]
	CastFrom(any) (T, error)
	GetRealType(T) GoenInheritType
}

type RepoForOther[T any] interface {
	New() T
	// 如果没有找到，则不必报错而是返回空值
	GetFromAllInstanceBy(member string, value any) T
	FindFromAllInstanceBy(member string, value any) []T
	AddInAllInstance(e T)
	RemoveFromAllInstance(e T)
	IsInAllInstance(e T) bool
	GetAll() []T
}

type RepoForEntity[PT any] interface {
	// Get 实际上不需要检查是否在allinstance里面
	Get(goenId int) (PT, error)
	FindFromMultiAssTable(tableName string, ownerId int) ([]PT, error)
	GetGoenId(PT) int
}
