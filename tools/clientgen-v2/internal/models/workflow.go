package models

import (
	"fmt"
	"go/token"
)

// WorkflowMethod represents a parsed workflow method
type WorkflowMethod struct {
	Name          string            `json:"name"`
	Signature     *MethodSignature  `json:"signature"`
	Documentation *Documentation    `json:"documentation,omitempty"`
	InputType     string            `json:"input_type"`
	OutputType    string            `json:"output_type"`
	Metadata      map[string]string `json:"metadata,omitempty"`
}

// MethodSignature contains method signature details
type MethodSignature struct {
	Name       string       `json:"name"`
	Parameters []*Parameter `json:"parameters"`
	Returns    []*Return    `json:"returns"`
	IsVariadic bool         `json:"is_variadic"`
}

// Parameter represents a method parameter
type Parameter struct {
	Name    string    `json:"name"`
	Type    *TypeInfo `json:"type"`
	IsCtx   bool      `json:"is_ctx"`
	IsInput bool      `json:"is_input"`
}

// Return represents a method return value
type Return struct {
	Type    *TypeInfo `json:"type"`
	IsError bool      `json:"is_error"`
}

// TypeInfo contains type information
type TypeInfo struct {
	Name      string `json:"name"`
	Package   string `json:"package,omitempty"`
	IsPointer bool   `json:"is_pointer"`
	IsSlice   bool   `json:"is_slice"`
}

// Documentation represents method documentation
type Documentation struct {
	Summary     string `json:"summary"`
	Description string `json:"description"`
}

// ParsedFile represents a parsed Go file
type ParsedFile struct {
	Path      string             `json:"path"`
	Package   string             `json:"package"`
	Workflows []*WorkflowMethod  `json:"workflows"`
	Position  token.Pos          `json:"-"`
	FileSet   *token.FileSet     `json:"-"`
}

// TemplateData represents data for template rendering
type TemplateData struct {
	PackageName     string             `json:"package_name"`
	ModulePath      string             `json:"module_path"`
	WorkflowMethods []*WorkflowMethod  `json:"workflow_methods"`
	GenerateDirective bool             `json:"generate_directive"`
	Imports         []string           `json:"imports"`
	Metadata        *GenerationMetadata `json:"metadata"`
}

// GenerationMetadata contains generation metadata
type GenerationMetadata struct {
	GeneratedAt string `json:"generated_at"`
	Version     string `json:"version"`
	Source      string `json:"source"`
	Generator   string `json:"generator"`
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Value   string `json:"value,omitempty"`
}

// Error implements the error interface
func (v *ValidationError) Error() string {
	return v.Message
}

// ValidationErrors represents multiple validation errors
type ValidationErrors []*ValidationError

// Error implements the error interface
func (v ValidationErrors) Error() string {
	if len(v) == 0 {
		return "no validation errors"
	}
	if len(v) == 1 {
		return v[0].Error()
	}
	return fmt.Sprintf("%s (and %d more errors)", v[0].Error(), len(v)-1)
}

// Add adds a validation error
func (v *ValidationErrors) Add(field, message string) {
	*v = append(*v, &ValidationError{
		Field:   field,
		Message: message,
	})
}

// IsEmpty returns true if there are no errors
func (v ValidationErrors) IsEmpty() bool {
	return len(v) == 0
}