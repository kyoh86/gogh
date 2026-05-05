// Package flags provides custom flag types and helpers for CLI command flags.
//
// This package contains type-safe flag implementations that follow the LocationFormat pattern,
// ensuring consistency across all CLI flags. Each flag type implements the pflag.Value interface
// and provides:
//
//   - A custom type (e.g., RepositoryStructure, LocationFormat, WorktreeFormat)
//   - pflag.Value interface implementation (String, Set, Type methods)
//   - A helper function (e.g., StructureFlag, LocationFormatFlag) for flag registration
//   - A completion function (e.g., CompleteStructure, CompleteLocationFormat) for shell autocompletion
//   - Short and long usage constants (e.g., StructureShortUsage, StructureLongUsage)
//
// # Architecture
//
// Flag types are defined in this UI layer package, while app/config uses plain string types
// for YAML/TOML serialization. This separation ensures:
//
//   - UI layer owns flag parsing and validation logic
//   - App layer remains independent of UI concerns
//   - Configuration files use simple strings for portability
//
// # Pattern Example
//
// To add a new flag type following this pattern:
//
//  1. Define the type with pflag.Value implementation:
//     type MyFlag string
//     func (f *MyFlag) Set(v string) error { ... }
//     func (f MyFlag) String() string { ... }
//     func (f MyFlag) Type() string { return "string" }
//
//  2. Create flag registration helper:
//     func MyFlagFlag(cmd *cobra.Command, flag *MyFlag, defaultValue string) error { ... }
//
//  3. Add shell completion function:
//     func CompleteMyFlag(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) { ... }
//
//  4. Define usage constants:
//     const MyFlagShortUsage = "..."
//     const MyFlagLongUsage = "..."
//
//  5. Use app/config with string type for YAML/TOML persistence
//     type Config struct {
//         MyFlag string `yaml:"myFlag"`
//     }
//
// # Existing Flag Types
//
//   - LocationFormat: Controls repository path output format (path, full-path, json, fields)
//   - RepositoryStructure: Controls repository storage structure (worktree, normal)
//   - WorktreeFormat: Controls worktree display format (default, full-path, json, fields)
//   - RepositoryFormat: Controls repository list output format (table, ref, url, json)
package flags
