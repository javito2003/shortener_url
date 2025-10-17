package user

import (
	"context"
	"time"

	"github.com/javito2003/shortener_url/internal/config"
	"github.com/javito2003/shortener_url/internal/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type repository struct {
	collection *mongo.Collection
}

type userModel struct {
	ID        bson.ObjectID `bson:"_id,omitempty"`
	FirstName string        `bson:"first_name"`
	LastName  string        `bson:"last_name"`
	Email     string        `bson:"email"`
	Password  string        `bson:"password"`
	CreatedAt time.Time     `bson:"created_at"`
}

func NewRepository(db *mongo.Client) *repository {
	return &repository{collection: db.Database(config.AppConfig.Mongo.Database).Collection("users")}
}

func (r *repository) Save(ctx context.Context, user *domain.User) (*domain.User, error) {
	createdAt := time.Now()

	result, err := r.collection.InsertOne(ctx, userModel{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: createdAt,
	})

	if err != nil {
		return nil, err
	}

	user.ID = result.InsertedID.(bson.ObjectID).Hex()
	user.CreatedAt = createdAt

	return user, nil
}

func (r *repository) FindByEmail(ctx context.Context, email string) (*domain.User, bool, error) {
	var model userModel
	err := r.collection.FindOne(ctx, map[string]interface{}{"email": email}).Decode(&model)

	if err == mongo.ErrNoDocuments {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	return &domain.User{
		ID:        model.ID.Hex(),
		FirstName: model.FirstName,
		LastName:  model.LastName,
		Email:     model.Email,
		Password:  model.Password,
		CreatedAt: model.CreatedAt,
	}, true, nil
}
