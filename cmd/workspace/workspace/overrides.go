package workspace

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/databricks/cli/libs/cmdio"
	"github.com/databricks/databricks-sdk-go/apierr"
	"github.com/databricks/databricks-sdk-go/service/workspace"
	"github.com/spf13/cobra"
)

func listOverride(listCmd *cobra.Command, listReq *workspace.ListWorkspaceRequest) {
	listReq.Path = "/"
	listCmd.Annotations["template"] = cmdio.Heredoc(`
	{{header "ID"}}	{{header "Type"}}	{{header "Language"}}	{{header "Path"}}
	{{range .}}{{green "%d" .ObjectId}}	{{blue "%s" .ObjectType}}	{{cyan "%s" .Language}}	{{.Path|cyan}}
	{{end}}`)
}

func exportOverride(exportCmd *cobra.Command, exportReq *workspace.ExportRequest) {
	// The export command prints the contents of the file to stdout by default.
	exportCmd.Annotations["template"] = `{{.Content | b64_decode}}`
}

func importOverride(importCmd *cobra.Command, importReq *workspace.Import) {
	var filePath string
	importCmd.Use = "import TARGET_PATH"
	importCmd.Flags().StringVar(&filePath, "file", "", `Path of local file to import`)
	importCmd.MarkFlagsMutuallyExclusive("content", "file")

	originalRunE := importCmd.RunE
	importCmd.RunE = func(cmd *cobra.Command, args []string) error {
		if filePath != "" {
			b, err := os.ReadFile(filePath)
			if err != nil {
				return err
			}
			importReq.Content = base64.StdEncoding.EncodeToString(b)
		}
		err := originalRunE(cmd, args)
		target := &apierr.APIError{}
		if errors.As(err, &target) &&
			(importReq.Format == workspace.ImportFormatSource || importReq.Format == "") &&
			target.StatusCode == http.StatusBadRequest &&
			strings.Contains(target.Message, "The zip file may not be valid or may be an unsupported version.") {
			return fmt.Errorf("%w Hint: Objects imported using source format are assumed to be zip encoded by default. Please specify a language using the --language flag if you are trying to import a notebook instead", err)
		}
		return err
	}

}

func init() {
	listOverrides = append(listOverrides, listOverride)
	exportOverrides = append(exportOverrides, exportOverride)
	importOverrides = append(importOverrides, importOverride)
}
