package core_handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	"github.com/shordem/api.thryvo/lib/config"
	"github.com/shordem/api.thryvo/lib/constants"
	"github.com/shordem/api.thryvo/lib/helper"
	"github.com/shordem/api.thryvo/payload/response"
)

type MediaHandlerInterface interface {
	UploadMedia(c *fiber.Ctx) error
	GetMedia(c *fiber.Ctx) error
}

type mediaHandler struct {
	mediaService config.FileConfigInterface
}

func NewMediaHandler(mediaService config.FileConfigInterface) MediaHandlerInterface {
	return &mediaHandler{mediaService: mediaService}
}

func (h *mediaHandler) UploadMedia(c *fiber.Ctx) error {
	var resp response.Response

	file, err := c.FormFile("file")
	if err != nil {
		resp.Status = constants.ClientUnProcessableEntity
		resp.Message = "File is required"

		return c.Status(http.StatusUnprocessableEntity).JSON(resp)
	}

	media, err := h.mediaService.UploadFile(file)
	if err != nil {
		resp.Status = constants.ServerErrorExternalService
		resp.Message = err.Error()

		return c.Status(http.StatusInternalServerError).JSON(resp)
	}

	resp.Status = constants.SuccessOperationCompleted
	resp.Message = "Media uploaded successfully"
	resp.Data = map[string]interface{}{"results": media}

	return c.Status(http.StatusOK).JSON(resp)
}

func (h *mediaHandler) GetMedia(c *fiber.Ctx) error {
	var resp response.Response
	mediaId := c.Params("key")

	media, err := h.mediaService.GetObject(mediaId)
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
