# Kastor Anthropic Plugin

External Anthropic Claude Managed Agents platform target for [Kastor](https://github.com/weirdGuy/kastor).

The plugin owns Anthropic SDK integration, target configuration, lifecycle reconciliation, normalization, drift comparison, MCP credential checks, and tool-permission readiness checks. Kastor core sends protocol-v1 resources and owns planning, ordering, state, and locking.

> Status: pre-release. The protocol-v1 implementation is available for integration testing but no stable binary has been released yet.

## Configuration

```hcl
target "production" {
  type   = "platform"
  plugin = "anthropic"

  config {
    api_key_env = "ANTHROPIC_API_KEY"
    vault_id    = "vault_..."
  }
}
```

## Development

```sh
go test ./...
go vet ./...
go build ./cmd/kastor-anthropic
```

Point a local Kastor checkout at the binary with:

```sh
export KASTOR_PLUGIN_ANTHROPIC=/absolute/path/to/kastor-anthropic
```

The source address is `github.com/getkastordev/kastor-anthropic`.

## License

Apache-2.0
