package docker

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetVersions(t *testing.T) {
	Convey("Get tags", t, func() {
		var packageName = "kpmtool/example"
		var dockerRegistry = DefaultDockerRegistry

		var ch = make(chan string, 1)

		var tags []string
		go func() {
			for tag := range ch {
				tags = append(tags, tag)
			}
		}()

		var err = GetImageTags(ch, packageName, dockerRegistry)
		So(err, ShouldBeNil)
		So(tags, ShouldContain, "1.0.0")
	})
}
