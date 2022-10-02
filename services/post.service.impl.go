package services

import (
	"context"
	"errors"
	"time"

	"github.com/wpcodevo/golang-mongodb/models"
	"github.com/wpcodevo/golang-mongodb/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PostServiceImpl struct {
	postCollection *mongo.Collection
	ctx            context.Context
}

func NewPostService(postCollection *mongo.Collection, ctx context.Context) PostService {
	return &PostServiceImpl{postCollection, ctx}
}

func (p *PostServiceImpl) CreatePost(post *models.CreatePostRequest) (*models.DBPost, error) {
	post.CreateAt = time.Now()
	post.UpdatedAt = post.CreateAt
	res, err := p.postCollection.InsertOne(p.ctx, post)

	if err != nil {
		if er, ok := err.(mongo.WriteException); ok && er.WriteErrors[0].Code == 11000 {
			return nil, errors.New("post with that title already exists")
		}
		return nil, err
	}

	opt := options.Index()
	opt.SetUnique(true)

	index := mongo.IndexModel{Keys: bson.M{"title": 1}, Options: opt}

	if _, err := p.postCollection.Indexes().CreateOne(p.ctx, index); err != nil {
		return nil, errors.New("could not create index for title")
	}

	var newPost *models.DBPost
	query := bson.M{"_id": res.InsertedID}
	if err = p.postCollection.FindOne(p.ctx, query).Decode(&newPost); err != nil {
		return nil, err
	}

	return newPost, nil
}

func (p *PostServiceImpl) UpdatePost(id string, data *models.UpdatePost) (*models.DBPost, error) {
	doc, err := utils.ToDoc(data)
	if err != nil {
		return nil, err
	}

	obId, _ := primitive.ObjectIDFromHex(id)
	query := bson.D{{Key: "_id", Value: obId}}
	update := bson.D{{Key: "$set", Value: doc}}
	res := p.postCollection.FindOneAndUpdate(p.ctx, query, update, options.FindOneAndUpdate().SetReturnDocument(1))

	var updatedPost *models.DBPost
	if err := res.Decode(&updatedPost); err != nil {
		return nil, errors.New("no post with that Id exists")
	}

	return updatedPost, nil
}

func (p *PostServiceImpl) FindPostById(id string) (*models.DBPost, error) {
	obId, _ := primitive.ObjectIDFromHex(id)

	query := bson.M{"_id": obId}

	var post *models.DBPost

	if err := p.postCollection.FindOne(p.ctx, query).Decode(&post); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("no document with that Id exists")
		}

		return nil, err
	}

	return post, nil
}

func (p *PostServiceImpl) FindPosts(page int, limit int) ([]*models.DBPost, error) {
	if page == 0 {
		page = 1
	}

	if limit == 0 {
		limit = 10
	}

	skip := (page - 1) * limit

	opt := options.FindOptions{}
	opt.SetLimit(int64(limit))
	opt.SetSkip(int64(skip))
	opt.SetSort(bson.M{"created_at": -1})

	query := bson.M{}

	cursor, err := p.postCollection.Find(p.ctx, query, &opt)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(p.ctx)

	var posts []*models.DBPost

	for cursor.Next(p.ctx) {
		post := &models.DBPost{}
		err := cursor.Decode(post)

		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	if len(posts) == 0 {
		return []*models.DBPost{}, nil
	}

	return posts, nil
}

func (p *PostServiceImpl) DeletePost(id string) error {
	obId, _ := primitive.ObjectIDFromHex(id)
	query := bson.M{"_id": obId}

	res, err := p.postCollection.DeleteOne(p.ctx, query)
	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return errors.New("no document with that Id exists")
	}

	return nil
}
