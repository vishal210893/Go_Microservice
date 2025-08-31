// Package db provides database seeding functionality for populating the database
// with test data including users, posts, and comments for development and testing purposes.
package db

import (
	"Go-Microservice/internal/repo"
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
)

// seedData contains all the sample data used for seeding the database
type seedData struct {
	usernames []string
	titles    []string
	contents  []string
	tags      []string
	comments  []string
}

// newSeedData returns a new instance of seedData with predefined sample data
func newSeedData() *seedData {
	return &seedData{
		usernames: []string{
			"alice", "bob", "charlie", "dave", "eve", "frank", "grace", "heidi",
			"ivan", "judy", "karl", "laura", "mallory", "nina", "oscar", "peggy",
			"quinn", "rachel", "steve", "trent", "ursula", "victor", "wendy", "xander",
			"yvonne", "zack", "amber", "brian", "carol", "doug", "eric", "fiona",
			"george", "hannah", "ian", "jessica", "kevin", "lisa", "mike", "natalie",
			"oliver", "peter", "queen", "ron", "susan", "tim", "uma", "vicky",
			"walter", "xenia", "yasmin", "zoe",
		},
		titles: []string{
			"The Power of Habit", "Embracing Minimalism", "Healthy Eating Tips",
			"Travel on a Budget", "Mindfulness Meditation", "Boost Your Productivity",
			"Home Office Setup", "Digital Detox", "Gardening Basics",
			"DIY Home Projects", "Yoga for Beginners", "Sustainable Living",
			"Mastering Time Management", "Exploring Nature", "Simple Cooking Recipes",
			"Fitness at Home", "Personal Finance Tips", "Creative Writing",
			"Mental Health Awareness", "Learning New Skills",
		},
		contents: []string{
			"In this post, we'll explore how to develop good habits that stick and transform your life.",
			"Discover the benefits of a minimalist lifestyle and how to declutter your home and mind.",
			"Learn practical tips for eating healthy on a budget without sacrificing flavor.",
			"Traveling doesn't have to be expensive. Here are some tips for seeing the world on a budget.",
			"Mindfulness meditation can reduce stress and improve your mental well-being. Here's how to get started.",
			"Increase your productivity with these simple and effective strategies.",
			"Set up the perfect home office to boost your work-from-home efficiency and comfort.",
			"A digital detox can help you reconnect with the real world and improve your mental health.",
			"Start your gardening journey with these basic tips for beginners.",
			"Transform your home with these fun and easy DIY projects.",
			"Yoga is a great way to stay fit and flexible. Here are some beginner-friendly poses to try.",
			"Sustainable living is good for you and the planet. Learn how to make eco-friendly choices.",
			"Master time management with these tips and get more done in less time.",
			"Nature has so much to offer. Discover the benefits of spending time outdoors.",
			"Whip up delicious meals with these simple and quick cooking recipes.",
			"Stay fit without leaving home with these effective at-home workout routines.",
			"Take control of your finances with these practical personal finance tips.",
			"Unleash your creativity with these inspiring writing prompts and exercises.",
			"Mental health is just as important as physical health. Learn how to take care of your mind.",
			"Learning new skills can be fun and rewarding. Here are some ideas to get you started.",
		},
		tags: []string{
			"Self Improvement", "Minimalism", "Health", "Travel", "Mindfulness",
			"Productivity", "Home Office", "Digital Detox", "Gardening", "DIY",
			"Yoga", "Sustainability", "Time Management", "Nature", "Cooking",
			"Fitness", "Personal Finance", "Writing", "Mental Health", "Learning",
		},
		comments: []string{
			"Great post! Thanks for sharing.",
			"I completely agree with your thoughts.",
			"Thanks for the tips, very helpful.",
			"Interesting perspective, I hadn't considered that.",
			"Thanks for sharing your experience.",
			"Well written, I enjoyed reading this.",
			"This is very insightful, thanks for posting.",
			"Great advice, I'll definitely try that.",
			"I love this, very inspirational.",
			"Thanks for the information, very useful.",
		},
	}
}

// Seed populates the database with test data including users, posts, and comments.
// It creates the specified number of users, posts, and comments using sample data.
//
// The seeding process:
//  1. Creates 100 users with unique usernames and emails
//  2. Creates 200 posts assigned to random users
//  3. Creates 500 comments on random posts by random users
//
// Parameters:
//   - repository: Repository interface for database operations
//
// If any error occurs during seeding, the function logs the error and returns early.
func Seed(repository repo.Repository, db *sql.DB) {
	ctx := context.Background()
	data := newSeedData()

	// Generate and create users
	users := data.generateUsers(100)

	tx, _ := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	for _, user := range users {
		if err := repository.Users.Create(ctx, tx, user); err != nil {
			_ = tx.Rollback()
			log.Printf("Error creating user %s: %v", user.Username, err)
			return
		}
	}
	log.Printf("Successfully created %d users", len(users))
	tx.Commit()

	// Generate and create posts
	posts := data.generatePosts(200, users)
	for _, post := range posts {
		if err := repository.Posts.Create(ctx, post); err != nil {
			log.Printf("Error creating post '%s': %v", post.Title, err)
			return
		}
	}
	log.Printf("Successfully created %d posts", len(posts))

	// Generate and create comments
	generatedComments := data.generateComments(500, users, posts)
	for _, comment := range generatedComments {
		if err := repository.Comments.Create(ctx, comment); err != nil {
			log.Printf("Error creating comment: %v", err)
			return
		}
	}
	log.Printf("Successfully created %d comments", len(generatedComments))

	log.Println("Database seeding completed successfully")
}

// generateUsers creates the specified number of user instances with unique usernames and emails.
//
// Parameters:
//   - num: Number of users to generate
//
// Returns:
//   - []*repo.User: Slice of generated user pointers
func (s *seedData) generateUsers(num int) []*repo.User {
	users := make([]*repo.User, num)

	for i := 0; i < num; i++ {
		baseUsername := s.usernames[i%len(s.usernames)]
		users[i] = &repo.User{
			Username: fmt.Sprintf("%s%d", baseUsername, i),
			Email:    fmt.Sprintf("%s%d@example.com", baseUsername, i),
			Role: repo.Role{
				Name: "user",
			},
		}
	}

	return users
}

// generatePosts creates the specified number of post instances assigned to random users.
//
// Parameters:
//   - num: Number of posts to generate
//   - users: Slice of users to assign posts to
//
// Returns:
//   - []*repo.Post: Slice of generated post pointers
func (s *seedData) generatePosts(num int, users []*repo.User) []*repo.Post {
	posts := make([]*repo.Post, num)

	for i := 0; i < num; i++ {
		user := users[rand.Intn(len(users))]
		posts[i] = &repo.Post{
			UserID:  user.ID,
			Title:   s.titles[rand.Intn(len(s.titles))],
			Content: s.contents[rand.Intn(len(s.contents))], // Fixed: was using titles instead of contents
			Tags: []string{
				s.tags[rand.Intn(len(s.tags))],
				s.tags[rand.Intn(len(s.tags))],
			},
		}
	}

	return posts
}

// generateComments creates the specified number of comment instances on random posts by random users.
//
// Parameters:
//   - num: Number of comments to generate
//   - users: Slice of users to assign comments to
//   - posts: Slice of posts to comment on
//
// Returns:
//   - []*repo.Comment: Slice of generated comment pointers
func (s *seedData) generateComments(num int, users []*repo.User, posts []*repo.Post) []*repo.Comment {
	generatedComments := make([]*repo.Comment, num)

	for i := 0; i < num; i++ {
		generatedComments[i] = &repo.Comment{
			PostID:  posts[rand.Intn(len(posts))].ID,
			UserID:  users[rand.Intn(len(users))].ID,
			Content: s.comments[rand.Intn(len(s.comments))],
		}
	}

	return generatedComments
}
