package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"

	"clientgen-v2/internal/config"
	"clientgen-v2/internal/models"
)

// ASTParser implements Parser using Go AST
type ASTParser struct {
	fileSet *token.FileSet
	config  *config.ParserConfig
}

// NewASTParser creates a new AST parser
func NewASTParser(config *config.ParserConfig) *ASTParser {
	return &ASTParser{
		fileSet: token.NewFileSet(),
		config:  config,
	}
}

// ParseFile parses a Go file and extracts workflow information
func (p *ASTParser) ParseFile(filePath string) (*models.ParsedFile, error) {
	src, err := parser.ParseFile(p.fileSet, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file %s: %w", filePath, err)
	}

	workflows, err := p.extractWorkflowsFromAST(src)
	if err != nil {
		return nil, fmt.Errorf("failed to extract workflows from %s: %w", filePath, err)
	}

	return &models.ParsedFile{
		Path:      filePath,
		Package:   src.Name.Name,
		Workflows: workflows,
		FileSet:   p.fileSet,
	}, nil
}

// ExtractWorkflows extracts workflow methods from a file
func (p *ASTParser) ExtractWorkflows(filePath string) ([]*models.WorkflowMethod, error) {
	parsedFile, err := p.ParseFile(filePath)
	if err != nil {
		return nil, err
	}
	return parsedFile.Workflows, nil
}

// Validate validates workflow methods
func (p *ASTParser) Validate(methods []*models.WorkflowMethod) error {
	var errors models.ValidationErrors

	for _, method := range methods {
		if method.Name == "" {
			errors.Add("name", "workflow method name cannot be empty")
		}

		if method.InputType == "" {
			errors.Add("input_type", fmt.Sprintf("workflow method %s must have an input type", method.Name))
		}

		// Output type is optional - workflows can return just error
		// if method.OutputType == "" {
		//	errors.Add("output_type", fmt.Sprintf("workflow method %s must have an output type", method.Name))
		// }

		if method.Signature == nil {
			errors.Add("signature", fmt.Sprintf("workflow method %s must have a signature", method.Name))
			continue
		}

		// Validate signature has context as first parameter
		if len(method.Signature.Parameters) < 2 {
			errors.Add("parameters", fmt.Sprintf("workflow method %s must have at least context and input parameters", method.Name))
		}

		// Validate first parameter is context
		if len(method.Signature.Parameters) > 0 {
			firstParam := method.Signature.Parameters[0]
			if !firstParam.IsCtx {
				errors.Add("context", fmt.Sprintf("workflow method %s first parameter must be workflow.Context", method.Name))
			}
		}
	}

	if !errors.IsEmpty() {
		return errors
	}

	return nil
}

// extractWorkflowsFromAST extracts workflow methods from AST
func (p *ASTParser) extractWorkflowsFromAST(file *ast.File) ([]*models.WorkflowMethod, error) {
	var workflows []*models.WorkflowMethod

	ast.Inspect(file, func(n ast.Node) bool {
		if typeSpec, ok := n.(*ast.TypeSpec); ok {
			// Look for interface named "Workflows"
			if typeSpec.Name.Name == "Workflows" {
				if iface, ok := typeSpec.Type.(*ast.InterfaceType); ok {
					for _, method := range iface.Methods.List {
						if len(method.Names) > 0 && method.Type != nil {
							if funcType, ok := method.Type.(*ast.FuncType); ok {
								workflow := p.parseWorkflowMethod(method.Names[0].Name, funcType, method.Doc)
								if workflow != nil && p.shouldIncludeMethod(workflow.Name) {
									workflows = append(workflows, workflow)
								}
							}
						}
					}
				}
			}
		}
		return true
	})

	return workflows, nil
}

// parseWorkflowMethod parses a workflow method from AST
func (p *ASTParser) parseWorkflowMethod(name string, funcType *ast.FuncType, doc *ast.CommentGroup) *models.WorkflowMethod {
	signature := p.parseMethodSignature(name, funcType)
	
	// Extract input and output types
	inputType := ""
	outputType := ""
	
	if funcType.Params != nil && len(funcType.Params.List) >= 2 {
		// Second parameter (after context) is the input type
		param := funcType.Params.List[1]
		if ident, ok := param.Type.(*ast.Ident); ok {
			inputType = ident.Name
		} else if selector, ok := param.Type.(*ast.SelectorExpr); ok {
			inputType = selector.Sel.Name
		}
	}
	
	if funcType.Results != nil && len(funcType.Results.List) >= 1 {
		// First return value is the output type
		if ident, ok := funcType.Results.List[0].Type.(*ast.Ident); ok {
			outputType = ident.Name
		}
	}
	

	return &models.WorkflowMethod{
		Name:          name,
		Signature:     signature,
		InputType:     inputType,
		OutputType:    outputType,
		Documentation: p.parseDocumentation(doc),
		Metadata:      make(map[string]string),
	}
}

// parseMethodSignature parses method signature
func (p *ASTParser) parseMethodSignature(name string, funcType *ast.FuncType) *models.MethodSignature {
	signature := &models.MethodSignature{
		Name:       name,
		Parameters: []*models.Parameter{},
		Returns:    []*models.Return{},
		IsVariadic: false,
	}

	// Parse parameters
	if funcType.Params != nil {
		for i, param := range funcType.Params.List {
			for _, paramName := range param.Names {
				parameter := &models.Parameter{
					Name:    paramName.Name,
					Type:    p.parseTypeInfo(param.Type),
					IsCtx:   i == 0, // First parameter is context
					IsInput: i == 1, // Second parameter is input
				}
				signature.Parameters = append(signature.Parameters, parameter)
			}
		}
	}

	// Parse return values
	if funcType.Results != nil {
		for _, result := range funcType.Results.List {
			returnVal := &models.Return{
				Type:    p.parseTypeInfo(result.Type),
				IsError: p.isErrorType(result.Type),
			}
			signature.Returns = append(signature.Returns, returnVal)
		}
	}

	return signature
}

// parseTypeInfo parses type information from AST
func (p *ASTParser) parseTypeInfo(expr ast.Expr) *models.TypeInfo {
	switch t := expr.(type) {
	case *ast.Ident:
		return &models.TypeInfo{
			Name: t.Name,
		}
	case *ast.SelectorExpr:
		pkg := ""
		if ident, ok := t.X.(*ast.Ident); ok {
			pkg = ident.Name
		}
		return &models.TypeInfo{
			Name:    t.Sel.Name,
			Package: pkg,
		}
	case *ast.StarExpr:
		inner := p.parseTypeInfo(t.X)
		inner.IsPointer = true
		return inner
	case *ast.ArrayType:
		inner := p.parseTypeInfo(t.Elt)
		inner.IsSlice = true
		return inner
	default:
		return &models.TypeInfo{
			Name: "unknown",
		}
	}
}

// parseDocumentation parses documentation from comment group
func (p *ASTParser) parseDocumentation(doc *ast.CommentGroup) *models.Documentation {
	if doc == nil {
		return nil
	}

	var lines []string
	for _, comment := range doc.List {
		line := strings.TrimPrefix(comment.Text, "//")
		line = strings.TrimSpace(line)
		if line != "" {
			lines = append(lines, line)
		}
	}

	if len(lines) == 0 {
		return nil
	}

	return &models.Documentation{
		Summary:     lines[0],
		Description: strings.Join(lines, " "),
	}
}

// isErrorType checks if a type is an error type
func (p *ASTParser) isErrorType(expr ast.Expr) bool {
	if ident, ok := expr.(*ast.Ident); ok {
		return ident.Name == "error"
	}
	return false
}

// shouldIncludeMethod checks if a method should be included based on patterns
func (p *ASTParser) shouldIncludeMethod(methodName string) bool {
	// Check exclude patterns first
	for _, pattern := range p.config.ExcludePatterns {
		if matched, _ := filepath.Match(pattern, methodName); matched {
			return false
		}
	}

	// Check include patterns
	for _, pattern := range p.config.WorkflowPatterns {
		if matched, _ := filepath.Match(pattern, methodName); matched {
			return true
		}
	}

	// Default: include if no patterns match
	return len(p.config.WorkflowPatterns) == 0
}