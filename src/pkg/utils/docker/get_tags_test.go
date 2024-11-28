package docker

import (
	"sync"
	"testing"

	"github.com/rohitramu/kpm/src/pkg/utils/log"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGetVersions(t *testing.T) {
	Convey("Get tags", t, func() {
		var packageName = "kpmtool/example"
		var dockerRegistry = DefaultDockerRegistry

		var ch = make(chan string, 1)

		var tags []string
		var wg sync.WaitGroup
		go func() {
			defer wg.Done()
			wg.Add(1)
			for tag := range ch {
				tags = append(tags, tag)
			}
		}()

		log.SetLevel(log.LevelDebug)

		var err = GetImageTags(ch, packageName, dockerRegistry)
		close(ch)
		wg.Wait()
		So(err, ShouldBeNil)
		So(tags, ShouldContain, "1.0.0")
	})
}
