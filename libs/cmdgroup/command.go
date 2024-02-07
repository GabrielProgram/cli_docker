package cmdgroup

import (
	"io"
	"strings"
	"text/template"
	"unicode"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type CommandWithGroupFlag struct {
	cmd        *cobra.Command
	flagGroups []*FlagGroup
}

func (c *CommandWithGroupFlag) RefreshFlags() {
	for _, fg := range c.flagGroups {
		c.cmd.Flags().AddFlagSet(fg.flagSet)
	}
}

func (c *CommandWithGroupFlag) Command() *cobra.Command {
	return c.cmd
}

func (c *CommandWithGroupFlag) FlagGroups() []*FlagGroup {
	return c.flagGroups
}

func (c *CommandWithGroupFlag) NonGroupedFlags() *pflag.FlagSet {
	nonGrouped := pflag.NewFlagSet("non-grouped", pflag.ContinueOnError)
	c.cmd.LocalFlags().VisitAll(func(f *pflag.Flag) {
		for _, fg := range c.flagGroups {
			if fg.Has(f) {
				return
			}
		}
		nonGrouped.AddFlag(f)
	})

	return nonGrouped
}

func (c *CommandWithGroupFlag) HasNonGroupedFlags() bool {
	return c.NonGroupedFlags().HasFlags()
}

func NewCommandWithGroupFlag(cmd *cobra.Command) *CommandWithGroupFlag {
	cmdWithFlagGroups := &CommandWithGroupFlag{cmd: cmd, flagGroups: make([]*FlagGroup, 0)}
	cmd.SetUsageFunc(func(c *cobra.Command) error {
		err := tmpl(c.OutOrStderr(), c.UsageTemplate(), cmdWithFlagGroups)
		if err != nil {
			c.PrintErrln(err)
		}
		return nil
	})
	cmd.SetUsageTemplate(usageTemplate)
	return cmdWithFlagGroups
}

func (c *CommandWithGroupFlag) AddFlagGroup(name string) *FlagGroup {
	fg := &FlagGroup{name: name, flagSet: pflag.NewFlagSet(name, pflag.ContinueOnError)}
	c.flagGroups = append(c.flagGroups, fg)
	return fg
}

type FlagGroup struct {
	name        string
	description string
	flagSet     *pflag.FlagSet
}

func (c *FlagGroup) Name() string {
	return c.name
}

func (c *FlagGroup) Description() string {
	return c.description
}

func (c *FlagGroup) SetDescription(description string) {
	c.description = description
}

func (c *FlagGroup) FlagSet() *pflag.FlagSet {
	return c.flagSet
}

func (c *FlagGroup) Has(f *pflag.Flag) bool {
	return c.flagSet.Lookup(f.Name) != nil
}

var templateFuncs = template.FuncMap{
	"trim":                    strings.TrimSpace,
	"trimRightSpace":          trimRightSpace,
	"trimTrailingWhitespaces": trimRightSpace,
}

func trimRightSpace(s string) string {
	return strings.TrimRightFunc(s, unicode.IsSpace)
}

// tmpl executes the given template text on data, writing the result to w.
func tmpl(w io.Writer, text string, data interface{}) error {
	t := template.New("top")
	t.Funcs(templateFuncs)
	template.Must(t.Parse(text))
	return t.Execute(w, data)
}
