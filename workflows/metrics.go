package workflows

import (
	"fmt"
	"strings"

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
	Metrics map[string]interface{} `json:"metrics"`
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

	ctx, ok := context.Get(ContextData{})
	if !ok {
		utils.Sugar.Errorf("Error retrieving ContextData from context: %v", context)
		return nil
	}
	ct, _ := ctx.(ContextData)

	metricsEndpoint := fmt.Sprintf(":%d%s",
		config.PoseidonConf.Workflows.Metrics.Port,
		strings.Replace(config.PoseidonConf.Workflows.Metrics.Endpoint, "{id}", ct.ID, -1),
	)

	response := make(map[string]interface{})
	err := clients.MetricsClient.Get(metricsEndpoint, &response)
	if err != nil {
		utils.Sugar.Errorf("Error fetching metrics from %s: %v", metricsEndpoint, err)
		return nil
	}

	utils.Sugar.Infof("Metrics data processed successfully: %v", response)
	return MetricsData{
		Metrics: response,
	}
}
