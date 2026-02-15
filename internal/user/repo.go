package user

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repo struct {
	coll *mongo.Collection
}

func NewRepo(db *mongo.Database) *Repo {
	return &Repo{coll: db.Collection("users")}
}

func (r *Repo) FindByEmail(ctx context.Context, email string) (User, error) {
	email = strings.ToLower(strings.TrimSpace(email))

	filter := bson.M{"email": email}

	var u User
	err := r.coll.FindOne(ctx, filter).Decode(&u)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return User{}, mongo.ErrNoDocuments
		}
		return User{}, fmt.Errorf("Find by email failed: %w", err)
	}
	return u, nil
}

func (r *Repo) Create(ctx context.Context, u User) (User, error) {

	res, err := r.coll.InsertOne(ctx, u)
	if err != nil {
		return User{}, fmt.Errorf("Insert User Failed: %w", err)
	}

	id, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return User{}, fmt.Errorf("Insert User Failed and Inserted id is not objectId: %w", err)
	}

	u.ID = id

	return u, nil
}
