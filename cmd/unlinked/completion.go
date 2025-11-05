package main

import (
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate shell completion scripts",
	Long: `Generate shell completion scripts for unlinked.

To load completions:

Bash:
  $ source <(unlinked completion bash)
  # To load completions for each session, add to your .bashrc:
  $ echo 'source <(unlinked completion bash)' >> ~/.bashrc

Zsh:
  # If shell completion is not already enabled, enable it:
  $ echo "autoload -U compinit; compinit" >> ~/.zshrc
  # To load completions for each session:
  $ unlinked completion zsh > "${fpath[1]}/_unlinked"
  $ source ~/.zshrc

Fish:
  $ unlinked completion fish | source
  # To load completions for each session:
  $ unlinked completion fish > ~/.config/fish/completions/unlinked.fish

PowerShell:
  PS> unlinked completion powershell | Out-String | Invoke-Expression
  # To load completions for every session, add to your PowerShell profile:
  PS> Add-Content $PROFILE "unlinked completion powershell | Out-String | Invoke-Expression"`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	RunE: func(cmd *cobra.Command, args []string) error {
		switch args[0] {
		case "bash":
			return rootCmd.GenBashCompletion(os.Stdout)
		case "zsh":
			return rootCmd.GenZshCompletion(os.Stdout)
		case "fish":
			return rootCmd.GenFishCompletion(os.Stdout, true)
		case "powershell":
			return rootCmd.GenPowerShellCompletionWithDesc(os.Stdout)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
