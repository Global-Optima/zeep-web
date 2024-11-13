package storage

import "github.com/gin-gonic/gin"

func RegisterStorageRoutes(router *gin.RouterGroup, handler *StorageHandler) {
	storageGroup := router.Group("/storage")
	{
		storageGroup.POST("/", handler.UploadFileHandler)
		storageGroup.DELETE("/", handler.DeleteFileHandler)
		storageGroup.GET("/", handler.GetFileURLHandler)
		storageGroup.GET("/download", handler.DownloadAndSaveFileHandler) // temp endpoint for testing purposes
		storageGroup.GET("/list-buckets", handler.ListBucketsHandler)     // temp endpoint for testing purposes
	}
}
