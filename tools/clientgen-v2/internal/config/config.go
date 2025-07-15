package config

import (
	"fmt"
	"os"
)

// Config represents the main configuration
type Config struct {
	Parser    ParserConfig    `yaml:"parser"`
	Generator GeneratorConfig `yaml:"generator"`
	Template  TemplateConfig  `yaml:"template"`
	Output    OutputConfig    `yaml:"output"`
}

// ParserConfig configures the interface parser
type ParserConfig struct {
	IncludePrivate   bool     `yaml:"include_private"`
	WorkflowPatterns []string `yaml:"workflow_patterns"`
	ActivityPatterns []string `yaml:"activity_patterns"`
	ExcludePatterns  []string `yaml:"exclude_patterns"`
}

// GeneratorConfig configures code generation
type GeneratorConfig struct {
	PackageName       string `yaml:"package_name"`
	ModulePath        string `yaml:"module_path"`
	ClientType        string `yaml:"client_type"`
	IncludeDirective  bool   `yaml:"include_directive"`
	GenerateTests     bool   `yaml:"generate_tests"`
}

// TemplateConfig configures template processing
type TemplateConfig struct {
	Directory   string            `yaml:"directory"`
	CustomFuncs map[string]string `yaml:"custom_funcs"`
	Variables   map[string]string `yaml:"variables"`
}

// OutputConfig configures output generation
type OutputConfig struct {
	Directory  string `yaml:"directory"`
	FileMode   string `yaml:"file_mode"`
	FormatCode bool   `yaml:"format_code"`
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		Parser: ParserConfig{
			IncludePrivate: false,
			WorkflowPatterns: []string{}, // Empty = include all methods
			ActivityPatterns: []string{"*Activity*", "*activities*"},
			ExcludePatterns:  []string{"*Test*", "*Mock*"},
		},
		Generator: GeneratorConfig{
			PackageName:      "",
			ModulePath:       "simple-temporal-workflow",
			ClientType:       "standard",
			IncludeDirective: true,
			GenerateTests:    false,
		},
		Template: TemplateConfig{
			Directory:   "templates",
			CustomFuncs: make(map[string]string),
			Variables:   make(map[string]string),
		},
		Output: OutputConfig{
			Directory:  ".",
			FileMode:   "0644",
			FormatCode: true,
		},
	}
}

// Load loads configuration from file
func Load(configPath string) (*Config, error) {
	config := DefaultConfig()
	
	if configPath == "" {
		return config, nil
	}
	
	_, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return config, nil // Use default if file doesn't exist
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	
	// For now, just use the default config
	// TODO: Implement JSON/YAML parsing when needed
	
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}
	
	return config, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Generator.PackageName == "" {
		return fmt.Errorf("generator.package_name is required")
	}
	
	if c.Generator.ModulePath == "" {
		return fmt.Errorf("generator.module_path is required")
	}
	
	if c.Output.Directory == "" {
		return fmt.Errorf("output.directory is required")
	}
	
	// Validate output directory exists or can be created
	if err := os.MkdirAll(c.Output.Directory, 0755); err != nil {
		return fmt.Errorf("cannot create output directory: %w", err)
	}
	
	return nil
}

// SetDefaults sets default values for missing fields
func (c *Config) SetDefaults() {
	defaults := DefaultConfig()
	
	if len(c.Parser.WorkflowPatterns) == 0 {
		c.Parser.WorkflowPatterns = defaults.Parser.WorkflowPatterns
	}
	
	if len(c.Parser.ActivityPatterns) == 0 {
		c.Parser.ActivityPatterns = defaults.Parser.ActivityPatterns
	}
	
	if c.Generator.ModulePath == "" {
		c.Generator.ModulePath = defaults.Generator.ModulePath
	}
	
	if c.Generator.ClientType == "" {
		c.Generator.ClientType = defaults.Generator.ClientType
	}
	
	if c.Output.FileMode == "" {
		c.Output.FileMode = defaults.Output.FileMode
	}
}