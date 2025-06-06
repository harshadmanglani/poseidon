package workflows

import (
	"fmt"
	"strings"

	"github.com/harshadmanglani/polaris"
	"github.com/harshadmanglani/poseidon/clients"
	"github.com/harshadmanglani/poseidon/config"
	"github.com/harshadmanglani/poseidon/utils"
)

type LogsBuilder struct {
}

type LogsData struct {
	RawLogs map[string]interface{} `json:"raw_logs"`
}

func (l LogsBuilder) GetBuilderInfo() polaris.BuilderInfo {
	return polaris.BuilderInfo{
		Consumes: []polaris.IData{
			ContextData{},
		},
		Produces:  LogsData{},
		Optionals: nil,
		Accesses:  nil,
	}
}

func (l LogsBuilder) Process(context polaris.BuilderContext) polaris.IData {
	if !config.PoseidonConf.Workflows.Logs.Enabled {
		return nil
	}

	ctx, ok := context.Get(ContextData{})
	if !ok {
		utils.Sugar.Errorf("Error retrieving ContextData from context: %v", context)
		return nil
	}
	ct, _ := ctx.(ContextData)

	logsEndpoint := fmt.Sprintf(":%d%s",
		config.PoseidonConf.Workflows.Logs.Port,
		strings.Replace(config.PoseidonConf.Workflows.Logs.Endpoint, "{id}", ct.ID, -1),
	)

	response := make(map[string]interface{})
	err := clients.LogsClient.Get(logsEndpoint, &response)
	if err != nil {
		utils.Sugar.Errorf("Error fetching logs from %s: %v", logsEndpoint, err)
		return nil
	}

	utils.Sugar.Infof("logs data processed successfully: %v", response)
	return LogsData{
		RawLogs: response,
	}
}
