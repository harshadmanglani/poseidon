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

type MetricsBuilder struct {
}

type UnitMetric struct {
	Uri       string `json:"uri"`
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Count     int64  `json:"count"`
}

type MetricsData struct {
	UnitMetrics []UnitMetric `json:"unit_metrics"`
}

func (m MetricsBuilder) GetBuilderInfo() polaris.BuilderInfo {
	return polaris.BuilderInfo{
		Consumes: []polaris.IData{
			ContextData{},
		},
		Produces:  MetricsData{},
		Optionals: nil,
		Accesses:  nil,
	}
}

func (m MetricsBuilder) Process(context polaris.BuilderContext) polaris.IData {
	if !config.PoseidonConf.Workflows.Metrics.Enabled {
		utils.Sugar.Info("Metrics workflow is disabled, skipping processing.")
		return nil
	}

	metricsEndpoint := fmt.Sprintf("%s:%d",
		config.PoseidonConf.Workflows.Metrics.Endpoint,
		config.PoseidonConf.Workflows.Metrics.Port,
	)

	var response interface{}
	err := clients.MetricsClient.Get(metricsEndpoint, response)
	if err != nil {
		utils.Sugar.Errorf("Error fetching metrics from %s: %v", metricsEndpoint, err)
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

	metricsData, err := clients.Anthropic.ConvertResponse(string(jsonStr), reflect.TypeOf(MetricsData{}))
	if err != nil {
		utils.Sugar.Errorf("Error converting response to MetricsData: %v", err)
		return nil
	}

	utils.Sugar.Infof("Metrics data processed successfully: %v", metricsData)
	return metricsData
}
