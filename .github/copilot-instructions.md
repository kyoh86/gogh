# Copilot Instruction for gogh Project

This document serves as a guide to understand the design and structure of the gogh project, focusing on the architecture and dependency rules of the program.

## Stance of Copilot

It is expected that you write appropriate code as a Go professional.

## Output Language

Describe mainly in Japanese.
Write comments in the code in English.

### Project Location

gogh is managed under the repository `github.com/kyoh86/gogh`.
The Go package name may have a version suffix, so please check the `module` line in `go.mod`.

### How to Write Tests

When writing tests for private functions or methods, write private tests.
Private tests should be placed in the same package with a file name ending in `_private_test.go`.

When writing tests for exported functions or methods, place them in a file named `*_test.go` with the package name suffixed with `_test` in the same directory.
When importing the exported test package, use the name `testtarget`; for example:

```go
import (
    testtarget "github.com/kyoh86/gogh/vx/core/foo"
)
```

When writing tests for exported functions or methods, use the standard `testing` package and do not use `testify/assert` or `testify/require`.

## Purpose of the Project

gogh is a repository management tool primarily targeting GitHub, providing functionalities such as cloning, creating, deleting, and listing repositories.
It is designed to structure and place repository references locally, process them in a consistent manner, and allow users to operate easily.

## Architecture of the Project

### 1. Four Main Concerns (Layers)

#### Core Functionality Layer (core)

- Location: `core/`

- Role
    - Defines the essential operations and concepts that the program performs.
- Example in this project
    - Processing repository references, defining repository operations.
- Definitions
    - Domain entity (concrete structure `Reference`, `BaseRepository` and etc.)
    - Interfaces needed for interaction with external systems (e.g., `RepositoryService`).
- Dependencies
    - None (innermost layer).
- Implements
    - Domain entities and pure logic.
    - Pure functionality that does not depend on external systems or UI.

#### Application Layer (app)

- Location: `app/`

- Role
    - Realizes use cases and coordinates core functionality.
- Example in this project
    - Repository cloning process, list retrieval logic.
- Dependencies
    - Only imports interfaces from the core layer.
- Implements
    - Combines core functionality to perform actual processing.
    - Implements the entire flow of use cases (e.g., `CloneService`).
    - Service objects, use case handlers.

#### External System Integration Layer (infra)

- Location: `infra/`

- Role
    - Integration with external systems and technical implementations.
- Example in this project
    - GitHub API calls, token storage, file operations.
- Dependencies
    - Imports core interfaces and types.
- Implements
    - API clients, data base connections, file operations.
    - Implements the interface defined in the core layer.

#### UI Layer (ui)

- Location: `ui/`

- Role
    - The method of interaction with the user.
- Example in this project
    - CLI commands, output formats.
- Features
    - Translates user requests to the application layer.
- Dependencies
    - Can import the application layer.
- Implements
    - Command line interface, output formats.
    - User interaction methods (e.g., `cli` commands).

### 2. Practical Dependency Rules

1. Core Functionality Layer Dependency Rule
   - Core functionality layer does not import any other layers.
   - (Other layers can import core functionality.)
   - Reason: Core operations should not depend on application-specific processing or external integrations.

2. Application Layer Dependency Rule
    - Application layer can import core functionality.
    - Does not import UI or external system layers.
    - Reason: Use cases use core functionality but should not depend on specific implementations or UI.

3. External System Integration Layer Dependency Rule
    - External system integration layer can import core functionality.
    - Does not import UI or application layers.
    - Reason: External integrations implement the interfaces defined in the core layer but do not need to know the details of UI or use cases.

4. UI Layer Dependency Rule
    - Can import the application layer.
    - Reason: UI calls and executes use cases.

### 3. Interface-Based Integration

- Core layer defines interfaces.
- External system integration layer implements the interfaces.
- Application layer uses the interfaces to call the external system.

#### Example of Interface Definition in Core Layer

```go
// core/repository/service.go
type RepositoryService interface {
   // Basic operations for repositories
   Get(name string) (Repository, error)
   List() ([]Repository, error)
}
```

#### Example of usecase in Application Layer

```go
// app/clone/service.go
type CloneService struct {
   repoService core.RepositoryService
}

func (s *CloneService) CloneRepository(name string) error {
   repo, err := s.repoService.Get(name)
   if err != nil {
       return err
   }
   // Clone processing implementation
}
```

#### Example of Implementation in External Integration Layer

```go
// infra/github/client.go
type githubClient struct {...}

// Implements the interface from the core layer
func (c *githubClient) Get(name string) (core.Repository, error) {
   // Implementation using GitHub API
}
```

#### Main Function (or Dependency Injection Container) for Dependency Injection

```go
// main.go
func main() {
    // 1. Create implementation of the infrastructure layer
    tokenStore := config.NewTokenStore()
    githubClient := infra.NewGithubClient(tokenStore)

    // 2. Create application layer service and inject the infrastructure implementation
    cloneService := clone.NewService(githubService)

    // 3. Create UI layer command and inject the application service
    cloneCommand := commands.NewCloneCommand(cloneService)

    // 4. Execute the application
    rootCmd.AddCommand(cloneCommand)
    rootCmd.Execute()
}
```

#### Factory for Dependency Injection

```go
// app/clone/factory.go
func NewCloneServiceWithGitHub(tokenStore core.TokenStore) *Service {
    githubService := github.NewRepositoryService(tokenStore)
    return NewService(githubService)
}
```

#### Example of Dependency Injection

| Layer             | Depends on                 | Implementation Injected | Injection Method                        |
|-------------------|----------------------------|-------------------------|-----------------------------------------|
| UI Layer          | Application Layer Services | `app.CloneService`      | main function/initilization of commands |
| Application Layer | Core Layer Interfaces      | `infra.GitHubService`   | main function/factory                   |
| External System   | Core Layer Interfaces      | - (itself)              | -                                       |
| Core Layer        | -                          | -                       | -                                       |

### 4. Recommended Directory Structure

```
gogh/
├── core/              # Core functionality - the essence of the program
│   ├── repository/    # Definitions related to repositories
│   └── auth/          # Definitions related to authentication
│
├── app/               # Application layer - use case implementations
│   ├── clone/         # Clone use case
│   └── list/          # List use case
│
├── infra/             # External system integration - technical implementations
│   ├── github/        # GitHub API client
│   └── storage/          # Data storage implementation
│
├── ui/                # User interface - interaction with the user
│   └── cli/           # CLI command implementations
```

### 5. Case Study

Example of Clone Operation

1. UI Layer: The `clone` command receives the repository name from the user.
2. Application Layer: The `CloneService` uses the `RepositoryService` to get the repository details.
3. External System Integration Layer: The `githubClient` implements the `RepositoryService` interface to fetch repository details from GitHub.
4. Core Layer: The core layer defines the `RepositoryService` interface and the `Repository` struct.

### 6. Two Types of Core Layer Elements

To clarify the core layer, we can categorize its elements into two types:

#### 6-1. Interfaces

- Interfaces that are implemented by the external system layer.
- Example: `RepositoryService`, `RemoteService`, etc.
- If the operation is dependent on the technical implementation (GitHub API, file system, etc.), it should be defined in the core layer.

#### 6-2. Comcrete Implementations

- Concrete implementations of the interfaces.
- Example: `githubClient`, `fileSystemClient`, etc.
- Domain logics that are not dependent on external systems should be defined in the core layer.

#### 6-3. Example

```go
// core/repository/repository.go
type reference struct { // Concrete type
    host string
    owner string
    name string
}

func (r *reference) Host() string { return r.host } // Implementation

// core/repository/parser.go
type Parse(s string) (Reference, error) { // Implementation
    // Parsing logic
    return &reference{...}, nil
}

// core/repository/service.go - External system layer implements this
type RepositoryService interface { // Interface only
    Get(ref Reference) (Repository, error)
    // ...
}
```

#### 6.4. Criteria for Determining

Whether to implement in the core layer or just define the interface:

1. Whether the operation is dependent on the technical implementation (e.g., GitHub API, file system)
    - Yes → Define only the interface (implementation in the external system layer).
    - No → Implement in the core layer.
2. Whether the operation is pure domain logic:
    - Yes → Implement in the core layer.
    - No → Define only the interface (implementation in the external system layer).

## Design of a Core of gogh

Core rules of gogh are as follows

### 1. Conceptual Model of Repository

- Definition of repository reference: Interpretation rules in the format `github.com/owner/repo`
- Components of a repository: Basic attributes such as name, owner, and URL format
- Relationships between repositories: Fork and origin, parent-child relationships, etc.

### 2. Structuring and Organizing Rules

- Grouping of repositories: How to manage multiple repositories
- Structure under the root: Path structure `<root>/<host>/<owner>/<name>`
- Namespace management: Rules to avoid name collisions

### 3. Basic Concept of Operations

- Definition of clone: What does it mean to duplicate from remote to local
- Abstract concept of search: How to search for repositories

### 4. Relationship between User and Repository

- User ownership and access rights: Relationship between user and repository
- Permission model: Concept of read/write permissions
