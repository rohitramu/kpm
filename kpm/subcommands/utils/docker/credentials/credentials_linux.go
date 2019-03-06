package credentials

import (
	"github.com/docker/docker-credential-helpers/pass"
	"github.com/docker/docker-credential-helpers/secretservice"
)

func getCredentialManagers() credentialManagers {
	return credentialManagers{
		"secretservice": secretservice.Secretservice{},
		"pass":          pass.Pass{},
	}
}
