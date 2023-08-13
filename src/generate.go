package src

import (
	"fmt"
	"github.com/Drakorgaur/jetono-api/src/storage"
	"github.com/labstack/echo/v4"
	nsc "github.com/nats-io/nsc/cmd"
	"github.com/nats-io/nsc/cmd/store"
	"github.com/spf13/cobra"
)

func init() {
	GetEchoRoot().GET("/generate/config", generateConfig)
	GetEchoRoot().POST("/generate/config", storeConfig)
}

// GenerateConfig now supports only `--nats-resolver` configuration
func GenerateConfig(operatorName string) ([]byte, error) {
	ctx, err := nsc.NewActx(&cobra.Command{}, []string{})
	if err != nil {
		return nil, err
	}
	s := ctx.StoreCtx().Store

	op, err := s.Read(store.JwtName(operatorName))
	opClaim, err := ctx.StoreCtx().Store.ReadOperatorClaim()
	if err != nil {
		return nil, err
	}

	generator := nsc.NewNatsResolverConfigBuilder(false)
	err = generator.Add(op)
	if err != nil {
		return nil, err
	}

	err = generator.SetSystemAccount(opClaim.SystemAccount)
	if err != nil {
		return nil, err
	}

	names, err := nsc.GetConfig().ListAccounts()
	if err != nil {
		return nil, err
	}

	if len(names) == 0 {
		return nil, fmt.Errorf("operator %q has no accounts", nsc.GetConfig().Operator)
	}

	for _, n := range names {
		d, err := s.Read(store.Accounts, n, store.JwtName(n))
		if err != nil {
			return nil, err
		}
		err = generator.Add(d)
		if err != nil {
			return nil, err
		}

		users, err := s.ListEntries(store.Accounts, n, store.Users)
		if err != nil {
			return nil, err
		}
		for _, u := range users {
			d, err := s.Read(store.Accounts, n, store.Users, store.JwtName(u))
			if err != nil {
				return nil, err
			}
			err = generator.Add(d)
			if err != nil {
				return nil, err
			}
		}
	}

	d, err := generator.Generate()
	if err != nil {
		return nil, err
	}
	return d, nil
}

//	@Tags		Generate
//	@Router		/generate/config [get]
//	@Param		operator	query	string	true	"Operator name"
//	@Summary	Sends configuration for nats server with resolver as this operator
//	@Success	200	{object}	string	"text/plain config file"
//	@Failure	500	{object}	string	"Internal error"
func generateConfig(c echo.Context) error {

	err, s := raiseForRequiredFlags(c.QueryParam, "operator")
	if err != nil {
		return badRequest(c, fmt.Errorf("required flag %s not set", s))
	}

	operator := c.QueryParam("operator")

	if err := nsc.GetConfig().SetOperator(operator); err != nil {
		return badRequest(c, err)
	}
	fmt.Println(nsc.GetConfig().Operator)
	config, err := GenerateConfig(operator)
	if err != nil {
		return badRequest(c, err)
	}
	return c.Blob(200, "text/plain", config)
}

type storeConfigForm struct {
	Operator string `json:"operator"`
	Name     string `json:"name"`
}

//	@Tags		Generate
//	@Router		/generate/config [post]
//	@Param		json	body	storeConfigForm	true	"json body"
//	@Summary	Stores configuration for nats server with resolver as this operator in kubernetes secret
//	@Success	200	{object}	SimpleJSONResponse	"200 ok"
//	@Failure	500	{object}	string				"Internal error"
func storeConfig(c echo.Context) error {
	s := storeConfigForm{}
	if err := c.Bind(&s); err != nil {
		return err
	}
	if err := nsc.GetConfig().SetOperator(s.Operator); err != nil {
		return badRequest(c, err)
	}

	storeT, err := storage.StoreType()
	if err != nil {
		return badRequest(c, err)
	}

	if kubeStore, ok := storeT.(*storage.KubernetesStore); ok {
		config, err := GenerateConfig(s.Operator)
		if err != nil {
			return badRequest(c, err)
		}
		err = kubeStore.StoreSecret(
			s.Name,
			map[string][]byte{
				".config": config,
			},
			kubeStore.ConfigNs,
		)
		if err != nil {
			return badRequest(c, err)
		}
		return c.JSON(200, SimpleJSONResponse{
			Status:  "200",
			Message: "config stored",
		})
	}
	return badRequest(c, fmt.Errorf("storage type %kubeStore does not support config storage", storeT))
}
