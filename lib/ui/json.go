package ui

import (
	"github.com/fatih/color"
	"github.com/nwidger/jsoncolor"
	"github.com/spf13/cobra"
)

func Render(cmd *cobra.Command, v any) error {
	// TODO: render in other formats
	pretty, err := MarshalJSON(v)
	if err != nil {
		return err
	}
	err = startPager(cmd)
	if err != nil {
		return err
	}
	cmd.OutOrStdout().Write(pretty)
	return nil
}

func MarshalJSON(v any) ([]byte, error) {
	// create custom formatter
	f := jsoncolor.NewFormatter()

	// set custom colors
	f.StringColor = color.New(color.FgGreen)
	f.TrueColor = color.New(color.FgGreen, color.Bold)
	f.FalseColor = color.New(color.FgRed)
	f.NumberColor = color.New(color.FgCyan)
	f.NullColor = color.New(color.FgMagenta)
	f.FieldColor = color.New(color.FgWhite, color.Bold)
	f.FieldQuoteColor = color.New(color.FgWhite)
	// KeyColor:        color.New(color.FgWhite),
	// StringColor:     color.New(color.FgGreen),
	// BoolColor:       color.New(),
	// NullColor:       color.New(),

	return jsoncolor.MarshalIndentWithFormatter(v, "", "  ", f)
}
