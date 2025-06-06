package workflows

import (
	"fmt"
	"strings"

	"github.com/harshadmanglani/polaris"
	"github.com/harshadmanglani/poseidon/clients"
	"github.com/harshadmanglani/poseidon/config"
	"github.com/harshadmanglani/poseidon/utils"
)

type IncidentsBuilder struct {
}

type Incident struct {
	IncidentID  string
	RootCause   string
	Impact      string
	Fix         string
	Description string
	Severity    string
	Status      string
	CreatedAt   string
}

type IncidentsData struct {
	Incidents map[string]interface{} `json:"incidents"`
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

	ctx, ok := context.Get(ContextData{})
	if !ok {
		utils.Sugar.Errorf("Error retrieving ContextData from context: %v", context)
		return nil
	}
	ct, _ := ctx.(ContextData)

	incidentEndpoint := fmt.Sprintf(":%d%s",
		config.PoseidonConf.Workflows.Incidents.Port,
		strings.Replace(config.PoseidonConf.Workflows.Incidents.Endpoint, "{id}", ct.ID, -1),
	)

	response := make(map[string]interface{})
	err := clients.IncidentsClient.Get(incidentEndpoint, &response)
	if err != nil {
		utils.Sugar.Errorf("Error fetching incidents from %s: %v", incidentEndpoint, err)
		return nil
	}

	utils.Sugar.Infof("incident data processed successfully: %v", response)
	return IncidentsData{
		Incidents: response,
	}
}
