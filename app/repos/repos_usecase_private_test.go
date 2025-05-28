package repos

import (
	"testing"

	"github.com/kyoh86/gogh/v4/core/hosting"
)

func TestConvertOpts(t *testing.T) {
	tests := []struct {
		name    string
		input   Options
		want    hosting.ListRepositoryOptions
		wantErr bool
	}{
		{
			name:  "default options",
			input: Options{},
			want: hosting.ListRepositoryOptions{
				Limit: 30, // Default limit
			},
			wantErr: false,
		},
		{
			name: "custom limit",
			input: Options{
				Limit: 50,
			},
			want: hosting.ListRepositoryOptions{
				Limit: 50,
			},
			wantErr: false,
		},
		{
			name: "no limit",
			input: Options{
				Limit: -1,
			},
			want: hosting.ListRepositoryOptions{
				Limit: 0, // No limit
			},
			wantErr: false,
		},
		{
			name: "privacy setting",
			input: Options{
				Privacy: "private",
			},
			want: hosting.ListRepositoryOptions{
				Limit:   30,
				Privacy: hosting.RepositoryPrivacyPrivate,
			},
			wantErr: false,
		},
		{
			name: "invalid privacy",
			input: Options{
				Privacy: "invalid",
			},
			wantErr: true,
		},
		{
			name: "fork setting",
			input: Options{
				Fork: "forked",
			},
			want: hosting.ListRepositoryOptions{
				Limit:  30,
				IsFork: hosting.TristateTrue,
			},
			wantErr: false,
		},
		{
			name: "invalid fork",
			input: Options{
				Fork: "invalid",
			},
			wantErr: true,
		},
		{
			name: "archive setting",
			input: Options{
				Archive: "archived",
			},
			want: hosting.ListRepositoryOptions{
				Limit:      30,
				IsArchived: hosting.TristateTrue,
			},
			wantErr: false,
		},
		{
			name: "invalid archive",
			input: Options{
				Archive: "invalid",
			},
			wantErr: true,
		},
		{
			name: "relation filters",
			input: Options{
				Relation: []string{"owner", "collaborator"},
			},
			want: hosting.ListRepositoryOptions{
				Limit: 30,
				OwnerAffiliations: []hosting.RepositoryAffiliation{
					hosting.RepositoryAffiliationOwner,
					hosting.RepositoryAffiliationCollaborator,
				},
			},
			wantErr: false,
		},
		{
			name: "invalid relation",
			input: Options{
				Relation: []string{"invalid"},
			},
			wantErr: true,
		},
		{
			name: "sort and order",
			input: Options{
				Sort:  "name",
				Order: "asc",
			},
			want: hosting.ListRepositoryOptions{
				Limit: 30,
				OrderBy: hosting.RepositoryOrder{
					Field:     hosting.RepositoryOrderFieldName,
					Direction: hosting.OrderDirectionAsc,
				},
			},
			wantErr: false,
		},
		{
			name: "invalid sort",
			input: Options{
				Sort: "invalid",
			},
			wantErr: true,
		},
		{
			name: "invalid order",
			input: Options{
				Order: "invalid",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convertOpts(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("convertOpts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}

			// Check the result fields
			if got.Limit != tt.want.Limit {
				t.Errorf("Limit = %v, want %v", got.Limit, tt.want.Limit)
			}
			if got.Privacy != tt.want.Privacy {
				t.Errorf("Privacy = %v, want %v", got.Privacy, tt.want.Privacy)
			}
			if got.IsFork != tt.want.IsFork {
				t.Errorf("IsFork = %v, want %v", got.IsFork, tt.want.IsFork)
			}
			if got.IsArchived != tt.want.IsArchived {
				t.Errorf("IsArchived = %v, want %v", got.IsArchived, tt.want.IsArchived)
			}
			// Test other fields as needed
		})
	}
}
