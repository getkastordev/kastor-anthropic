kastor {
  required_plugins {
    anthropic = {
      source  = "github.com/getkastordev/kastor-anthropic"
      version = "~> 0.1"
    }
  }
}

model "claude" {
  provider = "anthropic"
  id       = "claude-sonnet-4-5"
}

target "production" {
  type   = "platform"
  plugin = "anthropic"

  config {
    api_key_env = "ANTHROPIC_API_KEY"
  }
}
