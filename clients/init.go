package clients

import "github.com/harshadmanglani/poseidon/config"

var Anthropic *AnthropicClient

func Init() {
	Anthropic = NewAnthropicClient(config.PoseidonConf.Clients.Anthropic.Key)
}
