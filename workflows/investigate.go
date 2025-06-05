package workflows

import "github.com/harshadmanglani/polaris"

var InvestigateWorkflowKey string

func Init() {
	InvestigateWorkflowKey = "INVESTIGATION_WORKFLOW"
	polaris.RegisterWorkflow(InvestigateWorkflowKey, InvestigateWorkflow{})
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
