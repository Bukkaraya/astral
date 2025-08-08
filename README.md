# Astral

A Rails-like framework for Temporal workflows in Go. Build robust distributed applications with automatic code generation, type-safe clients, and convention-over-configuration patterns.

## Philosophy

Astral brings Rails' "convention over configuration" philosophy to Temporal workflows. Instead of writing boilerplate client code, activity definitions, and worker setup, Astral generates everything from your workflow definitions.

**Key Principles:**
- **Convention over Configuration** - Standard project structure and naming conventions eliminate boilerplate
- **Code Generation** - Type-safe clients and workers generated automatically from workflow definitions
- **Domain-Driven Design** - Organize workflows by business domains with clear bounded contexts
- **Developer Experience** - Focus on business logic, not infrastructure plumbing
- **Opinionated Defaults** - Sensible defaults that can be overridden when needed

## Design

Astral follows a structured approach combining Rails MVC with Domain-Driven Design:

- **Workflows** - Your business logic organized by domain (like Controllers)
- **Activities** - Reusable units of work within bounded contexts (like Models)  
- **Clients** - Auto-generated, type-safe interfaces per domain (like Views)
- **Workers** - Background processors that execute workflows
- **Domains** - Clear separation of business concerns with independent deployment

Projects are organized by business domains (e.g., `order/`, `payment/`, `inventory/`) with each domain containing its own workflows, activities, and generated clients. This promotes modularity and allows teams to own specific business capabilities.