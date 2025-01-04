package messages

type Register struct {
	Username string
}

type CreateSubreddit struct {
	Name        string
	Description string
}

type JoinSubreddit struct {
	SubredditName string
	Username      string
}

type LeaveSubreddit struct {
	SubredditName string
	Username      string
}

type Post struct {
	ID            string
	SubredditName string
	Username      string
	Title         string
	Content       string
	Karma         int
}

type Comment struct {
	ID       string
	PostID   string
	ParentID string // Empty if top-level comment
	Username string
	Content  string
	IsUpvote bool // Add this field
}
type Response struct {
	Success bool
	Message string
}
type PostInSubreddit struct {
	SubredditName string

	Username string

	Title string

	Content string
}
type Vote struct {
	PostID   string
	Username string
	IsUpvote bool
}
type SimulateAction struct {
	Action string
	PostID string
}
type GetKarma struct {
	PostID string
}

type KarmaResponse struct {
	PostID string
	Karma  int
}
type GetFeed struct {
	Username string
}

type FeedResponse struct {
	Posts []Post
}

type GetDirectMessages struct {
	Username string
}

type DirectMessageResponse struct {
	Messages []DirectMessage
}

type DirectMessage struct {
	From string

	To string

	Content string

	Timestamp int64
}

type ReplyDirectMessage struct {
	To string

	Content string

	Timestamp int64
}

type EngineReady struct{}
