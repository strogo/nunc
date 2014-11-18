package nunc

type Meta struct {
	ID    int64  `ql:"index xID"`
	Key   string `ql:"index xKey"`
	Value int64
}

// TODO functions Bool, String, etc.
