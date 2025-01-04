package simulation

import (
	"fmt"
	"math/rand"
	"protoactor-simulation/client"
	"protoactor-simulation/engine"
	"protoactor-simulation/messages"
	"time"

	"github.com/asynkron/protoactor-go/actor"
)

type Simulation struct {
	system         *actor.ActorSystem
	enginePID      *actor.PID
	numClients     int
	subreddits     []string
	posts          map[string][]Post
	comments       map[string][]Comment
	directMessages map[string][]DirectMessage
}

type Post struct {
	ID        string
	Title     string
	Content   string
	Username  string
	Subreddit string
}

type Comment struct {
	Content  string
	Username string
}

type DirectMessage struct {
	From    string
	Content string
}

func NewSimulation(system *actor.ActorSystem, numClients int) *Simulation {
	return &Simulation{
		system:         system,
		numClients:     numClients,
		subreddits:     []string{"AskReddit", "worldnews", "funny", "gaming", "aww", "todayilearned", "science"},
		posts:          make(map[string][]Post),
		comments:       make(map[string][]Comment),
		directMessages: make(map[string][]DirectMessage),
	}
}

func (s *Simulation) Run() {
	start := time.Now()

	engineProps := actor.PropsFromProducer(func() actor.Actor { return engine.NewEngine() })
	s.enginePID, _ = s.system.Root.SpawnNamed(engineProps, "engine")

	s.system.Root.RequestFuture(s.enginePID, &messages.EngineReady{}, 5*time.Second).Wait()

	clientPIDs := make([]*actor.PID, s.numClients)
	for i := 0; i < s.numClients; i++ {
		username := fmt.Sprintf("user%d", i)
		clientProps := actor.PropsFromProducer(func(i int) func() actor.Actor {
			return func() actor.Actor { return client.NewClient(username, s.enginePID) }
		}(i))
		clientPIDs[i] = s.system.Root.Spawn(clientProps)
	}

	s.simulateActions(clientPIDs)

	elapsed := time.Since(start)
	fmt.Printf("Simulation took %s\n", elapsed-5*time.Second)
}

func (s *Simulation) simulateActions(clientPIDs []*actor.PID) {
	for _, subreddit := range s.subreddits {
		s.system.Root.Send(s.enginePID, &messages.CreateSubreddit{Name: subreddit})
		fmt.Printf("Subreddit 'r/%s' created\n", subreddit)
	}

	for i, clientPID := range clientPIDs {
		subreddit := s.subreddits[rand.Intn(len(s.subreddits))]
		s.system.Root.Send(clientPID, &messages.JoinSubreddit{SubredditName: subreddit})
		fmt.Printf("u/%s joined r/%s\n", clientPID.Id, subreddit)

		// Simulate posting actions
		title, content := generateRealisticPost(subreddit)
		postID := fmt.Sprintf("post%d", rand.Intn(10000))
		post := Post{ID: postID, Title: title, Content: content, Username: clientPID.Id, Subreddit: subreddit}
		s.posts[subreddit] = append(s.posts[subreddit], post)
		fmt.Printf("u/%s posted in r/%s: %s\n", clientPID.Id, subreddit, title)

		// Simulate commenting actions
		commentContent := generateRealisticComment(subreddit)
		comment := Comment{Content: commentContent, Username: clientPID.Id}
		s.comments[postID] = append(s.comments[postID], comment)
		fmt.Printf("u/%s commented on post in r/%s: %s\n", clientPID.Id, subreddit, commentContent[:30]+"...")

		// Simulate direct messages
		if rand.Float32() < 0.1 { // 10% chance to send a direct message
			receiverIndex := rand.Intn(len(clientPIDs))
			if receiverIndex != i { // Ensure sender is not the receiver
				dmContent := fmt.Sprintf("Hello from %s to %s!", clientPID.Id, clientPIDs[receiverIndex].Id)
				dm := DirectMessage{From: clientPID.Id, Content: dmContent}
				s.directMessages[clientPIDs[receiverIndex].Id] = append(s.directMessages[clientPIDs[receiverIndex].Id], dm)
				fmt.Printf("Direct message sent from %s to %s\n", clientPID.Id, clientPIDs[receiverIndex].Id)
			}
		}
	}
}

// GetStatus returns the current status of the simulation
func (s *Simulation) GetStatus() map[string]interface{} {
	return map[string]interface{}{
		"active_users":   s.numClients,
		"subreddits":     len(s.subreddits),
		"total_posts":    len(s.posts),
		"total_comments": len(s.comments),
	}
}

func (s *Simulation) GetSubreddits() []string {
	return s.subreddits
}

func (s *Simulation) GetPosts() []Post {
	var allPosts []Post
	for _, posts := range s.posts {
		allPosts = append(allPosts, posts...)
	}
	return allPosts
}

func (s *Simulation) GetComments(postID string) []Comment {
	return s.comments[postID]
}

func (s *Simulation) GetFeed(username string) []Post {
	var userFeed []Post
	for _, posts := range s.posts {
		userFeed = append(userFeed, posts...)
	}
	return userFeed[:5] // Return top 5 posts as feed for simplicity
}

func (s *Simulation) GetDirectMessages(username string) []DirectMessage {
	return s.directMessages[username]
}

// Helper functions to generate realistic content

func generateRealisticPost(subreddit string) (string, string) {
	titlesMap := map[string][]string{
		"AskReddit":     {"What's the craziest thing you've ever done?", "If you could have dinner with any historical figure, who would it be?"},
		"worldnews":     {"Breaking: Major diplomatic breakthrough in Middle East", "New study shows alarming rate of climate change"},
		"funny":         {"My dog's reaction when I pretend to throw the ball", "Found this gem while cleaning out my grandpa's attic"},
		"gaming":        {"After 500 hours, I finally beat this boss", "New leak suggests GTA 6 release date"},
		"aww":           {"My rescue kitten's first day home", "This baby elephant learning to use its trunk"},
		"todayilearned": {"TIL the Great Wall of China is not visible from space", "TIL honey never spoils"},
		"science":       {"New breakthrough in quantum computing", "Scientists discover potential cure for common cold"},
	}

	contentMap := map[string][]string{
		"AskReddit":     {"I'm really curious to hear everyone's stories!", "Imagine the conversations you could have..."},
		"worldnews":     {"This could have significant implications for global politics.", "The study calls for immediate action to mitigate the effects."},
		"funny":         {"His face of betrayal is priceless!", "I can't believe this was just sitting there for years!"},
		"gaming":        {"The feeling of accomplishment is indescribable.", "If this is true, it's going to be a game-changer for the industry."},
		"aww":           {"She's already claimed my heart and the best spot on the couch.", "Nature is truly amazing. Look at that playfulness!"},
		"todayilearned": {"It's actually a common misconception. Here's why...", "Archaeologists have found pots of honey in ancient Egyptian tombs that are still perfectly edible."},
		"science":       {"This could revolutionize computing as we know it.", "The potential applications in medicine are enormous."},
	}

	titles, ok := titlesMap[subreddit]
	if !ok || len(titles) == 0 {
		return "Default Title", "Default Content"
	}

	content, ok := contentMap[subreddit]
	if !ok || len(content) == 0 {
		return titles[rand.Intn(len(titles))], "Default Content"
	}

	titleIdx := rand.Intn(len(titles))
	contentIdx := rand.Intn(len(content))

	return titles[titleIdx], content[contentIdx]
}
func generateRealisticComment(subreddit string) string {
	commentOptionsMap := map[string][]string{
		"AskReddit":     {"Wow, that's insane! I can't believe you actually did that.", "I'd choose Einstein. Imagine the mind-bending conversations!"},
		"worldnews":     {"This is huge if true. Hope it leads to lasting peace.", "We need to take this seriously and act now before it's too late."},
		"funny":         {"I can't stop laughing! The look on his face is priceless.", "Your grandpa must have been quite the character!"},
		"gaming":        {"Congrats! That boss gave me nightmares for weeks.", "Please let this be true. I've been waiting for so long!"},
		"aww":           {"She's adorable! You're so lucky to have found each other.", "Elephants are such intelligent and gentle creatures. This made my day!"},
		"todayilearned": {"Mind blown! I've been telling people this for years.", "That's fascinating! Nature never ceases to amaze me."},
		"science":       {"The implications of this are staggering. Can't wait to see where this leads.", "If this pans out, it could save millions of lives."},
	}

	commentOptions, ok := commentOptionsMap[subreddit]
	if !ok || len(commentOptions) == 0 {
		return "Default Comment"
	}

	commentIdx := rand.Intn(len(commentOptions))
	return commentOptions[commentIdx]
}
