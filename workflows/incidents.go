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

type IncidentsBuilder struct {
}

type IncidentsData struct {
	IncidentID  string
	RootCause   string
	Impact      string
	Fix         string
	Description string
	Severity    string
	Status      string
	CreatedAt   string
}

func (i IncidentsBuilder) GetBuilderInfo() polaris.BuilderInfo {
	return polaris.BuilderInfo{
		Consumes: []polaris.IData{
			ContextData{},
		},
		Produces:  IncidentsData{},
		Optionals: nil,
		Accesses:  nil,
	}
}

func (i IncidentsBuilder) Process(context polaris.BuilderContext) polaris.IData {
	if !config.PoseidonConf.Workflows.Incidents.Enabled {
		return nil
	}

	incidentEndpoint := fmt.Sprintf("%s:%d",
		config.PoseidonConf.Workflows.Incidents.Endpoint,
		config.PoseidonConf.Workflows.Incidents.Port,
	)

	var response interface{}
	err := clients.IncidentsClient.Get(incidentEndpoint, response)
	if err != nil {
		utils.Sugar.Errorf("Error fetching incidents from %s: %v", incidentEndpoint, err)
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

	incidentsData, err := clients.Anthropic.ConvertResponse(string(jsonStr), reflect.TypeOf(IncidentsData{}))
	if err != nil {
		utils.Sugar.Errorf("Error converting response to IncidentsData: %v", err)
		return nil
	}

	utils.Sugar.Infof("incident data processed successfully: %v", incidentsData)
	return incidentsData
}
