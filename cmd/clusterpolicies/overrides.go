package clusterpolicies

import "github.com/databricks/bricks/lib/ui"

func init() {
	listCmd.Annotations["template"] = ui.Heredoc(`
	{{range .}}{{.PolicyId | green}}	{{.Name}}
	{{end}}`)
}
