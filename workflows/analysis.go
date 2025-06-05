package workflows

import (
	"github.com/harshadmanglani/polaris"
)

type AnalysisBuilder struct {
}

type AnalysisData struct {
	RootCause          map[string]string
	Service            string
	StartTime          string
	EndTime            string
	Impact             string
	ImpactedComponents []string
}

func (m AnalysisBuilder) GetBuilderInfo() polaris.BuilderInfo {
	return polaris.BuilderInfo{
		Consumes:  []polaris.IData{},
		Produces:  AnalysisData{},
		Optionals: nil,
		Accesses:  nil,
	}
}

func (m AnalysisBuilder) Process(context polaris.BuilderContext) polaris.IData {
	return nil
}
