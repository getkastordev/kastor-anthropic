package schema

import protocol "github.com/weirdGuy/kastor/protocol/v1"

type Target = protocol.Target

const (
	SchemeEnv        = protocol.SchemeEnv
	SchemeConnection = protocol.SchemeConnection
)

var ParseCredentialRef = protocol.ParseCredentialRef
