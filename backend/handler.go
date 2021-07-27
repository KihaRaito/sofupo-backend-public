package backend

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/KihaRaito/sofupo-backend/model"
	"github.com/KihaRaito/sofupo-backend/repository"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rekognition"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// HandleGet 指定した投稿を表示させるhandler
func HandleGet(ctx *gin.Context) {
	db := DbConn()
	DB, _ := db.DB()
	defer DB.Close()

	// クエリパラメーターから投稿のid取得
	id, err := strconv.Atoi(ctx.Query("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request parameter",
		})
	}

	// 指定したidの投稿取得
	postRepository := repository.PostRepositoryImpl{DB: db}
	post, err := postRepository.RetrievePost(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "post not found",
		})
	}

	accessKey := os.Getenv("ACCESS_KEY_ID")
	privateKey := os.Getenv("SECRET_ACCESS_KEY")
	region := os.Getenv("S3_REGION")
	bucketName := os.Getenv("S3_BUCKET_NAME")
	fileName := filepath.Base(post.Image)

	url, err := GetPresignUrl(accessKey, privateKey, region, bucketName, fileName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to sign request",
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "retrieving post successfully",
		"post":    post,
		"url":     url,
	})
}

// HandleGetAll 投稿を全て表示させるhandler
func HandleGetAll(ctx *gin.Context) {
	db := DbConn()
	DB, _ := db.DB()
	defer DB.Close()

	// 投稿を全て取得
	postRepository := repository.PostRepositoryImpl{DB: db}
	posts, err := postRepository.RetrievePosts()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "retrieving posts failed",
		})
	}

	log.Info().Msgf("fetch %v posts", len(posts))

	ctx.JSON(http.StatusOK, gin.H{
		"message": "retrieving posts successfully",
		"posts":   posts,
	})
}

// HandleGetAllByUser 指定したユーザーの投稿を全て取得して表示させるhandler
func HandleGetAllByUser(ctx *gin.Context) {
	db := DbConn()
	DB, _ := db.DB()
	defer DB.Close()

	// クエリパラメーターからユーザーのIDを取得
	user_id := ctx.Query("id")

	// 指定ユーザーの投稿を全て取得
	postRepository := repository.PostRepositoryImpl{DB: db}
	posts, err := postRepository.RetrievePostsByUser(user_id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "posts not found",
			"posts":   []model.Shop{},
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "retrieving posts by user successfully",
		"posts":   posts,
	})
}

// HandleGetMyShops 指定したユーザーが投稿した店舗の取得をして表示させるhandler
func HandleGetMyShops(ctx *gin.Context) {
	db := DbConn()
	DB, _ := db.DB()
	defer DB.Close()

	// クエリパラメーターからユーザーのIDを取得
	user_id := ctx.Query("id")

	// 指定したユーザーが投稿した店舗の取得
	postRepository := repository.PostRepositoryImpl{DB: db}
	posts, err := postRepository.RetrievePostsByUser(user_id)
	myshops, err := GetShops(posts)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "shops not found",
			"myshops": []model.Shop{},
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "retrieving my shops successfully",
		"myshops": myshops,
	})
}

// HandleConfirm 投稿を確認する画面を表示させるhandler
func HandleConfirm(ctx *gin.Context) {
	var post model.Post
	var err error

	_ = ctx.Request.ParseForm()

	// 編集フォームから投稿IDを取得
	id := ctx.PostForm("id")

	// 投稿を構造体にセット
	if id != "" {
		id, err := strconv.Atoi(id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "invalid post id",
			})
		}
		post.ID = id
	}

	jsonStr := ctx.Request.FormValue("formData")
	log.Info().Msgf("jsonStr=%v", jsonStr)
	json.Unmarshal([]byte(jsonStr), &post)

	// フォームから画像を取得
	file, fileHeader, err := ctx.Request.FormFile("image")
	defer file.Close()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "invalid post image",
		})
	}

	// sessionの作成
	accessKey := os.Getenv("ACCESS_KEY_ID")
	privateKey := os.Getenv("SECRET_ACCESS_KEY")
	region := os.Getenv("S3_REGION")
	bucketName := os.Getenv("S3_BUCKET_NAME")
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(accessKey, privateKey, ""),
		Region:      aws.String(region),
	}))
	uploader := s3manager.NewUploader(sess)

	// S3にupload
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileHeader.Filename),
		Body:   file,
	})

	if result == nil || err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "uploading image failed",
		})
	}
	log.Info().Msgf("%v", result.Location)

	// ファイル名を取得
	post.Image = result.Location

	model := rekognition.New(sess)
	output, err := model.DetectLabels(
		&rekognition.DetectLabelsInput{
			Image: &rekognition.Image{
				S3Object: &rekognition.S3Object{
					Bucket: aws.String(bucketName),
					Name:   aws.String(fileHeader.Filename),
				}},
			// MinConfidence: aws.Float64(0.75)
		})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "recognizing image failed",
		})
	}
	// post.Score = "0.0"
	post.Score = 0.0
	for _, label := range output.Labels {
		if *label.Name == "Ice Cream" {
			post.Score = *label.Confidence
			// rv := reflect.ValueOf(*label.Confidence).Interface().(float64)
			// post.Score = strconv.FormatFloat(rv, 'f', 2, 64)
		}
	}
	log.Info().Msgf("post=%v", post)

	svc := s3.New(sess)
	fileName := filepath.Base(post.Image)
	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
	})
	url, err := req.Presign(2 * time.Minute)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to sign request",
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "posting successfully",
		"post":    post,
		"url":     url,
	})
}

// HandleEdit 投稿を編集するフォーム画面を表示させるhandler
func HandleEdit(ctx *gin.Context) {
	db := DbConn()
	DB, _ := db.DB()
	defer DB.Close()

	// クエリパラメーターから投稿IDを取得
	id, err := strconv.Atoi(ctx.Query("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request parameter",
		})
	}

	// 指定した投稿の取得
	postRepository := repository.PostRepositoryImpl{DB: db}
	post, err := postRepository.RetrievePost(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "post not found",
		})
	}

	accessKey := os.Getenv("ACCESS_KEY_ID")
	privateKey := os.Getenv("SECRET_ACCESS_KEY")
	region := os.Getenv("S3_REGION")
	bucketName := os.Getenv("S3_BUCKET_NAME")
	fileName := filepath.Base(post.Image)

	url, err := GetPresignUrl(accessKey, privateKey, region, bucketName, fileName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to sign request",
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "",
		"post":    post,
		"url":     url,
	})
}

// HandlePost 投稿を作成するhandler
func HandlePost(ctx *gin.Context) {
	db := DbConn()
	DB, _ := db.DB()
	defer DB.Close()

	// 投稿の作成
	var post model.Post
	post.UserID = ctx.PostForm("user_id")
	post.Comment = ctx.PostForm("comment")
	post.ShopName = ctx.PostForm("shop_name")
	post.ShopAddress = ctx.PostForm("shop_address")
	post.Image = ctx.PostForm("image")
	// post.Score = ctx.PostForm("score")
	score, _ := strconv.ParseFloat(ctx.PostForm("score"), 64)
	post.Score = score

	postRepository := repository.PostRepositoryImpl{DB: db}
	err := postRepository.Create(&post)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "creating post failed",
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "creating post successfully",
	})
}

// HandlePut 投稿を更新するhandler
func HandlePut(ctx *gin.Context) {
	db := DbConn()
	DB, _ := db.DB()
	defer DB.Close()

	// クエリパラメーターから投稿IDを取得
	id, err := strconv.Atoi(ctx.PostForm("ID"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request parameter",
		})
	}

	// 指定した投稿の取得
	postRepository := repository.PostRepositoryImpl{DB: db}
	post, err := postRepository.RetrievePost(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "post not found",
		})
	}
	post.UserID = ctx.PostForm("user_id")
	post.Comment = ctx.PostForm("comment")
	post.ShopName = ctx.PostForm("shop_name")
	post.ShopAddress = ctx.PostForm("shop_address")
	post.Image = ctx.PostForm("image")
	// post.Score = ctx.PostForm("score")
	score, _ := strconv.ParseFloat(ctx.PostForm("score"), 64)
	post.Score = score

	err = postRepository.Update(post)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "updating post failed",
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "updating post successfully",
	})
}

// HandleDelete 投稿を削除するhandler
func HandleDelete(ctx *gin.Context) {
	db := DbConn()
	DB, _ := db.DB()
	defer DB.Close()

	// IDを取得
	id, err := strconv.Atoi(ctx.PostForm("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request parameter",
		})
	}
	user_id := ctx.PostForm("user_id")

	// 指定した投稿の取得
	postRepository := repository.PostRepositoryImpl{DB: db}
	post, err := postRepository.RetrievePost(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "post not found",
		})
	}
	if post.UserID != user_id {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "deleting post failed",
		})
	}

	// 投稿を削除
	err = postRepository.Delete(post)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "deleting post failed",
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "deleting post successfully",
	})
}

// HandleMerge ユーザー間で投稿店舗のMergeを行うhandler
func HandleMerge(ctx *gin.Context) {
	db := DbConn()
	DB, _ := db.DB()
	defer DB.Close()

	// mergeするuserのIDを取得
	merge_user_id := ctx.PostForm("merge_user_id")
	log.Info().Msgf("merge_user_id = %v", merge_user_id)

	// mergeされるuserのIDを取得
	merged_user_id := ctx.PostForm("merged_user_id")
	log.Info().Msgf("merged_user_id = %v", merged_user_id)

	// ユーザー間での投稿店舗のMerge
	postRepository := repository.PostRepositoryImpl{DB: db}
	mergePosts, err := postRepository.RetrievePostsByUser(merge_user_id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "merge shops not found",
			"myshops": []model.Shop{},
		})
	}
	mergedPosts, err := postRepository.RetrievePostsByUser(merged_user_id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "merged shops not found",
			"myshops": []model.Shop{},
		})
	}
	mergeShops, _ := GetShops(mergePosts)
	mergedShops, _ := GetShops(mergedPosts)
	myshops := Merge(mergeShops, mergedShops)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "merging my shops successfully",
		"myshops": myshops,
	})
}