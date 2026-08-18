package plugin

import (
	"context"
	"testing"

	protocol "github.com/weirdGuy/kastor/protocol/v1"
)

func TestValidateRequiresVaultForConnectionCredential(t *testing.T) {
	response, err := (Handler{}).Validate(context.Background(), &protocol.ValidateRequest{
		Target: &protocol.Target{Name: "production", Type: "platform"},
		Module: &protocol.Module{
			Tools:      []*protocol.Tool{{Name: "search", Source: &protocol.ToolSource{Kind: "mcp", URI: "mcp://search/query"}}},
			MCPServers: []*protocol.MCPServer{{Name: "search", Transport: "http", Auth: []*protocol.MCPAuth{{Ref: "connection://cred_123"}}}},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(response.Diagnostics) != 1 || response.Diagnostics[0].Addr != "target.production" {
		t.Fatalf("diagnostics = %#v", response.Diagnostics)
	}
}

func TestDiffDoesNotRequireCredentials(t *testing.T) {
	_, err := (Handler{}).Diff(context.Background(), &protocol.DiffRequest{
		Target:  &protocol.Target{Name: "production", Type: "platform"},
		Desired: &protocol.Resource{Addr: "agent.weather", Config: protocol.Object{}},
	})
	if err == nil {
		return
	}
	// The empty resource is invalid, but the failure must come from config
	// normalization rather than ambient API-key discovery.
	if err.Error() == "Claude Managed Agents API key is missing; set ANTHROPIC_API_KEY" {
		t.Fatal(err)
	}
}
