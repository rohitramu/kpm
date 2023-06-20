package subcommands

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/rohitramu/kpm/subcommands/utils/docker"
	"github.com/rohitramu/kpm/subcommands/utils/log"
	"github.com/rohitramu/kpm/subcommands/utils/validation"
)

type tagsResponse struct {
	NextPage string `json:"next"`
	Results  []struct {
		TagName     string `json:"name"`
		ContentType string `json:"content_type"`
		TagStatus   string `json:"tag_status"`
	} `json:"results"`
}

// GetPackageVersionsCmd gets the list of tags on a Docker repository in a remote registry.
func GetPackageVersionsCmd(packageNameArg *string, dockerRegistryArg *string) error {
	// Get docker registry location.
	dockerRegistry := validation.GetStringOrDefault(dockerRegistryArg, docker.DefaultDockerRegistry)

	// Get package name.
	packageName, err := validation.GetStringOrError(packageNameArg, "packageName")
	if err != nil {
		return err
	}

	// Validate package name.
	err = validation.ValidatePackageName(packageName)
	if err != nil {
		return err
	}

	// Parse and validate the repository name, and get the namespace and image name.
	dockerRepository := packageName
	if !strings.Contains(dockerRepository, "/") {
		// For first-party images, we need to set the namespace to "library".
		dockerRepository = "library/" + dockerRepository
	}

	// Get image name (combining package name with docker registry)
	docker.GetImageNameWithoutTag(dockerRegistry, packageName)

	// Create a channel to send the results.
	ch := make(chan string, 1)

	// Set up the receiver
	go func() {
		// Print the results.
		numTags := 0
		for tag := range ch {
			numTags++
			fmt.Println(tag)
		}

		if numTags == 0 {
			log.Warning(`There were no tags in repository "%s"`, packageName)
		}
	}()

	// Make the API call and get the results.
	err = getTagsFromDockerV2Api(ch, "hub.docker.com", dockerRepository)
	if err != nil {
		return fmt.Errorf("Failed to get tags from Docker V2 API: %s", err)
	}

	// Channel will be garbage collected when it drops out of scope.
	return nil
}

func getTagsFromDockerV2Api(ch chan<- string, baseUrl string, dockerRepository string) error {
	// Construct the initial URL
	requestUrl := fmt.Sprintf("https://%s/v2/repositories/%s/tags?ordering=-name", baseUrl, dockerRepository)

	for requestUrl != "" {
		// Make the HTTP request
		httpResponse, err := http.Get(requestUrl)
		if err != nil {
			return fmt.Errorf("Failed to call the Docker registry: %s", err)
		}
		defer func() {
			err := httpResponse.Body.Close()
			if err != nil {
				log.Error("Failed to close response stream: %s", err)
			}
		}()

		// Parse the response
		jsonDecoder := json.NewDecoder(httpResponse.Body)
		response := tagsResponse{}
		err = jsonDecoder.Decode(&response)
		if err != nil {
			return fmt.Errorf("Failed to read the HTTP response from the Docker registry: %s", err)
		}

		// Extract the tags
		for _, obj := range response.Results {
			if obj.ContentType == "image" && obj.TagStatus == "active" {
				ch <- obj.TagName
			}
		}

		// Set the next URL
		requestUrl = response.NextPage
	}

	return nil
}
