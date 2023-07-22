package storage

type Storage interface {
	Store(usm *AccountServerMap) error
	Read(usm *AccountServerMap) error
}

type AccountServerMap struct {
	Operator string `objectbox:"-" json:"operator"`
	Account  string `objectbox:"-" json:"account"`
	Servers  string `json:"servers"`
}
