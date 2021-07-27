package repository

import (
	"github.com/KihaRaito/sofupo-backend/model"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type PostRepositoryImpl struct {
	DB *gorm.DB
}

type PostRepository interface {
	RetrievePostsByUser(user_id string)
	RetrievePosts()
	RetrievePost(id int)
	Create(post *model.Post)
	Update(post *model.Post)
	Delete(post *model.Post)
}

// RetrievePostsByUser ユーザーの投稿を全て取得するFunction
func (postRepo PostRepositoryImpl) RetrievePostsByUser(user_id string) (posts []model.Post, err error) {
	err = postRepo.DB.Where("user_id = ?", user_id).Find(&posts).Error
	log.Info().Msg("retrieve all posts by user_id")
	return posts, err
}

// RetrievePosts 投稿を全て取得するFunction
func (postRepo PostRepositoryImpl) RetrievePosts() (posts []model.Post, err error) {
	err = postRepo.DB.Find(&posts).Error
	log.Info().Msg("retrieve all posts")
	return posts, err
}

// RetrievePost IDにマッチした投稿を取得するFunction
func (postRepo PostRepositoryImpl) RetrievePost(id int) (*model.Post, error) {
	var post model.Post
	err := postRepo.DB.Where("id = ?", id).Find(&post).Error
	log.Info().Msgf("retrieve post(id=%v)", id)
	return &post, err
}

// Create 投稿を作成するFunction
func (postRepo PostRepositoryImpl) Create(post *model.Post) (err error) {
	err = postRepo.DB.Create(post).Error
	log.Info().Msgf("create post(id=%v)", post.ID)
	return err
}

// Update 投稿を更新するFunction
func (postRepo PostRepositoryImpl) Update(post *model.Post) (err error) {
	err = postRepo.DB.Save(post).Error
	log.Info().Msgf("update post(id=%v)", post.ID)
	return err
}

// Delete 投稿を削除するFunction
func (postRepo PostRepositoryImpl) Delete(post *model.Post) (err error) {
	err = postRepo.DB.Delete(post).Error
	log.Info().Msgf("delete post(id=%v)", post.ID)
	return err
}

