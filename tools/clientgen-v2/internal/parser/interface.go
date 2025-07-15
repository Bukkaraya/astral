package parser

import (
	"clientgen-v2/internal/models"
)

// Parser defines the interface for parsing workflow interfaces
type Parser interface {
	ParseFile(filePath string) (*models.ParsedFile, error)
	ExtractWorkflows(filePath string) ([]*models.WorkflowMethod, error)
	Validate(methods []*models.WorkflowMethod) error
}

// TypeAnalyzer analyzes Go types and extracts metadata
type TypeAnalyzer interface {
	AnalyzeMethod(method interface{}) (*models.MethodSignature, error)
	ResolveType(typeName string) (*models.TypeInfo, error)
}

// ParseOptions configures parsing behavior
type ParseOptions struct {
	IncludePrivate   bool
	WorkflowPatterns []string
	ExcludePatterns  []string
}