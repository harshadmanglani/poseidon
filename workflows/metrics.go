package workflows

import "github.com/harshadmanglani/polaris"

type MetricsBuilder struct {
}

func (m MetricsBuilder) GetBuilderInfo() polaris.BuilderInfo {
	return polaris.BuilderInfo{
		Consumes:  []polaris.IData{},
		Produces:  nil,
		Optionals: nil,
		Accesses:  nil,
	}
}

func (m MetricsBuilder) Process(context polaris.BuilderContext) polaris.IData {
	return nil
}
