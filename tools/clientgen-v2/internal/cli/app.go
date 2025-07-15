package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"clientgen-v2/internal/config"
	"clientgen-v2/internal/generator"
	"clientgen-v2/internal/models"
	"clientgen-v2/internal/parser"
)

// App represents the CLI application
type App struct {
	config *config.Config
}

// NewApp creates a new CLI application
func NewApp() *App {
	return &App{}
}

// Execute executes the CLI application
func (a *App) Execute() error {
	if len(os.Args) < 2 {
		return a.showHelp()
	}

	command := os.Args[1]
	switch command {
	case "generate":
		return a.runGenerate()
	case "help", "-h", "--help":
		return a.showHelp()
	case "version":
		return a.showVersion()
	default:
		return fmt.Errorf("unknown command: %s", command)
	}
}

// runGenerate runs the generate command
func (a *App) runGenerate() error {
	// Parse flags
	flags, err := a.parseGenerateFlags()
	if err != nil {
		return err
	}

	// Load configuration
	cfg, err := config.Load(flags.ConfigFile)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Override config with flags
	if flags.Domain != "" {
		cfg.Generator.PackageName = flags.Domain
	}
	if flags.OutputFile != "" {
		cfg.Output.Directory = filepath.Dir(flags.OutputFile)
	}
	if flags.ModulePath != "" {
		cfg.Generator.ModulePath = flags.ModulePath
	}

	cfg.SetDefaults()
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	a.config = cfg

	// Parse interface file
	interfaceFile := flags.InterfaceFile
	if interfaceFile == "" {
		interfaceFile = "interfaces.go"
	}

	workflowParser := parser.NewASTParser(&cfg.Parser)
	workflows, err := workflowParser.ExtractWorkflows(interfaceFile)
	if err != nil {
		return fmt.Errorf("failed to parse workflows: %w", err)
	}

	if err := workflowParser.Validate(workflows); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Prepare template data
	templateData := &models.TemplateData{
		PackageName:       cfg.Generator.PackageName,
		ModulePath:        cfg.Generator.ModulePath,
		WorkflowMethods:   workflows,
		GenerateDirective: cfg.Generator.IncludeDirective,
		Imports:          []string{},
	}

	// Generate client code
	gen := generator.NewClientGenerator(&cfg.Generator)
	
	// Load template
	if err := gen.SetTemplate(createTemplate()); err != nil {
		return fmt.Errorf("failed to set template: %w", err)
	}

	generated, err := gen.Generate(templateData)
	if err != nil {
		return fmt.Errorf("failed to generate client: %w", err)
	}

	// Write output file
	outputFile := flags.OutputFile
	if outputFile == "" {
		outputFile = "client.go"
	}

	if err := os.WriteFile(outputFile, generated, 0644); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	fmt.Printf("Generated %s client code in %s (%d methods)\n", 
		cfg.Generator.PackageName, outputFile, len(workflows))

	return nil
}

// GenerateFlags represents generate command flags
type GenerateFlags struct {
	Domain        string
	OutputFile    string
	InterfaceFile string
	ModulePath    string
	ConfigFile    string
}

// parseGenerateFlags parses generate command flags
func (a *App) parseGenerateFlags() (*GenerateFlags, error) {
	flags := &GenerateFlags{}
	
	args := os.Args[2:] // Skip program name and command
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-d", "--domain":
			if i+1 >= len(args) {
				return nil, fmt.Errorf("missing value for %s", args[i])
			}
			flags.Domain = args[i+1]
			i++
		case "-o", "--output":
			if i+1 >= len(args) {
				return nil, fmt.Errorf("missing value for %s", args[i])
			}
			flags.OutputFile = args[i+1]
			i++
		case "-i", "--interface":
			if i+1 >= len(args) {
				return nil, fmt.Errorf("missing value for %s", args[i])
			}
			flags.InterfaceFile = args[i+1]
			i++
		case "-m", "--module":
			if i+1 >= len(args) {
				return nil, fmt.Errorf("missing value for %s", args[i])
			}
			flags.ModulePath = args[i+1]
			i++
		case "-c", "--config":
			if i+1 >= len(args) {
				return nil, fmt.Errorf("missing value for %s", args[i])
			}
			flags.ConfigFile = args[i+1]
			i++
		default:
			return nil, fmt.Errorf("unknown flag: %s", args[i])
		}
	}

	return flags, nil
}

// showHelp shows help information
func (a *App) showHelp() error {
	fmt.Println(`Temporal Workflow Client Generator v2.0

USAGE:
    clientgen generate [OPTIONS]

COMMANDS:
    generate    Generate workflow client code
    help        Show this help message
    version     Show version information

GENERATE OPTIONS:
    -d, --domain STRING         Domain name (required)
    -o, --output STRING         Output file path (default: client.go)
    -i, --interface STRING      Interface file path (default: interfaces.go)
    -m, --module STRING         Module path (default: simple-temporal-workflow)
    -c, --config STRING         Configuration file path

EXAMPLES:
    clientgen generate -d order -o client.go
    clientgen generate -d payment -i payment_interfaces.go -o payment_client.go
    clientgen generate -d shipping -c config.yaml

For more information, visit: https://github.com/your-org/temporal-client-generator`)

	return nil
}

// showVersion shows version information
func (a *App) showVersion() error {
	fmt.Println("Temporal Workflow Client Generator v2.0.0")
	return nil
}

// createTemplate creates and loads the template
func createTemplate() generator.Template {
	goTemplate := generator.NewGoTemplate()
	err := goTemplate.Load(generator.ClientTemplate)
	if err != nil {
		panic(fmt.Sprintf("Failed to load template: %v", err))
	}
	return goTemplate
}