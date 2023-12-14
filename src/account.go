package src

import (
	"github.com/labstack/echo/v4"
	nsc "github.com/nats-io/nsc/v2/cmd"
)

func init() {
	root := GetEchoRoot()
	root.POST("operator/:operator/account", addAccount)

	root.GET("operator/:operator/accounts", listAccounts)

	root.GET("operator/:operator/account/:name", describeAccount)

	root.PATCH("operator/:operator/account/:name", updateAccount)
}

type addAccountForm struct {
	Name             string `json:"name"`
	PublicKey        string `json:"public_key"`
	ResponseTTL      string `json:"response_ttl,omitempty"`
	AllowPubResponse string `json:"allow_pub_response,omitempty"`
	AllowPub         string `json:"allow_pub,omitempty"`
	AllowPubsub      string `json:"allow_pubsub,omitempty"`
	AllowSub         string `json:"allow_sub,omitempty"`
	DenyPub          string `json:"deny_pub,omitempty"`
	DenyPubsub       string `json:"deny_pubsub,omitempty"`
	DenySub          string `json:"deny_sub,omitempty"`
	Start            string `json:"start,omitempty"`
	Expiry           string `json:"expiry,omitempty"`
}

// @Tags			Account
// @Router			/operator/{operator}/account [post]
// @Summary		Add an account
// @Description	Add an account with given operator to the store
// @Param			operator	path		string				true	"Operator name"
// @Param			json		body		addAccountForm		true	"Account data in json format"
// @Success		200			{object}	SimpleJSONResponse	"Account added"
// @Failure		400			{object}	SimpleJSONResponse	"Bad request"
// @Failure		500			{object}	string				"Internal error"
func addAccount(c echo.Context) error {

	err := runNsc(nil, nil, "select", "operator", c.Param("operator"))
	if err != nil {
		return badRequest(c, err)
	}

	a := addAccountForm{}

	err = runNsc(&a, c, "add", "account")
	if err != nil {
		return badRequest(c, err)
	}

	return c.JSON(200, &SimpleJSONResponse{
		Status:  "200",
		Message: "Account added",
	})
}

// @Tags			Account
// @Router			/operator/{operator}/accounts [get]
// @Summary		List accounts
// @Param			operator	path	string	true	"Operator name"
// @Description	Returns json list of existing accounts for given operator
// @Success		200	{object}	[]string	"Operator's accounts list"
// @Failure		500	{object}	string		"Internal error"
func listAccounts(c echo.Context) error {
	config := nsc.GetConfig()

	// TODO:
	//  {
	//   "status": "400",
	//   "message": "`/nsc/store` is not a valid data directory: stat /nsc/store/.nsc: no such file or directory"
	//  }

	config.Operator = c.Param("operator")

	if accounts, err := config.ListAccounts(); err != nil {
		return badRequest(c, err)
	} else {
		return c.JSON(200, map[string][]string{"accounts": accounts})
	}
}

// @Tags			Account
// @Router			/operator/{operator}/account/{name} [get]
// @Param			name		path	string	true	"Account name"
// @Param			operator	path	string	true	"Operator name"
// @Summary		Describes an account
// @Description	Returns json object with account description
// @Success		200	{object}	AccountDescription	"Operator description"
// @Failure		500	{object}	string				"Internal error"
func describeAccount(c echo.Context) error {
	store, err := nsc.GetStoreForOperator(c.Param("operator"))
	if err != nil {
		return badRequest(c, err)
	}

	claim, err := store.ReadRawAccountClaim(c.Param("name"))
	if err != nil {
		return badRequest(c, err)
	}

	body, err := bodyAsJson(claim)
	if err != nil {
		return badRequest(c, err)
	}

	return c.JSONBlob(200, body)
}

type updateAccountForm struct {
	Tag                string `json:"tag,omitempty"`
	RmTag              string `json:"rm_tag,omitempty"`
	Conns              string `json:"conns,omitempty"`
	LeafConns          string `json:"leaf_conns,omitempty"`
	Data               string `json:"data,omitempty"`
	Exports            string `json:"exports,omitempty"`
	Imports            string `json:"imports,omitempty"`
	Payload            string `json:"payload,omitempty"`
	Subscriptions      string `json:"subscriptions,omitempty"`
	WildcardExports    bool   `json:"wildcard_exports,omitempty"`
	DisallowBearer     bool   `json:"disallow_bearer,omitempty"`
	RmSk               string `json:"rm_sk,omitempty"`
	Description        string `json:"description,omitempty"`
	InfoUrl            string `json:"info_url,omitempty"`
	JsTier             string `json:"js_tier,omitempty"`
	RmJsTier           string `json:"rm_js_tier,omitempty"`
	JsMemStorage       string `json:"js_mem_storage,omitempty"`
	JsDiskStorage      string `json:"js_disk_storage,omitempty"`
	JsStreams          string `json:"js_streams,omitempty"`
	JsConsumer         string `json:"js_consumer,omitempty"`
	JsMaxMemStream     string `json:"js_max_mem_stream,omitempty"`
	JsMaxDiskStream    string `json:"js_max_disk_stream,omitempty"`
	JsMaxBytesRequired string `json:"js_max_bytes_required,omitempty"`
	JsMaxAckPending    string `json:"js_max_ack_pending,omitempty"`
}

// @Tags			Account
// @Router			/operator/{operator}/account/{name} [patch]
// @Param			name		path	string				true	"Account name"
// @Param			operator	path	string				true	"Operator name"
// @Param			json		body	updateAccountForm	true	"Account data in json format"
// @Summary		Updates an account
// @Description	Updates an account and returns json with status ok if successful
// @Success		200	{object}	SimpleJSONResponse	"Status ok"
// @Failure		500	{object}	string				"Internal error"
func updateAccount(c echo.Context) error {
	s := &updateAccountForm{}

	err := runNsc(nil, nil, "select", "operator", c.Param("operator"))
	if err != nil {
		return badRequest(c, err)
	}

	err = runNsc(s, c, "edit", "account", c.Param("name"), "--operator", c.Param("operator"))
	if err != nil {
		return badRequest(c, err)
	}

	return c.JSON(200, &SimpleJSONResponse{
		Status:  "200",
		Message: "Account updated",
	})
}
