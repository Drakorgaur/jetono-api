package storage

import (
	"github.com/objectbox/objectbox-go/objectbox"
)

type ObjectBoxStore struct{}

type ErrNotFound struct{}

func (e ErrNotFound) Error() string { return "object(s) was not found" }

func initObjectBox() (*objectbox.ObjectBox, error) {
	objectBox, err := objectbox.NewBuilder().Model(ObjectBoxModel()).Build()
	if err != nil {
		return nil, err
	}

	return objectBox, nil
}

func (s *ObjectBoxStore) Store(asm *AccountServerMap) error {
	ob, err := initObjectBox()
	if err != nil {
		return err
	}
	defer ob.Close() // In a server app, you would just keep ob and close on shutdown

	box := BoxForAccountServerMapB(ob)

	asmb, err := getAccountServerMapB(box, asm)
	if _, ok := err.(ErrNotFound); ok {
		_, err = box.Put(&AccountServerMapB{
			AccountServerMap: *asm,
			Uid:              GenerateUid(asm.Operator, asm.Account),
		})
		if err != nil {
			return err
		}
	} else {
		asmb.AccountServerMap = *asm
		err = box.Update(asmb)
		if err != nil {
			return err
		}
	}

	return nil
}

func getAccountServerMapB(box *AccountServerMapBBox, asm *AccountServerMap) (*AccountServerMapB, error) {
	res, err := box.Query(AccountServerMapB_.Uid.Equals(GenerateUid(
		asm.Operator, asm.Account,
	), true)).Find()
	if err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, ErrNotFound{}
	}

	return res[0], nil
}

func (s *ObjectBoxStore) Read(asm *AccountServerMap) error {
	ob, err := initObjectBox()
	if err != nil {
		return err
	}
	defer ob.Close() // In a server app, you would just keep ob and close on shutdown

	box := BoxForAccountServerMapB(ob)

	asmb, err := getAccountServerMapB(box, asm)
	if err != nil {
		return err
	}

	asm.Servers = asmb.Servers

	return nil
}
