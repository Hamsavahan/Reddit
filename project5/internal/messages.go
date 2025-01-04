package messages

// User-related messages

// Subreddit-related messages
type CreateSubreddit struct {
	Name        string
	Description string
	CreatorID   string
}

type JoinSubreddit struct {
	UserID      string
	SubredditID string
}

type LeaveSubreddit struct {
	UserID      string
	SubredditID string
}

// Post-related messages

// Comment-related messages
type CreateComment struct {
	Content  string
	AuthorID string
	PostID   string
	ParentID string
}

// Voting messages
type Vote struct {
	UserID   string
	TargetID string
	IsUpvote bool
}

type VoteRecorded struct{}

// Feed-related messages
type GetFeed struct {
	UserID string
}

// Utility messages
type GetSubredditPID struct {
	ID string
}

type GetUserPID struct {
	ID string
}

type GetPostPID struct {
	ID string
}

type GetCommentPID struct {
	ID string
}

type GetPost struct {
	PostID string
}

type GetComments struct {
	PostID string
}

type GetComment struct {
	CommentID string
}

type Error struct {
	Message string
}
