package client

import (
	"github.com/pteropackages/soar/logger"
	"github.com/pteropackages/soar/util"
	"github.com/spf13/cobra"
)

var log = logger.New()

func GroupCommands() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "client subcommand [options]",
		Short: "client api management",
		Long:  "Commands for interacting with the CLient API.",
	}

	util.ApplyDefaultFlags(getAccountCmd)
	util.ApplyDefaultFlags(getPermissionsCmd)
	util.ApplyDefaultFlags(getServersCmd)
	util.ApplyDefaultFlags(getTwoFactorCodeCmd)
	util.ApplyDefaultFlags(enableTwoFactorCmd)
	util.ApplyDefaultFlags(disableTwoFactorCmd)
	util.ApplyDefaultFlags(getAccountActivityCmd)
	util.ApplyDefaultFlags(getAPIKeysCmd)
	util.ApplyDefaultFlags(deleteAPIKeyCmd)
	util.ApplyDefaultFlags(getServerWSCmd)
	util.ApplyDefaultFlags(getServerResourcesCmd)
	util.ApplyDefaultFlags(getServerActivityCmd)
	util.ApplyDefaultFlags(sendServerCommandCmd)
	util.ApplyDefaultFlags(setServerPowerStateCmd)
	util.ApplyDefaultFlags(getDatabasesCmd)
	util.ApplyDefaultFlags(listFilesCmd)
	util.ApplyDefaultFlags(getFileInfoCmd)
	util.ApplyDefaultFlags(getFileContentsCmd)
	util.ApplyDefaultFlags(downloadFileCmd)
	util.ApplyDefaultFlags(renameFileCmd)
	util.ApplyDefaultFlags(copyFileCmd)
	util.ApplyDefaultFlags(writeFileCmd)
	util.ApplyDefaultFlags(createFileCmd)
	util.ApplyDefaultFlags(compressFilesCmd)
	util.ApplyDefaultFlags(decompressFileCmd)
	util.ApplyDefaultFlags(deleteFilesCmd)
	util.ApplyDefaultFlags(createFolderCmd)
	util.ApplyDefaultFlags(chmodFileCmd)
	util.ApplyDefaultFlags(pullFileCmd)

	listFilesCmd.Flags().BoolP("dir", "d", false, "only list directories")
	listFilesCmd.Flags().BoolP("file", "f", false, "only list files")
	listFilesCmd.Flags().String("root", "/", "the root directory")
	downloadFileCmd.Flags().String("dest", "", "the path to save the file at")
	downloadFileCmd.Flags().BoolP("url-only", "U", false, "only return the url")
	renameFileCmd.Flags().String("root", "/", "the root directory of the file")
	compressFilesCmd.Flags().String("root", "/", "the root directory of the files")
	decompressFileCmd.Flags().String("root", "/", "the root directory of the file")
	deleteFilesCmd.Flags().String("root", "/", "the root directory of the files")
	createFolderCmd.Flags().String("root", "/", "the root directory for the folder")
	chmodFileCmd.Flags().String("root", "/", "the root directory of the file")
	pullFileCmd.Flags().String("dest", "", "the destination directory for the file")
	pullFileCmd.Flags().String("name", "", "the name to save the file as")
	pullFileCmd.Flags().Bool("use-header", false, "use the source content header")
	pullFileCmd.Flags().BoolP("foreground", "f", false, "pull the file in the foreground")

	cmd.AddCommand(getAccountCmd)
	cmd.AddCommand(getPermissionsCmd)
	cmd.AddCommand(getServersCmd)
	cmd.AddCommand(getTwoFactorCodeCmd)
	cmd.AddCommand(enableTwoFactorCmd)
	cmd.AddCommand(disableTwoFactorCmd)
	cmd.AddCommand(getAccountActivityCmd)
	cmd.AddCommand(getAPIKeysCmd)
	cmd.AddCommand(deleteAPIKeyCmd)
	cmd.AddCommand(getServerWSCmd)
	cmd.AddCommand(getServerResourcesCmd)
	cmd.AddCommand(getServerActivityCmd)
	cmd.AddCommand(sendServerCommandCmd)
	cmd.AddCommand(setServerPowerStateCmd)
	cmd.AddCommand(getDatabasesCmd)
	cmd.AddCommand(listFilesCmd)
	cmd.AddCommand(getFileInfoCmd)
	cmd.AddCommand(getFileContentsCmd)
	cmd.AddCommand(downloadFileCmd)
	cmd.AddCommand(renameFileCmd)
	cmd.AddCommand(copyFileCmd)
	cmd.AddCommand(writeFileCmd)
	cmd.AddCommand(createFileCmd)
	cmd.AddCommand(compressFilesCmd)
	cmd.AddCommand(decompressFileCmd)
	cmd.AddCommand(decompressFileCmd)
	cmd.AddCommand(deleteFilesCmd)
	cmd.AddCommand(createFolderCmd)
	cmd.AddCommand(chmodFileCmd)
	cmd.AddCommand(pullFileCmd)

	return cmd
}
