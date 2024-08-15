package core_handler

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/shordem/api.thryvo/dto"
	"github.com/shordem/api.thryvo/handler"
	"github.com/shordem/api.thryvo/lib/constants"
	"github.com/shordem/api.thryvo/lib/helper"
	"github.com/shordem/api.thryvo/payload/response"
	core_repository "github.com/shordem/api.thryvo/repository/core"
	core_service "github.com/shordem/api.thryvo/service/core"
)

type FileHandlerInterface interface {
	UploadFile(c *fiber.Ctx) error
	GetUserFiles(c *fiber.Ctx) error
	GetFile(c *fiber.Ctx) error
}

type fileHandler struct {
	fileService core_service.FileServiceInterface
}

func NewFileHandler(fileService core_service.FileServiceInterface) FileHandlerInterface {
	return &fileHandler{fileService: fileService}
}

func (h *fileHandler) GeneratePageable(c *fiber.Ctx) (filePageable core_repository.FilePageable) {
	var resp response.Response

	filePageable.Pageable = handler.GeneratePageable(c)
	filePageable.UserId = handler.GetUserId(c)

	if folderId := c.Query("folder_id"); folderId != "" {
		folderIdParsed, err := uuid.Parse(folderId)
		if err != nil {
			resp.Status = constants.ClientErrorBadRequest
			resp.Message = err.Error()

			c.Status(http.StatusBadRequest).JSON(resp)
		}

		filePageable.FolderId = folderIdParsed
	}

	if hasFolder := c.Query("has_folder"); hasFolder != "" {
		hasFolderConv, err := strconv.ParseBool(hasFolder)
		if err != nil {
			resp.Status = constants.ClientUnProcessableEntity
			resp.Message = "has_folder query is not a valid boolean format"

			c.Status(http.StatusUnprocessableEntity).JSON(resp)
		}

		filePageable.HasFolder = hasFolderConv
	}

	return filePageable
}

func (h *fileHandler) UploadFile(c *fiber.Ctx) error {
	var resp response.Response
	var fileDto dto.FileDTO

	userId := handler.GetUserId(c)
	folderId := c.FormValue("folder_id")
	file, err := c.FormFile("file")

	if err != nil {
		resp.Status = constants.ClientUnProcessableEntity
		resp.Message = "File is required"

		return c.Status(http.StatusUnprocessableEntity).JSON(resp)
	}

	fileDto.UserID = userId
	if folderId != "" {
		folderUUID, err := uuid.Parse(folderId)
		if err != nil {
			resp.Status = constants.ClientUnProcessableEntity
			resp.Message = "Invalid folder ID"

			return c.Status(http.StatusUnprocessableEntity).JSON(resp)
		}
		fileDto.FolderID = &folderUUID
	}
	fileDto.OriginalName = file.Filename
	fileDto.MimeType = file.Header.Get("Content-Type")
	fileDto.Size = file.Size
	fileDto.Visibility = core_service.FileVisibilityPublic

	fileKey, err := h.fileService.UploadFile(fileDto, file)
	if err != nil {
		resp.Status = constants.ServerErrorExternalService
		resp.Message = err.Error()

		return c.Status(http.StatusInternalServerError).JSON(resp)
	}

	resp.Status = constants.SuccessOperationCompleted
	resp.Message = "file uploaded successfully"
	resp.Data = map[string]interface{}{"result": fileKey}

	return c.Status(http.StatusOK).JSON(resp)
}

func (h *fileHandler) GetUserFiles(c *fiber.Ctx) error {
	var resp response.Response
	pageable := h.GeneratePageable(c)

	files, pagination, err := h.fileService.FindAllFiles(pageable)
	if err != nil {
		resp.Status = constants.ServerErrorExternalService
		resp.Message = "Failed to get user files"

		return c.Status(http.StatusInternalServerError).JSON(resp)
	}

	resp.Status = constants.SuccessOperationCompleted
	resp.Message = "User files fetched successfully"
	resp.Data = map[string]interface{}{"pagination": pagination, "result": files}

	return c.Status(http.StatusOK).JSON(resp)
}

func (h *fileHandler) GetFile(c *fiber.Ctx) error {
	var resp response.Response
	mediaId := c.Params("key")

	media, err := h.fileService.GetFile(mediaId)
	if err != nil {
		resp.Status = constants.ServerErrorExternalService
		resp.Message = "Failed to get media"
		c.Status(http.StatusInternalServerError).JSON(resp)
		return err
	}

	// stream media
	c.Set("Content-Type", *media.ContentType)
	c.Set("Content-Disposition", "inline")
	c.Set("Content-Length", helper.Int64ToString(*media.ContentLength))
	c.SendStream(media.Body)

	return nil
}
