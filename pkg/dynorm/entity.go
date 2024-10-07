package dynorm

type Entity interface {
	MetaInfo() MetaInfo
}

type MetaInfo struct {
	Table string

	PartitionKey      string
	PartitionkeyValue string

	SortKey      string
	SortKeyValue interface{}

	IndexName     string
	IndexKey      string
	IndexKeyValue string
}
