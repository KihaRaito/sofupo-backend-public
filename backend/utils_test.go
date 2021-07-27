package backend

import (
	"github.com/KihaRaito/sofupo-backend/model"
	"testing"
)

func TestUnique(t *testing.T) {
	tests := []struct{
		message string
		shops []model.Shop
		want []model.Shop
	}{
		{
			message: "normal",
			shops: []model.Shop{
				{
					Name: "test",
					Address: "test",
				},
				{
					Name: "test",
					Address: "test",
				},
			},
			want: []model.Shop{
				{
					Name: "test",
					Address: "test",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.message, func(t *testing.T) {
			got := Unique(tt.shops)
			for i, expected := range tt.want {
				if expected != got[i] {
					t.Errorf("GetShops() = %v, want %v", expected, got[i])
				}
			}
		})
	}
}

func TestGetShops(t *testing.T) {
	tests := []struct{
		message string
		posts []model.Post
		want []model.Shop
	}{
		{
			message: "normal",
			posts: []model.Post{
				{
					ShopName: "test",
					ShopAddress: "test",
				},
			},
			want: []model.Shop{
				{
					Name: "test",
					Address: "test",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.message, func(t *testing.T) {
			got, _ := GetShops(tt.posts)
			for i, expected := range tt.want {
				if expected != got[i] {
					t.Errorf("GetShops() = %v, want %v", expected, got[i])
				}
			}
		})
	}
}

func TestMerge(t *testing.T) {
	tests := []struct{
		message string
		mergeShops []model.Shop
		mergedShops []model.Shop
		want []model.Shop
	}{
		{
			message: "normal",
			mergeShops: []model.Shop{
				{
					Name: "test1",
					Address: "test1",
				},
			},
			mergedShops: []model.Shop{
				{
					Name: "test2",
					Address: "test2",
				},
			},
			want: []model.Shop{
				{
					Name: "test2",
					Address: "test2",
				},
				{
					Name: "test1",
					Address: "test1",
				},
			},
		},
		{
			message: "mergedShops length is zero",
			mergeShops: []model.Shop{
				{
					Name: "test1",
					Address: "test1",
				},
			},
			mergedShops: []model.Shop{},
			want: []model.Shop{
				{
					Name: "test1",
					Address: "test1",
				},
			},
		},
		{
			message: "duplicate shop",
			mergeShops: []model.Shop{
				{
					Name: "test1",
					Address: "test1",
				},
				{
					Name: "test3",
					Address: "test3",
				},
			},
			mergedShops: []model.Shop{
				{
					Name: "test2",
					Address: "test2",
				},
				{
					Name: "test3",
					Address: "test3",
				},
			},
			want: []model.Shop{
				{
					Name: "test2",
					Address: "test2",
				},
				{
					Name: "test3",
					Address: "test3",
				},
				{
					Name: "test1",
					Address: "test1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.message, func(t *testing.T) {
			got := Merge(tt.mergeShops, tt.mergedShops)
			for i, expected := range tt.want {
				if expected != got[i] {
					t.Errorf("Merge() = %v, want %v", expected, got[i])
				}
			}
		})
	}
}