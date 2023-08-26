package storage

import (
	"errors"
	"github.com/objectbox/objectbox-go/objectbox"
	"os"
)

type ObjectBoxStore struct{}

type ErrNotFound struct{}

func (e ErrNotFound) Error() string { return "object(s) was not found" }

var errNotFound ErrNotFound

func initObjectBox() (*objectbox.ObjectBox, error) {
	builder := objectbox.NewBuilder().Model(ObjectBoxModel())
	volume := os.Getenv("OBJECTBOX_VOL")
	if volume == "" {
		volume = "./objectbox"
	}
	builder.Directory(volume)
	objectBox, err := builder.Build()
	if err != nil {
		return nil, err
	}

	return objectBox, nil
}

func (s *ObjectBoxStore) StoreCtx(asm *AccountServerMap) error {
	ob, err := initObjectBox()
	if err != nil {
		return err
	}
	defer ob.Close() // In a server app, you would just keep ob and close on shutdown

	box := BoxForAccountServerMapB(ob)

	asmb, err := getAccountServerMapB(box, asm)
	if errors.As(err, &errNotFound) {
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

func (s *ObjectBoxStore) ReadCtx(asm *AccountServerMap) error {
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

	asm.Server = asmb.Server

	return nil
}
