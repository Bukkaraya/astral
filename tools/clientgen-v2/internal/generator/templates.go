package generator

const ClientTemplate = `{{if .GenerateDirective}}//go:generate go run ../tools/clientgen-v2/cmd/main.go generate -d {{.PackageName}} -o client.go

{{end}}package client

import (
	"context"
	"fmt"

	"{{.ModulePath}}/common"
	"{{.ModulePath}}/{{.PackageName}}/workflows"
)

{{range .WorkflowMethods}}
// {{.Name}}Request represents a request to {{.Name | toDescription | toLower}}
type {{.Name}}Request struct {
	{{.InputType | extractIDField}} string ` + "`json:\"{{.InputType | extractIDField | toJSONTag}}\"`" + `
}
{{end}}

// Client provides methods to execute {{.PackageName}} workflows
type Client interface {
{{range .WorkflowMethods}}	{{.Name}}(ctx context.Context, req {{.Name}}Request) (*common.WorkflowResult, error)
{{end}}}

// {{.PackageName}}Client implements the Client interface
type {{.PackageName}}Client struct {
	executor common.TemporalExecutor
}

// NewClient creates a new {{.PackageName}} workflow client
func NewClient(executor common.TemporalExecutor) Client {
	return &{{.PackageName}}Client{
		executor: executor,
	}
}

{{range .WorkflowMethods}}
// {{.Name}} starts a {{.Name}} workflow
func (c *{{$.PackageName}}Client) {{.Name}}(ctx context.Context, req {{.Name}}Request) (*common.WorkflowResult, error) {
	workflowInput := workflows.{{.InputType}}{{"{"}}{{.InputType | extractIDField}}: req.{{.InputType | extractIDField}}{{"}"}}
	
	return c.executor.ExecuteWorkflow(ctx, common.WorkflowExecutionParams{
		WorkflowType:     "{{.Name}}.v1",
		WorkflowIDPrefix: "{{.Name | toKebabCase}}",
		WorkflowInput:    workflowInput,
		SuccessMessage:   fmt.Sprintf("{{.Name | toDescription}} workflow started for {{.InputType | extractEntity}} %s", req.{{.InputType | extractIDField}}),
	})
}
{{end}}`