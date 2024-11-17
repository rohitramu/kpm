package docker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/rohitramu/kpm/src/pkg/utils/log"
	"github.com/rohitramu/kpm/src/pkg/utils/validation"
)

type tagsResponse struct {
	NextPage string `json:"next"`
	Results  []struct {
		TagName     string `json:"name"`
		ContentType string `json:"content_type"`
		TagStatus   string `json:"tag_status"`
	} `json:"results"`
}

// GetImageTags gets the list of tags on a Docker repository in a remote registry.
func GetImageTags(
	ch chan<- string,
	imageName string,
	dockerRegistry string,
) (err error) {
	// Validate package name.
	err = validation.ValidatePackageName(imageName)
	if err != nil {
		return err
	}

	// Parse and validate the repository name, and get the namespace and image name.
	dockerRepository := imageName
	if !strings.Contains(dockerRepository, "/") {
		// For first-party images, we need to set the namespace to "library".
		dockerRepository = "library/" + dockerRepository
	}

	// Get image name (combining package name with docker registry)
	GetImageNameWithoutTag(dockerRegistry, imageName)

	// Make the API call and get the results.
	err = getTagsFromDockerRegistryV2Api(ch, "hub.docker.com", dockerRepository)
	if err != nil {
		return fmt.Errorf("failed to get tags from Docker V2 API: %s", err)
	}

	// Channel will be garbage collected when it drops out of scope.
	return nil
}

func getTagsFromDockerRegistryV2Api(ch chan<- string, baseUrl string, dockerRepository string) error {
	// Construct the initial URL
	requestUrl := fmt.Sprintf("https://%s/v2/repositories/%s/tags?ordering=-name", baseUrl, dockerRepository)

	for requestUrl != "" {
		// Make the HTTP request
		httpResponse, err := http.Get(requestUrl)
		if err != nil {
			return fmt.Errorf("failed to call the Docker registry: %s", err)
		}
		defer func() {
			err := httpResponse.Body.Close()
			if err != nil {
				log.Errorf("failed to close response stream: %s", err)
			}
		}()

		// Parse the response
		jsonDecoder := json.NewDecoder(httpResponse.Body)
		response := tagsResponse{}
		err = jsonDecoder.Decode(&response)
		if err != nil {
			return fmt.Errorf("failed to parse the HTTP response from the Docker registry: %s", err)
		}

		// Extract the tags
		for _, obj := range response.Results {
			if obj.ContentType == "image" {
				if obj.TagStatus == "active" {
					ch <- obj.TagName
				} else {
					log.Debugf("Inactive tag: %s:%s", dockerRepository, obj.TagName)
				}
			}
		}

		// Set the next URL
		requestUrl = response.NextPage
	}

	return nil
}
