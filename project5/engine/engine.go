package engine

import (
	"fmt"
	"protoactor-simulation/messages"
	"sort"
	"sync"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/google/uuid"
)

type Engine struct {
	users          map[string]*actor.PID
	subreddits     map[string]map[string]bool
	posts          map[string]messages.Post
	comments       map[string][]messages.Comment
	directMessages map[string][]messages.DirectMessage
	mu             sync.RWMutex
}

func NewEngine() *Engine {
	return &Engine{
		users:      make(map[string]*actor.PID),
		subreddits: make(map[string]map[string]bool),
		posts:      make(map[string]messages.Post),
		comments:   make(map[string][]messages.Comment),
	}
}

func (e *Engine) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *actor.Started:
		fmt.Println("Engine started")
	case *messages.Register:
		e.handleRegister(context, msg)
	case *messages.CreateSubreddit:
		e.handleCreateSubreddit(context, msg)
	case *messages.JoinSubreddit:
		e.handleJoinSubreddit(context, msg)
	case *messages.LeaveSubreddit:
		e.handleLeaveSubreddit(context, msg)
	case *messages.Post:
		e.handlePost(context, msg)
	case *messages.Comment:
		e.handleComment(context, msg)
	case *messages.Vote:
		e.handleVote(context, msg)
	case *messages.GetFeed:
		e.handleGetFeed(context, msg)
	case *messages.GetDirectMessages:
		e.handleGetDirectMessages(context, msg)
	case *messages.ReplyDirectMessage:
		e.handleReplyDirectMessage(context, msg)

	}
}

func (e *Engine) handleRegister(context actor.Context, msg *messages.Register) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if _, exists := e.users[msg.Username]; !exists {
		e.users[msg.Username] = context.Sender()
		context.Respond(&messages.Response{Success: true, Message: "Registered successfully"})
	} else {
		context.Respond(&messages.Response{Success: false, Message: "Username already exists"})
	}
}

func (e *Engine) handleCreateSubreddit(context actor.Context, msg *messages.CreateSubreddit) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if _, exists := e.subreddits[msg.Name]; !exists {
		e.subreddits[msg.Name] = make(map[string]bool)
		context.Respond(&messages.Response{Success: true, Message: "Subreddit created successfully"})
	} else {
		context.Respond(&messages.Response{Success: false, Message: "Subreddit already exists"})
	}
}

func (e *Engine) handleJoinSubreddit(context actor.Context, msg *messages.JoinSubreddit) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if subreddit, exists := e.subreddits[msg.SubredditName]; exists {
		subreddit[msg.Username] = true
		context.Respond(&messages.Response{Success: true, Message: "Joined subreddit successfully"})
	} else {
		context.Respond(&messages.Response{Success: false, Message: "Subreddit does not exist"})
	}
}

func (e *Engine) handleLeaveSubreddit(context actor.Context, msg *messages.LeaveSubreddit) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if subreddit, exists := e.subreddits[msg.SubredditName]; exists {
		delete(subreddit, msg.Username)
		context.Respond(&messages.Response{Success: true, Message: "Left subreddit successfully"})
	} else {
		context.Respond(&messages.Response{Success: false, Message: "Subreddit does not exist"})
	}
}

func (e *Engine) handlePost(context actor.Context, msg *messages.Post) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if subreddit, exists := e.subreddits[msg.SubredditName]; exists {
		if subreddit[msg.Username] {
			postID := uuid.New().String()
			msg.ID = postID
			e.posts[postID] = *msg
			context.Respond(&messages.Response{Success: true, Message: fmt.Sprintf("Posted successfully. Post ID: %s", postID)})
		} else {
			context.Respond(&messages.Response{Success: false, Message: "You are not a member of this subreddit"})
		}
	} else {
		context.Respond(&messages.Response{Success: false, Message: "Subreddit does not exist"})
	}
}

func (e *Engine) handleComment(context actor.Context, msg *messages.Comment) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if post, exists := e.posts[msg.PostID]; exists {
		commentID := uuid.New().String()
		msg.ID = commentID
		e.comments[msg.PostID] = append(e.comments[msg.PostID], *msg)

		// Handle voting
		if msg.IsUpvote {
			post.Karma++
		} else {
			post.Karma--
		}
		e.posts[msg.PostID] = post

		voteType := map[bool]string{true: "upvoted", false: "downvoted"}[msg.IsUpvote]
		response := fmt.Sprintf("Commented and %s successfully. Comment ID: %s. New post karma: %d.", voteType, commentID, post.Karma)
		context.Respond(&messages.Response{Success: true, Message: response})
	} else {
		context.Respond(&messages.Response{Success: false, Message: "Post does not exist"})
	}
}

func (e *Engine) handleVote(context actor.Context, msg *messages.Vote) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if post, exists := e.posts[msg.PostID]; exists {
		if subreddit, subredditExists := e.subreddits[post.SubredditName]; subredditExists {
			if subreddit[msg.Username] {
				if msg.IsUpvote {
					post.Karma++
				} else {
					post.Karma--
				}
				context.Respond(&messages.Response{Success: true, Message: fmt.Sprintf("Voted successfully. New karma: %d", post.Karma)})
				e.posts[msg.PostID] = post
				voteType := map[bool]string{true: "upvoted", false: "downvoted"}[msg.IsUpvote]
				response := fmt.Sprintf("User %s %s post %s. New karma: %d", msg.Username, voteType, msg.PostID, post.Karma)
				context.Respond(&messages.Response{Success: true, Message: response})
			} else {
				context.Respond(&messages.Response{Success: false, Message: "You are not a member of this subreddit"})
			}
		} else {
			context.Respond(&messages.Response{Success: false, Message: "Subreddit does not exist"})
		}
	} else {
		context.Respond(&messages.Response{Success: false, Message: "Voted Successfully"})
	}
}

func (e *Engine) handleGetFeed(context actor.Context, msg *messages.GetFeed) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	var feed []messages.Post
	for _, post := range e.posts {
		if e.subreddits[post.SubredditName][msg.Username] {
			feed = append(feed, post)
		}
	}

	sort.Slice(feed, func(i, j int) bool {
		return feed[i].Karma > feed[j].Karma
	})

	context.Respond(&messages.FeedResponse{Posts: feed})
}

func (e *Engine) handleGetDirectMessages(context actor.Context, msg *messages.GetDirectMessages) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	userMessages := e.directMessages[msg.Username]
	context.Respond(&messages.DirectMessageResponse{Messages: userMessages})
}

func (e *Engine) handleReplyDirectMessage(context actor.Context, msg *messages.ReplyDirectMessage) {
	e.mu.Lock()
	defer e.mu.Unlock()

	newMessage := messages.DirectMessage{
		From:      context.Sender().Id,
		To:        msg.To,
		Content:   msg.Content,
		Timestamp: msg.Timestamp,
	}

	e.directMessages[msg.To] = append(e.directMessages[msg.To], newMessage)
	context.Respond(&messages.Response{Success: true, Message: "Direct message sent"})
}
