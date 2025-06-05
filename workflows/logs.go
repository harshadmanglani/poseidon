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

type LogsBuilder struct {
}

type LogsData struct {
	RawLogs []string `json:"raw_logs"`
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

	logsEndpoint := fmt.Sprintf("%s:%d",
		config.PoseidonConf.Workflows.Logs.Endpoint,
		config.PoseidonConf.Workflows.Logs.Port,
	)

	var response interface{}
	err := clients.LogsClient.Get(logsEndpoint, response)
	if err != nil {
		utils.Sugar.Errorf("Error fetching logs from %s: %v", logsEndpoint, err)
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

	logsData, err := clients.Anthropic.ConvertResponse(string(jsonStr), reflect.TypeOf(LogsData{}))
	if err != nil {
		utils.Sugar.Errorf("Error converting response to LogsData: %v", err)
		return nil
	}

	utils.Sugar.Infof("logs data processed successfully: %v", logsData)
	return logsData
}
