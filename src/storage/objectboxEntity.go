package storage

import "fmt"

//go:generate go run github.com/objectbox/objectbox-go/cmd/objectbox-gogen
type AccountServerMapB struct {
	AccountServerMap
	Id  uint64 `objectbox:"id"`
	Uid string `objectbox:"index:hash"`
}

func GenerateUid(op string, acc string) string {
	return fmt.Sprintf("%s.%s", op, acc)
}

func (usmb *AccountServerMapB) generateUid() {
	usmb.Uid = GenerateUid(usmb.Operator, usmb.Account)
}
