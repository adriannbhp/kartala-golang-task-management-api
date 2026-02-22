package main

import (
	"context"
	"github.com/adriannbhp/kartala-golang-task-management-api/internal/config"
	"github.com/adriannbhp/kartala-golang-task-management-api/internal/tasks"
	"github.com/adriannbhp/kartala-golang-task-management-api/internal/users"
	"github.com/adriannbhp/kartala-golang-task-management-api/pkg/database"
	"github.com/adriannbhp/kartala-golang-task-management-api/pkg/logger"
	"github.com/adriannbhp/kartala-golang-task-management-api/pkg/security"
	"time"
	"os"
	"github.com/google/uuid"
)

var seedUsers = []struct {
	Name     string
	Username string
	Email    string
	Password string
}{
	{
		Name:     "Admin User",
		Username: "admin",
		Email:    "admin@example.com",
		Password: "admin123",
	},
	{
		Name:     "Test User",
		Username: "testuser",
		Email:    "user@example.com",
		Password: "user123",
	},
}

var seedTasks = []struct {
	Title       string
	Description string
	UserEmail   string
}{
	{
		Title:       "Setup Project",
		Description: "Configure the initial project structure and dependencies",
		UserEmail:   "admin@example.com",
	},
	{
		Title:       "Implement Auth",
		Description: "Create JWT authentication and authorization",
		UserEmail:   "admin@example.com",
	},
	{
		Title:       "Create Crud",
		Description: "Build CRUD operations for users and tasks",
		UserEmail:   "user@example.com",
	},
}

var exitFunc = os.Exit

func main() {
	if err := BootstrapSeeder(); err != nil {
		logger.Logger.Errorf("Seeding failed: %v", err)
		exitFunc(1)
	}
	logger.Logger.Info("Seeding completed!")
}

func BootstrapSeeder() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}
	logger.SetupLogger()
	
	logger.Logger.Info("Starting database seeder...")
	
	// Initialize Database
	db, err := database.NewDatabaseConnection(cfg.Database)
	if err != nil {
		return err
	}
	
	userRepo := users.NewRepository(db)
	taskRepo := tasks.NewRepository(db)

	return RunSeed(context.Background(), userRepo, taskRepo)
}

func RunSeed(ctx context.Context, userRepo users.Repository, taskRepo tasks.Repository) error {
	// Seed Users
	userMap := make(map[string]uuid.UUID)
	for _, seedUser := range seedUsers {
		existingUser, err := userRepo.FindByEmail(ctx, seedUser.Email)
		if err == nil && existingUser != nil {
			logger.Logger.Warnf("User %s already exists, skipping...", seedUser.Email)
			userMap[seedUser.Email] = existingUser.ID
			continue
		}
		
		hashedPassword, _ := security.HashPassword(seedUser.Password)
		user := &users.User{
			ID:        uuid.New(),
			Username:  seedUser.Username,
			Email:     seedUser.Email,
			Password:  hashedPassword,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		
		userID, err := userRepo.InsertNewUser(ctx, user)
		if err != nil {
			logger.Logger.Errorf("Failed to create user %s: %v", seedUser.Email, err)
			continue
		}
		userMap[seedUser.Email] = user.ID
		logger.Logger.Infof("Created user: %s - ID: %s", seedUser.Name, userID)
	}

	// Seed Tasks
	for _, seedTask := range seedTasks {
		userID, ok := userMap[seedTask.UserEmail]
		if !ok {
			logger.Logger.Errorf("Failed to find user ID for email %s, skipping task: %s", seedTask.UserEmail, seedTask.Title)
			continue
		}

		task := &tasks.Task{
			ID:          uuid.New(),
			Title:       seedTask.Title,
			Description: seedTask.Description,
			Status:      tasks.TaskStatusTodo,
			UserID:      userID,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		if err := taskRepo.Insert(ctx, task); err != nil {
			logger.Logger.Errorf("Failed to create task %s: %v", seedTask.Title, err)
			continue
		}
		logger.Logger.Infof("Created task: %s for user: %s", seedTask.Title, seedTask.UserEmail)
	}

	return nil
}
