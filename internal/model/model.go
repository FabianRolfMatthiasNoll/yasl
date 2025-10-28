package model

type List struct {
	ID   int64
	Name string
}

type Item struct {
	ID       int64
	ListID   int64
	Name     string
	Category string
	Checked  bool
}