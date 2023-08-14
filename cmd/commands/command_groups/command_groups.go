package command_groups

import "github.com/spf13/cobra"

var PackageManagement = &cobra.Group{
	ID:    "package-management",
	Title: "Template Package Management",
}

var PackageEditing = &cobra.Group{
	ID:    "package-editing",
	Title: "Template Package Editing",
}
