package nunc

type Context struct {
	ID        int64 `ql:"index xID"`
	Name      string
	ShortName string `ql:"uindex xShortName"`
	Inactive  bool
}
