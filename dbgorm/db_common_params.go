package dbgorm

type PagingParams struct {
	Offset    int
	Limit     int
	NeedTotal bool
}

type PagingResult struct {
	Offset int
	Limit  int
	Total  int64
}
