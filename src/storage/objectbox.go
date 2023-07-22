package storage

import (
	"fmt"
	"github.com/objectbox/objectbox-go/objectbox"
)

type ObjectBoxStore struct{}

func initObjectBox() (*objectbox.ObjectBox, error) {
	objectBox, err := objectbox.NewBuilder().Model(ObjectBoxModel()).Build()
	if err != nil {
		return nil, err
	}

	return objectBox, nil
}

func (s *ObjectBoxStore) Store(usm *AccountServerMap) error {
	ob, err := initObjectBox()
	if err != nil {
		return err
	}
	defer ob.Close() // In a server app, you would just keep ob and close on shutdown

	box := BoxForAccountServerMapB(ob)

	// Create
	_, err = box.Put(&AccountServerMapB{
		AccountServerMap: *usm,
		Uid:              GenerateUid(usm.Operator, usm.Account),
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *ObjectBoxStore) Read(usm *AccountServerMap) error {
	ob, err := initObjectBox()
	if err != nil {
		return err
	}
	defer ob.Close() // In a server app, you would just keep ob and close on shutdown

	box := BoxForAccountServerMapB(ob)

	res, err := box.Query(AccountServerMapB_.Uid.Equals(GenerateUid(
		usm.Operator, usm.Account,
	), true)).Find()
	if err != nil {
		return err
	}

	if len(res) == 0 {
		return fmt.Errorf("no result")
	}

	usmb := res[0]

	usm.Servers = usmb.Servers

	return nil
}
