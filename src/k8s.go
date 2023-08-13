package src

import (
	"fmt"
	"github.com/Drakorgaur/jetono-api/src/storage"
	"github.com/labstack/echo/v4"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

type createSecretForm struct {
	Operator   string `json:"operator"`
	Account    string `json:"account"`
	User       string `json:"user"`
	SecretName string `json:"secret_name"`
	Namespace  string `json:"namespace"`
}

//	@Tags		Secret
//	@Router		/secret [post]
//	@Param		json	body	createSecretForm	true	"json body"
//	@Summary	Creates secret with credentials
//	@Success	200	{object}	SecretResponse	"200 ok"
//	@Failure	500	{object}	string			"Internal error"
func createSecret(c echo.Context) error {
	s := &createSecretForm{
		Namespace:  "default",
		SecretName: "",
	}
	if err := c.Bind(&s); err != nil {
		return badRequest(c, err)
	}

	creds, err := GetUserCreds(s.Operator, s.Account, s.User)

	if s.SecretName == "" {
		s.SecretName = fmt.Sprintf("%s.%s.%s.creds", storage.Slugify(s.Operator), storage.Slugify(s.Account), storage.Slugify(s.User))
	}

	if err != nil {
		return badRequest(c, err)
	}

	secret, err := createSecretWithCredentials(s.SecretName, s.Namespace, map[string][]byte{"creds": creds})
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

func createSecretWithCredentials(secretName string, ns string, data map[string][]byte) (*v1.Secret, error) {
	kube, ctx, err := storage.InitKubeWithCtx()
	if err != nil {
		return nil, err
	}

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
