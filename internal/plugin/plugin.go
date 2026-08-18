package plugin

import (
	"context"
	"fmt"
	"sort"

	"github.com/getkastordev/kastor-anthropic/internal/claude"
	"github.com/getkastordev/kastor-anthropic/internal/provider"
	protocol "github.com/weirdGuy/kastor/protocol/v1"
)

const Source = "github.com/getkastordev/kastor-anthropic"

var Version = "0.1.0-dev"

type Handler struct{}

func (Handler) Metadata(context.Context) (protocol.Metadata, error) {
	return protocol.Metadata{
		Protocol: protocol.Version,
		Source:   Source,
		Version:  Version,
		Kinds:    []protocol.Kind{protocol.KindPlatform},
		Capabilities: protocol.Capabilities{
			CredentialSchemes: []string{protocol.SchemeConnection},
			Scaffold:          true,
			Config: map[string]protocol.ConfigAttribute{
				"api_key_env": {Type: "string", Description: "Environment variable containing the Anthropic API key"},
				"vault_id":    {Type: "string", Description: "Claude credential vault used by readiness checks"},
			},
			Check: true,
		},
	}, nil
}

func (Handler) Validate(_ context.Context, request *protocol.ValidateRequest) (*protocol.ValidateResponse, error) {
	if request == nil || request.Module == nil || request.Target == nil {
		return nil, fmt.Errorf("anthropic: validate request is incomplete")
	}
	var diagnostics []protocol.Diagnostic
	if request.Target.Type != string(protocol.KindPlatform) {
		diagnostics = append(diagnostics, diagnostic(request.Target.Addr(), "target must have type platform", ""))
	}
	keys := make([]string, 0, len(request.Target.Config))
	for key := range request.Target.Config {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		if key != "api_key_env" && key != "vault_id" {
			diagnostics = append(diagnostics, diagnostic(request.Target.Addr(), "unsupported config attribute", key))
			continue
		}
		if _, ok := request.Target.Config[key].(string); !ok {
			diagnostics = append(diagnostics, diagnostic(request.Target.Addr(), "config attribute must be a string", key))
		}
	}
	referenced := referencedServers(request.Module)
	for _, server := range request.Module.MCPServers {
		if !referenced[server.Name] {
			continue
		}
		if server.Transport == "stdio" {
			diagnostics = append(diagnostics, diagnostic(server.Addr(), "stdio transport is not supported", "Claude Managed Agents requires an HTTP endpoint"))
		}
		auth, ok := server.AuthFor(request.Target.Addr())
		if !ok {
			continue
		}
		scheme, _, err := protocol.ParseCredentialRef(auth.Ref)
		if err != nil {
			continue
		}
		if scheme != protocol.SchemeConnection {
			diagnostics = append(diagnostics, diagnostic(server.Addr(), "unsupported credential scheme", scheme+"://"))
			continue
		}
		vault, ok := request.Target.ConfigString("vault_id")
		if !ok || vault == "" {
			diagnostics = append(diagnostics, diagnostic(request.Target.Addr(), "vault_id is required for connection credentials", auth.Ref))
		}
	}
	return &protocol.ValidateResponse{Diagnostics: diagnostics}, nil
}

func (Handler) Read(ctx context.Context, request *protocol.ReadRequest) (*protocol.ReadResponse, error) {
	implementation, err := claude.Factory(request.Target)
	if err != nil {
		return nil, err
	}
	remote, found, err := implementation.Read(ctx, request.ID)
	if err != nil {
		return nil, err
	}
	return &protocol.ReadResponse{Remote: remote, Found: found}, nil
}

func (Handler) Create(ctx context.Context, request *protocol.CreateRequest) (*protocol.CreateResponse, error) {
	implementation, err := claude.Factory(request.Target)
	if err != nil {
		return nil, err
	}
	id, err := implementation.Create(ctx, request.Desired)
	if err != nil {
		return nil, err
	}
	return &protocol.CreateResponse{ID: id}, nil
}

func (Handler) Update(ctx context.Context, request *protocol.UpdateRequest) error {
	implementation, err := claude.Factory(request.Target)
	if err != nil {
		return err
	}
	return implementation.Update(ctx, request.ID, request.Desired)
}

func (Handler) Delete(ctx context.Context, request *protocol.DeleteRequest) error {
	implementation, err := claude.Factory(request.Target)
	if err != nil {
		return err
	}
	return implementation.Delete(ctx, request.ID)
}

func (Handler) Diff(_ context.Context, request *protocol.DiffRequest) (*protocol.DiffResponse, error) {
	diffs, err := claude.New().Diff(request.Desired, request.Remote)
	if err != nil {
		return nil, err
	}
	return &protocol.DiffResponse{Diffs: diffs}, nil
}

func (Handler) Check(ctx context.Context, request *protocol.CheckRequest) (*protocol.CheckResponse, error) {
	implementation, err := claude.Factory(request.Target)
	if err != nil {
		return nil, err
	}
	checker, ok := implementation.(provider.Checker)
	if !ok {
		return nil, fmt.Errorf("anthropic provider does not implement readiness checks")
	}
	checks, err := checker.Check(ctx, request.Desired, request.Remote)
	if err != nil {
		return nil, err
	}
	return &protocol.CheckResponse{Checks: checks}, nil
}

func diagnostic(addr, summary, detail string) protocol.Diagnostic {
	return protocol.Diagnostic{Severity: protocol.SeverityError, Addr: addr, Summary: summary, Detail: detail}
}

func referencedServers(module *protocol.Module) map[string]bool {
	referenced := map[string]bool{}
	for _, tool := range module.Tools {
		if tool.Source == nil || tool.Source.Kind != "mcp" {
			continue
		}
		server, _, err := protocol.ParseMCPURI(tool.Source.URI)
		if err == nil {
			referenced[server] = true
		}
	}
	return referenced
}
