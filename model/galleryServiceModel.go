package model

import "time"

type GalleryModel struct {
	Id     int    `json:"id" binding:"required"`
	UserId int    `json:"user_id" binding:"required"`
	Name   string `json:"gallery_name" binding:"required"`
	Status int    `json:"gallery_status_id" binding:"required"`
}

type CreateGallery struct {
	Name   string `json:"gallery_name" binding:"required"`
	Status int    `json:"gallery_status_id"`
	UserId int    `json:"user_id"`
}

type GalleryStatus struct {
	Id     int    `json:"id" binding:"required"`
	Status string `json:"status" binding:"required"`
}

type GetGallery struct {
	GalleryName string    `json:"gallery_name"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	Files       []string  `json:"files"`
}
