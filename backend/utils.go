package backend

import (
	"errors"
	"github.com/KihaRaito/sofupo-backend/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/rs/zerolog/log"
	"time"
)

// Unique Shopのリストから重複を取り除くFunction
func Unique(shops []model.Shop) (myshops []model.Shop) {
	m := make(map[model.Shop]bool)
	for _, shop := range shops {
		if !m[shop] {
			m[shop] = true
			myshops = append(myshops, shop)
		}
	}
	return myshops
}

// Merge ユーザー間で投稿店舗をMergeした差分の店舗を取得するFunction
func Merge(mergeShops []model.Shop, mergedShops []model.Shop) (shops []model.Shop) {
	for _, mergedShop := range mergedShops {
		for _, mergeShop := range mergeShops {
			if mergedShop == mergeShop {
				break
			} else {
				shops = append(shops, mergeShop)
			}
		}
	}

	if len(mergedShops) == 0 {
		shops = mergeShops
	}

	shops = append(mergedShops, shops...)
	shops = Unique(shops)

	for _, shop := range shops {
		log.Info().Msgf("shop_name=%v shop_address=%v", shop.Name, shop.Address)
	}
	return shops
}

// GetShops 指定ユーザーの投稿店舗を取得するFunction
func GetShops(posts []model.Post) (shops []model.Shop, err error) {
	var shop model.Shop
	for _, post := range posts {
		shop.Name = post.ShopName
		shop.Address = post.ShopAddress
		shops = append(shops, shop)
	}
	return shops, err
}

// GetPresignUrl 署名付きURLを取得するFunction
func GetPresignUrl(accessKey string, privateKey string, region string, bucketName string, fileName string) (url string, err error) {
	if accessKey == "" || privateKey == "" || region == "" || bucketName == "" || fileName == "" {
		err = errors.New("aws keys not found")
	}
	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(accessKey, privateKey, ""),
		Region: aws.String(region),
	})
	svc := s3.New(sess)
	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key: aws.String(fileName),
	})
	url, err = req.Presign(2 * time.Minute)
	return url, err
}
