package pkg

import (
	"fmt"

	"github.com/rohitramu/kpm/src/pkg/utils/files"
	"github.com/rohitramu/kpm/src/pkg/utils/local_package_repo"
	"github.com/rohitramu/kpm/src/pkg/utils/log"
	"github.com/rohitramu/kpm/src/pkg/utils/user_prompts"
)

// PurgeCmd removes all versions of a template package from the local KPM repository.
func PurgeCmd(
	packageName string,
	userHasConfirmed bool,
	kpmHomeDirPath string,
) error {
	// Get KPM home directory.
	var kpmHomeDir string
	kpmHomeDir, err := files.GetAbsolutePath(kpmHomeDirPath)
	if err != nil {
		return err
	}

	if packageName != "" {
		// Get user confirmation.
		if !userHasConfirmed {
			if userHasConfirmed, err = user_prompts.ConfirmWithUser("All versions of package '%s' will be deleted from the local KPM repository.", packageName); err != nil {
				log.Panicf("Failed to get user confirmation. \n%s", err)
			}

			if !userHasConfirmed {
				return fmt.Errorf("purge operation cancelled - user did not confirm the delete action")
			}
		}

		return local_package_repo.RemoveAllVersionsOfPackages(packageName)
	} else {
		// Get user confirmation.
		if !userHasConfirmed {
			if userHasConfirmed, err = user_prompts.ConfirmWithUser("All versions of all packages will be deleted from the local KPM repository."); err != nil {
				log.Panicf("Failed to get user confirmation. \n%s", err)
			}

			if !userHasConfirmed {
				return fmt.Errorf("purge operation cancelled - user did not confirm the delete action")
			}
		}

		return local_package_repo.RemoveAllVersionsOfAllPackages(kpmHomeDir)
	}
}
