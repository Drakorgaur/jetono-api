package src

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"strings"
)

func init() {
	GetEchoRoot().POST("secret", createSecret)
}

type SecretInfo struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

type SecretResponse struct {
	Message string     `json:"message"`
	Secret  SecretInfo `json:"secret"`
}

func createSecret(c echo.Context) error {
	if err, s := raiseForRequiredFlags(c.FormValue, "operator", "account", "user"); err != nil {
		return c.JSON(400, map[string]string{"error": err.Error(), "field required": s})
	}

	ns := c.FormValue("namespace")
	if ns == "" {
		ns = "default"
	}

	secretName := c.FormValue("secret_name")

	operator := c.FormValue("operator")
	account := c.FormValue("account")
	user := c.FormValue("user")
	creds, err := GetUserCreds(operator, account, user)

	if secretName == "" {
		secretName = fmt.Sprintf("%s.%s.%s.creds", slugify(operator), slugify(account), slugify(user))
	}

	if err != nil {
		return badRequest(c, err)
	}

	secret, err := createSecretWithCredentials(secretName, ns, map[string][]byte{"creds": creds})
	if err != nil {
		return badRequest(c, err)
	}

	return c.JSON(
		200,
		SecretResponse{
			Message: "secret was successfully created",
			Secret: SecretInfo{
				Name:      secret.Name,
				Namespace: secret.Namespace,
			},
		},
	)
}

func slugify(s string) string {
	// slugify operator name for k8s secret name
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "_", "-")
	return s
}

func createSecretWithCredentials(secretName string, ns string, data map[string][]byte) (*v1.Secret, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	kube, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	ctx := context.TODO()
	secret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: secretName,
		},
		Data: data,
	}

	if secret, err = kube.CoreV1().Secrets(ns).Create(ctx, secret, metav1.CreateOptions{}); err != nil {
		return nil, err
	}

	return secret, nil
}
