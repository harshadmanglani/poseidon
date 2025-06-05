package workflows

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/harshadmanglani/polaris"
	"github.com/harshadmanglani/poseidon/clients"
	"github.com/harshadmanglani/poseidon/config"
	"github.com/harshadmanglani/poseidon/utils"
)

type OperationsBuilder struct {
}

type Operation struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	Name        string `json:"name"`
	Service     string `json:"service"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
	Description string `json:"description"`
}

type OperationsData struct {
}

func (o OperationsBuilder) GetBuilderInfo() polaris.BuilderInfo {
	return polaris.BuilderInfo{
		Consumes:  []polaris.IData{ContextData{}},
		Produces:  nil,
		Optionals: nil,
		Accesses:  nil,
	}
}

func (o OperationsBuilder) Process(context polaris.BuilderContext) polaris.IData {
	if !config.PoseidonConf.Workflows.Metrics.Enabled {
		return nil
	}

	operationsEndpoint := fmt.Sprintf("%s:%d",
		config.PoseidonConf.Workflows.Operations.Endpoint,
		config.PoseidonConf.Workflows.Operations.Port,
	)

	var response interface{}
	err := clients.OperationsClient.Get(operationsEndpoint, response)
	if err != nil {
		utils.Sugar.Errorf("Error fetching operations from %s: %v", operationsEndpoint, err)
		return nil
	}

	responseMap, ok := response.(map[string]interface{})
	if !ok {
		utils.Sugar.Errorf("Error converting response to map[string]interface{}: %v", err)
		return nil
	}

	jsonStr, err := json.Marshal(responseMap)
	if err != nil {
		utils.Sugar.Errorf("Error marshalling response to JSON: %v", err)
		return nil
	}

	operationsData, err := clients.Anthropic.ConvertResponse(string(jsonStr), reflect.TypeOf(OperationsData{}))
	if err != nil {
		utils.Sugar.Errorf("Error converting response to OperationsData: %v", err)
		return nil
	}

	utils.Sugar.Infof("operations data processed successfully: %v", operationsData)
	return operationsData

	return nil
}
