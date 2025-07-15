# Clientgen v2 - Clean Architecture Implementation

This is a refactored version of the Temporal workflow client generator following clean architecture principles.

## Architecture Overview

The refactored generator implements the following design patterns and principles:

### 🏗️ **Clean Architecture**
- **Separation of Concerns**: Each package has a single responsibility
- **Dependency Inversion**: High-level modules don't depend on low-level modules
- **Interface Segregation**: Small, focused interfaces
- **Single Responsibility**: Each class/function has one reason to change

### 📦 **Package Structure**
```
clientgen-v2/
├── cmd/main.go              # CLI entry point
├── internal/
│   ├── config/              # Configuration management
│   │   └── config.go
│   ├── models/              # Domain models
│   │   └── workflow.go
│   ├── parser/              # Interface parsing
│   │   ├── interface.go
│   │   └── ast_parser.go
│   ├── generator/           # Code generation
│   │   ├── interface.go
│   │   ├── client_generator.go
│   │   └── templates.go
│   └── cli/                 # CLI handling
│       └── app.go
└── go.mod
```

### 🎯 **Design Patterns Implemented**

1. **Strategy Pattern** - Different generation strategies (client vs worker)
2. **Template Method Pattern** - Common generation workflow 
3. **Builder Pattern** - Configuration and template data building
4. **Factory Pattern** - Creating appropriate generators
5. **Dependency Injection** - Interface-based dependencies

### ✨ **Key Improvements**

#### **Before (Original)**
- ❌ 299-line monolithic main.go
- ❌ Mixed responsibilities in single functions
- ❌ Hard-coded templates as strings
- ❌ No separation of concerns
- ❌ Difficult to test and extend

#### **After (Refactored)**
- ✅ Clean package structure with focused responsibilities
- ✅ Interface-based design for testability
- ✅ Configurable template system
- ✅ Comprehensive error handling
- ✅ Extensible architecture for new features

### 🔧 **Components**

#### **Parser Package**
- `Parser` interface for different parsing strategies
- `ASTParser` implementation using Go AST
- `TypeAnalyzer` for type information extraction
- Comprehensive validation with structured errors

#### **Generator Package**
- `Generator` interface for different generation targets
- `ClientGenerator` for workflow client generation
- `Template` interface for different template engines
- `Formatter` interface for code formatting

#### **Configuration Package**
- Type-safe configuration structures
- Support for multiple configuration sources
- Validation and default value handling
- Environment-specific overrides

#### **Models Package**
- Clean domain models for workflow metadata
- Template data structures
- Validation error types
- Generation metadata

### 🚀 **Usage Examples**

#### **Basic Usage**
```bash
# Generate client for order domain
./clientgen generate -d order -o client.go

# Generate with custom config
./clientgen generate -d payment -c config.yaml -o payment_client.go
```

#### **Programmatic Usage**
```go
// Create configuration
config := config.DefaultConfig()
config.Generator.PackageName = "order"

// Parse workflows
parser := parser.NewASTParser(&config.Parser)
workflows, err := parser.ExtractWorkflows("interfaces.go")

// Generate client
generator := generator.NewClientGenerator(&config.Generator)
data := &models.TemplateData{
    PackageName:     "order",
    WorkflowMethods: workflows,
}
code, err := generator.Generate(data)
```

### 🧪 **Testing Strategy**

The clean architecture enables comprehensive testing:

1. **Unit Tests** - Each component tested in isolation
2. **Integration Tests** - End-to-end workflow testing
3. **Mock Objects** - Interface-based mocking
4. **Test Data** - Comprehensive test scenarios

### 🔮 **Future Extensions**

The architecture supports easy extension:

1. **New Parsers** - Add support for different input formats
2. **New Generators** - Support for worker, activity, or test generation
3. **New Templates** - Custom templates for different frameworks
4. **Plugins** - Plugin system for custom transformations

## Design Benefits

### **Maintainability**
- Clear separation of concerns
- Easy to understand and modify
- Consistent error handling
- Comprehensive documentation

### **Testability**
- Interface-based design
- Dependency injection
- Mock-friendly architecture
- Isolated components

### **Extensibility**
- Strategy pattern for new features
- Plugin-ready architecture
- Configuration-driven behavior
- Template-based customization

### **Reliability**
- Comprehensive validation
- Structured error handling
- Type-safe operations
- Graceful failure handling

This refactored version transforms the original monolithic script into a well-architected, maintainable, and extensible code generation tool following Go best practices and clean architecture principles.