package credentials

import (
	"github.com/docker/docker-credential-helpers/osxkeychain"
)

func getCredentialManagers() credentialManagers {
	return credentialManagers{
		"osxkeychain": osxkeychain.Osxkeychain{},
	}
}
