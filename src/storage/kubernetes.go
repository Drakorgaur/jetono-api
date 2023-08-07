package storage

import (
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"strings"
)

type KubernetesStore struct {
	CtxNs    string
	ConfigNs string
}

func NewKubernetesStore(ctxNs string, configNs string) *KubernetesStore {
	if ctxNs == "" {
		ctxNs = "default"
	}
	if configNs == "" {
		configNs = "default"
	}

	return &KubernetesStore{
		CtxNs:    ctxNs,
		ConfigNs: configNs,
	}
}

func InitKubeWithCtx() (*kubernetes.Clientset, context.Context, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, nil, err
	}
	kube, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, err
	}
	ctx := context.TODO()
	return kube, ctx, nil
}

func Slugify(s string) string {
	// slugify operator name for k8s secret name
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "_", "-")
	return s
}

func GenerateCMName(op string, acc string) string {
	return fmt.Sprintf("jetono.%s.%s", Slugify(op), Slugify(acc))
}

func (s *KubernetesStore) StoreCtx(usm *AccountServerMap) error {
	return s.StoreCM(
		GenerateCMName(usm.Operator, usm.Account),
		map[string]string{
			"servers": usm.ServersList,
		},
		s.CtxNs,
	)
}

// StoreCM TODO: store/update
func (s *KubernetesStore) StoreCM(name string, data map[string]string, ns string) error {
	kube, ctx, err := InitKubeWithCtx()
	if err != nil {
		return err
	}

	_, err = kube.CoreV1().ConfigMaps(ns).Create(ctx, &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Data: data,
	}, metav1.CreateOptions{})

	return err
}

// StoreSecret TODO: store/update
func (s *KubernetesStore) StoreSecret(name string, data map[string][]byte, ns string) error {
	kube, ctx, err := InitKubeWithCtx()
	if err != nil {
		return err
	}

	_, err = kube.CoreV1().Secrets(ns).Create(ctx, &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Data: data,
	}, metav1.CreateOptions{})

	return err
}

func (s *KubernetesStore) ReadCtx(usm *AccountServerMap) error {
	name := GenerateCMName(usm.Operator, usm.Account)

	kube, ctx, err := InitKubeWithCtx()
	if err != nil {
		return err
	}

	cm, err := kube.CoreV1().ConfigMaps("default").Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	usm.ServersList = cm.Data["servers"]

	return nil
}
