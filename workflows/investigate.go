package workflows

import (
	"errors"

	"github.com/harshadmanglani/polaris"
	"github.com/harshadmanglani/poseidon/utils"
)

var InvestigateWorkflowKey string
var executor polaris.Executor

func Init() {
	InvestigateWorkflowKey = "INVESTIGATION_WORKFLOW"
	polaris.RegisterWorkflow(InvestigateWorkflowKey, InvestigateWorkflow{})

	executor = polaris.Executor{}
}

func Invoke(id string, contextData ContextData) (AnalysisData, error) {
	response, err := executor.Sequential(InvestigateWorkflowKey, id, contextData)
	if err != nil {
		utils.Sugar.Errorf("Error executing workflow %s: %v", InvestigateWorkflowKey, err)
		return AnalysisData{}, err
	}

	analysis, ok := response.Get(AnalysisData{})
	if !ok {
		utils.Sugar.Errorf("Error retrieving AnalysisData from response: %v", response)
		return AnalysisData{}, errors.New("failed to retrieve AnalysisData from response")
	}
	return analysis.(AnalysisData), nil
}

type ContextData struct {
	Service   string `json:"service"`
	Type      string `json:"type"`
	ID        string `json:"id"`
	Timestamp string `json:"timestamp"`
}

type InvestigateWorkflow struct {
}

func (iW InvestigateWorkflow) GetWorkflowMeta() polaris.WorkflowMeta {
	return polaris.WorkflowMeta{
		Builders: []polaris.IBuilder{
			LogsBuilder{},
			IncidentsBuilder{},
			MetricsBuilder{},
			OperationsBuilder{},
			AnalysisBuilder{},
		},
		TargetData: AnalysisData{},
	}
}
