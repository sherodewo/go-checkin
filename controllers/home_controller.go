package controllers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/xuri/excelize/v2"
	"go-checkin/dto"
	"go-checkin/models"
	"go-checkin/service"
	"go-checkin/utils/session"
	"net/http"
	"os"
	"strconv"
)

var downloadName string
var count = 1

type HomeController struct {
	BaseBackendController
	service *service.HomeService
}

func NewHomeController(service *service.HomeService) HomeController {
	return HomeController{
		BaseBackendController: BaseBackendController{
			Menu:        "Home",
			BreadCrumbs: []map[string]interface{}{},
		},
		service: service,
	}
}

func (c *HomeController) Index(ctx echo.Context) error {
	breadCrumbs := map[string]interface{}{
		"menu": "Home",
		"link": "/check/admin/home",
	}
	userInfo, _ := session.Manager.Get(ctx, session.SessionId)
	return Render(ctx, "Home", "index", c.Menu, session.GetFlashMessage(ctx),
		append(c.BreadCrumbs, breadCrumbs), userInfo)
}

func (c *HomeController) List(ctx echo.Context) error {

	draw, err := strconv.Atoi(ctx.Request().URL.Query().Get("draw"))
	search := ctx.Request().URL.Query().Get("search[value]")
	start, err := strconv.Atoi(ctx.Request().URL.Query().Get("start"))
	length, err := strconv.Atoi(ctx.Request().URL.Query().Get("length"))
	order, err := strconv.Atoi(ctx.Request().URL.Query().Get("order[0][column]"))
	orderName := ctx.Request().URL.Query().Get("columns[" + strconv.Itoa(order) + "][name]")
	//orderAscDesc := ctx.Request().URL.Query().Get("order[0][dir]")
	var req dto.Excel
	fmt.Println("REQUEST : ", req)
	recordTotal, recordFiltered, data, err := c.service.QueryDatatable(search, "desc", orderName, length, start)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	listOfData := make([]map[string]interface{}, len(data))
	for k, v := range data {
		in := v.Checkin.Format("Monday ,02 January")
		inTime := v.Checkin.Format(" 15:04 ")
		out := v.Checkout.Format("Monday ,02 January")
		outTime := v.Checkout.Format(" 15:04 ")

		if v.Checkout.After(v.Checkin) {
			listOfData[k] = map[string]interface{}{
				"name":      v.Name,
				"location":  v.LocationName,
				"check_in":  in + ", " + inTime,
				"check_out": out + ", " + outTime,
			}
		} else {
			listOfData[k] = map[string]interface{}{
				"name":      v.Name,
				"location":  v.LocationName,
				"check_in":  in + ", " + inTime,
				"check_out": "WORKING",
			}
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

func (c *HomeController) DownloadExcel(ctx echo.Context) error {
	var req dto.Excel
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(400, echo.Map{"message": "error binding data"})
	}
	// Get All
	data, err := c.service.GetAllPresence(req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	f := excelize.NewFile()
	_, _ = f.NewConditionalStyle("center")
	// Set value of a cell.
	_ = f.SetCellValue("Sheet1", "A1", "Name")
	_ = f.SetCellValue("Sheet1", "B1", "Check IN Time")
	_ = f.SetCellValue("Sheet1", "C1", "Check Out Time")
	_ = f.SetCellValue("Sheet1", "D1", "Location Name")
	_ = f.SetCellValue("Sheet1", "E1", "Created At")

	for i, v := range data {
		var checkout string
		if v.Checkout.After(v.Checkin) {
			checkout = v.Checkout.Format("Monday ,02 January")
		} else {
			checkout = "WORKING"
		}
		inTime := v.Checkin.Format(" 15:04 ")
		outTime := v.Checkout.Format(" 15:04 ")
		_ = f.SetCellValue("Sheet1", "A"+strconv.Itoa(i+2), v.Name)
		_ = f.SetCellValue("Sheet1", "B"+strconv.Itoa(i+2), v.Checkin.Format("Monday ,02 January")+", "+inTime)
		_ = f.SetCellValue("Sheet1", "C"+strconv.Itoa(i+2), checkout+", "+outTime)
		_ = f.SetCellValue("Sheet1", "D"+strconv.Itoa(i+2), v.LocationName)
		_ = f.SetCellValue("Sheet1", "E"+strconv.Itoa(i+2), v.CreatedAt.Format("Monday ,02 January"))
	}
	// Save spreadsheet by the given path.
	count++
	name := fmt.Sprintf("Report_%d.xlsx", count)
	downloadName = name
	fmt.Println("NAME DOWNLOAD : ", downloadName)

	if err := f.SaveAs(name); err != nil {
		fmt.Println(err)
	}

	listOfData := make([]map[string]interface{}, len(data))
	dataI := int64(0)
	for k, v := range data {
		dataI++
		in := v.Checkin.Format("Monday ,02 January")
		inTime := v.Checkin.Format(" 15:04 ")
		out := v.Checkout.Format("Monday ,02 January")
		outTime := v.Checkout.Format(" 15:04 ")

		if v.Checkout.After(v.Checkin) {
			listOfData[k] = map[string]interface{}{
				"name":      v.Name,
				"location":  v.LocationName,
				"check_in":  in + ", " + inTime,
				"check_out": out + ", " + outTime,
			}
		} else {
			listOfData[k] = map[string]interface{}{
				"name":      v.Name,
				"location":  v.LocationName,
				"check_in":  in + ", " + inTime,
				"check_out": "WORKING",
			}
		}

	}

	result := models.ResponseDatatable{
		Draw:            10,
		RecordsTotal:    dataI,
		RecordsFiltered: dataI,
		Data:            listOfData,
	}
	return ctx.JSON(http.StatusOK, &result)
}

func (c *HomeController) ExportExcel(ctx echo.Context) error {
	defer os.Remove(downloadName)
	return ctx.File(downloadName)
}
