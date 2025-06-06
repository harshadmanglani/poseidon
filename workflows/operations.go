package workflows

import (
	"fmt"
	"strings"

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
	OperationsHistory map[string]interface{} `json:"operations_history"`
}

func (o OperationsBuilder) GetBuilderInfo() polaris.BuilderInfo {
	return polaris.BuilderInfo{
		Consumes:  []polaris.IData{ContextData{}},
		Produces:  OperationsData{},
		Optionals: nil,
		Accesses:  nil,
	}
}

func (o OperationsBuilder) Process(context polaris.BuilderContext) polaris.IData {
	if !config.PoseidonConf.Workflows.Metrics.Enabled {
		return nil
	}
	ctx, ok := context.Get(ContextData{})
	if !ok {
		utils.Sugar.Errorf("Error retrieving ContextData from context: %v", context)
		return nil
	}
	ct, _ := ctx.(ContextData)

	operationsEndpoint := fmt.Sprintf(":%d%s",
		config.PoseidonConf.Workflows.Operations.Port,
		strings.Replace(config.PoseidonConf.Workflows.Operations.Endpoint, "{id}", ct.ID, -1),
	)

	response := make(map[string]interface{})
	err := clients.OperationsClient.Get(operationsEndpoint, &response)
	if err != nil {
		utils.Sugar.Errorf("Error fetching operations from %s: %v", operationsEndpoint, err)
		return nil
	}

	utils.Sugar.Infof("operations data processed successfully: %v", response)
	return OperationsData{
		OperationsHistory: response,
	}
}
