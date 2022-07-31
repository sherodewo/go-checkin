package controllers

import (
	"errors"
	"github.com/labstack/echo/v4"
	"go-checkin/dto"
	"go-checkin/models"
	"go-checkin/service"
	"go-checkin/utils/session"
	"gopkg.in/go-playground/validator.v9"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type JabatanController struct {
	BaseBackendController
	service *service.JabatanService
}

func NewJabatanController(service *service.JabatanService) JabatanController {
	return JabatanController{
		BaseBackendController: BaseBackendController{
			Menu:        "Jabatan",
			BreadCrumbs: []map[string]interface{}{},
		},
		service: service,
	}
}
func (c *JabatanController) Index(ctx echo.Context) error {
	breadCrumbs := map[string]interface{}{
		"menu": "List Role",
		"link": "/check/admin/jabatan",
	}
	return Render(ctx, "Jabatan List", "jabatan/index", c.Menu, session.GetFlashMessage(ctx),
		append(c.BreadCrumbs, breadCrumbs), nil)
}
func (c *JabatanController) List(ctx echo.Context) error {

	draw, err := strconv.Atoi(ctx.Request().URL.Query().Get("draw"))
	search := ctx.Request().URL.Query().Get("search[value]")
	start, err := strconv.Atoi(ctx.Request().URL.Query().Get("start"))
	length, err := strconv.Atoi(ctx.Request().URL.Query().Get("length"))
	order, err := strconv.Atoi(ctx.Request().URL.Query().Get("order[0][column]"))
	orderName := ctx.Request().URL.Query().Get("columns[" + strconv.Itoa(order) + "][name]")
	//orderAscDesc := ctx.Request().URL.Query().Get("order[0][dir]")

	recordTotal, recordFiltered, data, err := c.service.QueryDatatable(search, "desc", orderName, length, start)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	var action string

	//var role string
	listOfData := make([]map[string]interface{}, len(data))
	for k, v := range data {

		action = `<div class="btn-group open">
					<button class="btn btn-xs dropdown-toggle" type="button" data-toggle="dropdown" aria-expanded="true"> Actions</button>
                      <ul class="dropdown-menu" role="menu">
                      <li>
                      <a href="/check/admin/jabatan/edit/` + v.ID + `" style="text-decoration: none;font-weight: 400; color: #333;">
                      <i class="fa fa-edit"></i>Edit </a>
                      </li>
                      <li>
                      <a href="/check/admin/jabatan/detail/` + v.ID + `"style="text-decoration: none;font-weight: 400; color: #333;">
                      <i class="fa fa-user"></i>Detail </a>
                      </li>
                      <li>
                      <a href="javascript:;" onclick="Delete('` + v.ID + `')" data-toggle="tooltip" data-placement="right" title="Delete" style="text-decoration: none;font-weight: 400; color: #333;">
                      <i class="fa fa-trash" style="color: #ff4d65;"></i> Delete </a>
                      </li>
                      </ul>
                      </div>`

		listOfData[k] = map[string]interface{}{
			"id":          v.ID,
			"name":        v.Name,
			"description": v.Description,
			"action":      action,
		}
	}

	result := models.ResponseDatatable{
		Draw:            draw,
		RecordsTotal:    recordTotal,
		RecordsFiltered: recordFiltered,
		Data:            listOfData,
	}
	return ctx.JSON(http.StatusOK, &result)
}

func (c *JabatanController) Add(ctx echo.Context) error {
	breadCrumbs := map[string]interface{}{
		"menu": "Add",
		"link": "/check/admin/jabatan/add",
	}
	return Render(ctx, "Add Jabatan", "jabatan/add", c.Menu, session.GetFlashMessage(ctx),
		append(c.BreadCrumbs, breadCrumbs), nil)
}
func (c *JabatanController) Store(ctx echo.Context) error {
	var jabatanDto dto.JabatanDto
	if err := ctx.Bind(&jabatanDto); err != nil {
		return ctx.JSON(400, echo.Map{"message": "error binding data"})
	}
	if err := ctx.Validate(&jabatanDto); err != nil {
		var validationErrors []models.ValidationError
		if errs, ok := err.(validator.ValidationErrors); ok {
			validationErrors = models.WrapValidationErrors(errs)
		}
		return ctx.JSON(400, echo.Map{"message": "error validation", "errors": validationErrors})
	}

	_, err := c.service.SaveJabatan(jabatanDto)
	if err != nil {
		return ctx.JSON(400, echo.Map{"message": "error save data user"})
	}

	session.SetFlashMessage(ctx, "save data success", "success", nil)
	return ctx.Redirect(302, "/check/admin/jabatan")
}

func (c *JabatanController) Edit(ctx echo.Context) error {
	id := ctx.Param("id")
	data, err := c.service.FindUserById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			session.SetFlashMessage(ctx, err.Error(), "error", nil)
			return ctx.Redirect(302, "/check/admin/jabatan")
		}
		session.SetFlashMessage(ctx, err.Error(), "error", nil)
		return ctx.Redirect(302, "/check/admin/jabatan")
	}
	breadCrumbs := map[string]interface{}{
		"menu": "Edit",
		"link": "/check/admin/jabatan/edit",
	}

	dataJabatan := models.Jabatan{
		ID:          data.ID,
		Name:        data.Name,
		Description: data.Description,
	}
	return Render(ctx, "Edit Jabatan", "jabatan/edit", c.Menu, session.GetFlashMessage(ctx),
		append(c.BreadCrumbs, breadCrumbs), dataJabatan)
}

func (c *JabatanController) Delete(ctx echo.Context) error {
	id := ctx.Param("id")
	err := c.service.DeleteJabatan(id)
	if err != nil {
		return ctx.JSON(500, echo.Map{"message": "error when trying delete data"})
	}
	return ctx.JSON(200, echo.Map{"message": "delete data has been deleted"})
}

func (c *JabatanController) Update(ctx echo.Context) error {
	id := ctx.Param("id")
	var jabatanDto dto.JabatanDto
	if err := ctx.Bind(&jabatanDto); err != nil {
		return ctx.JSON(400, echo.Map{"message": "error binding data"})
	}

	if err := ctx.Validate(&jabatanDto); err != nil {
		var validationErrors []models.ValidationError
		if errs, ok := err.(validator.ValidationErrors); ok {
			validationErrors = models.WrapValidationErrors(errs)
		}
		return ctx.JSON(400, echo.Map{"message": "error validation", "errors": validationErrors})
	}
	_, err := c.service.UpdateJabatan(id, jabatanDto)
	if err != nil {
		return ctx.JSON(400, echo.Map{"message": "error update data user"})
	}
	session.SetFlashMessage(ctx, "update data success", "success", nil)
	return ctx.Redirect(302, "/check/admin/jabatan")
}

func (c *JabatanController) View(ctx echo.Context) error {
	id := ctx.Param("id")
	var data models.Jabatan
	err := c.service.GetDbInstance().First(&data, "id =?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			session.SetFlashMessage(ctx, err.Error(), "error", nil)
			return ctx.Redirect(302, "/check/admin/jabatan")
		}
		session.SetFlashMessage(ctx, err.Error(), "error", nil)
		return ctx.Redirect(302, "/check/admin/jabatan")
	}
	breadCrumbs := map[string]interface{}{
		"menu": "Detail Jabatan",
		"link": "/check/admin/jabatan/detail",
	}
	return Render(ctx, "Detail Jabatan ", "jabatan/view", c.Menu, session.FlashMessage{},
		append(c.BreadCrumbs, breadCrumbs), echo.Map{"Jabatan": data})
}
