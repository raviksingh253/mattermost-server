//go:generate go run layer_generators/main.go

// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package store

import (
	"context"
	"time"

	"github.com/mattermost/mattermost-server/v5/model"
)

type StoreResult struct {
	Data interface{}

	// NErr a temporary field used by the new code for the AppError migration. This will later become Err when the entire store is migrated.
	NErr error
}

type Store interface {
	Team() TeamStore
	Channel() ChannelStore
	Post() PostStore
	Thread() ThreadStore
	User() UserStore
	Bot() BotStore
	Audit() AuditStore
	ClusterDiscovery() ClusterDiscoveryStore
	Compliance() ComplianceStore
	Session() SessionStore
	OAuth() OAuthStore
	System() SystemStore
	Webhook() WebhookStore
	Command() CommandStore
	CommandWebhook() CommandWebhookStore
	Preference() PreferenceStore
	License() LicenseStore
	Token() TokenStore
	Emoji() EmojiStore
	Status() StatusStore
	FileInfo() FileInfoStore
	UploadSession() UploadSessionStore
	Reaction() ReactionStore
	Role() RoleStore
	Scheme() SchemeStore
	Job() JobStore
	UserAccessToken() UserAccessTokenStore
	ChannelMemberHistory() ChannelMemberHistoryStore
	Plugin() PluginStore
	TermsOfService() TermsOfServiceStore
	ProductNotices() ProductNoticesStore
	Group() GroupStore
	UserTermsOfService() UserTermsOfServiceStore
	LinkMetadata() LinkMetadataStore
	MarkSystemRanUnitTests()
	Close()
	LockToMaster()
	UnlockFromMaster()
	DropAllTables()
	RecycleDBConnections(d time.Duration)
	GetCurrentSchemaVersion() string
	GetDbVersion(numerical bool) (string, error)
	TotalMasterDbConnections() int
	TotalReadDbConnections() int
	TotalSearchDbConnections() int
	ReplicaLagTime() error
	ReplicaLagAbs() error
	CheckIntegrity() <-chan model.IntegrityCheckResult
	SetContext(context context.Context)
	Context() context.Context
}

type TeamStore interface {
	Save(team *model.Team) (*model.Team, error)
	Update(team *model.Team) (*model.Team, error)
	Get(id string) (*model.Team, error)
	GetByName(name string) (*model.Team, error)
	GetByNames(name []string) ([]*model.Team, error)
	SearchAll(term string, opts *model.TeamSearch) ([]*model.Team, error)
	SearchAllPaged(term string, opts *model.TeamSearch) ([]*model.Team, int64, error)
	SearchOpen(term string) ([]*model.Team, error)
	SearchPrivate(term string) ([]*model.Team, error)
	GetAll() ([]*model.Team, error)
	GetAllPage(offset int, limit int) ([]*model.Team, error)
	GetAllPrivateTeamListing() ([]*model.Team, error)
	GetAllPrivateTeamPageListing(offset int, limit int) ([]*model.Team, error)
	GetAllPublicTeamPageListing(offset int, limit int) ([]*model.Team, error)
	GetAllTeamListing() ([]*model.Team, error)
	GetAllTeamPageListing(offset int, limit int) ([]*model.Team, error)
	GetTeamsByUserId(userId string) ([]*model.Team, error)
	GetByInviteId(inviteID string) (*model.Team, error)
	PermanentDelete(teamID string) error
	AnalyticsTeamCount(includeDeleted bool) (int64, error)
	AnalyticsPublicTeamCount() (int64, error)
	AnalyticsPrivateTeamCount() (int64, error)
	SaveMultipleMembers(members []*model.TeamMember, maxUsersPerTeam int) ([]*model.TeamMember, error)
	SaveMember(member *model.TeamMember, maxUsersPerTeam int) (*model.TeamMember, error)
	UpdateMember(member *model.TeamMember) (*model.TeamMember, error)
	UpdateMultipleMembers(members []*model.TeamMember) ([]*model.TeamMember, error)
	GetMember(ctx context.Context, teamID string, userId string) (*model.TeamMember, error)
	GetMembers(teamID string, offset int, limit int, teamMembersGetOptions *model.TeamMembersGetOptions) ([]*model.TeamMember, error)
	GetMembersByIds(teamID string, userIds []string, restrictions *model.ViewUsersRestrictions) ([]*model.TeamMember, error)
	GetTotalMemberCount(teamID string, restrictions *model.ViewUsersRestrictions) (int64, error)
	GetActiveMemberCount(teamID string, restrictions *model.ViewUsersRestrictions) (int64, error)
	GetTeamsForUser(ctx context.Context, userId string) ([]*model.TeamMember, error)
	GetTeamsForUserWithPagination(userId string, page, perPage int) ([]*model.TeamMember, error)
	GetChannelUnreadsForAllTeams(excludeTeamID, userId string) ([]*model.ChannelUnread, error)
	GetChannelUnreadsForTeam(teamID, userId string) ([]*model.ChannelUnread, error)
	RemoveMember(teamID string, userId string) error
	RemoveMembers(teamID string, userIds []string) error
	RemoveAllMembersByTeam(teamID string) error
	RemoveAllMembersByUser(userId string) error
	UpdateLastTeamIconUpdate(teamID string, curTime int64) error
	GetTeamsByScheme(schemeID string, offset int, limit int) ([]*model.Team, error)
	MigrateTeamMembers(fromTeamID string, fromUserId string) (map[string]string, error)
	ResetAllTeamSchemes() error
	ClearAllCustomRoleAssignments() error
	AnalyticsGetTeamCountForScheme(schemeID string) (int64, error)
	GetAllForExportAfter(limit int, afterID string) ([]*model.TeamForExport, error)
	GetTeamMembersForExport(userId string) ([]*model.TeamMemberForExport, error)
	UserBelongsToTeams(userId string, teamIds []string) (bool, error)
	GetUserTeamIds(userId string, allowFromCache bool) ([]string, error)
	InvalidateAllTeamIdsForUser(userId string)
	ClearCaches()

	// UpdateMembersRole sets all of the given team members to admins and all of the other members of the team to
	// non-admin members.
	UpdateMembersRole(teamID string, userIDs []string) error

	// GroupSyncedTeamCount returns the count of non-deleted group-constrained teams.
	GroupSyncedTeamCount() (int64, error)
}

type ChannelStore interface {
	Save(channel *model.Channel, maxChannelsPerTeam int64) (*model.Channel, error)
	CreateDirectChannel(userId *model.User, otherUserId *model.User) (*model.Channel, error)
	SaveDirectChannel(channel *model.Channel, member1 *model.ChannelMember, member2 *model.ChannelMember) (*model.Channel, error)
	Update(channel *model.Channel) (*model.Channel, error)
	UpdateSidebarChannelCategoryOnMove(channel *model.Channel, newTeamID string) error
	ClearSidebarOnTeamLeave(userId, teamID string) error
	Get(id string, allowFromCache bool) (*model.Channel, error)
	InvalidateChannel(id string)
	InvalidateChannelByName(teamID, name string)
	GetFromMaster(id string) (*model.Channel, error)
	Delete(channelID string, time int64) error
	Restore(channelID string, time int64) error
	SetDeleteAt(channelID string, deleteAt int64, updateAt int64) error
	PermanentDelete(channelID string) error
	PermanentDeleteByTeam(teamID string) error
	GetByName(team_id string, name string, allowFromCache bool) (*model.Channel, error)
	GetByNames(team_id string, names []string, allowFromCache bool) ([]*model.Channel, error)
	GetByNameIncludeDeleted(team_id string, name string, allowFromCache bool) (*model.Channel, error)
	GetDeletedByName(team_id string, name string) (*model.Channel, error)
	GetDeleted(team_id string, offset int, limit int, userId string) (*model.ChannelList, error)
	GetChannels(teamID string, userId string, includeDeleted bool, lastDeleteAt int) (*model.ChannelList, error)
	GetAllChannels(page, perPage int, opts ChannelSearchOpts) (*model.ChannelListWithTeamData, error)
	GetAllChannelsCount(opts ChannelSearchOpts) (int64, error)
	GetMoreChannels(teamID string, userId string, offset int, limit int) (*model.ChannelList, error)
	GetPrivateChannelsForTeam(teamID string, offset int, limit int) (*model.ChannelList, error)
	GetPublicChannelsForTeam(teamID string, offset int, limit int) (*model.ChannelList, error)
	GetPublicChannelsByIdsForTeam(teamID string, channelIds []string) (*model.ChannelList, error)
	GetChannelCounts(teamID string, userId string) (*model.ChannelCounts, error)
	GetTeamChannels(teamID string) (*model.ChannelList, error)
	GetAll(teamID string) ([]*model.Channel, error)
	GetChannelsByIds(channelIds []string, includeDeleted bool) ([]*model.Channel, error)
	GetForPost(postID string) (*model.Channel, error)
	SaveMultipleMembers(members []*model.ChannelMember) ([]*model.ChannelMember, error)
	SaveMember(member *model.ChannelMember) (*model.ChannelMember, error)
	UpdateMember(member *model.ChannelMember) (*model.ChannelMember, error)
	UpdateMultipleMembers(members []*model.ChannelMember) ([]*model.ChannelMember, error)
	GetMembers(channelID string, offset, limit int) (*model.ChannelMembers, error)
	GetMember(ctx context.Context, channelID string, userId string) (*model.ChannelMember, error)
	GetChannelMembersTimezones(channelID string) ([]model.StringMap, error)
	GetAllChannelMembersForUser(userId string, allowFromCache bool, includeDeleted bool) (map[string]string, error)
	InvalidateAllChannelMembersForUser(userId string)
	IsUserInChannelUseCache(userId string, channelID string) bool
	GetAllChannelMembersNotifyPropsForChannel(channelID string, allowFromCache bool) (map[string]model.StringMap, error)
	InvalidateCacheForChannelMembersNotifyProps(channelID string)
	GetMemberForPost(postID string, userId string) (*model.ChannelMember, error)
	InvalidateMemberCount(channelID string)
	GetMemberCountFromCache(channelID string) int64
	GetMemberCount(channelID string, allowFromCache bool) (int64, error)
	GetMemberCountsByGroup(ctx context.Context, channelID string, includeTimezones bool) ([]*model.ChannelMemberCountByGroup, error)
	InvalidatePinnedPostCount(channelID string)
	GetPinnedPostCount(channelID string, allowFromCache bool) (int64, error)
	InvalidateGuestCount(channelID string)
	GetGuestCount(channelID string, allowFromCache bool) (int64, error)
	GetPinnedPosts(channelID string) (*model.PostList, error)
	RemoveMember(channelID string, userId string) error
	RemoveMembers(channelID string, userIds []string) error
	PermanentDeleteMembersByUser(userId string) error
	PermanentDeleteMembersByChannel(channelID string) error
	UpdateLastViewedAt(channelIds []string, userId string, updateThreads bool) (map[string]int64, error)
	UpdateLastViewedAtPost(unreadPost *model.Post, userID string, mentionCount, mentionCountRoot int, updateThreads bool) (*model.ChannelUnreadAt, error)
	CountPostsAfter(channelID string, timestamp int64, userId string) (int, int, error)
	IncrementMentionCount(channelID string, userId string, updateThreads, isRoot bool) error
	AnalyticsTypeCount(teamID string, channelType string) (int64, error)
	GetMembersForUser(teamID string, userId string) (*model.ChannelMembers, error)
	GetMembersForUserWithPagination(teamID, userId string, page, perPage int) (*model.ChannelMembers, error)
	AutocompleteInTeam(teamID string, term string, includeDeleted bool) (*model.ChannelList, error)
	AutocompleteInTeamForSearch(teamID string, userId string, term string, includeDeleted bool) (*model.ChannelList, error)
	SearchAllChannels(term string, opts ChannelSearchOpts) (*model.ChannelListWithTeamData, int64, error)
	SearchInTeam(teamID string, term string, includeDeleted bool) (*model.ChannelList, error)
	SearchArchivedInTeam(teamID string, term string, userId string) (*model.ChannelList, error)
	SearchForUserInTeam(userId string, teamID string, term string, includeDeleted bool) (*model.ChannelList, error)
	SearchMore(userId string, teamID string, term string) (*model.ChannelList, error)
	SearchGroupChannels(userId, term string) (*model.ChannelList, error)
	GetMembersByIds(channelID string, userIds []string) (*model.ChannelMembers, error)
	GetMembersByChannelIds(channelIds []string, userId string) (*model.ChannelMembers, error)
	AnalyticsDeletedTypeCount(teamID string, channelType string) (int64, error)
	GetChannelUnread(channelID, userId string) (*model.ChannelUnread, error)
	ClearCaches()
	GetChannelsByScheme(schemeID string, offset int, limit int) (model.ChannelList, error)
	MigrateChannelMembers(fromChannelId string, fromUserId string) (map[string]string, error)
	ResetAllChannelSchemes() error
	ClearAllCustomRoleAssignments() error
	MigratePublicChannels() error
	CreateInitialSidebarCategories(userId, teamID string) (*model.OrderedSidebarCategories, error)
	GetSidebarCategories(userId, teamID string) (*model.OrderedSidebarCategories, error)
	GetSidebarCategory(categoryID string) (*model.SidebarCategoryWithChannels, error)
	GetSidebarCategoryOrder(userId, teamID string) ([]string, error)
	CreateSidebarCategory(userId, teamID string, newCategory *model.SidebarCategoryWithChannels) (*model.SidebarCategoryWithChannels, error)
	UpdateSidebarCategoryOrder(userId, teamID string, categoryOrder []string) error
	UpdateSidebarCategories(userId, teamID string, categories []*model.SidebarCategoryWithChannels) ([]*model.SidebarCategoryWithChannels, []*model.SidebarCategoryWithChannels, error)
	UpdateSidebarChannelsByPreferences(preferences *model.Preferences) error
	DeleteSidebarChannelsByPreferences(preferences *model.Preferences) error
	DeleteSidebarCategory(categoryID string) error
	GetAllChannelsForExportAfter(limit int, afterID string) ([]*model.ChannelForExport, error)
	GetAllDirectChannelsForExportAfter(limit int, afterID string) ([]*model.DirectChannelForExport, error)
	GetChannelMembersForExport(userId string, teamID string) ([]*model.ChannelMemberForExport, error)
	RemoveAllDeactivatedMembers(channelID string) error
	GetChannelsBatchForIndexing(startTime, endTime int64, limit int) ([]*model.Channel, error)
	UserBelongsToChannels(userId string, channelIds []string) (bool, error)

	// UpdateMembersRole sets all of the given team members to admins and all of the other members of the team to
	// non-admin members.
	UpdateMembersRole(channelID string, userIDs []string) error

	// GroupSyncedChannelCount returns the count of non-deleted group-constrained channels.
	GroupSyncedChannelCount() (int64, error)
}

type ChannelMemberHistoryStore interface {
	LogJoinEvent(userId string, channelID string, joinTime int64) error
	LogLeaveEvent(userId string, channelID string, leaveTime int64) error
	GetUsersInChannelDuring(startTime int64, endTime int64, channelID string) ([]*model.ChannelMemberHistoryResult, error)
	PermanentDeleteBatch(endTime int64, limit int64) (int64, error)
}
type ThreadStore interface {
	SaveMultiple(thread []*model.Thread) ([]*model.Thread, int, error)
	Save(thread *model.Thread) (*model.Thread, error)
	Update(thread *model.Thread) (*model.Thread, error)
	Get(id string) (*model.Thread, error)
	GetThreadsForUser(userId, teamId string, opts model.GetUserThreadsOpts) (*model.Threads, error)
	GetThreadForUser(userId, teamId, threadId string, extended bool) (*model.ThreadResponse, error)
	Delete(postId string) error
	GetPosts(threadId string, since int64) ([]*model.Post, error)

	MarkAllAsRead(userId, teamID string) error
	MarkAsRead(userId, threadID string, timestamp int64) error

	SaveMembership(membership *model.ThreadMembership) (*model.ThreadMembership, error)
	UpdateMembership(membership *model.ThreadMembership) (*model.ThreadMembership, error)
	GetMembershipsForUser(userId, teamID string) ([]*model.ThreadMembership, error)
	GetMembershipForUser(userId, postID string) (*model.ThreadMembership, error)
	DeleteMembershipForUser(userId, postID string) error
	MaintainMembership(userId, postID string, following, incrementMentions, updateFollowing, updateViewedTimestamp bool) error
	CollectThreadsWithNewerReplies(userId string, channelIds []string, timestamp int64) ([]string, error)
	UpdateUnreadsByChannel(userId string, changedThreads []string, timestamp int64, updateViewedTimestamp bool) error
}

type PostStore interface {
	SaveMultiple(posts []*model.Post) ([]*model.Post, int, error)
	Save(post *model.Post) (*model.Post, error)
	Update(newPost *model.Post, oldPost *model.Post) (*model.Post, error)
	Get(ctx context.Context, id string, skipFetchThreads, collapsedThreads, collapsedThreadsExtended bool, userID string) (*model.PostList, error)
	GetSingle(id string) (*model.Post, error)
	Delete(postID string, time int64, deleteByID string) error
	PermanentDeleteByUser(userId string) error
	PermanentDeleteByChannel(channelID string) error
	GetPosts(options model.GetPostsOptions, allowFromCache bool) (*model.PostList, error)
	GetFlaggedPosts(userId string, offset int, limit int) (*model.PostList, error)
	// @openTracingParams userId, teamID, offset, limit
	GetFlaggedPostsForTeam(userId, teamID string, offset int, limit int) (*model.PostList, error)
	GetFlaggedPostsForChannel(userId, channelID string, offset int, limit int) (*model.PostList, error)
	GetPostsBefore(options model.GetPostsOptions) (*model.PostList, error)
	GetPostsAfter(options model.GetPostsOptions) (*model.PostList, error)
	GetPostsSince(options model.GetPostsSinceOptions, allowFromCache bool) (*model.PostList, error)
	GetPostAfterTime(channelID string, time int64, collapsedThreads bool) (*model.Post, error)
	GetPostIdAfterTime(channelID string, time int64, collapsedThreads bool) (string, error)
	GetPostIdBeforeTime(channelID string, time int64, collapsedThreads bool) (string, error)
	GetEtag(channelID string, allowFromCache bool, collapsedThreads bool) string
	Search(teamID string, userId string, params *model.SearchParams) (*model.PostList, error)
	AnalyticsUserCountsWithPostsByDay(teamID string) (model.AnalyticsRows, error)
	AnalyticsPostCountsByDay(options *model.AnalyticsPostCountsOptions) (model.AnalyticsRows, error)
	AnalyticsPostCount(teamID string, mustHaveFile bool, mustHaveHashtag bool) (int64, error)
	ClearCaches()
	InvalidateLastPostTimeCache(channelID string)
	GetPostsCreatedAt(channelID string, time int64) ([]*model.Post, error)
	Overwrite(post *model.Post) (*model.Post, error)
	OverwriteMultiple(posts []*model.Post) ([]*model.Post, int, error)
	GetPostsByIds(postIds []string) ([]*model.Post, error)
	GetPostsBatchForIndexing(startTime int64, endTime int64, limit int) ([]*model.PostForIndexing, error)
	PermanentDeleteBatch(endTime int64, limit int64) (int64, error)
	GetOldest() (*model.Post, error)
	GetMaxPostSize() int
	GetParentsForExportAfter(limit int, afterID string) ([]*model.PostForExport, error)
	GetRepliesForExport(parentID string) ([]*model.ReplyForExport, error)
	GetDirectPostParentsForExportAfter(limit int, afterID string) ([]*model.DirectPostForExport, error)
	SearchPostsInTeamForUser(paramsList []*model.SearchParams, userId, teamID string, page, perPage int) (*model.PostSearchResults, error)
	GetOldestEntityCreationTime() (int64, error)
}

type UserStore interface {
	Save(user *model.User) (*model.User, error)
	Update(user *model.User, allowRoleUpdate bool) (*model.UserUpdate, error)
	UpdateLastPictureUpdate(userId string) error
	ResetLastPictureUpdate(userId string) error
	UpdatePassword(userId, newPassword string) error
	UpdateUpdateAt(userId string) (int64, error)
	UpdateAuthData(userId string, service string, authData *string, email string, resetMfa bool) (string, error)
	UpdateMfaSecret(userId, secret string) error
	UpdateMfaActive(userId string, active bool) error
	Get(ctx context.Context, id string) (*model.User, error)
	GetMany(ctx context.Context, ids []string) ([]*model.User, error)
	GetAll() ([]*model.User, error)
	ClearCaches()
	InvalidateProfilesInChannelCacheByUser(userId string)
	InvalidateProfilesInChannelCache(channelID string)
	GetProfilesInChannel(options *model.UserGetOptions) ([]*model.User, error)
	GetProfilesInChannelByStatus(options *model.UserGetOptions) ([]*model.User, error)
	GetAllProfilesInChannel(ctx context.Context, channelId string, allowFromCache bool) (map[string]*model.User, error)
	GetProfilesNotInChannel(teamId string, channelId string, groupConstrained bool, offset int, limit int, viewRestrictions *model.ViewUsersRestrictions) ([]*model.User, error)
	GetProfilesWithoutTeam(options *model.UserGetOptions) ([]*model.User, error)
	GetProfilesByUsernames(usernames []string, viewRestrictions *model.ViewUsersRestrictions) ([]*model.User, error)
	GetAllProfiles(options *model.UserGetOptions) ([]*model.User, error)
	GetProfiles(options *model.UserGetOptions) ([]*model.User, error)
	GetProfileByIds(ctx context.Context, userIds []string, options *UserGetByIdsOpts, allowFromCache bool) ([]*model.User, error)
	GetProfileByGroupChannelIdsForUser(userId string, channelIds []string) (map[string][]*model.User, error)
	InvalidateProfileCacheForUser(userId string)
	GetByEmail(email string) (*model.User, error)
	GetByAuth(authData *string, authService string) (*model.User, error)
	GetAllUsingAuthService(authService string) ([]*model.User, error)
	GetAllNotInAuthService(authServices []string) ([]*model.User, error)
	GetByUsername(username string) (*model.User, error)
	GetForLogin(loginID string, allowSignInWithUsername, allowSignInWithEmail bool) (*model.User, error)
	VerifyEmail(userId, email string) (string, error)
	GetEtagForAllProfiles() string
	GetEtagForProfiles(teamID string) string
	UpdateFailedPasswordAttempts(userId string, attempts int) error
	GetSystemAdminProfiles() (map[string]*model.User, error)
	PermanentDelete(userId string) error
	AnalyticsActiveCount(time int64, options model.UserCountOptions) (int64, error)
	AnalyticsActiveCountForPeriod(startTime int64, endTime int64, options model.UserCountOptions) (int64, error)
	GetUnreadCount(userId string) (int64, error)
	GetUnreadCountForChannel(userId string, channelID string) (int64, error)
	GetAnyUnreadPostCountForChannel(userId string, channelID string) (int64, error)
	GetRecentlyActiveUsersForTeam(teamID string, offset, limit int, viewRestrictions *model.ViewUsersRestrictions) ([]*model.User, error)
	GetNewUsersForTeam(teamID string, offset, limit int, viewRestrictions *model.ViewUsersRestrictions) ([]*model.User, error)
	Search(teamID string, term string, options *model.UserSearchOptions) ([]*model.User, error)
	SearchNotInTeam(notInTeamId string, term string, options *model.UserSearchOptions) ([]*model.User, error)
	SearchInChannel(channelID string, term string, options *model.UserSearchOptions) ([]*model.User, error)
	SearchNotInChannel(teamID string, channelID string, term string, options *model.UserSearchOptions) ([]*model.User, error)
	SearchWithoutTeam(term string, options *model.UserSearchOptions) ([]*model.User, error)
	SearchInGroup(groupID string, term string, options *model.UserSearchOptions) ([]*model.User, error)
	AnalyticsGetInactiveUsersCount() (int64, error)
	AnalyticsGetExternalUsers(hostDomain string) (bool, error)
	AnalyticsGetSystemAdminCount() (int64, error)
	AnalyticsGetGuestCount() (int64, error)
	GetProfilesNotInTeam(teamID string, groupConstrained bool, offset int, limit int, viewRestrictions *model.ViewUsersRestrictions) ([]*model.User, error)
	GetEtagForProfilesNotInTeam(teamID string) string
	ClearAllCustomRoleAssignments() error
	InferSystemInstallDate() (int64, error)
	GetAllAfter(limit int, afterID string) ([]*model.User, error)
	GetUsersBatchForIndexing(startTime, endTime int64, limit int) ([]*model.UserForIndexing, error)
	Count(options model.UserCountOptions) (int64, error)
	GetTeamGroupUsers(teamID string) ([]*model.User, error)
	GetChannelGroupUsers(channelID string) ([]*model.User, error)
	PromoteGuestToUser(userID string) error
	DemoteUserToGuest(userID string) (*model.User, error)
	DeactivateGuests() ([]string, error)
	AutocompleteUsersInChannel(teamID, channelID, term string, options *model.UserSearchOptions) (*model.UserAutocompleteInChannel, error)
	GetKnownUsers(userID string) ([]string, error)
}

type BotStore interface {
	Get(userId string, includeDeleted bool) (*model.Bot, error)
	GetAll(options *model.BotGetOptions) ([]*model.Bot, error)
	Save(bot *model.Bot) (*model.Bot, error)
	Update(bot *model.Bot) (*model.Bot, error)
	PermanentDelete(userId string) error
}

type SessionStore interface {
	Get(ctx context.Context, sessionIDOrToken string) (*model.Session, error)
	Save(session *model.Session) (*model.Session, error)
	GetSessions(userId string) ([]*model.Session, error)
	GetSessionsWithActiveDeviceIds(userId string) ([]*model.Session, error)
	GetSessionsExpired(thresholdMillis int64, mobileOnly bool, unnotifiedOnly bool) ([]*model.Session, error)
	UpdateExpiredNotify(sessionid string, notified bool) error
	Remove(sessionIDOrToken string) error
	RemoveAllSessions() error
	PermanentDeleteSessionsByUser(teamID string) error
	UpdateExpiresAt(sessionID string, time int64) error
	UpdateLastActivityAt(sessionID string, time int64) error
	UpdateRoles(userId string, roles string) (string, error)
	UpdateDeviceId(id string, deviceID string, expiresAt int64) (string, error)
	UpdateProps(session *model.Session) error
	AnalyticsSessionCount() (int64, error)
	Cleanup(expiryTime int64, batchSize int64)
}

type AuditStore interface {
	Save(audit *model.Audit) error
	Get(user_id string, offset int, limit int) (model.Audits, error)
	PermanentDeleteByUser(userId string) error
}

type ClusterDiscoveryStore interface {
	Save(discovery *model.ClusterDiscovery) error
	Delete(discovery *model.ClusterDiscovery) (bool, error)
	Exists(discovery *model.ClusterDiscovery) (bool, error)
	GetAll(discoveryType, clusterName string) ([]*model.ClusterDiscovery, error)
	SetLastPingAt(discovery *model.ClusterDiscovery) error
	Cleanup() error
}

type ComplianceStore interface {
	Save(compliance *model.Compliance) (*model.Compliance, error)
	Update(compliance *model.Compliance) (*model.Compliance, error)
	Get(id string) (*model.Compliance, error)
	GetAll(offset, limit int) (model.Compliances, error)
	ComplianceExport(compliance *model.Compliance) ([]*model.CompliancePost, error)
	MessageExport(after int64, limit int) ([]*model.MessageExport, error)
}

type OAuthStore interface {
	SaveApp(app *model.OAuthApp) (*model.OAuthApp, error)
	UpdateApp(app *model.OAuthApp) (*model.OAuthApp, error)
	GetApp(id string) (*model.OAuthApp, error)
	GetAppByUser(userId string, offset, limit int) ([]*model.OAuthApp, error)
	GetApps(offset, limit int) ([]*model.OAuthApp, error)
	GetAuthorizedApps(userId string, offset, limit int) ([]*model.OAuthApp, error)
	DeleteApp(id string) error
	SaveAuthData(authData *model.AuthData) (*model.AuthData, error)
	GetAuthData(code string) (*model.AuthData, error)
	RemoveAuthData(code string) error
	PermanentDeleteAuthDataByUser(userId string) error
	SaveAccessData(accessData *model.AccessData) (*model.AccessData, error)
	UpdateAccessData(accessData *model.AccessData) (*model.AccessData, error)
	GetAccessData(token string) (*model.AccessData, error)
	GetAccessDataByUserForApp(userId, clientId string) ([]*model.AccessData, error)
	GetAccessDataByRefreshToken(token string) (*model.AccessData, error)
	GetPreviousAccessData(userId, clientId string) (*model.AccessData, error)
	RemoveAccessData(token string) error
	RemoveAllAccessData() error
}

type SystemStore interface {
	Save(system *model.System) error
	SaveOrUpdate(system *model.System) error
	Update(system *model.System) error
	Get() (model.StringMap, error)
	GetByName(name string) (*model.System, error)
	PermanentDeleteByName(name string) (*model.System, error)
	InsertIfExists(system *model.System) (*model.System, error)
	SaveOrUpdateWithWarnMetricHandling(system *model.System) error
}

type WebhookStore interface {
	SaveIncoming(webhook *model.IncomingWebhook) (*model.IncomingWebhook, error)
	GetIncoming(id string, allowFromCache bool) (*model.IncomingWebhook, error)
	GetIncomingList(offset, limit int) ([]*model.IncomingWebhook, error)
	GetIncomingListByUser(userId string, offset, limit int) ([]*model.IncomingWebhook, error)
	GetIncomingByTeam(teamID string, offset, limit int) ([]*model.IncomingWebhook, error)
	GetIncomingByTeamByUser(teamID string, userId string, offset, limit int) ([]*model.IncomingWebhook, error)
	UpdateIncoming(webhook *model.IncomingWebhook) (*model.IncomingWebhook, error)
	GetIncomingByChannel(channelID string) ([]*model.IncomingWebhook, error)
	DeleteIncoming(webhookID string, time int64) error
	PermanentDeleteIncomingByChannel(channelID string) error
	PermanentDeleteIncomingByUser(userId string) error

	SaveOutgoing(webhook *model.OutgoingWebhook) (*model.OutgoingWebhook, error)
	GetOutgoing(id string) (*model.OutgoingWebhook, error)
	GetOutgoingByChannel(channelID string, offset, limit int) ([]*model.OutgoingWebhook, error)
	GetOutgoingByChannelByUser(channelID string, userId string, offset, limit int) ([]*model.OutgoingWebhook, error)
	GetOutgoingList(offset, limit int) ([]*model.OutgoingWebhook, error)
	GetOutgoingListByUser(userId string, offset, limit int) ([]*model.OutgoingWebhook, error)
	GetOutgoingByTeam(teamID string, offset, limit int) ([]*model.OutgoingWebhook, error)
	GetOutgoingByTeamByUser(teamID string, userId string, offset, limit int) ([]*model.OutgoingWebhook, error)
	DeleteOutgoing(webhookID string, time int64) error
	PermanentDeleteOutgoingByChannel(channelID string) error
	PermanentDeleteOutgoingByUser(userId string) error
	UpdateOutgoing(hook *model.OutgoingWebhook) (*model.OutgoingWebhook, error)

	AnalyticsIncomingCount(teamID string) (int64, error)
	AnalyticsOutgoingCount(teamID string) (int64, error)
	InvalidateWebhookCache(webhook string)
	ClearCaches()
}

type CommandStore interface {
	Save(webhook *model.Command) (*model.Command, error)
	GetByTrigger(teamID string, trigger string) (*model.Command, error)
	Get(id string) (*model.Command, error)
	GetByTeam(teamID string) ([]*model.Command, error)
	Delete(commandID string, time int64) error
	PermanentDeleteByTeam(teamID string) error
	PermanentDeleteByUser(userId string) error
	Update(hook *model.Command) (*model.Command, error)
	AnalyticsCommandCount(teamID string) (int64, error)
}

type CommandWebhookStore interface {
	Save(webhook *model.CommandWebhook) (*model.CommandWebhook, error)
	Get(id string) (*model.CommandWebhook, error)
	TryUse(id string, limit int) error
	Cleanup()
}

type PreferenceStore interface {
	Save(preferences *model.Preferences) error
	GetCategory(userId string, category string) (model.Preferences, error)
	Get(userId string, category string, name string) (*model.Preference, error)
	GetAll(userId string) (model.Preferences, error)
	Delete(userId, category, name string) error
	DeleteCategory(userId string, category string) error
	DeleteCategoryAndName(category string, name string) error
	PermanentDeleteByUser(userId string) error
	CleanupFlagsBatch(limit int64) (int64, error)
}

type LicenseStore interface {
	Save(license *model.LicenseRecord) (*model.LicenseRecord, error)
	Get(id string) (*model.LicenseRecord, error)
}

type TokenStore interface {
	Save(recovery *model.Token) error
	Delete(token string) error
	GetByToken(token string) (*model.Token, error)
	Cleanup()
	RemoveAllTokensByType(tokenType string) error
}

type EmojiStore interface {
	Save(emoji *model.Emoji) (*model.Emoji, error)
	Get(ctx context.Context, id string, allowFromCache bool) (*model.Emoji, error)
	GetByName(ctx context.Context, name string, allowFromCache bool) (*model.Emoji, error)
	GetMultipleByName(names []string) ([]*model.Emoji, error)
	GetList(offset, limit int, sort string) ([]*model.Emoji, error)
	Delete(emoji *model.Emoji, time int64) error
	Search(name string, prefixOnly bool, limit int) ([]*model.Emoji, error)
}

type StatusStore interface {
	SaveOrUpdate(status *model.Status) error
	Get(userId string) (*model.Status, error)
	GetByIds(userIds []string) ([]*model.Status, error)
	ResetAll() error
	GetTotalActiveUsersCount() (int64, error)
	UpdateLastActivityAt(userId string, lastActivityAt int64) error
}

type FileInfoStore interface {
	Save(info *model.FileInfo) (*model.FileInfo, error)
	Upsert(info *model.FileInfo) (*model.FileInfo, error)
	Get(id string) (*model.FileInfo, error)
	GetByIds(ids []string) ([]*model.FileInfo, error)
	GetByPath(path string) (*model.FileInfo, error)
	GetForPost(postID string, readFromMaster, includeDeleted, allowFromCache bool) ([]*model.FileInfo, error)
	GetForUser(userId string) ([]*model.FileInfo, error)
	GetWithOptions(page, perPage int, opt *model.GetFileInfosOptions) ([]*model.FileInfo, error)
	InvalidateFileInfosForPostCache(postID string, deleted bool)
	AttachToPost(fileID string, postID string, creatorId string) error
	DeleteForPost(postID string) (string, error)
	PermanentDelete(fileID string) error
	PermanentDeleteBatch(endTime int64, limit int64) (int64, error)
	PermanentDeleteByUser(userId string) (int64, error)
	SetContent(fileID, content string) error
	Search(paramsList []*model.SearchParams, userId, teamID string, page, perPage int) (*model.FileInfoList, error)
	CountAll() (int64, error)
	GetFilesBatchForIndexing(startTime, endTime int64, limit int) ([]*model.FileForIndexing, error)
	ClearCaches()
}

type UploadSessionStore interface {
	Save(session *model.UploadSession) (*model.UploadSession, error)
	Update(session *model.UploadSession) error
	Get(id string) (*model.UploadSession, error)
	GetForUser(userId string) ([]*model.UploadSession, error)
	Delete(id string) error
}

type ReactionStore interface {
	Save(reaction *model.Reaction) (*model.Reaction, error)
	Delete(reaction *model.Reaction) (*model.Reaction, error)
	GetForPost(postID string, allowFromCache bool) ([]*model.Reaction, error)
	DeleteAllWithEmojiName(emojiName string) error
	PermanentDeleteBatch(endTime int64, limit int64) (int64, error)
	BulkGetForPosts(postIds []string) ([]*model.Reaction, error)
}

type JobStore interface {
	Save(job *model.Job) (*model.Job, error)
	UpdateOptimistically(job *model.Job, currentStatus string) (bool, error)
	UpdateStatus(id string, status string) (*model.Job, error)
	UpdateStatusOptimistically(id string, currentStatus string, newStatus string) (bool, error)
	Get(id string) (*model.Job, error)
	GetAllPage(offset int, limit int) ([]*model.Job, error)
	GetAllByType(jobType string) ([]*model.Job, error)
	GetAllByTypePage(jobType string, offset int, limit int) ([]*model.Job, error)
	GetAllByStatus(status string) ([]*model.Job, error)
	GetNewestJobByStatusAndType(status string, jobType string) (*model.Job, error)
	GetNewestJobByStatusesAndType(statuses []string, jobType string) (*model.Job, error)
	GetCountByStatusAndType(status string, jobType string) (int64, error)
	Delete(id string) (string, error)
}

type UserAccessTokenStore interface {
	Save(token *model.UserAccessToken) (*model.UserAccessToken, error)
	DeleteAllForUser(userId string) error
	Delete(tokenID string) error
	Get(tokenID string) (*model.UserAccessToken, error)
	GetAll(offset int, limit int) ([]*model.UserAccessToken, error)
	GetByToken(tokenString string) (*model.UserAccessToken, error)
	GetByUser(userId string, page, perPage int) ([]*model.UserAccessToken, error)
	Search(term string) ([]*model.UserAccessToken, error)
	UpdateTokenEnable(tokenID string) error
	UpdateTokenDisable(tokenID string) error
}

type PluginStore interface {
	SaveOrUpdate(keyVal *model.PluginKeyValue) (*model.PluginKeyValue, error)
	CompareAndSet(keyVal *model.PluginKeyValue, oldValue []byte) (bool, error)
	CompareAndDelete(keyVal *model.PluginKeyValue, oldValue []byte) (bool, error)
	SetWithOptions(pluginID string, key string, value []byte, options model.PluginKVSetOptions) (bool, error)
	Get(pluginID, key string) (*model.PluginKeyValue, error)
	Delete(pluginID, key string) error
	DeleteAllForPlugin(PluginID string) error
	DeleteAllExpired() error
	List(pluginID string, page, perPage int) ([]string, error)
}

type RoleStore interface {
	Save(role *model.Role) (*model.Role, error)
	Get(roleID string) (*model.Role, error)
	GetAll() ([]*model.Role, error)
	GetByName(name string) (*model.Role, error)
	GetByNames(names []string) ([]*model.Role, error)
	Delete(roleID string) (*model.Role, error)
	PermanentDeleteAll() error

	// HigherScopedPermissions retrieves the higher-scoped permissions of a list of role names. The higher-scope
	// (either team scheme or system scheme) is determined based on whether the team has a scheme or not.
	ChannelHigherScopedPermissions(roleNames []string) (map[string]*model.RolePermissions, error)

	// AllChannelSchemeRoles returns all of the roles associated to channel schemes.
	AllChannelSchemeRoles() ([]*model.Role, error)

	// ChannelRolesUnderTeamRole returns all of the non-deleted roles that are affected by updates to the
	// given role.
	ChannelRolesUnderTeamRole(roleName string) ([]*model.Role, error)
}

type SchemeStore interface {
	Save(scheme *model.Scheme) (*model.Scheme, error)
	Get(schemeID string) (*model.Scheme, error)
	GetByName(schemeName string) (*model.Scheme, error)
	GetAllPage(scope string, offset int, limit int) ([]*model.Scheme, error)
	Delete(schemeID string) (*model.Scheme, error)
	PermanentDeleteAll() error
	CountByScope(scope string) (int64, error)
	CountWithoutPermission(scope, permissionID string, roleScope model.RoleScope, roleType model.RoleType) (int64, error)
}

type TermsOfServiceStore interface {
	Save(termsOfService *model.TermsOfService) (*model.TermsOfService, error)
	GetLatest(allowFromCache bool) (*model.TermsOfService, error)
	Get(id string, allowFromCache bool) (*model.TermsOfService, error)
}

type ProductNoticesStore interface {
	View(userId string, notices []string) error
	Clear(notices []string) error
	ClearOldNotices(currentNotices *model.ProductNotices) error
	GetViews(userId string) ([]model.ProductNoticeViewState, error)
}

type UserTermsOfServiceStore interface {
	GetByUser(userId string) (*model.UserTermsOfService, error)
	Save(userTermsOfService *model.UserTermsOfService) (*model.UserTermsOfService, error)
	Delete(userId, termsOfServiceId string) error
}

type GroupStore interface {
	Create(group *model.Group) (*model.Group, error)
	Get(groupID string) (*model.Group, error)
	GetByName(name string, opts model.GroupSearchOpts) (*model.Group, error)
	GetByIDs(groupIDs []string) ([]*model.Group, error)
	GetByRemoteID(remoteID string, groupSource model.GroupSource) (*model.Group, error)
	GetAllBySource(groupSource model.GroupSource) ([]*model.Group, error)
	GetByUser(userId string) ([]*model.Group, error)
	Update(group *model.Group) (*model.Group, error)
	Delete(groupID string) (*model.Group, error)

	GetMemberUsers(groupID string) ([]*model.User, error)
	GetMemberUsersPage(groupID string, page int, perPage int) ([]*model.User, error)
	GetMemberCount(groupID string) (int64, error)

	GetMemberUsersInTeam(groupID string, teamID string) ([]*model.User, error)
	GetMemberUsersNotInChannel(groupID string, channelID string) ([]*model.User, error)

	UpsertMember(groupID string, userID string) (*model.GroupMember, error)
	DeleteMember(groupID string, userID string) (*model.GroupMember, error)
	PermanentDeleteMembersByUser(userId string) error

	CreateGroupSyncable(groupSyncable *model.GroupSyncable) (*model.GroupSyncable, error)
	GetGroupSyncable(groupID string, syncableID string, syncableType model.GroupSyncableType) (*model.GroupSyncable, error)
	GetAllGroupSyncablesByGroupId(groupID string, syncableType model.GroupSyncableType) ([]*model.GroupSyncable, error)
	UpdateGroupSyncable(groupSyncable *model.GroupSyncable) (*model.GroupSyncable, error)
	DeleteGroupSyncable(groupID string, syncableID string, syncableType model.GroupSyncableType) (*model.GroupSyncable, error)

	// TeamMembersToAdd returns a slice of UserTeamIDPair that need newly created memberships
	// based on the groups configurations. The returned list can be optionally scoped to a single given team.
	//
	// Typically since will be the last successful group sync time.
	TeamMembersToAdd(since int64, teamID *string) ([]*model.UserTeamIDPair, error)

	// ChannelMembersToAdd returns a slice of UserChannelIDPair that need newly created memberships
	// based on the groups configurations. The returned list can be optionally scoped to a single given channel.
	//
	// Typically since will be the last successful group sync time.
	ChannelMembersToAdd(since int64, channelID *string) ([]*model.UserChannelIDPair, error)

	// TeamMembersToRemove returns all team members that should be removed based on group constraints.
	TeamMembersToRemove(teamID *string) ([]*model.TeamMember, error)

	// ChannelMembersToRemove returns all channel members that should be removed based on group constraints.
	ChannelMembersToRemove(channelID *string) ([]*model.ChannelMember, error)

	GetGroupsByChannel(channelID string, opts model.GroupSearchOpts) ([]*model.GroupWithSchemeAdmin, error)
	CountGroupsByChannel(channelID string, opts model.GroupSearchOpts) (int64, error)

	GetGroupsByTeam(teamID string, opts model.GroupSearchOpts) ([]*model.GroupWithSchemeAdmin, error)
	GetGroupsAssociatedToChannelsByTeam(teamID string, opts model.GroupSearchOpts) (map[string][]*model.GroupWithSchemeAdmin, error)
	CountGroupsByTeam(teamID string, opts model.GroupSearchOpts) (int64, error)

	GetGroups(page, perPage int, opts model.GroupSearchOpts) ([]*model.Group, error)

	TeamMembersMinusGroupMembers(teamID string, groupIDs []string, page, perPage int) ([]*model.UserWithGroups, error)
	CountTeamMembersMinusGroupMembers(teamID string, groupIDs []string) (int64, error)
	ChannelMembersMinusGroupMembers(channelID string, groupIDs []string, page, perPage int) ([]*model.UserWithGroups, error)
	CountChannelMembersMinusGroupMembers(channelID string, groupIDs []string) (int64, error)

	// AdminRoleGroupsForSyncableMember returns the IDs of all of the groups that the user is a member of that are
	// configured as SchemeAdmin: true for the given syncable.
	AdminRoleGroupsForSyncableMember(userID, syncableID string, syncableType model.GroupSyncableType) ([]string, error)

	// PermittedSyncableAdmins returns the IDs of all of the user who are permitted by the group syncable to have
	// the admin role for the given syncable.
	PermittedSyncableAdmins(syncableID string, syncableType model.GroupSyncableType) ([]string, error)

	// GroupCount returns the total count of records in the UserGroups table.
	GroupCount() (int64, error)

	// GroupTeamCount returns the total count of records in the GroupTeams table.
	GroupTeamCount() (int64, error)

	// GroupChannelCount returns the total count of records in the GroupChannels table.
	GroupChannelCount() (int64, error)

	// GroupMemberCount returns the total count of records in the GroupMembers table.
	GroupMemberCount() (int64, error)

	// DistinctGroupMemberCount returns the count of records in the GroupMembers table with distinct UserId values.
	DistinctGroupMemberCount() (int64, error)

	// GroupCountWithAllowReference returns the count of records in the Groups table with AllowReference set to true.
	GroupCountWithAllowReference() (int64, error)
}

type LinkMetadataStore interface {
	Save(linkMetadata *model.LinkMetadata) (*model.LinkMetadata, error)
	Get(url string, timestamp int64) (*model.LinkMetadata, error)
}

// ChannelSearchOpts contains options for searching channels.
//
// NotAssociatedToGroup will exclude channels that have associated, active GroupChannels records.
// IncludeDeleted will include channel records where DeleteAt != 0.
// ExcludeChannelNames will exclude channels from the results by name.
// Paginate whether to paginate the results.
// Page page requested, if results are paginated.
// PerPage number of results per page, if paginated.
//
type ChannelSearchOpts struct {
	NotAssociatedToGroup    string
	IncludeDeleted          bool
	Deleted                 bool
	ExcludeChannelNames     []string
	TeamIds                 []string
	GroupConstrained        bool
	ExcludeGroupConstrained bool
	Public                  bool
	Private                 bool
	Page                    *int
	PerPage                 *int
}

func (c *ChannelSearchOpts) IsPaginated() bool {
	return c.Page != nil && c.PerPage != nil
}

type UserGetByIdsOpts struct {
	// IsAdmin tracks whether or not the request is being made by an administrator. Does nothing when provided by a client.
	IsAdmin bool

	// Restrict to search in a list of teams and channels. Does nothing when provided by a client.
	ViewRestrictions *model.ViewUsersRestrictions

	// Since filters the users based on their UpdateAt timestamp.
	Since int64
}
