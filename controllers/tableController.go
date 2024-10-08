package controller

import (
	"context"
	"log"
	"net/http"
	"restaurant-management/database"
	"restaurant-management/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/x/mongo/driver/mongocrypt/options"
)

var tableCollection *mongo.Collection = database.OpenCollection(database.Client, "table")

func GetTables() gin.HandlerFunc {

	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		result, err := tableCollection.Find(context.TODO(), bson.M{})
		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while feteching tables"})
		}
		var allTables []bson.M
		if err = result.All(ctx, &allTables); err != nil {
			log.Fatal(err)
		}

		c.JSON(http.StatusOK, allTables)

	}
}

func GetTable() gin.HandlerFunc {

	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		tableId := c.Param("table_id")

		var table models.Table
		err := orderCollection.FindOne(ctx, bson.M{"table_id": tableId}).Decode(&table)
		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while fetching the table"})
		}
		c.JSON(http.StatusOK, table)

	}
}

func CreateTable() gin.HandlerFunc {

	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var table models.Table
		if err := c.BindJSON(&table); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		validationErr := validate.Struct(table)

		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		table.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		table.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		table.ID = primitive.NewObjectID()
		table.Table_id = table.ID.Hex()

		result, resultErr := tableCollection.InsertOne(ctx, table)

		if resultErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Create Table Failed"})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, result)
	}
}

func UpdateTable() gin.HandlerFunc {

	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var table models.Table
		if err := c.BindJSON(&table); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		}
		tableId := c.Param("table_id")

		var updateObj primitive.D

		if table.Number_of_guests != nil {
			updateObj = append(updateObj, bson.E{Key: "number_of_guest", Value: table.Number_of_guests})
		}
		if table.Table_number != nil {
			updateObj = append(updateObj, bson.E{Key: "table_number", Value: table.Table_number})
		}
		table.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{Key: "updated_at", Value: table.Updated_at})
		upsert := true
		filter := bson.M{"table_id": tableId}
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}
		result, err := tableCollection.UpdateOne(ctx, filter, bson.D{
			{Key: "$set", Value: updateObj},
		}, &opt)

		if err != nil {
			msg := "Table update failed"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, result)

	}
}
