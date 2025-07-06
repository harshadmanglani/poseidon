package workflows

import (
	"github.com/harshadmanglani/polaris"
	"github.com/harshadmanglani/poseidon/clients"
	"github.com/harshadmanglani/poseidon/utils"
)

type AnalysisBuilder struct {
}

type AnalysisData struct {
	RootCause   map[string]interface{} `json:"rootCause"`
	Service     string                 `json:"service"`
	StartTime   string                 `json:"startTime"`
	Summary     string                 `json:"summary"`
	PossibleFix string                 `json:"possibleFix"`
}

func (m AnalysisBuilder) GetBuilderInfo() polaris.BuilderInfo {
	return polaris.BuilderInfo{
		Consumes: []polaris.IData{
			ContextData{},
		},
		Produces: AnalysisData{},
		Optionals: []polaris.IData{
			LogsData{},
			IncidentsData{},
			OperationsData{},
			MetricsData{},
		},
		Accesses: nil,
	}
}

func (m AnalysisBuilder) Process(context polaris.BuilderContext) polaris.IData {
	trigger, _ := context.Get(ContextData{})
	logsData, logsPresent := context.Get(LogsData{})
	metricsData, metricsPresent := context.Get(MetricsData{})
	operationsData, operationsPresent := context.Get(OperationsData{})

	if !logsPresent && !metricsPresent && !operationsPresent {
		utils.Sugar.Error("Required data not present in context for analysis")
		return nil
	}

	contextData := map[string]interface{}{
		"trigger":    trigger.(ContextData),
		"logs":       logsData.(LogsData).RawLogs,
		"metrics":    metricsData.(MetricsData).Metrics,
		"operations": operationsData.(OperationsData).OperationsHistory,
	}

	outputJson := `{
		"rootCause": {"key": "value"},
		"service": "serviceName",
		"startTime": "2023-01-01T00:00:00Z",
		"summary": "This is a summary of the analysis, keep this limited to 100 words.",
		"possibleFix": "This is a possible fix to stop loss and bring the system to a normal state." 
	}`

	prompt := `
	You are an expert in incident analysis. Given the following context data, analyze the situation and strictly follow the outputJson to return the analysis.
	You are given data about the incident and the system in the following format:
	{
		"trigger":    "the trigger (manual invocation, grafana alert, etc)",
		"logs":       "raw logs",
		"metrics":    "metrics by grafana APIs",
		"operations": "history of operations such as config changes, releases, etc",
	}`

	result, err := clients.Anthropic.Analyze(prompt, contextData, outputJson)
	if err != nil {
		return nil
	}

	utils.Sugar.Infof("Analysis result: %v", result)
	resultMap := result.(map[string]interface{})
	analysis := AnalysisData{
		RootCause:   resultMap["rootCause"].(map[string]interface{}),
		Service:     resultMap["service"].(string),
		StartTime:   resultMap["startTime"].(string),
		Summary:     resultMap["summary"].(string),
		PossibleFix: resultMap["possibleFix"].(string),
	}

	return analysis
}
