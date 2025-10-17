package link

import (
	"context"

	"github.com/javito2003/shortener_url/internal/config"
	link "github.com/javito2003/shortener_url/internal/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type repository struct {
	collection *mongo.Collection
}

type linkModel struct {
	ID         bson.ObjectID `bson:"_id,omitempty"`
	ShortCode  string        `bson:"short_code"`
	Url        string        `bson:"url"`
	ClickCount int           `bson:"click_count"`
}

func NewRepository(db *mongo.Client) *repository {
	return &repository{collection: db.Database(config.AppConfig.Mongo.Database).Collection("links")}
}

func (r *repository) Save(ctx context.Context, link *link.Link) (*link.Link, error) {
	result, err := r.collection.InsertOne(ctx, linkModel{
		ShortCode:  link.ShortCode,
		Url:        link.URL,
		ClickCount: link.ClickCount,
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

	return &link.Link{
		ID:         model.ID.Hex(),
		ShortCode:  model.ShortCode,
		URL:        model.Url,
		ClickCount: model.ClickCount,
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

	return &link.Link{
		ID:         model.ID.Hex(),
		ShortCode:  model.ShortCode,
		URL:        model.Url,
		ClickCount: model.ClickCount,
	}, true, nil
}

// func (r *repository) IncrementBulk(ctx context.Context, links []*ports.LinkData) error {
// 	models := make([]mongo.WriteModel, 0, len(links))

// 	batch := 20

// 	for i := 0; i < len(links); i += batch {
// 		end := i + batch
// 		if end > len(links) {
// 			end = len(links)
// 		}

// 		models = models[:0]
// 		bulk := []mongo.WriteModel{}
// 		for _, link := range links[i:end] {
// 			bulk = append(bulk, mongo.NewUpdateOneModel().SetFilter(bson.M{"short_code": link.Id}).SetUpdate(bson.M{"$inc": bson.M{"click_count": link.Increment}}))
// 		}

// 		models = append(models, bulk...)
// 		_, err := r.collection.BulkWrite(ctx, models)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	_, err := r.collection.BulkWrite(ctx, models)
// 	return err
// }
