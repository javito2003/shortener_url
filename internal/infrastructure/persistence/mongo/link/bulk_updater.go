package link

import (
	"context"

	"github.com/javito2003/shortener_url/internal/config"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// LinkBulkUpdater implementa la interfaz de actualización en lote.
type LinkBulkUpdater struct {
	collection *mongo.Collection
}

func NewLinkBulkUpdater(db *mongo.Client) *LinkBulkUpdater {
	return &LinkBulkUpdater{collection: db.Database(config.AppConfig.Mongo.Database).Collection("links")}
}

// IncrementClickCounts usa BulkWrite para máxima eficiencia.
func (u *LinkBulkUpdater) IncrementClickCounts(ctx context.Context, counts map[string]int64) error {
	if len(counts) == 0 {
		return nil
	}

	writeModels := make([]mongo.WriteModel, 0, len(counts))

	for shortCode, increment := range counts {
		model := mongo.NewUpdateOneModel().
			SetFilter(bson.M{"short_code": shortCode}).
			SetUpdate(bson.M{"$inc": bson.M{"click_count": increment}})
		writeModels = append(writeModels, model)
	}

	// Unordered: si una actualización falla, las otras continúan.
	_, err := u.collection.BulkWrite(ctx, writeModels)
	return err
}
