package repository

import (
	"lamhat/core"
	customErrors "lamhat/errors"
	"lamhat/model"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func GetGalleryStatus(ctx *gin.Context, status string, tx pgx.Tx) (model.GalleryStatus, error) {
	// Look for status
	var status_id model.GalleryStatus
	const query string = "SELECT id, status FROM app.gallery_status WHERE status = $1;"
	core.Sugar.Infof("Running %s", query)
	row := tx.QueryRow(ctx, query, status)
	err := row.Scan(
		&status_id.Id,
		&status_id.Status,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return status_id, customErrors.ErrStatusNotFound
		}
		// Handle other errors
		Sugar_logger.Errorf("Error querying status table: %v", err)
		return status_id, err
	}
	return status_id, nil
}

func CreateGallery(ctx *gin.Context, gallery model.CreateGallery, tx pgx.Tx) (model.GalleryModel, error) {
	var gallery_model model.GalleryModel
	const query string = "INSERT INTO app.gallery (user_id, gallery_name, gallery_status_id) VALUES ($1, $2, $3) RETURNING id, user_id, gallery_name, gallery_status_id;"
	core.Sugar.Infof("Running %s", query)

	err := tx.QueryRow(ctx, query, gallery.UserId, gallery.Name, gallery.Status).Scan(
		&gallery_model.Id,
		&gallery_model.UserId,
		&gallery_model.Name,
		&gallery_model.Status,
	)

	if err != nil {
		return gallery_model, err
	}

	// Commit the transaction after successful insertion
	if err := tx.Commit(ctx); err != nil {
		return gallery_model, err
	}
	return gallery_model, nil
}
