# Architecture of the Project

## 1. Four Main Concerns (Layers)

### Core Functionality Layer (core)

- **Role**: Defines the essential operations and concepts that the program performs.
- **Example in this project**: Processing repository references, defining repository operations.
- **Content**: Data structures, interfaces, basic operations.
- **Features**: Pure functionality that does not depend on external systems or UI.

### Application Layer (app)

- **Role**: Realizes use cases and coordinates core functionality.
- **Example in this project**: Repository cloning process, list retrieval logic.
- **Content**: Service objects, use case handlers.
- **Features**: Combines core functionality to perform actual processing.

### External System Integration Layer (infra)

- **Role**: Integration with external systems and technical implementations.
- **Example in this project**: GitHub API calls, token storage, file operations.
- **Content**: API clients, data storage, file handling.
- **Features**: Implementation details that are not part of the core functionality.

### UI Layer (ui)

- **Role**: The method of interaction with the user.
- **Example in this project**: CLI commands, output formats.
- **Content**: Command implementations, display logic, input handling.
- **Features**: Translates user requests to the application layer.

## 2. Practical Dependency Rules

1. **Core Functionality Layer Dependency Rule**
   - Core functionality layer does not **import** any other layers.
   - (Other layers can import core functionality.)
   - **Reason**: Core operations should not depend on application-specific processing or external integrations.

2. **Application Layer Dependency Rule**
    - Application layer can **import** core functionality.
    - Does not import UI or external system layers.
    - **Reason**: Use cases use core functionality but should not depend on specific implementations or UI.

3. **External System Integration Layer Dependency Rule**
    - External system integration layer can **import** core functionality.
    - Does not import UI or application layers.
    - **Reason**: External integrations implement the interfaces defined in the core layer but do not need to know the details of UI or use cases.

4. **UI Layer Dependency Rule**
    - Can **import** the application layer.
    - **Reason**: UI calls and executes use cases.

## 3. Interface-Based Integration

1. **Interface Ownership**
   - The user layer defines the interface.
   - **Example**: The interface for remote repositories is defined in the core layer.

2. **Example of Interface Definition in Core Layer**
   ```go
   // core/repository/service.go
   type RepositoryService interface {
       // Basic operations for repositories
       Get(name string) (Repository, error)
       List() ([]Repository, error)
   }
   ```

3. **Example of usecase in Application Layer**
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

4. **Example of Implementation in External Integration Layer**
   ```go
   // infra/github/client.go
   type githubClient struct {...}
   
   // Implements the interface from the core layer
   func (c *githubClient) Get(name string) (core.Repository, error) {
       // Implementation using GitHub API
   }
   ```

## 4. Recommended Directory Structure

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

## 5. Case Study

Example of Clone Operation

1. **UI Layer**: The `clone` command receives the repository name from the user.
2. **Application Layer**: The `CloneService` uses the `RepositoryService` to get the repository details.
3. **External System Integration Layer**: The `githubClient` implements the `RepositoryService` interface to fetch repository details from GitHub.
4. **Core Layer**: The core layer defines the `RepositoryService` interface and the `Repository` struct.
