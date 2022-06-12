package main

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func loginHandler(dbClient *mongo.Client) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		userUuid := c.PostForm("uuid")
		_, err := uuid.Parse(userUuid)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid uuid",
			})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		col := dbClient.Database("MEDODS-interview").Collection("users")
		col.FindOneAndDelete(ctx, bson.M{"uuid": userUuid})

		user := User{Uuid: userUuid}

		accessToken, err := generateAccessToken(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "error generating access token",
			})
			return
		}

		refreshToken, err := generateRefreshToken(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "error generating refresh token",
			})
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "error generating hash",
			})
			return
		}

		doc := bson.D{
			{Key: "uuid", Value: userUuid},
			{Key: "token", Value: string(hash)},
		}

		_, err = col.InsertOne(ctx, doc)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "error inserting in db",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"access-token":  accessToken,
			"refresh-token": refreshToken,
		})
	}
	return fn
}
