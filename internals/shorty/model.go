package shorty

import "go.mongodb.org/mongo-driver/bson/primitive"

type Record struct {
	ID       primitive.ObjectID `bson:"_id"`
	LongURL  string             `bson:"long_url" json:"long_url"`
	ShortURL string             `bson:"short_url" json:"short_url"`
}
