# Design of a Core of gogh

Core rules of gogh are as follows

## 1. Conceptual Model of Repository

- Definition of repository reference: Interpretation rules in the format `github.com/owner/repo`
- Components of a repository: Basic attributes such as name, owner, and URL format
- Relationships between repositories: Fork and origin, parent-child relationships, etc.

## 2. Structuring and Organizing Rules

- Grouping of repositories: How to manage multiple repositories
- Structure under the root: Path structure `<root>/<host>/<owner>/<name>`
- Namespace management: Rules to avoid name collisions

## 3. Basic Concept of Operations

- Definition of clone: What does it mean to duplicate from remote to local
- Abstract concept of search: How to search for repositories

## 4. Relationship between User and Repository

- User ownership and access rights: Relationship between user and repository
- Permission model: Concept of read/write permissions
