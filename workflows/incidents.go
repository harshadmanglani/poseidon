package workflows

import "github.com/harshadmanglani/polaris"

type IncidentsBuilder struct {
}

func (i IncidentsBuilder) GetBuilderInfo() polaris.BuilderInfo {
	return polaris.BuilderInfo{
		Consumes:  []polaris.IData{},
		Produces:  nil,
		Optionals: nil,
		Accesses:  nil,
	}
}

func (i IncidentsBuilder) Process(context polaris.BuilderContext) polaris.IData {
	return nil
}
