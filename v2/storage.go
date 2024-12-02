package bf

type Storage interface {
	NewBlock(capacity uint32)

	Count() int

	Block(index int) StorageBlock

	Equals(other Storage) bool
}

type StorageBlock interface {
	Set(index uint32)

	Clear(index uint32)

	Get(index uint32) bool

	Capacity() uint32

	Equals(other StorageBlock) bool
}

type BatchIntersect interface {
	Intersect(other StorageBlock)
}

type BatchUnion interface {
	Union(other StorageBlock)
}
