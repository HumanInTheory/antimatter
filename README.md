# Antimatter

To build:
'''
$env:GOOS="linux"
$env:GOARCH="arm"
go build plugin_linux_arm.go
'''

Reference Snippets:
'''
type SearchParams struct {
	Terms                  string
	ExcludedTerms          string
	IsHashtag              bool
	InChannels             []string
	ExcludedChannels       []string
	FromUsers              []string
	ExcludedUsers          []string
	AfterDate              string
	ExcludedAfterDate      string
	BeforeDate             string
	ExcludedBeforeDate     string
	OnDate                 string
	ExcludedDate           string
	OrTerms                bool
	IncludeDeletedChannels bool
	TimeZoneOffset         int
	// True if this search doesn't originate from a "current user".
	SearchWithoutUserId bool
}
SearchPostsInTeam(teamId string, paramsList []*model.SearchParams) ([]*model.Post, *model.AppError)
DeletePost(postId string) *model.AppError
'''

Pseudo-Code:

Create /antimatter command
On execution, if mode1 is clear and mode2 is channel:
Confirm user is admin of channel, else return error
Confirm user really wants to nuke channel, else return cancel
Get all posts in channel
Delete those posts
Return ephemeral confirmation