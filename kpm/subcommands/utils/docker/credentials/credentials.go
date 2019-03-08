package credentials

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	dockerCreds "github.com/docker/docker-credential-helpers/credentials"
	dockerTypes "github.com/docker/docker/api/types"

	"../../files"
	"../../logger"
)

// DockerConfig represents the structure of the Docker ~/.docker/config.json file.
type DockerConfig struct {
	Auths      *map[string]dockerConfigAuth `json:"auths"`
	CredsStore *string                      `json:"credsStore"`
}

type dockerConfigAuth struct {
	Auth     *string `json:"auth"`
	Username *string `json:"username"`
	Email    *string `json:"email"`
}

type credentialManagers map[string]dockerCreds.Helper

// GetRegistryURLs returns all of the Docker registry URLs in the Docker configuration file.
func (config *DockerConfig) GetRegistryURLs() ([]string, error) {
	// See if there are any Docker credentials at all
	var auths map[string]dockerConfigAuth
	if config.Auths == nil {
		return nil, fmt.Errorf("No Docker credentials found in Docker configuration")
	}
	auths = *config.Auths

	var registryURLs = make([]string, len(auths))
	var i = 0
	for key := range auths {
		registryURLs[i] = key
		i++
	}

	return registryURLs, nil
}

// GetCredentials retrieves credentials from a Docker config.
func (config *DockerConfig) GetCredentials(dockerRegistryURL string) (string, error) {
	var err error
	var ok bool

	// See if there are any Docker credentials at all
	var auths map[string]dockerConfigAuth
	if config.Auths == nil {
		return "", fmt.Errorf("No Docker credentials found in Docker configuration")
	}
	auths = *config.Auths

	// See if there are any credentials for the given Docker registry URL
	var auth dockerConfigAuth
	auth, ok = auths[dockerRegistryURL]
	if !ok {
		return "", fmt.Errorf("No Docker credentials found for registry URL: %s", dockerRegistryURL)
	}

	// Get auth config
	var authConfig dockerTypes.AuthConfig
	authConfig, err = getAuthConfig(dockerRegistryURL, (*config).CredsStore, auth)
	if err != nil {
		return "", err
	}

	// Convert auth config to credentials
	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		logger.Default.Error.Panicln(err)
	}
	var authString = base64.URLEncoding.EncodeToString(encodedJSON)

	logger.Default.Verbose.Println("Auth:\n" + string(encodedJSON))

	return authString, nil
}

// GetDockerConfig retrieves the Docker configuration from the ~/.docker/config.json file.
func GetDockerConfig() (*DockerConfig, error) {
	var err error

	// Get absolute path of Docker config file
	var path string
	path, err = files.GetAbsolutePath("~/.docker/config.json")
	if err != nil {
		return nil, err
	}

	// Read the Docker config file
	var configJSON []byte
	configJSON, err = files.ReadBytes(path)
	if err != nil {
		return nil, err
	}

	// Parse the Docker config file
	var config = new(DockerConfig)
	err = json.Unmarshal(configJSON, config)
	if err != nil {
		return nil, fmt.Errorf("Failed to deserialize Docker config file: %s", err)
	}

	return config, nil
}

func getAuthConfig(dockerRegistryURL string, credsStoreTypeArg *string, auth dockerConfigAuth) (dockerTypes.AuthConfig, error) {
	var err error

	// Create the result object
	var result = dockerTypes.AuthConfig{
		// ServerAddress: dockerRegistryURL,
	}

	// If the auth token was provided, use the details in the Docker config file
	if auth.Auth != nil {
		result.Auth = *auth.Auth

		// Get the auth token if it was provided
		if auth.Username != nil {
			result.Username = *auth.Username
		}

		// Get the email address if it was provided
		if auth.Email != nil {
			result.Email = *auth.Email
		}

		return result, nil
	}

	// Since the auth token hasn't been provided, see if there is a credential store defined
	var credsStoreType string
	if credsStoreTypeArg == nil {
		return result, fmt.Errorf("No auth token found for Docker credentials for registry URL: %s", dockerRegistryURL)
	}
	credsStoreType = *credsStoreTypeArg

	// Get the appropriate credential manager based on the credential store type
	var credsStore dockerCreds.Helper
	credsStore, err = getCredsStore(credsStoreType)
	if err != nil {
		return result, err
	}

	// Get credential from the credential store
	var username string
	var password string
	username, password, err = credsStore.Get(dockerRegistryURL)
	if err != nil {
		return result, err
	}

	result.Username = username
	result.Password = password

	return result, nil
}

func getCredsStore(credsStoreType string) (dockerCreds.Helper, error) {
	var credentialManagers = getCredentialManagers()

	var ok bool
	var credsStore dockerCreds.Helper
	credsStore, ok = credentialManagers[credsStoreType]
	if !ok {
		return nil, fmt.Errorf("Failed to get credential storage manager: %s", credsStoreType)
	}

	return credsStore, nil
}
