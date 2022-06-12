package main

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func tokenHandler(dbClient *mongo.Client) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		access := c.PostForm("access-token")
		accessTkn, err := authorizeAccessToken(access)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "access token not authorized",
			})
			return
		}

		refresh := c.PostForm("refresh-token")
		refreshTkn, err := authorizeRefreshToken(refresh)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "refresh token not authorized",
			})
			return
		}

		if refreshTkn.Uuid != accessTkn.Uuid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "refresh and access token uuids don't match",
			})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var res bson.M
		col := dbClient.Database("MEDODS-interview").Collection("users")
		err = col.FindOne(ctx, bson.M{"uuid": accessTkn.Uuid}).Decode(&res)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "couldn't find document with given uuid",
			})
			return
		}
		err = bcrypt.CompareHashAndPassword([]byte(res["token"].(string)), []byte(refresh))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "token doesn't match to one in the database",
			})
			return
		}

		user := User{Uuid: accessTkn.Uuid}
		newAccessToken, err := generateAccessToken(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "error generating access token",
			})
			return
		}

		newRefreshToken, err := generateRefreshToken(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "error generating refresh token",
			})
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(newRefreshToken), bcrypt.DefaultCost)
		doc := bson.D{
			{Key: "uuid", Value: accessTkn.Uuid},
			{Key: "token", Value: string(hash)},
		}

		col.FindOneAndDelete(ctx, bson.M{"uuid": accessTkn.Uuid})

		_, err = col.InsertOne(ctx, doc)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "error inserting refresh token in db",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"access-token":  newAccessToken,
			"refresh-token": newRefreshToken,
		})
	}
	return fn
}