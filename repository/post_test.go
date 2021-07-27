package repository

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/KihaRaito/sofupo-backend/model"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"regexp"
	"testing"
)

// PostRepositoryTestSuite テストスイートの構造体
type PostRepositoryTestSuite struct {
	suite.Suite
	postRepository PostRepositoryImpl
	mock sqlmock.Sqlmock
}

// SetupTest テストのセットアップ
func (suite *PostRepositoryTestSuite) SetupTest() {
	db, mock, _ := sqlmock.New()
	DB, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	suite.mock = mock
	suite.postRepository.DB = DB
}

// TearDownTest テスト終了時の処理
func (suite *PostRepositoryTestSuite) TearDownTest() {
	db, _ := suite.postRepository.DB.DB()
	db.Close()
}

// TestPostRepositoryTestSuite テストスイートの実行
func TestPostRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(PostRepositoryTestSuite))
}

// TestCreate Createのテスト
func (suite *PostRepositoryTestSuite) TestCreate() {
	suite.Run("create a post", func() {
		newId := 1
		rows := sqlmock.NewRows([]string{"id"}).AddRow(newId)
		suite.mock.ExpectBegin()
		query := regexp.QuoteMeta(
			`INSERT INTO "posts" ("created_at","updated_at","deleted_at","user_id","comment","shop_name","shop_address","image","score") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING "id"`)
		suite.mock.ExpectQuery(query).WillReturnRows(rows)
		suite.mock.ExpectCommit()
		post := &model.Post{
			UserID: "sample@gmail.com",
			Comment: "sample comment",
			ShopName: "sample shop",
			ShopAddress: "sample shop address",
			Image: "sample image path",
			Score: 100.0,
		}
		err := suite.postRepository.Create(post)

		if err != nil {
			suite.Fail(err.Error())
		}
	})
}
