package workflows

import (
	"fmt"

	"github.com/harshadmanglani/polaris"
	"github.com/harshadmanglani/poseidon/clients"
	"github.com/harshadmanglani/poseidon/utils"
)

type AnalysisBuilder struct {
	anthropicClient *clients.AnthropicClient
}

type AnalysisData struct {
	RootCause string `json:"rootCause"`
	Service   string `json:"service"`
	StartTime string `json:"startTime"`
	Summary   string `json:"summary"`
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
	logsData, logsPresent := context.Get(LogsData{})
	incidentsData, incidentsPresent := context.Get(IncidentsData{})
	metricsData, metricsPresent := context.Get(MetricsData{})
	operationsData, operationsPresent := context.Get(OperationsData{})

	if !logsPresent || !incidentsPresent || !metricsPresent || !operationsPresent {
		utils.Sugar.Error("Required data not present in context for analysis")
		return nil
	}

	contextData := []string{
		fmt.Sprintf("Logs: %v", logsData),
		fmt.Sprintf("Incidents: %v", incidentsData),
		fmt.Sprintf("Metrics: %v", metricsData),
		fmt.Sprintf("Operations: %v", operationsData),
	}

	outputJson := `{
		"rootCause": {"key": "value"},
		"service": "serviceName",
		"startTime": "2023-01-01T00:00:00Z",
		"summary": "This is a summary of the analysis",
	}`

	prompt := "Analyze this incident data and determine the root cause, impact, affected service and components"

	result, err := m.anthropicClient.Analyze(prompt, contextData, outputJson)
	if err != nil {
		return nil
	}

	analysis := AnalysisData{
		RootCause: result["rootCause"],
		Service:   result["service"],
		StartTime: result["startTime"],
		Summary:   result["summary"],
	}

	return analysis
}
