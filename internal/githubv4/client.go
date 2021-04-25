// Code generated by github.com/Yamashou/gqlgenc, DO NOT EDIT.

package githubv4

import (
	"context"
	"net/http"

	"github.com/Yamashou/gqlgenc/client"
)

type Client struct {
	Client *client.Client
}

func NewClient(cli *http.Client, baseURL string, options ...client.HTTPRequestOption) *Client {
	return &Client{Client: client.NewClient(cli, baseURL, options...)}
}

type Query struct {
	CodeOfConduct                            *CodeOfConduct                     "json:\"codeOfConduct\" graphql:\"codeOfConduct\""
	CodesOfConduct                           []*CodeOfConduct                   "json:\"codesOfConduct\" graphql:\"codesOfConduct\""
	Enterprise                               *Enterprise                        "json:\"enterprise\" graphql:\"enterprise\""
	EnterpriseAdministratorInvitation        *EnterpriseAdministratorInvitation "json:\"enterpriseAdministratorInvitation\" graphql:\"enterpriseAdministratorInvitation\""
	EnterpriseAdministratorInvitationByToken *EnterpriseAdministratorInvitation "json:\"enterpriseAdministratorInvitationByToken\" graphql:\"enterpriseAdministratorInvitationByToken\""
	License                                  *License                           "json:\"license\" graphql:\"license\""
	Licenses                                 []*License                         "json:\"licenses\" graphql:\"licenses\""
	MarketplaceCategories                    []*MarketplaceCategory             "json:\"marketplaceCategories\" graphql:\"marketplaceCategories\""
	MarketplaceCategory                      *MarketplaceCategory               "json:\"marketplaceCategory\" graphql:\"marketplaceCategory\""
	MarketplaceListing                       *MarketplaceListing                "json:\"marketplaceListing\" graphql:\"marketplaceListing\""
	MarketplaceListings                      MarketplaceListingConnection       "json:\"marketplaceListings\" graphql:\"marketplaceListings\""
	Meta                                     GitHubMetadata                     "json:\"meta\" graphql:\"meta\""
	Node                                     Node                               "json:\"node\" graphql:\"node\""
	Nodes                                    []Node                             "json:\"nodes\" graphql:\"nodes\""
	Organization                             *Organization                      "json:\"organization\" graphql:\"organization\""
	RateLimit                                *RateLimit                         "json:\"rateLimit\" graphql:\"rateLimit\""
	Relay                                    *Query                             "json:\"relay\" graphql:\"relay\""
	Repository                               *Repository                        "json:\"repository\" graphql:\"repository\""
	RepositoryOwner                          RepositoryOwner                    "json:\"repositoryOwner\" graphql:\"repositoryOwner\""
	Resource                                 UniformResourceLocatable           "json:\"resource\" graphql:\"resource\""
	Search                                   SearchResultItemConnection         "json:\"search\" graphql:\"search\""
	SecurityAdvisories                       SecurityAdvisoryConnection         "json:\"securityAdvisories\" graphql:\"securityAdvisories\""
	SecurityAdvisory                         *SecurityAdvisory                  "json:\"securityAdvisory\" graphql:\"securityAdvisory\""
	SecurityVulnerabilities                  SecurityVulnerabilityConnection    "json:\"securityVulnerabilities\" graphql:\"securityVulnerabilities\""
	Sponsorables                             SponsorableItemConnection          "json:\"sponsorables\" graphql:\"sponsorables\""
	SponsorsListing                          *SponsorsListing                   "json:\"sponsorsListing\" graphql:\"sponsorsListing\""
	Topic                                    *Topic                             "json:\"topic\" graphql:\"topic\""
	User                                     *User                              "json:\"user\" graphql:\"user\""
	Viewer                                   User                               "json:\"viewer\" graphql:\"viewer\""
}
type Mutation struct {
	AcceptEnterpriseAdministratorInvitation                     *AcceptEnterpriseAdministratorInvitationPayload                     "json:\"acceptEnterpriseAdministratorInvitation\" graphql:\"acceptEnterpriseAdministratorInvitation\""
	AcceptTopicSuggestion                                       *AcceptTopicSuggestionPayload                                       "json:\"acceptTopicSuggestion\" graphql:\"acceptTopicSuggestion\""
	AddAssigneesToAssignable                                    *AddAssigneesToAssignablePayload                                    "json:\"addAssigneesToAssignable\" graphql:\"addAssigneesToAssignable\""
	AddComment                                                  *AddCommentPayload                                                  "json:\"addComment\" graphql:\"addComment\""
	AddEnterpriseSupportEntitlement                             *AddEnterpriseSupportEntitlementPayload                             "json:\"addEnterpriseSupportEntitlement\" graphql:\"addEnterpriseSupportEntitlement\""
	AddLabelsToLabelable                                        *AddLabelsToLabelablePayload                                        "json:\"addLabelsToLabelable\" graphql:\"addLabelsToLabelable\""
	AddProjectCard                                              *AddProjectCardPayload                                              "json:\"addProjectCard\" graphql:\"addProjectCard\""
	AddProjectColumn                                            *AddProjectColumnPayload                                            "json:\"addProjectColumn\" graphql:\"addProjectColumn\""
	AddPullRequestReview                                        *AddPullRequestReviewPayload                                        "json:\"addPullRequestReview\" graphql:\"addPullRequestReview\""
	AddPullRequestReviewComment                                 *AddPullRequestReviewCommentPayload                                 "json:\"addPullRequestReviewComment\" graphql:\"addPullRequestReviewComment\""
	AddPullRequestReviewThread                                  *AddPullRequestReviewThreadPayload                                  "json:\"addPullRequestReviewThread\" graphql:\"addPullRequestReviewThread\""
	AddReaction                                                 *AddReactionPayload                                                 "json:\"addReaction\" graphql:\"addReaction\""
	AddStar                                                     *AddStarPayload                                                     "json:\"addStar\" graphql:\"addStar\""
	AddVerifiableDomain                                         *AddVerifiableDomainPayload                                         "json:\"addVerifiableDomain\" graphql:\"addVerifiableDomain\""
	ApproveVerifiableDomain                                     *ApproveVerifiableDomainPayload                                     "json:\"approveVerifiableDomain\" graphql:\"approveVerifiableDomain\""
	ArchiveRepository                                           *ArchiveRepositoryPayload                                           "json:\"archiveRepository\" graphql:\"archiveRepository\""
	CancelEnterpriseAdminInvitation                             *CancelEnterpriseAdminInvitationPayload                             "json:\"cancelEnterpriseAdminInvitation\" graphql:\"cancelEnterpriseAdminInvitation\""
	ChangeUserStatus                                            *ChangeUserStatusPayload                                            "json:\"changeUserStatus\" graphql:\"changeUserStatus\""
	ClearLabelsFromLabelable                                    *ClearLabelsFromLabelablePayload                                    "json:\"clearLabelsFromLabelable\" graphql:\"clearLabelsFromLabelable\""
	CloneProject                                                *CloneProjectPayload                                                "json:\"cloneProject\" graphql:\"cloneProject\""
	CloneTemplateRepository                                     *CloneTemplateRepositoryPayload                                     "json:\"cloneTemplateRepository\" graphql:\"cloneTemplateRepository\""
	CloseIssue                                                  *CloseIssuePayload                                                  "json:\"closeIssue\" graphql:\"closeIssue\""
	ClosePullRequest                                            *ClosePullRequestPayload                                            "json:\"closePullRequest\" graphql:\"closePullRequest\""
	ConvertProjectCardNoteToIssue                               *ConvertProjectCardNoteToIssuePayload                               "json:\"convertProjectCardNoteToIssue\" graphql:\"convertProjectCardNoteToIssue\""
	ConvertPullRequestToDraft                                   *ConvertPullRequestToDraftPayload                                   "json:\"convertPullRequestToDraft\" graphql:\"convertPullRequestToDraft\""
	CreateBranchProtectionRule                                  *CreateBranchProtectionRulePayload                                  "json:\"createBranchProtectionRule\" graphql:\"createBranchProtectionRule\""
	CreateCheckRun                                              *CreateCheckRunPayload                                              "json:\"createCheckRun\" graphql:\"createCheckRun\""
	CreateCheckSuite                                            *CreateCheckSuitePayload                                            "json:\"createCheckSuite\" graphql:\"createCheckSuite\""
	CreateEnterpriseOrganization                                *CreateEnterpriseOrganizationPayload                                "json:\"createEnterpriseOrganization\" graphql:\"createEnterpriseOrganization\""
	CreateIPAllowListEntry                                      *CreateIPAllowListEntryPayload                                      "json:\"createIpAllowListEntry\" graphql:\"createIpAllowListEntry\""
	CreateIssue                                                 *CreateIssuePayload                                                 "json:\"createIssue\" graphql:\"createIssue\""
	CreateProject                                               *CreateProjectPayload                                               "json:\"createProject\" graphql:\"createProject\""
	CreatePullRequest                                           *CreatePullRequestPayload                                           "json:\"createPullRequest\" graphql:\"createPullRequest\""
	CreateRef                                                   *CreateRefPayload                                                   "json:\"createRef\" graphql:\"createRef\""
	CreateRepository                                            *CreateRepositoryPayload                                            "json:\"createRepository\" graphql:\"createRepository\""
	CreateTeamDiscussion                                        *CreateTeamDiscussionPayload                                        "json:\"createTeamDiscussion\" graphql:\"createTeamDiscussion\""
	CreateTeamDiscussionComment                                 *CreateTeamDiscussionCommentPayload                                 "json:\"createTeamDiscussionComment\" graphql:\"createTeamDiscussionComment\""
	DeclineTopicSuggestion                                      *DeclineTopicSuggestionPayload                                      "json:\"declineTopicSuggestion\" graphql:\"declineTopicSuggestion\""
	DeleteBranchProtectionRule                                  *DeleteBranchProtectionRulePayload                                  "json:\"deleteBranchProtectionRule\" graphql:\"deleteBranchProtectionRule\""
	DeleteDeployment                                            *DeleteDeploymentPayload                                            "json:\"deleteDeployment\" graphql:\"deleteDeployment\""
	DeleteIPAllowListEntry                                      *DeleteIPAllowListEntryPayload                                      "json:\"deleteIpAllowListEntry\" graphql:\"deleteIpAllowListEntry\""
	DeleteIssue                                                 *DeleteIssuePayload                                                 "json:\"deleteIssue\" graphql:\"deleteIssue\""
	DeleteIssueComment                                          *DeleteIssueCommentPayload                                          "json:\"deleteIssueComment\" graphql:\"deleteIssueComment\""
	DeleteProject                                               *DeleteProjectPayload                                               "json:\"deleteProject\" graphql:\"deleteProject\""
	DeleteProjectCard                                           *DeleteProjectCardPayload                                           "json:\"deleteProjectCard\" graphql:\"deleteProjectCard\""
	DeleteProjectColumn                                         *DeleteProjectColumnPayload                                         "json:\"deleteProjectColumn\" graphql:\"deleteProjectColumn\""
	DeletePullRequestReview                                     *DeletePullRequestReviewPayload                                     "json:\"deletePullRequestReview\" graphql:\"deletePullRequestReview\""
	DeletePullRequestReviewComment                              *DeletePullRequestReviewCommentPayload                              "json:\"deletePullRequestReviewComment\" graphql:\"deletePullRequestReviewComment\""
	DeleteRef                                                   *DeleteRefPayload                                                   "json:\"deleteRef\" graphql:\"deleteRef\""
	DeleteTeamDiscussion                                        *DeleteTeamDiscussionPayload                                        "json:\"deleteTeamDiscussion\" graphql:\"deleteTeamDiscussion\""
	DeleteTeamDiscussionComment                                 *DeleteTeamDiscussionCommentPayload                                 "json:\"deleteTeamDiscussionComment\" graphql:\"deleteTeamDiscussionComment\""
	DeleteVerifiableDomain                                      *DeleteVerifiableDomainPayload                                      "json:\"deleteVerifiableDomain\" graphql:\"deleteVerifiableDomain\""
	DisablePullRequestAutoMerge                                 *DisablePullRequestAutoMergePayload                                 "json:\"disablePullRequestAutoMerge\" graphql:\"disablePullRequestAutoMerge\""
	DismissPullRequestReview                                    *DismissPullRequestReviewPayload                                    "json:\"dismissPullRequestReview\" graphql:\"dismissPullRequestReview\""
	EnablePullRequestAutoMerge                                  *EnablePullRequestAutoMergePayload                                  "json:\"enablePullRequestAutoMerge\" graphql:\"enablePullRequestAutoMerge\""
	FollowUser                                                  *FollowUserPayload                                                  "json:\"followUser\" graphql:\"followUser\""
	InviteEnterpriseAdmin                                       *InviteEnterpriseAdminPayload                                       "json:\"inviteEnterpriseAdmin\" graphql:\"inviteEnterpriseAdmin\""
	LinkRepositoryToProject                                     *LinkRepositoryToProjectPayload                                     "json:\"linkRepositoryToProject\" graphql:\"linkRepositoryToProject\""
	LockLockable                                                *LockLockablePayload                                                "json:\"lockLockable\" graphql:\"lockLockable\""
	MarkFileAsViewed                                            *MarkFileAsViewedPayload                                            "json:\"markFileAsViewed\" graphql:\"markFileAsViewed\""
	MarkPullRequestReadyForReview                               *MarkPullRequestReadyForReviewPayload                               "json:\"markPullRequestReadyForReview\" graphql:\"markPullRequestReadyForReview\""
	MergeBranch                                                 *MergeBranchPayload                                                 "json:\"mergeBranch\" graphql:\"mergeBranch\""
	MergePullRequest                                            *MergePullRequestPayload                                            "json:\"mergePullRequest\" graphql:\"mergePullRequest\""
	MinimizeComment                                             *MinimizeCommentPayload                                             "json:\"minimizeComment\" graphql:\"minimizeComment\""
	MoveProjectCard                                             *MoveProjectCardPayload                                             "json:\"moveProjectCard\" graphql:\"moveProjectCard\""
	MoveProjectColumn                                           *MoveProjectColumnPayload                                           "json:\"moveProjectColumn\" graphql:\"moveProjectColumn\""
	PinIssue                                                    *PinIssuePayload                                                    "json:\"pinIssue\" graphql:\"pinIssue\""
	RegenerateEnterpriseIdentityProviderRecoveryCodes           *RegenerateEnterpriseIdentityProviderRecoveryCodesPayload           "json:\"regenerateEnterpriseIdentityProviderRecoveryCodes\" graphql:\"regenerateEnterpriseIdentityProviderRecoveryCodes\""
	RegenerateVerifiableDomainToken                             *RegenerateVerifiableDomainTokenPayload                             "json:\"regenerateVerifiableDomainToken\" graphql:\"regenerateVerifiableDomainToken\""
	RemoveAssigneesFromAssignable                               *RemoveAssigneesFromAssignablePayload                               "json:\"removeAssigneesFromAssignable\" graphql:\"removeAssigneesFromAssignable\""
	RemoveEnterpriseAdmin                                       *RemoveEnterpriseAdminPayload                                       "json:\"removeEnterpriseAdmin\" graphql:\"removeEnterpriseAdmin\""
	RemoveEnterpriseIdentityProvider                            *RemoveEnterpriseIdentityProviderPayload                            "json:\"removeEnterpriseIdentityProvider\" graphql:\"removeEnterpriseIdentityProvider\""
	RemoveEnterpriseOrganization                                *RemoveEnterpriseOrganizationPayload                                "json:\"removeEnterpriseOrganization\" graphql:\"removeEnterpriseOrganization\""
	RemoveEnterpriseSupportEntitlement                          *RemoveEnterpriseSupportEntitlementPayload                          "json:\"removeEnterpriseSupportEntitlement\" graphql:\"removeEnterpriseSupportEntitlement\""
	RemoveLabelsFromLabelable                                   *RemoveLabelsFromLabelablePayload                                   "json:\"removeLabelsFromLabelable\" graphql:\"removeLabelsFromLabelable\""
	RemoveOutsideCollaborator                                   *RemoveOutsideCollaboratorPayload                                   "json:\"removeOutsideCollaborator\" graphql:\"removeOutsideCollaborator\""
	RemoveReaction                                              *RemoveReactionPayload                                              "json:\"removeReaction\" graphql:\"removeReaction\""
	RemoveStar                                                  *RemoveStarPayload                                                  "json:\"removeStar\" graphql:\"removeStar\""
	ReopenIssue                                                 *ReopenIssuePayload                                                 "json:\"reopenIssue\" graphql:\"reopenIssue\""
	ReopenPullRequest                                           *ReopenPullRequestPayload                                           "json:\"reopenPullRequest\" graphql:\"reopenPullRequest\""
	RequestReviews                                              *RequestReviewsPayload                                              "json:\"requestReviews\" graphql:\"requestReviews\""
	RerequestCheckSuite                                         *RerequestCheckSuitePayload                                         "json:\"rerequestCheckSuite\" graphql:\"rerequestCheckSuite\""
	ResolveReviewThread                                         *ResolveReviewThreadPayload                                         "json:\"resolveReviewThread\" graphql:\"resolveReviewThread\""
	SetEnterpriseIdentityProvider                               *SetEnterpriseIdentityProviderPayload                               "json:\"setEnterpriseIdentityProvider\" graphql:\"setEnterpriseIdentityProvider\""
	SetOrganizationInteractionLimit                             *SetOrganizationInteractionLimitPayload                             "json:\"setOrganizationInteractionLimit\" graphql:\"setOrganizationInteractionLimit\""
	SetRepositoryInteractionLimit                               *SetRepositoryInteractionLimitPayload                               "json:\"setRepositoryInteractionLimit\" graphql:\"setRepositoryInteractionLimit\""
	SetUserInteractionLimit                                     *SetUserInteractionLimitPayload                                     "json:\"setUserInteractionLimit\" graphql:\"setUserInteractionLimit\""
	SubmitPullRequestReview                                     *SubmitPullRequestReviewPayload                                     "json:\"submitPullRequestReview\" graphql:\"submitPullRequestReview\""
	TransferIssue                                               *TransferIssuePayload                                               "json:\"transferIssue\" graphql:\"transferIssue\""
	UnarchiveRepository                                         *UnarchiveRepositoryPayload                                         "json:\"unarchiveRepository\" graphql:\"unarchiveRepository\""
	UnfollowUser                                                *UnfollowUserPayload                                                "json:\"unfollowUser\" graphql:\"unfollowUser\""
	UnlinkRepositoryFromProject                                 *UnlinkRepositoryFromProjectPayload                                 "json:\"unlinkRepositoryFromProject\" graphql:\"unlinkRepositoryFromProject\""
	UnlockLockable                                              *UnlockLockablePayload                                              "json:\"unlockLockable\" graphql:\"unlockLockable\""
	UnmarkFileAsViewed                                          *UnmarkFileAsViewedPayload                                          "json:\"unmarkFileAsViewed\" graphql:\"unmarkFileAsViewed\""
	UnmarkIssueAsDuplicate                                      *UnmarkIssueAsDuplicatePayload                                      "json:\"unmarkIssueAsDuplicate\" graphql:\"unmarkIssueAsDuplicate\""
	UnminimizeComment                                           *UnminimizeCommentPayload                                           "json:\"unminimizeComment\" graphql:\"unminimizeComment\""
	UnpinIssue                                                  *UnpinIssuePayload                                                  "json:\"unpinIssue\" graphql:\"unpinIssue\""
	UnresolveReviewThread                                       *UnresolveReviewThreadPayload                                       "json:\"unresolveReviewThread\" graphql:\"unresolveReviewThread\""
	UpdateBranchProtectionRule                                  *UpdateBranchProtectionRulePayload                                  "json:\"updateBranchProtectionRule\" graphql:\"updateBranchProtectionRule\""
	UpdateCheckRun                                              *UpdateCheckRunPayload                                              "json:\"updateCheckRun\" graphql:\"updateCheckRun\""
	UpdateCheckSuitePreferences                                 *UpdateCheckSuitePreferencesPayload                                 "json:\"updateCheckSuitePreferences\" graphql:\"updateCheckSuitePreferences\""
	UpdateEnterpriseAdministratorRole                           *UpdateEnterpriseAdministratorRolePayload                           "json:\"updateEnterpriseAdministratorRole\" graphql:\"updateEnterpriseAdministratorRole\""
	UpdateEnterpriseAllowPrivateRepositoryForkingSetting        *UpdateEnterpriseAllowPrivateRepositoryForkingSettingPayload        "json:\"updateEnterpriseAllowPrivateRepositoryForkingSetting\" graphql:\"updateEnterpriseAllowPrivateRepositoryForkingSetting\""
	UpdateEnterpriseDefaultRepositoryPermissionSetting          *UpdateEnterpriseDefaultRepositoryPermissionSettingPayload          "json:\"updateEnterpriseDefaultRepositoryPermissionSetting\" graphql:\"updateEnterpriseDefaultRepositoryPermissionSetting\""
	UpdateEnterpriseMembersCanChangeRepositoryVisibilitySetting *UpdateEnterpriseMembersCanChangeRepositoryVisibilitySettingPayload "json:\"updateEnterpriseMembersCanChangeRepositoryVisibilitySetting\" graphql:\"updateEnterpriseMembersCanChangeRepositoryVisibilitySetting\""
	UpdateEnterpriseMembersCanCreateRepositoriesSetting         *UpdateEnterpriseMembersCanCreateRepositoriesSettingPayload         "json:\"updateEnterpriseMembersCanCreateRepositoriesSetting\" graphql:\"updateEnterpriseMembersCanCreateRepositoriesSetting\""
	UpdateEnterpriseMembersCanDeleteIssuesSetting               *UpdateEnterpriseMembersCanDeleteIssuesSettingPayload               "json:\"updateEnterpriseMembersCanDeleteIssuesSetting\" graphql:\"updateEnterpriseMembersCanDeleteIssuesSetting\""
	UpdateEnterpriseMembersCanDeleteRepositoriesSetting         *UpdateEnterpriseMembersCanDeleteRepositoriesSettingPayload         "json:\"updateEnterpriseMembersCanDeleteRepositoriesSetting\" graphql:\"updateEnterpriseMembersCanDeleteRepositoriesSetting\""
	UpdateEnterpriseMembersCanInviteCollaboratorsSetting        *UpdateEnterpriseMembersCanInviteCollaboratorsSettingPayload        "json:\"updateEnterpriseMembersCanInviteCollaboratorsSetting\" graphql:\"updateEnterpriseMembersCanInviteCollaboratorsSetting\""
	UpdateEnterpriseMembersCanMakePurchasesSetting              *UpdateEnterpriseMembersCanMakePurchasesSettingPayload              "json:\"updateEnterpriseMembersCanMakePurchasesSetting\" graphql:\"updateEnterpriseMembersCanMakePurchasesSetting\""
	UpdateEnterpriseMembersCanUpdateProtectedBranchesSetting    *UpdateEnterpriseMembersCanUpdateProtectedBranchesSettingPayload    "json:\"updateEnterpriseMembersCanUpdateProtectedBranchesSetting\" graphql:\"updateEnterpriseMembersCanUpdateProtectedBranchesSetting\""
	UpdateEnterpriseMembersCanViewDependencyInsightsSetting     *UpdateEnterpriseMembersCanViewDependencyInsightsSettingPayload     "json:\"updateEnterpriseMembersCanViewDependencyInsightsSetting\" graphql:\"updateEnterpriseMembersCanViewDependencyInsightsSetting\""
	UpdateEnterpriseOrganizationProjectsSetting                 *UpdateEnterpriseOrganizationProjectsSettingPayload                 "json:\"updateEnterpriseOrganizationProjectsSetting\" graphql:\"updateEnterpriseOrganizationProjectsSetting\""
	UpdateEnterpriseProfile                                     *UpdateEnterpriseProfilePayload                                     "json:\"updateEnterpriseProfile\" graphql:\"updateEnterpriseProfile\""
	UpdateEnterpriseRepositoryProjectsSetting                   *UpdateEnterpriseRepositoryProjectsSettingPayload                   "json:\"updateEnterpriseRepositoryProjectsSetting\" graphql:\"updateEnterpriseRepositoryProjectsSetting\""
	UpdateEnterpriseTeamDiscussionsSetting                      *UpdateEnterpriseTeamDiscussionsSettingPayload                      "json:\"updateEnterpriseTeamDiscussionsSetting\" graphql:\"updateEnterpriseTeamDiscussionsSetting\""
	UpdateEnterpriseTwoFactorAuthenticationRequiredSetting      *UpdateEnterpriseTwoFactorAuthenticationRequiredSettingPayload      "json:\"updateEnterpriseTwoFactorAuthenticationRequiredSetting\" graphql:\"updateEnterpriseTwoFactorAuthenticationRequiredSetting\""
	UpdateIPAllowListEnabledSetting                             *UpdateIPAllowListEnabledSettingPayload                             "json:\"updateIpAllowListEnabledSetting\" graphql:\"updateIpAllowListEnabledSetting\""
	UpdateIPAllowListEntry                                      *UpdateIPAllowListEntryPayload                                      "json:\"updateIpAllowListEntry\" graphql:\"updateIpAllowListEntry\""
	UpdateIssue                                                 *UpdateIssuePayload                                                 "json:\"updateIssue\" graphql:\"updateIssue\""
	UpdateIssueComment                                          *UpdateIssueCommentPayload                                          "json:\"updateIssueComment\" graphql:\"updateIssueComment\""
	UpdateNotificationRestrictionSetting                        *UpdateNotificationRestrictionSettingPayload                        "json:\"updateNotificationRestrictionSetting\" graphql:\"updateNotificationRestrictionSetting\""
	UpdateProject                                               *UpdateProjectPayload                                               "json:\"updateProject\" graphql:\"updateProject\""
	UpdateProjectCard                                           *UpdateProjectCardPayload                                           "json:\"updateProjectCard\" graphql:\"updateProjectCard\""
	UpdateProjectColumn                                         *UpdateProjectColumnPayload                                         "json:\"updateProjectColumn\" graphql:\"updateProjectColumn\""
	UpdatePullRequest                                           *UpdatePullRequestPayload                                           "json:\"updatePullRequest\" graphql:\"updatePullRequest\""
	UpdatePullRequestReview                                     *UpdatePullRequestReviewPayload                                     "json:\"updatePullRequestReview\" graphql:\"updatePullRequestReview\""
	UpdatePullRequestReviewComment                              *UpdatePullRequestReviewCommentPayload                              "json:\"updatePullRequestReviewComment\" graphql:\"updatePullRequestReviewComment\""
	UpdateRef                                                   *UpdateRefPayload                                                   "json:\"updateRef\" graphql:\"updateRef\""
	UpdateRepository                                            *UpdateRepositoryPayload                                            "json:\"updateRepository\" graphql:\"updateRepository\""
	UpdateSubscription                                          *UpdateSubscriptionPayload                                          "json:\"updateSubscription\" graphql:\"updateSubscription\""
	UpdateTeamDiscussion                                        *UpdateTeamDiscussionPayload                                        "json:\"updateTeamDiscussion\" graphql:\"updateTeamDiscussion\""
	UpdateTeamDiscussionComment                                 *UpdateTeamDiscussionCommentPayload                                 "json:\"updateTeamDiscussionComment\" graphql:\"updateTeamDiscussionComment\""
	UpdateTopics                                                *UpdateTopicsPayload                                                "json:\"updateTopics\" graphql:\"updateTopics\""
	VerifyVerifiableDomain                                      *VerifyVerifiableDomainPayload                                      "json:\"verifyVerifiableDomain\" graphql:\"verifyVerifiableDomain\""
}
type Repo struct {
	ID    string "json:\"id\" graphql:\"id\""
	URL   string "json:\"url\" graphql:\"url\""
	Owner struct {
		Login string "json:\"login\" graphql:\"login\""
	} "json:\"owner\" graphql:\"owner\""
	Name             string  "json:\"name\" graphql:\"name\""
	Description      *string "json:\"description\" graphql:\"description\""
	CreatedAt        string  "json:\"createdAt\" graphql:\"createdAt\""
	IsArchived       bool    "json:\"isArchived\" graphql:\"isArchived\""
	IsDisabled       bool    "json:\"isDisabled\" graphql:\"isDisabled\""
	IsEmpty          bool    "json:\"isEmpty\" graphql:\"isEmpty\""
	IsFork           bool    "json:\"isFork\" graphql:\"isFork\""
	IsInOrganization bool    "json:\"isInOrganization\" graphql:\"isInOrganization\""
	IsLocked         bool    "json:\"isLocked\" graphql:\"isLocked\""
	IsPrivate        bool    "json:\"isPrivate\" graphql:\"isPrivate\""
	IsTemplate       bool    "json:\"isTemplate\" graphql:\"isTemplate\""
	LicenseInfo      *struct {
		Name string "json:\"name\" graphql:\"name\""
		Key  string "json:\"key\" graphql:\"key\""
	} "json:\"licenseInfo\" graphql:\"licenseInfo\""
	Parent *struct {
		Owner struct {
			ID    string "json:\"id\" graphql:\"id\""
			Login string "json:\"login\" graphql:\"login\""
		} "json:\"owner\" graphql:\"owner\""
		Name string "json:\"name\" graphql:\"name\""
	} "json:\"parent\" graphql:\"parent\""
	PrimaryLanguage *struct {
		Name string "json:\"name\" graphql:\"name\""
	} "json:\"primaryLanguage\" graphql:\"primaryLanguage\""
	PushedAt         *string "json:\"pushedAt\" graphql:\"pushedAt\""
	RepositoryTopics struct {
		Edges []*struct {
			Node *struct {
				Topic struct {
					Name string "json:\"name\" graphql:\"name\""
				} "json:\"topic\" graphql:\"topic\""
			} "json:\"node\" graphql:\"node\""
		} "json:\"edges\" graphql:\"edges\""
	} "json:\"repositoryTopics\" graphql:\"repositoryTopics\""
	ResourcePath   string "json:\"resourcePath\" graphql:\"resourcePath\""
	StargazerCount int64  "json:\"stargazerCount\" graphql:\"stargazerCount\""
	UpdatedAt      string "json:\"updatedAt\" graphql:\"updatedAt\""
}
type GetUserID struct {
	User *struct {
		ID string "json:\"id\" graphql:\"id\""
	} "json:\"user\" graphql:\"user\""
}
type GetViewerID struct {
	Viewer struct {
		ID string "json:\"id\" graphql:\"id\""
	} "json:\"viewer\" graphql:\"viewer\""
}
type ListMyOrganizations struct {
	Viewer struct {
		Organizations struct {
			Edges []*struct {
				Node *struct {
					Login string "json:\"login\" graphql:\"login\""
					ID    string "json:\"id\" graphql:\"id\""
				} "json:\"node\" graphql:\"node\""
				Cursor string "json:\"cursor\" graphql:\"cursor\""
			} "json:\"edges\" graphql:\"edges\""
		} "json:\"organizations\" graphql:\"organizations\""
	} "json:\"viewer\" graphql:\"viewer\""
}
type ListRepoUrls struct {
	Viewer struct {
		Repositories struct {
			Edges []*struct {
				Node *struct {
					URL string "json:\"url\" graphql:\"url\""
				} "json:\"node\" graphql:\"node\""
				Cursor string "json:\"cursor\" graphql:\"cursor\""
			} "json:\"edges\" graphql:\"edges\""
		} "json:\"repositories\" graphql:\"repositories\""
	} "json:\"viewer\" graphql:\"viewer\""
}
type ListRepos struct {
	Viewer struct {
		Repositories struct {
			Edges []*struct {
				Node   *Repo  "json:\"node\" graphql:\"node\""
				Cursor string "json:\"cursor\" graphql:\"cursor\""
			} "json:\"edges\" graphql:\"edges\""
		} "json:\"repositories\" graphql:\"repositories\""
	} "json:\"viewer\" graphql:\"viewer\""
}

const GetUserIDDocument = `query GetUserId ($login: String!) {
	user(login: $login) {
		id
	}
}
`

func (c *Client) GetUserID(ctx context.Context, login string, httpRequestOptions ...client.HTTPRequestOption) (*GetUserID, error) {
	vars := map[string]interface{}{
		"login": login,
	}

	var res GetUserID
	if err := c.Client.Post(ctx, "GetUserId", GetUserIDDocument, &res, vars, httpRequestOptions...); err != nil {
		return nil, err
	}

	return &res, nil
}

const GetViewerIDDocument = `query GetViewerId {
	viewer {
		id
	}
}
`

func (c *Client) GetViewerID(ctx context.Context, httpRequestOptions ...client.HTTPRequestOption) (*GetViewerID, error) {
	vars := map[string]interface{}{}

	var res GetViewerID
	if err := c.Client.Post(ctx, "GetViewerId", GetViewerIDDocument, &res, vars, httpRequestOptions...); err != nil {
		return nil, err
	}

	return &res, nil
}

const ListMyOrganizationsDocument = `query ListMyOrganizations ($organizationCursor: String) {
	viewer {
		organizations(after: $organizationCursor, first: 100) {
			edges {
				node {
					login
					id
				}
				cursor
			}
		}
	}
}
`

func (c *Client) ListMyOrganizations(ctx context.Context, organizationCursor *string, httpRequestOptions ...client.HTTPRequestOption) (*ListMyOrganizations, error) {
	vars := map[string]interface{}{
		"organizationCursor": organizationCursor,
	}

	var res ListMyOrganizations
	if err := c.Client.Post(ctx, "ListMyOrganizations", ListMyOrganizationsDocument, &res, vars, httpRequestOptions...); err != nil {
		return nil, err
	}

	return &res, nil
}

const ListRepoUrlsDocument = `query ListRepoUrls ($repositoryCursor: String) {
	viewer {
		repositories(first: 30, after: $repositoryCursor) {
			edges {
				node {
					url
				}
				cursor
			}
		}
	}
}
`

func (c *Client) ListRepoUrls(ctx context.Context, repositoryCursor *string, httpRequestOptions ...client.HTTPRequestOption) (*ListRepoUrls, error) {
	vars := map[string]interface{}{
		"repositoryCursor": repositoryCursor,
	}

	var res ListRepoUrls
	if err := c.Client.Post(ctx, "ListRepoUrls", ListRepoUrlsDocument, &res, vars, httpRequestOptions...); err != nil {
		return nil, err
	}

	return &res, nil
}

const ListReposDocument = `query ListRepos ($repositoryCursor: String) {
	viewer {
		repositories(first: 30, after: $repositoryCursor) {
			edges {
				node {
					... Repo
				}
				cursor
			}
		}
	}
}
fragment Repo on Repository {
	id
	url
	owner {
		login
	}
	name
	description
	createdAt
	isArchived
	isDisabled
	isEmpty
	isFork
	isInOrganization
	isLocked
	isPrivate
	isTemplate
	licenseInfo {
		name
		key
	}
	parent {
		owner {
			id
			login
		}
		name
	}
	primaryLanguage {
		name
	}
	pushedAt
	repositoryTopics(first: 30) {
		edges {
			node {
				topic {
					name
				}
			}
		}
	}
	resourcePath
	stargazerCount
	updatedAt
}
`

func (c *Client) ListRepos(ctx context.Context, repositoryCursor *string, httpRequestOptions ...client.HTTPRequestOption) (*ListRepos, error) {
	vars := map[string]interface{}{
		"repositoryCursor": repositoryCursor,
	}

	var res ListRepos
	if err := c.Client.Post(ctx, "ListRepos", ListReposDocument, &res, vars, httpRequestOptions...); err != nil {
		return nil, err
	}

	return &res, nil
}
