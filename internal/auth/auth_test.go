package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func TestPasswordHashing(t *testing.T) {
	password := "testpassword123"

	// Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	// Verify hash is not empty and different from plain
	if len(hashed) == 0 {
		t.Error("Hashed password is empty")
	}
	if string(hashed) == password {
		t.Error("Hashed password is the same as plain password")
	}

	// Compare hash with correct password
	err = bcrypt.CompareHashAndPassword(hashed, []byte(password))
	if err != nil {
		t.Errorf("Password comparison failed for correct password: %v", err)
	}

	// Compare hash with incorrect password
	err = bcrypt.CompareHashAndPassword(hashed, []byte("wrongpassword"))
	if err == nil {
		t.Error("Password comparison should fail for incorrect password")
	}
}

func TestJWTGenerationAndParsing(t *testing.T) {
	testSecret := []byte("test-jwt-secret-for-testing")
	userID := uint(123)

	// Generate token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString(testSecret)
	if err != nil {
		t.Fatalf("Failed to generate JWT: %v", err)
	}

	// Parse token
	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return testSecret, nil
	})
	if err != nil {
		t.Fatalf("Failed to parse JWT: %v", err)
	}
	if !parsedToken.Valid {
		t.Fatal("Token is not valid")
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok {
		if userIDFloat, ok := claims["user_id"].(float64); ok {
			parsedUserID := uint(userIDFloat)
			if parsedUserID != userID {
				t.Errorf("Parsed user ID does not match: got %d, want %d", parsedUserID, userID)
			}
		} else {
			t.Error("user_id claim not found or not float64")
		}
	} else {
		t.Error("Claims are not MapClaims")
	}
}
