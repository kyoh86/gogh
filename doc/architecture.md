# Architecture of the Project

## 1. Three Main Concerns (Layers)

### Core Functionality Layer (core)

- **Role**: The essential operations and concepts that the program performs.
- **Example in this project**: Processing repository references, defining repository operations.
- **Content**: Data structures, interfaces, basic operations.
- **Features**: Pure functionality that does not depend on external systems or UI.

### External System Integration Layer (infra)

- **Role**: Integration with external systems and technical implementations.
- **Example in this project**: GitHub API calls, token storage, file operations.
- **Content**: API clients, data storage, file handling.
- **Features**: Implementation details that are not part of the core functionality.

### UI Layer (ui)

- **Role**: The method of interaction with the user.
- **Example in this project**: CLI commands, output formats.
- **Content**: Command implementations, display logic, input handling.
- **Features**: Translates user requests to core functionality or external integrations.

## 2. Practical Dependency Rules

1. **Core Functionality Layer Dependency Rule**
   - Does **NOT** import from UI or external system integration.
   - **Reason**: Core operations should not depend on how they are triggered or how data is obtained.

2. **External System Integration Layer Dependency Rule**
   - Can **import** from core functionality layer.
   - Does **NOT** import from UI.
   - **Reason**: GitHub integration needs to know core concepts but does not need to know CLI details.

3. **UI Layer Dependency Rule**
    - Can **import** from both core functionality and external system integration layers.
    - **Reason**: Commands need both functionality and data.

## 3. Interface-Based Integration

1. **Interface Ownership**
   - The user defines the interface.
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

3. **Example of Implementation in External Integration Layer**
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
├── infra/             # External system integration - technical implementations
│   ├── github/        # GitHub API client
│   └── storage/          # Data storage implementation
│
├── ui/                # User interface - interaction with the user
│   └── cli/           # CLI command implementations
```
