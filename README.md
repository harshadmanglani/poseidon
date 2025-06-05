# Poseidon
An AI agent to investigate production outages (and soon minimize impact!)

## Current Scope
This is literally v0.0.1, and the concept works. I plan to add MCP, a more flexible schema for the API, and overall just make it more universal - so feedback is critical. Please create issues or email me directly: harshad.gm@gmail.com

Kind of vibe coded my way through this, so it just works. Lot of things I'll clean up in the next version.

## How to run
Use the config.yaml file to set up
- a Postgres store
- HTTP info for getting incidents, operations, logs and metrics
- API key for Anthropic and configure the model for analysis
- configure a webhook (like your own service) which can receive a payload (forward to Slack, etc)

## Tests