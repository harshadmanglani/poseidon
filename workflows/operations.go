package workflows

import "github.com/harshadmanglani/polaris"

type OperationsBuilder struct {
}

func (o OperationsBuilder) GetBuilderInfo() polaris.BuilderInfo {
	return polaris.BuilderInfo{
		Consumes:  []polaris.IData{},
		Produces:  nil,
		Optionals: nil,
		Accesses:  nil,
	}
}

func (o OperationsBuilder) Process(context polaris.BuilderContext) polaris.IData {
	return nil
}
