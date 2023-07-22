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
}

func InitK8sWithCtx() (*kubernetes.Clientset, context.Context, error) {
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

func (s *KubernetesStore) Store(usm *AccountServerMap) error {
	kube, ctx, err := InitK8sWithCtx()
	if err != nil {
		return err
	}

	name := GenerateCMName(usm.Operator, usm.Account)

	_, err = kube.CoreV1().ConfigMaps("default").Create(ctx, &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Data: map[string]string{
			"servers": usm.Servers,
		},
	}, metav1.CreateOptions{})

	if err != nil {
		return err
	}

	return nil
}

func (s *KubernetesStore) Read(usm *AccountServerMap) error {
	name := GenerateCMName(usm.Operator, usm.Account)

	kube, ctx, err := InitK8sWithCtx()
	if err != nil {
		return err
	}

	cm, err := kube.CoreV1().ConfigMaps("default").Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	usm.Servers = cm.Data["servers"]

	return nil
}
