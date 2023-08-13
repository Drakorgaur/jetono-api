package storage

import (
	"fmt"
	"os"
)

type Storage interface {
	StoreCtx(usm *AccountServerMap) error
	ReadCtx(usm *AccountServerMap) error
}

type AccountServerMap struct {
	Operator    string `objectbox:"-" json:"operator"`
	Account     string `objectbox:"-" json:"account"`
	ServersList string `json:"servers"`
}

func StoreType() (Storage, error) {
	storeType := os.Getenv("JETONO_STORE_TYPE")
	switch storeType {
	case "kubernetes":
		return NewKubernetesStore(os.Getenv("JETONO_K8S_CTX_NS"), os.Getenv("JETONO_K8S_CONFIG_NS")), nil
	case "objectbox":
		return &ObjectBoxStore{}, nil
	default:
		return nil, fmt.Errorf("invalid store type: %s", storeType)
	}
}

func FillAccCtxFromStorage(accCtx *AccountServerMap) error {
	strg, err := StoreType()
	if err != nil {
		return err
	}

	if err := strg.ReadCtx(accCtx); err != nil {
		return err
	}
	return nil
}
