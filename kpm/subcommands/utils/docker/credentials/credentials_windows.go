package credentials

import (
	"github.com/docker/docker-credential-helpers/wincred"
)

func getCredentialManagers() credentialManagers {
	return credentialManagers{
		"wincred": wincred.Wincred{},
	}
}
