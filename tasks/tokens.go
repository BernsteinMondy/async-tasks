package tasks

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"sync"
)

const (
	tokenContextKey = "token"
	tokenSecretKey  = "token_secret"
)

func main() {
	ctx := context.Background()
	ctx, err := AddJWTToContext(ctx, 123)
	if err != nil {
		fmt.Println(err)
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		userId, err := ExtractUserIDFromContext(ctx)
		if err != nil {
			fmt.Println(err)
			return
		}

		if userId != 123 {
			fmt.Println("userId != 123")
			return
		}

		fmt.Println("userId == 123")
	}()

	wg.Wait()
}

func AddJWTToContext(ctx context.Context, userID int) (context.Context, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
	},
	)

	tokenString, err := token.SignedString([]byte(tokenSecretKey))
	if err != nil {
		return nil, fmt.Errorf("failed to sign token: %w", err)
	}

	ctx = context.WithValue(ctx, tokenContextKey, tokenString)
	return ctx, nil
}

func ExtractUserIDFromContext(ctx context.Context) (int, error) {
	tokenValue := ctx.Value(tokenContextKey)
	if tokenValue == nil {
		return 0, fmt.Errorf("token not found in context")
	}

	tokenString, ok := tokenValue.(string)
	if !ok {
		return 0, fmt.Errorf("token is not a string")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(tokenSecretKey), nil
	})

	if err != nil {
		return 0, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return 0, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("invalid token claims")
	}

	userIDClaim, exists := claims["user_id"]
	if !exists {
		return 0, fmt.Errorf("user_id claim not found")
	}

	userID, ok := userIDClaim.(float64)
	if !ok {
		return 0, fmt.Errorf("user_id is not a number")
	}

	userIdInt := int(userID)

	return userIdInt, nil
}
