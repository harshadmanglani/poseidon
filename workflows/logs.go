package workflows

import "github.com/harshadmanglani/polaris"

type LogsBuilder struct {
}

func (l LogsBuilder) GetBuilderInfo() polaris.BuilderInfo {
	return polaris.BuilderInfo{
		Consumes:  []polaris.IData{},
		Produces:  nil,
		Optionals: nil,
		Accesses:  nil,
	}
}

func (l LogsBuilder) Process(context polaris.BuilderContext) polaris.IData {
	return nil
}
