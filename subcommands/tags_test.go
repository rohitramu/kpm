package subcommands

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetVersions(t *testing.T) {
	Convey("Get tags", t, func() {
		packageName := "kpmtool/example"
		So(GetPackageVersionsCmd(&packageName, nil), ShouldBeNil)
	})
}
