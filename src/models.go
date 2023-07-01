package src

type SimpleJSONResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type OperatorDescription struct {
	Iat  int    `json:"iat"`
	Iss  string `json:"iss"`
	Jti  string `json:"jti"`
	Name string `json:"name"`
	Nats struct {
		Type    string `json:"type"`
		Version int    `json:"version"`
	} `json:"nats"`
	Sub string `json:"sub"`
}

type AccountDescription struct {
	Iat  int    `json:"iat"`
	Iss  string `json:"iss"`
	Jti  string `json:"jti"`
	Name string `json:"name"`
	Nats struct {
		Authorization struct {
			AuthUsers interface{} `json:"auth_users"`
		} `json:"authorization"`
		DefaultPermissions struct {
			Pub map[string]interface{} `json:"pub"`
			Sub map[string]interface{} `json:"sub"`
		} `json:"default_permissions"`
		Limits struct {
			Conn      int  `json:"conn"`
			Data      int  `json:"data"`
			Exports   int  `json:"exports"`
			Imports   int  `json:"imports"`
			Leaf      int  `json:"leaf"`
			Payload   int  `json:"payload"`
			Subs      int  `json:"subs"`
			Wildcards bool `json:"wildcards"`
		} `json:"limits"`
		Type    string `json:"type"`
		Version int    `json:"version"`
	} `json:"nats"`
	Sub string `json:"sub"`
}

type UserDescription struct {
	Iat  int    `json:"iat"`
	Iss  string `json:"iss"`
	Jti  string `json:"jti"`
	Name string `json:"name"`
	Nats struct {
		Data    int                    `json:"data"`
		Payload int                    `json:"payload"`
		Pub     map[string]interface{} `json:"pub"`
		Sub     map[string]interface{} `json:"sub"`
		Subs    int                    `json:"subs"`
		Type    string                 `json:"type"`
		Version int                    `json:"version"`
	} `json:"nats"`
	Sub string `json:"sub"`
}
