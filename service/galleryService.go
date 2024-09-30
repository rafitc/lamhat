package service

import (
	"lamhat/core"
	"lamhat/model"
	"lamhat/repository"
	"lamhat/utils"

	"github.com/gin-gonic/gin"
)

func FetchGallery(ctx *gin.Context, gallery_id int, user_id int) model.Response {
	var result model.Response
	Sugar_logger.Infof("Fetching gallery details with ID  %d", gallery_id)

	// Acquire a connection from the pool
	connection, err := repository.ConObjOfDB.Acquire(ctx)
	if err != nil {
		Sugar_logger.Errorf("Error while acquiring connection from the database pool!! %v", err.Error())

		result.Status = false
		result.Data = nil
		result.Code = 500
		result.ErrorMsg = err.Error()
		return result
	}
	defer connection.Release()

	// Get transaction from connection and use it till the end.
	// If any err, do rollback else do commit
	tx, err := connection.Begin(ctx)
	if err != nil {
		Sugar_logger.Errorf("Error in DB connection %v", err.Error())
		defer tx.Rollback(ctx)

		result.Status = false
		result.Data = nil
		result.Code = 500
		result.ErrorMsg = err.Error()
		return result
	}

	// Fetch the whole Gallery data
	gallery, err := repository.GetGalleryDetails(ctx, gallery_id, user_id, tx)
	if err != nil {
		Sugar_logger.Errorf("%v", err.Error())

		result.Status = false
		result.Data = nil
		result.Code = 400
		result.ErrorMsg = err.Error()
		return result
	}

	Sugar_logger.Errorf("Gallery %v", gallery)

	result.Status = true
	result.Data = gallery
	result.Code = 200
	result.ErrorMsg = ""
	return result
}

func CreateGallery(ctx *gin.Context, body model.CreateGallery) model.Response {
	var result model.Response
	Sugar_logger.Info("Gallery Creation started with body %s", body)

	// Acquire a connection from the pool
	connection, err := repository.ConObjOfDB.Acquire(ctx)
	if err != nil {
		Sugar_logger.Errorf("Error while acquiring connection from the database pool!! %v", err.Error())

		result.Status = false
		result.Data = nil
		result.Code = 500
		result.ErrorMsg = err.Error()
		return result
	}
	defer connection.Release()

	// Get transaction from connection and use it till the end.
	// If any err, do rollback else do commit
	tx, err := connection.Begin(ctx)
	if err != nil {
		Sugar_logger.Errorf("Error in DB connection %v", err.Error())
		defer tx.Rollback(ctx)

		result.Status = false
		result.Data = nil
		result.Code = 500
		result.ErrorMsg = err.Error()
		return result
	}
	// Create a gallery with the given name
	// Get status_id of the gallery
	const draft_status string = "DRAFT"
	status, err := repository.GetGalleryStatus(ctx, draft_status, tx)
	if err != nil {
		Sugar_logger.Errorf("Can't find gallery status %v", err.Error())

		result.Status = false
		result.Data = nil
		result.Code = 400
		result.ErrorMsg = err.Error()
		return result
	}
	// By default make everything DRAFT
	body.UserId = ctx.GetInt("userId")
	body.Status = status.Id // Assign Status ID
	gallery, err := repository.CreateGallery(ctx, body, tx)
	if err != nil {
		Sugar_logger.Errorf("Error while creating gallery %v", err.Error())

		result.Status = false
		result.Data = nil
		result.Code = 400
		result.ErrorMsg = err.Error()
		return result
	}

	result.Status = true
	result.Data = gallery
	result.Code = 200
	result.ErrorMsg = ""
	return result
}

func UploadIntoGallery(ctx *gin.Context, gallery_id int, user_id int) model.Response {
	var result model.Response

	Sugar_logger.Infof("Uploading files into gallery %d", gallery_id)
	data, err := utils.UploadIntoGallery(ctx, gallery_id, user_id)

	if err != nil {
		result.Status = false
		result.Data = nil
		result.Code = 400
		result.ErrorMsg = err.Error()
		return result
	}

	core.Sugar.Debug("upload done. starting DB update")
	// Acquire a connection from the pool
	connection, err := repository.ConObjOfDB.Acquire(ctx)
	if err != nil {
		Sugar_logger.Errorf("Error while acquiring connection from the database pool!! %v", err.Error())

		result.Status = false
		result.Data = nil
		result.Code = 500
		result.ErrorMsg = err.Error()
		return result
	}
	defer connection.Release()

	// Get transaction from connection and use it till the end.
	// If any err, do rollback else do commit
	tx, err := connection.Begin(ctx)
	if err != nil {
		Sugar_logger.Errorf("Error in DB connection %v", err.Error())
		defer tx.Rollback(ctx)

		result.Status = false
		result.Data = nil
		result.Code = 500
		result.ErrorMsg = err.Error()
		return result
	}

	err = repository.InsertFileInfo(ctx, data, tx)

	if err != nil {
		core.Sugar.Infof("Error while updating files")

		result.Status = false
		result.Data = nil
		result.Code = 402
		result.ErrorMsg = err.Error()
		return result
	}
	tx.Commit(ctx)

	result.Status = true
	result.Data = nil
	result.Code = 200
	result.ErrorMsg = ""
	return result
}

func PublishGallery(ctx *gin.Context, user_id int, gallery_status model.PublishGallery) model.Response {
	var result model.Response

	Sugar_logger.Infof("Publising gallery %d", gallery_status.GalleryId)

	// Acquire a connection from the pool
	connection, err := repository.ConObjOfDB.Acquire(ctx)
	if err != nil {
		Sugar_logger.Errorf("Error while acquiring connection from the database pool!! %v", err.Error())

		result.Status = false
		result.Data = nil
		result.Code = 500
		result.ErrorMsg = err.Error()
		return result
	}
	defer connection.Release()

	// Get transaction from connection and use it till the end.
	// If any err, do rollback else do commit
	tx, err := connection.Begin(ctx)
	if err != nil {
		Sugar_logger.Errorf("Error in DB connection %v", err.Error())
		defer tx.Rollback(ctx)

		result.Status = false
		result.Data = nil
		result.Code = 500
		result.ErrorMsg = err.Error()
		return result
	}

	status, err := repository.GetGalleryStatus(ctx, gallery_status.Status, tx)

	if err != nil {
		result.Status = false
		result.Data = nil
		result.Code = 500
		result.ErrorMsg = err.Error()
		return result
	}

	err = repository.Publish(ctx, gallery_status.GalleryId, status.Id, tx)
	if err != nil {
		result.Status = false
		result.Data = nil
		result.Code = 500
		result.ErrorMsg = err.Error()
		return result
	}

	result.Status = true
	result.Data = nil
	result.Code = 200
	result.ErrorMsg = ""
	return result

}
