package src

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"os"
	"time"
)

const fileName = "dataflows.json"

func init() {
	GetEchoRoot().GET("dataflows", listDataFlows)

	GetEchoRoot().POST("dataflows", addDataFlow)

	GetEchoRoot().DELETE("dataflow/:id", deleteDataFlow)
}

type addDataFlowForm struct {
	Name    string `json:"name" validate:"required"`
	Server  string `json:"server" validate:"required"`
	Created string `json:"created,omitempty"`
}

func readJsonFile(fileName string, obj interface{}) error {
	bytedJson, err := os.ReadFile(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			// create empty file
			_, err = os.Create(fileName)
			if err != nil {
				return err
			}
			// write empty array
			bytedJson, err = json.Marshal([]string{})
			if err != nil {
				return err
			}

			err = os.WriteFile(fileName, bytedJson, 0644)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	return json.Unmarshal(bytedJson, obj)
}

func readDataFlows() ([]addDataFlowForm, error) {
	var dataFlows []addDataFlowForm
	err := readJsonFile(fileName, &dataFlows)
	return dataFlows, err
}

func storeJson(filename string, dataFlows interface{}) error {
	jsonData, err := json.Marshal(dataFlows)
	if err != nil {
		return err
	}

	return os.WriteFile(filename, jsonData, 0644)
}

// @Tags			DataFlow
// @Router			/dataflow [post]
// @Summary		Add a dataflow
// @Description	Add a dataflow to the store
// @Param			json	body		addDataFlowForm		true	"request body"
// @Success		200		{object}	SimpleJSONResponse	"DataFlow added"
// @Failure		400		{object}	SimpleJSONResponse	"Bad request"
// @Failure		500		{object}	string				"Internal error"
func addDataFlow(c echo.Context) error {
	var dataFlow addDataFlowForm
	err := c.Bind(&dataFlow)
	if err != nil {
		return badRequest(c, err)
	}

	dataFlow.Created = GetTimeNow()

	dataFlows, err := readDataFlows()
	if err != nil {
		return badRequest(c, err)
	}

	dataFlows = append(dataFlows, dataFlow)

	err = storeJson(fileName, dataFlows)
	if err != nil {
		return badRequest(c, err)
	}

	return c.JSON(200, &SimpleJSONResponse{
		Status:  "200",
		Message: "DataFlow added",
	})
}

func GetTimeNow() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// @Tags			DataFlow
// @Router			/dataflows [get]
// @Summary		List dataflows
// @Description	Returns json list of existing dataflows
// @Success		200	{object}	[]string	"DataFlows list"
// @Failure		500	{object}	string		"Internal error"
func listDataFlows(c echo.Context) error {
	dataFlows, err := readDataFlows()
	if err != nil {
		return badRequest(c, err)
	}

	return c.JSON(200, map[string][]addDataFlowForm{"dataflows": dataFlows})
}

// @Tags			DataFlow
// @Router			/dataflow/{id} [delete]
// @Summary		Delete a dataflow
// @Description	Delete a dataflow from the store
// @Param			id	path	string				true	"DataFlow ID"
// @Success		200		{object}	SimpleJSONResponse	"DataFlow deleted"
// @Failure		400		{object}	SimpleJSONResponse	"Bad request"
// @Failure		500		{object}	string				"Internal error"
func deleteDataFlow(c echo.Context) error {
	id := c.Param("id")

	dataFlows, err := readDataFlows()
	if err != nil {
		return badRequest(c, err)
	}

	var newDataFlows []addDataFlowForm
	for _, dataFlow := range dataFlows {
		if dataFlow.Name != id {
			newDataFlows = append(newDataFlows, dataFlow)
		}
	}

	err = storeJson(fileName, newDataFlows)
	if err != nil {
		return badRequest(c, err)
	}

	return c.JSON(200, &SimpleJSONResponse{
		Status:  "200",
		Message: "DataFlow deleted",
	})
}
