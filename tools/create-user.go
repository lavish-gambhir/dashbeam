package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/lavish-gambhir/dashbeam/shared/config"
	"github.com/lavish-gambhir/dashbeam/shared/database/postgres"
	"github.com/lavish-gambhir/dashbeam/shared/database/repositories"
	"github.com/lavish-gambhir/dashbeam/shared/models"
)

func main() {
	var (
		email         = flag.String("email", "", "Email for the user")
		name          = flag.String("name", "", "Name of the user")
		role          = flag.String("role", "student", "Role for the user (student, teacher)")
		schoolName    = flag.String("school-name", "Test School", "Name of the school to create")
		classroomName = flag.String("classroom-name", "Test Classroom", "Name of the classroom to create")
		appType       = flag.String("app", "whiteboard", "App type (whiteboard, notebook)")
	)
	flag.Parse()

	if *email == "" || *name == "" {
		fmt.Println("Usage: create-user -email=<email> -name=<name> [-role=<role>] [-school-name=<name>] [-classroom-name=<name>] [-app=<app-type>] [-jwt]")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Validate role
	if *role != "student" && *role != "teacher" {
		fmt.Printf("Invalid role: %s. Valid roles are: student, teacher\n", *role)
		os.Exit(1)
	}

	ctx := context.Background()

	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: failed to load .env file: %v", err)
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	pool, err := postgres.Connect(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	pgdb := postgres.New(pool, nil)

	schoolID, err := createSchool(ctx, pgdb, *schoolName)
	if err != nil {
		log.Fatalf("Failed to create school: %v", err)
	}

	classroomID, err := createClassroom(ctx, pgdb, *classroomName, schoolID)
	if err != nil {
		log.Fatalf("Failed to create classroom: %v", err)
	}

	// Create user
	userID, err := createUser(ctx, pgdb, *email, *name, *role, schoolID, classroomID)
	if err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}

	fmt.Printf("=== CREATED SUCCESSFULLY ===\n")
	fmt.Printf("School ID: %s (Name: %s)\n", schoolID.String(), *schoolName)
	fmt.Printf("Classroom ID: %s (Name: %s)\n", classroomID.String(), *classroomName)
	fmt.Printf("User ID: %s\n", userID.String())
	fmt.Printf("Email: %s\n", *email)
	fmt.Printf("Name: %s\n", *name)
	fmt.Printf("Role: %s\n", *role)
	fmt.Printf("App Type: %s\n", *appType)

	token, err := generateJWT(userID, *email, *name, *role, schoolID, &classroomID, *appType, cfg.Auth)
	if err != nil {
		log.Fatalf("Failed to generate JWT: %v", err)
	}

	fmt.Printf("\n=== JWT TOKEN ===\n")
	fmt.Printf("%s\n", token)
	fmt.Printf("\n=== CURL EXAMPLES ===\n")
	fmt.Printf("# Test batch events:\n")
	fmt.Printf("curl -X POST -H \"Authorization: Bearer %s\" \\\n", token)
	fmt.Printf("  -H \"Content-Type: application/json\" \\\n")
	fmt.Printf("  -d '{\"events\":[{\"type\":\"user_login\",\"user_id\":\"%s\",\"school_id\":\"%s\",\"timestamp\":\"%s\",\"payload\":{}}]}' \\\n",
		userID.String(), schoolID.String(), time.Now().UTC().Format(time.RFC3339))
	fmt.Printf("  http://localhost:8080/events/batch\n")
}

func createSchool(ctx context.Context, db *postgres.DB, name string) (uuid.UUID, error) {
	schoolID := uuid.New()
	now := time.Now().UTC()

	query := `
		INSERT INTO schools (
			id, name, district, address, timezone,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7
		)`

	_, err := db.Conn(ctx).Exec(ctx, query,
		schoolID,
		name,
		"123 Test Street, Test City", // district
		"123 Test Street, Test City", // default address
		"UTC",                        // tz

		now,
		now,
	)

	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create school: %w", err)
	}

	fmt.Printf("School created: %s (%s)\n", schoolID.String(), name)
	return schoolID, nil
}

func createClassroom(ctx context.Context, db *postgres.DB, name string, schoolID uuid.UUID) (uuid.UUID, error) {
	classroomID := uuid.New()
	now := time.Now().UTC()

	query := `
		INSERT INTO classrooms (
			id, school_id, name,
			created_at
		) VALUES (
			$1, $2, $3, $4
		)`

	_, err := db.Conn(ctx).Exec(ctx, query,
		classroomID,
		schoolID,
		name,
		now,
	)

	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create classroom: %w", err)
	}

	fmt.Printf("Classroom created: %s (%s)\n", classroomID.String(), name)
	return classroomID, nil
}

func createUser(ctx context.Context, db *postgres.DB, email, name, role string, schoolID, classroomID uuid.UUID) (uuid.UUID, error) {
	userRepo := repositories.NewUserRepository(db)
	userID := uuid.New()
	now := time.Now().UTC()

	user := &models.User{
		ID:                     userID.String(),
		Email:                  email,
		Name:                   name,
		Role:                   role,
		SchoolID:               schoolID.String(),
		FirstSeenAt:            now,
		LastSeenAt:             now,
		TotalQuizSessions:      0,
		TotalQuestionsAnswered: 0,
		AverageResponseTimeMS:  0,
		CreatedAt:              now,
		UpdatedAt:              now,
	}

	err := userRepo.CreateUser(ctx, user)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create user: %w", err)
	}

	err = createClassroomMembership(ctx, db, userID, classroomID, role)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create classroom membership: %w", err)
	}

	fmt.Printf("User created: %s (%s)\n", userID.String(), email)
	return userID, nil
}

func createClassroomMembership(ctx context.Context, db *postgres.DB, userID, classroomID uuid.UUID, role string) error {
	now := time.Now().UTC()

	query := `
		INSERT INTO user_classroom_memberships (
			id, user_id, classroom_id, role, joined_at
		) VALUES (
			$1, $2, $3, $4, $5
		)`

	_, err := db.Conn(ctx).Exec(ctx, query,
		uuid.New(), // membership ID
		userID,
		classroomID,
		role,
		now, // joined_at

	)

	if err != nil {
		return fmt.Errorf("failed to create classroom membership: %w", err)
	}

	fmt.Printf("Classroom membership created for user %s in classroom %s\n", userID.String(), classroomID.String())
	return nil
}

func generateJWT(userID uuid.UUID, email, name, role string, schoolID uuid.UUID, classroomID *uuid.UUID, appType string, authConfig config.AuthConfig) (string, error) {
	now := time.Now().UTC()
	expiry := now.Add(authConfig.AccessTokenExpiry)

	claims := &models.JWTClaims{
		UserID:      userID,
		Email:       email,
		Name:        name,
		Role:        role,
		SchoolID:    schoolID,
		ClassroomID: classroomID,
		AppType:     appType,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiry),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "dashbeam-test",
			Subject:   userID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(authConfig.JWTSecretKey))
}
