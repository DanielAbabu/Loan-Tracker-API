package middleware

import (
	"loan-tracker/domain"
	"loan-tracker/infrastructure"
	"loan-tracker/repositories"
	"log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// AuthMiddleware returns a middleware function that checks for the presence of a valid access token and handles token refreshing if needed.
func AuthMiddleware(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		authParts := strings.Split(authHeader, " ")
		if len(authParts) != 2 || strings.ToLower(authParts[0]) != "bearer" {
			c.JSON(401, gin.H{"error": "Invalid authorization header"})
			c.Abort()
			return
		}

		tokenString := authParts[1]
		log.Println("Token string: ", tokenString)
		token, err := infrastructure.TokenClaimer(tokenString)
		if err != nil {
			log.Println("Token parsing error:", err.Error())
			c.JSON(401, gin.H{"error": "Invalid JWT"})
			c.Abort()
			return
		}

		log.Println("Token: ", token)

		claims, ok := token.Claims.(*domain.JWTClaim)
		if !ok || !token.Valid {
			log.Println("Token parsing error:", claims)
			c.JSON(401, gin.H{"error": "Invalid JWT"})
			c.Abort()
			return
		}

		userRepo := repositories.NewUserRepository(client)
		var tem domain.User
		tem.ID, err = primitive.ObjectIDFromHex(claims.UserID)

		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid user ID in JWT"})
			c.Abort()
			return
		}

		user, err := userRepo.FindByID(tem)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(401, gin.H{"error": "User not found"})
			} else {
				c.JSON(500, gin.H{"error": "Database error"})
			}
			c.Abort()
			return
		}

		if claims.Exp < time.Now().Unix() {
			refreshToken := user.RefreshToken
			if refreshToken == "" {
				c.JSON(401, gin.H{"error": "Access token expired and no refresh token provided"})
				c.Abort()
				return
			}

			newAccessToken, err := infrastructure.RefreshAccessToken(refreshToken)
			if err != nil {
				c.JSON(401, gin.H{"error": "Refresh token invalid or expired"})
				c.Abort()
				return
			}

			c.Header("New-Access-Token", newAccessToken)
		}

		// Optionally, verify additional claims with the database values
		if user.Email != claims.Email {
			c.JSON(401, gin.H{"error": "Invalid JWT claims"})
			c.Abort()
			return
		}

		c.Set("isadmin", user.IsAdmin)
		c.Set("userid", user.ID.Hex())
		log.Println(c.GetString("userid"), claims.UserID)

		c.Next()
	}
}

// The middleware for Authentication
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		isAdmin, exists := c.Get("isadmin")
		if !exists || !isAdmin.(bool) {
			c.JSON(403, gin.H{"error": "Forbidden: You don't have admin privileges"})
			c.Abort()
			return
		}
		c.Next()
	}
}
