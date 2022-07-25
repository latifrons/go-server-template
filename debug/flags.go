package debug

type Flags struct {
	ReturnDetailError bool
	DbLog             bool
	RequestLog        bool
	ResponseLog       bool
	RpcLog            bool
	GinDebug          bool
	LogLevel          string
	SkipVerifyAdmin   bool
	SkipEmail         bool
	Swagger           bool
}
