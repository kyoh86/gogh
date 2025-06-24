package remove_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	testtarget "github.com/kyoh86/gogh/v4/app/extra/remove"
	"github.com/kyoh86/gogh/v4/core/extra_mock"
	"github.com/kyoh86/gogh/v4/core/repository"
	"github.com/kyoh86/gogh/v4/core/repository_mock"
	"go.uber.org/mock/gomock"
)

func TestUseCase_Execute(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		name      string
		opts      testtarget.Options
		setupMock func(*gomock.Controller) (*extra_mock.MockExtraService, *repository_mock.MockReferenceParser)
		wantErr   bool
		errMsg    string
	}{
		{
			name: "Successfully remove extra by ID",
			opts: testtarget.Options{
				ID: uuid.New().String(),
			},
			setupMock: func(ctrl *gomock.Controller) (*extra_mock.MockExtraService, *repository_mock.MockReferenceParser) {
				es := extra_mock.NewMockExtraService(ctrl)
				rp := repository_mock.NewMockReferenceParser(ctrl)

				es.EXPECT().Remove(ctx, gomock.Any()).DoAndReturn(
					func(ctx context.Context, id string) error {
						// Validate UUID format
						if _, err := uuid.Parse(id); err != nil {
							t.Errorf("Expected valid UUID, got %s", id)
						}
						return nil
					},
				)

				return es, rp
			},
			wantErr: false,
		},
		{
			name: "Remove extra by ID - not found",
			opts: testtarget.Options{
				ID: uuid.New().String(),
			},
			setupMock: func(ctrl *gomock.Controller) (*extra_mock.MockExtraService, *repository_mock.MockReferenceParser) {
				es := extra_mock.NewMockExtraService(ctrl)
				rp := repository_mock.NewMockReferenceParser(ctrl)

				es.EXPECT().Remove(ctx, gomock.Any()).Return(errors.New("extra not found"))

				return es, rp
			},
			wantErr: true,
			errMsg:  "removing extra by ID",
		},
		{
			name: "Successfully remove named extra",
			opts: testtarget.Options{
				Name: "my-extra",
			},
			setupMock: func(ctrl *gomock.Controller) (*extra_mock.MockExtraService, *repository_mock.MockReferenceParser) {
				es := extra_mock.NewMockExtraService(ctrl)
				rp := repository_mock.NewMockReferenceParser(ctrl)

				es.EXPECT().RemoveNamedExtra(ctx, "my-extra").Return(nil)

				return es, rp
			},
			wantErr: false,
		},
		{
			name: "Remove named extra - not found",
			opts: testtarget.Options{
				Name: "non-existent",
			},
			setupMock: func(ctrl *gomock.Controller) (*extra_mock.MockExtraService, *repository_mock.MockReferenceParser) {
				es := extra_mock.NewMockExtraService(ctrl)
				rp := repository_mock.NewMockReferenceParser(ctrl)

				es.EXPECT().RemoveNamedExtra(ctx, "non-existent").Return(errors.New("named extra not found"))

				return es, rp
			},
			wantErr: true,
			errMsg:  "removing named extra",
		},
		{
			name: "Successfully remove auto extra",
			opts: testtarget.Options{
				Repository: "github.com/owner/repo",
			},
			setupMock: func(ctrl *gomock.Controller) (*extra_mock.MockExtraService, *repository_mock.MockReferenceParser) {
				es := extra_mock.NewMockExtraService(ctrl)
				rp := repository_mock.NewMockReferenceParser(ctrl)

				ref := repository.NewReference("github.com", "owner", "repo")
				rp.EXPECT().Parse("github.com/owner/repo").Return(&ref, nil)
				es.EXPECT().RemoveAutoExtra(ctx, ref).Return(nil)

				return es, rp
			},
			wantErr: false,
		},
		{
			name: "Remove auto extra - invalid repository reference",
			opts: testtarget.Options{
				Repository: "invalid-ref",
			},
			setupMock: func(ctrl *gomock.Controller) (*extra_mock.MockExtraService, *repository_mock.MockReferenceParser) {
				es := extra_mock.NewMockExtraService(ctrl)
				rp := repository_mock.NewMockReferenceParser(ctrl)

				rp.EXPECT().Parse("invalid-ref").Return(nil, errors.New("invalid reference"))

				return es, rp
			},
			wantErr: true,
			errMsg:  "invalid repository reference",
		},
		{
			name: "Remove auto extra - not found",
			opts: testtarget.Options{
				Repository: "github.com/owner/repo",
			},
			setupMock: func(ctrl *gomock.Controller) (*extra_mock.MockExtraService, *repository_mock.MockReferenceParser) {
				es := extra_mock.NewMockExtraService(ctrl)
				rp := repository_mock.NewMockReferenceParser(ctrl)

				ref := repository.NewReference("github.com", "owner", "repo")
				rp.EXPECT().Parse("github.com/owner/repo").Return(&ref, nil)
				es.EXPECT().RemoveAutoExtra(ctx, ref).Return(errors.New("auto extra not found"))

				return es, rp
			},
			wantErr: true,
			errMsg:  "removing auto extra",
		},
		{
			name: "No option specified",
			opts: testtarget.Options{},
			setupMock: func(ctrl *gomock.Controller) (*extra_mock.MockExtraService, *repository_mock.MockReferenceParser) {
				es := extra_mock.NewMockExtraService(ctrl)
				rp := repository_mock.NewMockReferenceParser(ctrl)
				return es, rp
			},
			wantErr: true,
			errMsg:  "one of --id, --name, or --repository must be specified",
		},
		{
			name: "Multiple options specified (ID and Name)",
			opts: testtarget.Options{
				ID:   uuid.New().String(),
				Name: "my-extra",
			},
			setupMock: func(ctrl *gomock.Controller) (*extra_mock.MockExtraService, *repository_mock.MockReferenceParser) {
				es := extra_mock.NewMockExtraService(ctrl)
				rp := repository_mock.NewMockReferenceParser(ctrl)

				// ID takes precedence
				es.EXPECT().Remove(ctx, gomock.Any()).Return(nil)

				return es, rp
			},
			wantErr: false,
		},
		{
			name: "Multiple options specified (Name and Repository)",
			opts: testtarget.Options{
				Name:       "my-extra",
				Repository: "github.com/owner/repo",
			},
			setupMock: func(ctrl *gomock.Controller) (*extra_mock.MockExtraService, *repository_mock.MockReferenceParser) {
				es := extra_mock.NewMockExtraService(ctrl)
				rp := repository_mock.NewMockReferenceParser(ctrl)

				// Name takes precedence over Repository
				es.EXPECT().RemoveNamedExtra(ctx, "my-extra").Return(nil)

				return es, rp
			},
			wantErr: false,
		},
		{
			name: "Remove with empty ID",
			opts: testtarget.Options{
				ID: "",
			},
			setupMock: func(ctrl *gomock.Controller) (*extra_mock.MockExtraService, *repository_mock.MockReferenceParser) {
				es := extra_mock.NewMockExtraService(ctrl)
				rp := repository_mock.NewMockReferenceParser(ctrl)
				return es, rp
			},
			wantErr: true,
			errMsg:  "one of --id, --name, or --repository must be specified",
		},
		{
			name: "Remove with empty name",
			opts: testtarget.Options{
				Name: "",
			},
			setupMock: func(ctrl *gomock.Controller) (*extra_mock.MockExtraService, *repository_mock.MockReferenceParser) {
				es := extra_mock.NewMockExtraService(ctrl)
				rp := repository_mock.NewMockReferenceParser(ctrl)
				return es, rp
			},
			wantErr: true,
			errMsg:  "one of --id, --name, or --repository must be specified",
		},
		{
			name: "Remove with empty repository",
			opts: testtarget.Options{
				Repository: "",
			},
			setupMock: func(ctrl *gomock.Controller) (*extra_mock.MockExtraService, *repository_mock.MockReferenceParser) {
				es := extra_mock.NewMockExtraService(ctrl)
				rp := repository_mock.NewMockReferenceParser(ctrl)
				return es, rp
			},
			wantErr: true,
			errMsg:  "one of --id, --name, or --repository must be specified",
		},
		{
			name: "Service returns unexpected error for ID",
			opts: testtarget.Options{
				ID: uuid.New().String(),
			},
			setupMock: func(ctrl *gomock.Controller) (*extra_mock.MockExtraService, *repository_mock.MockReferenceParser) {
				es := extra_mock.NewMockExtraService(ctrl)
				rp := repository_mock.NewMockReferenceParser(ctrl)

				es.EXPECT().Remove(ctx, gomock.Any()).Return(errors.New("storage error"))

				return es, rp
			},
			wantErr: true,
			errMsg:  "removing extra by ID",
		},
		{
			name: "Service returns unexpected error for named extra",
			opts: testtarget.Options{
				Name: "my-extra",
			},
			setupMock: func(ctrl *gomock.Controller) (*extra_mock.MockExtraService, *repository_mock.MockReferenceParser) {
				es := extra_mock.NewMockExtraService(ctrl)
				rp := repository_mock.NewMockReferenceParser(ctrl)

				es.EXPECT().RemoveNamedExtra(ctx, "my-extra").Return(errors.New("storage error"))

				return es, rp
			},
			wantErr: true,
			errMsg:  "removing named extra",
		},
		{
			name: "Service returns unexpected error for auto extra",
			opts: testtarget.Options{
				Repository: "github.com/owner/repo",
			},
			setupMock: func(ctrl *gomock.Controller) (*extra_mock.MockExtraService, *repository_mock.MockReferenceParser) {
				es := extra_mock.NewMockExtraService(ctrl)
				rp := repository_mock.NewMockReferenceParser(ctrl)

				ref := repository.NewReference("github.com", "owner", "repo")
				rp.EXPECT().Parse("github.com/owner/repo").Return(&ref, nil)
				es.EXPECT().RemoveAutoExtra(ctx, ref).Return(errors.New("storage error"))

				return es, rp
			},
			wantErr: true,
			errMsg:  "removing auto extra",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			es, rp := tc.setupMock(ctrl)
			uc := testtarget.NewUseCase(es, rp)

			err := uc.Execute(ctx, tc.opts)
			if (err != nil) != tc.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tc.wantErr)
			}

			if err != nil && tc.errMsg != "" {
				if !errors.Is(err, errors.New(tc.errMsg)) {
					// Check if error message contains expected substring
					if err.Error() == tc.errMsg || contains(err.Error(), tc.errMsg) {
						// Error message matches expectation
					} else {
						t.Errorf("Execute() error message = %v, want to contain %v", err.Error(), tc.errMsg)
					}
				}
			}
		})
	}
}

func TestUseCase_Execute_PriorityOrder(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	es := extra_mock.NewMockExtraService(ctrl)
	rp := repository_mock.NewMockReferenceParser(ctrl)
	uc := testtarget.NewUseCase(es, rp)

	// Test priority: ID > Name > Repository
	opts := testtarget.Options{
		ID:         uuid.New().String(),
		Name:       "my-extra",
		Repository: "github.com/owner/repo",
	}

	// Only ID removal should be called
	es.EXPECT().Remove(ctx, opts.ID).Return(nil)
	// Name and Repository should NOT be called

	err := uc.Execute(ctx, opts)
	if err != nil {
		t.Errorf("Execute() unexpected error = %v", err)
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && len(substr) > 0 && s[:len(substr)] == substr || len(s) > len(substr) && s[len(s)-len(substr):] == substr || (len(substr) > 0 && len(s) > len(substr) && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
