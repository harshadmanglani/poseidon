package clients

import (
	"github.com/harshadmanglani/poseidon/config"
)

var Anthropic *AnthropicClient
var LogsClient *OutboundClient
var MetricsClient *OutboundClient
var IncidentsClient *OutboundClient
var OperationsClient *OutboundClient

func Init() {
	cfg := config.PoseidonConf
	Anthropic = NewAnthropicClient(cfg.Clients.Anthropic.Key)
	if cfg.Workflows.Logs.Enabled {
		LogsClient = NewClient(config.PoseidonConf.Workflows.Logs.Host)
	}
	if cfg.Workflows.Metrics.Enabled {
		MetricsClient = NewClient(config.PoseidonConf.Workflows.Metrics.Host)
	}
	if cfg.Workflows.Incidents.Enabled {
		IncidentsClient = NewClient(config.PoseidonConf.Workflows.Incidents.Host)
	}
	if cfg.Workflows.Operations.Enabled {
		OperationsClient = NewClient(config.PoseidonConf.Workflows.Operations.Host)
	}
}
