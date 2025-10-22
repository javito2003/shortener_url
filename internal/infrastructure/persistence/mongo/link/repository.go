package link

import (
	"context"

	"github.com/javito2003/shortener_url/internal/config"
	link "github.com/javito2003/shortener_url/internal/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type repository struct {
	collection *mongo.Collection
}

type linkModel struct {
	ID         bson.ObjectID `bson:"_id,omitempty"`
	ShortCode  string        `bson:"short_code"`
	Url        string        `bson:"url"`
	ClickCount int           `bson:"click_count"`
	UserID     bson.ObjectID `bson:"user_id"`
	ExpiresAt  bson.DateTime `bson:"expires_at"`
}

func NewRepository(db *mongo.Client) *repository {
	return &repository{collection: db.Database(config.AppConfig.Mongo.Database).Collection("links")}
}

func (r *repository) Save(ctx context.Context, link *link.Link) (*link.Link, error) {
	userId, err := bson.ObjectIDFromHex(link.UserID)
	if err != nil {
		return nil, err
	}

	result, err := r.collection.InsertOne(ctx, linkModel{
		ShortCode:  link.ShortCode,
		Url:        link.URL,
		ClickCount: link.ClickCount,
		UserID:     userId,
		ExpiresAt:  bson.DateTime(link.ExpiresAt.UnixMilli()),
	})

	if err != nil {
		return nil, err
	}

	link.ID = result.InsertedID.(bson.ObjectID).Hex()

	return link, nil
}

func (r *repository) FindByShortCode(ctx context.Context, shortCode string) (*link.Link, bool, error) {
	var model linkModel
	err := r.collection.FindOne(ctx, map[string]interface{}{"short_code": shortCode}).Decode(&model)

	if err == mongo.ErrNoDocuments {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	expiresAtTime := model.ExpiresAt.Time()
	return &link.Link{
		ID:         model.ID.Hex(),
		ShortCode:  model.ShortCode,
		URL:        model.Url,
		ClickCount: model.ClickCount,
		ExpiresAt:  &expiresAtTime,
	}, true, nil
}

func (r *repository) GetByUrl(ctx context.Context, url string) (*link.Link, bool, error) {
	var model linkModel
	err := r.collection.FindOne(ctx, map[string]interface{}{"url": url}).Decode(&model)
	if err == mongo.ErrNoDocuments {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	expiresAtTime := model.ExpiresAt.Time()

	return &link.Link{
		ID:         model.ID.Hex(),
		ShortCode:  model.ShortCode,
		URL:        model.Url,
		ClickCount: model.ClickCount,
		ExpiresAt:  &expiresAtTime,
	}, true, nil
}

func (r *repository) GetByUser(ctx context.Context, userID string, limit, skip int32) ([]*link.Link, error) {
	userIdObjectID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	filter := map[string]interface{}{"user_id": userIdObjectID}
	limitInt64 := int64(limit)
	skipInt64 := int64(skip)
	options := options.Find().SetLimit(limitInt64).SetSkip(skipInt64)

	cursor, err := r.collection.Find(ctx, filter, options)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var links []*link.Link
	for cursor.Next(ctx) {
		var model linkModel
		if err := cursor.Decode(&model); err != nil {
			return nil, err
		}

		expiresAtTime := model.ExpiresAt.Time()
		links = append(links, &link.Link{
			ID:         model.ID.Hex(),
			ShortCode:  model.ShortCode,
			URL:        model.Url,
			ClickCount: model.ClickCount,
			ExpiresAt:  &expiresAtTime,
		})
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return links, nil
}
