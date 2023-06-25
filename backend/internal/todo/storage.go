package todo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type todoDB struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description" json:"description"`
	Completed   bool               `bson:"completed" json:"completed"`
}

type TodoStorage struct {
	db *mongo.Database
}

func NewTodoStorage(db *mongo.Database) *TodoStorage {
	return &TodoStorage{
		db: db,
	}
}

func (s *TodoStorage) getAllTodos(ctx context.Context) ([]todoDB, error) {
	collection := s.db.Collection("todos")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	todos := make([]todoDB, 0)
	if err = cursor.All(ctx, &todos); err != nil {
		return nil, err
	}

	return todos, nil
}

func (s *TodoStorage) createTodo(title string, description string, completed bool, ctx context.Context) (string, error) {
	collection := s.db.Collection("todos")

	result, err := collection.InsertOne(ctx, bson.M{
		"title": title, "description": description, "completed": completed,
	})
	if err != nil {
		return "", err
	}

	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (s *TodoStorage) completeTodo(id string, ctx context.Context) error {
	collection := s.db.Collection("todos")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": bson.M{"completed": true}}

	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (s *TodoStorage) editTodo(id, title, description string, completed bool, ctx context.Context) error {
	collection := s.db.Collection("todos")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": bson.M{
		"title":       title,
		"description": description,
		"completed":   completed,
	}}

	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (s *TodoStorage) delete(id string, ctx context.Context) error {
	collection := s.db.Collection("todos")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}

	return nil
}
