package core_handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/shordem/api.thryvo/dto"
	"github.com/shordem/api.thryvo/handler"
	"github.com/shordem/api.thryvo/lib/constants"
	"github.com/shordem/api.thryvo/payload/request"
	"github.com/shordem/api.thryvo/payload/response"
	core_service "github.com/shordem/api.thryvo/service/core"
)

type FolderHandlerInterface interface {
	CreateFolder(c *fiber.Ctx) error
	GetUserFolders(c *fiber.Ctx) error
	GetFoldersByParent(c *fiber.Ctx) error
	UpdateFolder(c *fiber.Ctx) error
	DeleteFolder(c *fiber.Ctx) error
}

type folderHandler struct {
	folderService core_service.FolderServiceInterface
}

func NewFolderHandler(folderService core_service.FolderServiceInterface) FolderHandlerInterface {
	return &folderHandler{folderService: folderService}
}

func (h *folderHandler) CreateFolder(c *fiber.Ctx) error {
	var resp response.Response
	var folderDto dto.FolderDTO
	var createFolderReq request.CreateFolderRequest

	userId := handler.GetUserId(c)
	if err := c.BodyParser(&createFolderReq); err != nil {
		resp.Status = constants.ClientUnProcessableEntity
		resp.Message = "Invalid request"

		return c.Status(http.StatusUnprocessableEntity).JSON(resp)
	}

	folderDto.UserID = userId
	folderDto.Name = createFolderReq.Name
	folderDto.ParentID = createFolderReq.ParentID

	_, err := h.folderService.CreateFolder(folderDto)
	if err != nil {
		resp.Status = constants.ServerErrorInternal
		resp.Message = err.Error()

		return c.Status(http.StatusInternalServerError).JSON(resp)
	}

	resp.Status = constants.SuccessOperationCompleted
	resp.Message = "Folder created successfully"

	return c.JSON(resp)
}

func (h *folderHandler) GetUserFolders(c *fiber.Ctx) error {
	var resp response.Response

	userId := handler.GetUserId(c)

	folders, err := h.folderService.FindFoldersByUserId(userId)
	if err != nil {
		resp.Status = constants.ServerErrorInternal
		resp.Message = "Failed to fetch folders"

		return c.Status(http.StatusInternalServerError).JSON(resp)
	}

	resp.Status = constants.SuccessOperationCompleted
	resp.Message = "Folders fetched successfully"
	resp.Data = map[string]interface{}{"result": folders}

	return c.JSON(resp)
}

func (h *folderHandler) GetFoldersByParent(c *fiber.Ctx) error {
	var resp response.Response

	userId := handler.GetUserId(c)
	parentId, err := uuid.Parse(c.Params("parent_id"))
	if err != nil {
		resp.Status = constants.ClientErrorBadRequest
		resp.Message = err.Error()

		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	folders, err := h.folderService.FindFoldersByParentId(userId, parentId)
	if err != nil {
		resp.Status = constants.ServerErrorInternal
		resp.Message = "Failed to fetch folders"

		return c.Status(http.StatusInternalServerError).JSON(resp)
	}

	resp.Status = constants.SuccessOperationCompleted
	resp.Message = "Folders fetched successfully"
	resp.Data = map[string]interface{}{"result": folders}

	return c.JSON(resp)
}

func (h *folderHandler) UpdateFolder(c *fiber.Ctx) error {
	var resp response.Response
	var updateFolderReq request.UpdateFolderRequest
	var folderDto dto.FolderDTO

	folderId, err := uuid.Parse(c.Params("id"))
	if err != nil {
		resp.Status = constants.ClientErrorBadRequest
		resp.Message = err.Error()

		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	if err := c.BodyParser(&updateFolderReq); err != nil {
		resp.Status = constants.ClientErrorBadRequest
		resp.Message = err.Error()

		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	userId := handler.GetUserId(c)

	folderDto.ID = folderId
	folderDto.UserID = userId
	folderDto.Name = updateFolderReq.Name
	folderDto.ParentID = updateFolderReq.ParentID

	if _, err = h.folderService.UpdateFolder(folderDto); err != nil {
		resp.Status = constants.ServerErrorInternal
		resp.Message = "Failed to update folder"

		return c.Status(http.StatusInternalServerError).JSON(resp)
	}

	resp.Status = constants.SuccessOperationCompleted
	resp.Message = "Folder updated successfully"

	return c.JSON(resp)
}

func (h *folderHandler) DeleteFolder(c *fiber.Ctx) error {
	var resp response.Response

	folderId, err := uuid.Parse(c.Params("id"))
	if err != nil {
		resp.Status = constants.ClientErrorBadRequest
		resp.Message = err.Error()

		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	userId := handler.GetUserId(c)

	if err = h.folderService.DeleteFolder(folderId, userId); err != nil {
		resp.Status = constants.ServerErrorInternal
		resp.Message = "Failed to delete folder"

		return c.Status(http.StatusInternalServerError).JSON(resp)
	}

	resp.Status = constants.SuccessOperationCompleted
	resp.Message = "Folder deleted successfully"

	return c.JSON(resp)
}
