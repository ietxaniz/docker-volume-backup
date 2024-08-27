package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}

func addListFlags(cmd *cobra.Command) {
	cmd.Flags().String("prefix", "", "Prefix of bucket")
	cmd.Flags().String("delimiter", "", "List delimiter. '' for recursive '/' for local items only")
}

func addKeyFlags(cmd *cobra.Command) {
	cmd.Flags().String("salt", "", "Salt for key generation")
	cmd.Flags().Int("iterations", 0, "Iterations for key generation")
}

func init() {
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(derivekeyCmd)
	rootCmd.AddCommand(volumebackupCmd)
	rootCmd.AddCommand(volumerestoreCmd)
	rootCmd.AddCommand(fileEncryptCmd)
	rootCmd.AddCommand(fileDecryptCmd)
	rootCmd.AddCommand(keyGenerateCmd)
	rootCmd.AddCommand(keyEncryptCmd)
	rootCmd.AddCommand(keyDecryptCmd)
	rootCmd.AddCommand(keyDecrypt2Cmd)

	volumebackupCmd.Flags().BoolP("no-compression", "n", false, "Create backup without compression")

	addListFlags(listCmd)
	addKeyFlags(derivekeyCmd)
}
