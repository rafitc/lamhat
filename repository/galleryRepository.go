package repository

import (
	"lamhat/core"
	customErrors "lamhat/errors"
	"lamhat/model"
	"lamhat/utils"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func GetGalleryDetails(ctx *gin.Context, gallery_id int, user_id int, tx pgx.Tx) (model.GetGallery, error) {
	var gallery_files model.GetGallery
	const query = `select gallery_name, gs.status, g.created_at,
					jsonb_agg(gf.file_path) as files from app.gallery g 
					join app.gallery_status gs 
					on g.gallery_status_id = gs.id
					join app.gallery_files gf 
					on gf.gallery_id = g.id 
					where g.user_id = $1 and g.id = $2
					group by 1,2,3`
	core.Sugar.Infof("Running %v", query, user_id, gallery_id)
	row := tx.QueryRow(ctx, query, user_id, gallery_id)
	err := row.Scan(
		&gallery_files.GalleryName,
		&gallery_files.Status,
		&gallery_files.CreatedAt,
		&gallery_files.Files,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return gallery_files, customErrors.ErrNoGalleryDetailsFound
		}
		// Handle other errors
		Sugar_logger.Errorf("Error querying status table: %v", err)
		return gallery_files, err
	}

	return gallery_files, nil
}

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

func InsertFileInfo(ctx *gin.Context, upload_status []utils.UploadStatus, tx pgx.Tx) error {

	var uploadFiles [][]interface{}
	for _, each := range upload_status {
		uploadFiles = append(uploadFiles, []interface{}{each.Gallery_id, each.Objectname, each.Status, each.Bucketname})
	}

	copyCount, err := tx.CopyFrom(ctx,
		pgx.Identifier{"app", "gallery_files"},
		[]string{"gallery_id", "file_path", "is_active", "bucket_name"},
		pgx.CopyFromRows(uploadFiles))

	if err != nil {
		return err
	}
	core.Sugar.Infof("Inserted Row of %d count", copyCount)

	return err

}

func Publish(ctx *gin.Context, gallery_id int, status_id int, tx pgx.Tx) error {
	const query string = "UPDATE app.gallery SET gallery_status_id = $1 WHERE id = $2"
	core.Sugar.Infof("Running %s", query)

	_, err := tx.Exec(ctx, query, status_id, gallery_id)

	if err != nil {
		return err
	}
	err = tx.Commit(ctx)

	if err != nil {
		return err
	}
	return nil
}
