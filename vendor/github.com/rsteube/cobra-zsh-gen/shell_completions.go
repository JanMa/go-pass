package cobra

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// MarkFlagRequired adds the BashCompOneRequiredFlag annotation to the named flag if it exists,
// and causes your command to report an error if invoked without the flag.
func (c *ZshCommand) MarkFlagRequired(name string) error {
	return MarkFlagRequired(c.Command.Flags(), name)
}

// MarkPersistentFlagRequired adds the BashCompOneRequiredFlag annotation to the named persistent flag if it exists,
// and causes your command to report an error if invoked without the flag.
func (c *ZshCommand) MarkPersistentFlagRequired(name string) error {
	return MarkFlagRequired(c.Command.PersistentFlags(), name)
}

// MarkFlagRequired adds the BashCompOneRequiredFlag annotation to the named flag if it exists,
// and causes your command to report an error if invoked without the flag.
func MarkFlagRequired(flags *pflag.FlagSet, name string) error {
	return flags.SetAnnotation(name, cobra.BashCompOneRequiredFlag, []string{"true"})
}

// MarkFlagFilename adds the BashCompFilenameExt annotation to the named flag, if it exists.
// Generated bash autocompletion will select filenames for the flag, limiting to named extensions if provided.
func (c *ZshCommand) MarkFlagFilename(name string, extensions ...string) error {
	return MarkFlagFilename(c.Command.Flags(), name, extensions...)
}

// MarkFlagCustom adds the BashCompCustom annotation to the named flag, if it exists.
// Generated bash autocompletion will call the bash function f for the flag.
func (c *ZshCommand) MarkFlagCustom(name string, f string) error {
	return MarkFlagCustom(c.Command.Flags(), name, f)
}

// MarkPFlagCustom adds the BashCompCustom annotation to the named pflag, if it exists.
// Generated shell autocompletion will call the function f for the flag.
func (c *ZshCommand) MarkPFlagCustom(name string, f string) error {
	return MarkFlagCustom(c.Command.PersistentFlags(), name, f)
}

// MarkPersistentFlagFilename instructs the various shell completion
// implementations to limit completions for this persistent flag to the
// specified extensions (patterns).
//
// Shell Completion compatibility matrix: bash, zsh
func (c *ZshCommand) MarkPersistentFlagFilename(name string, extensions ...string) error {
	return MarkFlagFilename(c.Command.PersistentFlags(), name, extensions...)
}

// MarkFlagFilename instructs the various shell completion implementations to
// limit completions for this flag to the specified extensions (patterns).
//
// Shell Completion compatibility matrix: bash, zsh
func MarkFlagFilename(flags *pflag.FlagSet, name string, extensions ...string) error {
	return flags.SetAnnotation(name, cobra.BashCompFilenameExt, extensions)
}

// MarkFlagCustom instructs the various shell completion implementations to
// limit completions for this flag to the specified extensions (patterns).
//
// Shell Completion compatibility matrix: bash, zsh
func MarkFlagCustom(flags *pflag.FlagSet, name string, f string) error {
	return flags.SetAnnotation(name, cobra.BashCompCustom, []string{f})
}

// MarkFlagDirname instructs the various shell completion implementations to
// complete only directories with this named flag.
//
// Shell Completion compatibility matrix: zsh
func (c *ZshCommand) MarkFlagDirname(name string) error {
	return MarkFlagDirname(c.Command.Flags(), name)
}

// MarkPersistentFlagDirname instructs the various shell completion
// implementations to complete only directories with this persistent named flag.
//
// Shell Completion compatibility matrix: zsh
func (c *ZshCommand) MarkPersistentFlagDirname(name string) error {
	return MarkFlagDirname(c.Command.PersistentFlags(), name)
}

// MarkFlagDirname instructs the various shell completion implementations to
// complete only directories with this specified flag.
//
// Shell Completion compatibility matrix: zsh
func MarkFlagDirname(flags *pflag.FlagSet, name string) error {
	zshPattern := "-(/)"
	return flags.SetAnnotation(name, zshCompDirname, []string{zshPattern})
}
