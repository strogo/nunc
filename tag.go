package nunc

type Tag struct {
	ID   int64 `ql:"index xID"`
	Name string
}

type Tagged struct {
	ID   int64 `ql:"index xID"`
	Tag  int64
	Task int64
}
