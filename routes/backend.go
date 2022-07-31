package routes

import (
	"errors"
	"github.com/foolin/goview/supports/echoview-v4"
	"go-checkin/config"
	"go-checkin/repository"
	"go-checkin/service"
	"html/template"
	"strconv"
	"time"

	"github.com/foolin/goview"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"go-checkin/controllers"
	"go-checkin/middleware"
	"go-checkin/models"
	"go-checkin/utils"
	"go-checkin/utils/session"
	"gorm.io/gorm"
)

func BackendRoute(e *echo.Echo, db *gorm.DB) {

	//=========== Backend ===========//
	var userInfo session.UserInfo
	//new middleware
	mv := echoview.NewMiddleware(goview.Config{
		Root:      "views",
		Extension: ".html",
		Master:    "layouts/master",
		Partials: []string{
			"partials/asside",
			"partials/brand-menu",
			"partials/chart",
			"partials/chat",
			"partials/demo-panel",
			"partials/header-mobile",
			"partials/language",
			"partials/notification",
			"partials/quick-action",
			"partials/quick-panel",
			"partials/quick-panel-toogle",
			"partials/scrolltop",
			"partials/search",
			"partials/sticky-toolbar",
			"partials/sub-header",
			"partials/user-bar",
		},
		Funcs: template.FuncMap{
			"user": func(ctx echo.Context) models.User {
				userModel, err := utils.GetUserInfoFromContext(ctx, db)
				if err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						return models.User{}
					}
					return models.User{}
				}
				return userModel
			},
			"floatToString": func(value *float64) string {
				if value == nil {
					return ""
				}
				return strconv.FormatFloat(*value, 'f', -1, 64)
			},
			"floatNotPointerToString": func(value float64) string {
				return strconv.FormatFloat(value, 'f', -1, 64)
			},
			"formatDate": func(date time.Time, layout string) string {
				return date.Format(layout)
			},
			"maritalStatus": func(data string) string {
				var result string
				if data == "S" {
					result = "Belum Pernah Menikah"
				}
				return result
			},
			"MenuParent": func(ctx echo.Context) []map[string]interface{} {
				var dataMenu map[string]interface{}
				var listOfMenu []map[string]interface{}
				result, err := session.Manager.Get(ctx, session.SessionId)
				if err != nil {
					panic(err)
				}
				userInfo = result.(session.UserInfo)
				//var menu models.Menu
				//var menuRole []models.MenuRole
				//if err := db.Raw("select * from web_menu_role where user_role_id = ? and active = ?", userInfo.UserRoleID, 1 ).
				//Scan(&menuRole); err != nil {
				//}

				var listParentMenu []models.Menu
				if err := db.Raw("select * from web_menu where is_active = ? and menu_type = ? ",
					1, "parent_menu").Scan(&listParentMenu).Error; err != nil {
				}

				for _, v := range listParentMenu {
					var menus []models.Menu
					if err := db.Raw("select * from web_menu where parent_id = ? and is_active = ? ",
						v.ID, 1).Scan(&menus).Error; err != nil {
						log.Info("ERROR GET MENU BY ROLE ", err.Error())
					}
					dataMenu = map[string]interface{}{
						"Name":  v.Name,
						"Icon":  v.IconClass,
						"Menus": menus,
					}
					listOfMenu = append(listOfMenu, dataMenu)
				}

				return listOfMenu
			},
			"Menu": func(ctx echo.Context) []map[string]interface{} {
				var dataMenu map[string]interface{}
				var listOfMenu []map[string]interface{}
				result, err := session.Manager.Get(ctx, session.SessionId)
				if err != nil {
					panic(err)
				}
				userInfo = result.(session.UserInfo)

				var menus []models.Menu
				if err := db.Raw("select * from web_menu where parent_id = ? and is_active = ? and menu_type = ? ",
					0, 1, "menu").Scan(&menus).Error; err != nil {
					log.Info("ERROR GET MENU BY ROLE ", err.Error())
				}
				dataMenu = map[string]interface{}{
					"Menus": menus,
				}
				listOfMenu = append(listOfMenu, dataMenu)
				return listOfMenu
			},
		},
		DisableCache: true,
	})

	// You should use helper func `Middleware()` to set the supplied
	// TemplateEngine and make `HTML()` work validly.
	bGroup := e.Group("/check")
	backendGroup := bGroup.Group("/admin", mv, middleware.SessionMiddleware(session.Manager))
	authorizationMiddleware := middleware.NewAuthorizationMiddleware(db)

	var menus []models.Menu
	if err := db.Raw("select * from web_menu where is_active = ? ",
		1).Scan(&menus).Error; err != nil {
		log.Info("ERROR GET MENU BY ROLE ", err.Error())
	}

	Repository := repository.NewHomeRepository(db)
	Service := service.NewHomeService(Repository)
	homeController := controllers.NewHomeController(Service)
	backendGroup.GET("/home", homeController.Index)
	backendGroup.GET("/home/datatable", homeController.List)
	bGroup.POST("/export/excel", homeController.DownloadExcel)
	bGroup.GET("/download", homeController.ExportExcel)

	authController := config.InjectAuthController(db)
	backendGroup.POST("/logout", authController.Logout)

	//userController
	userController := config.InjectUserController(db)
	userGroup := backendGroup.Group("/register", authorizationMiddleware.AuthorizationMiddleware(menus, "user"))
	userGroup.GET("", userController.Index)
	userGroup.GET("/all", userController.GetAllUser)
	userGroup.POST("/store", userController.Store)
	userGroup.GET("/add", userController.Add)
	userGroup.GET("/profile", userController.Profile)
	userGroup.GET("/datatable", userController.List)
	userGroup.DELETE("/delete/:id", userController.Delete)
	userGroup.GET("/detail/:id", userController.View)
	userGroup.GET("/edit/:id", userController.Edit)
	userGroup.POST("/update/:id", userController.Update)

	//ConfigController
	configController := config.InjectConfigController(db)
	configGroup := backendGroup.Group("/config", authorizationMiddleware.AuthorizationMiddleware(menus, "config"))
	configGroup.GET("", configController.Index)
	configGroup.POST("/store", configController.Store)
	configGroup.GET("/datatable", configController.Datatable)
	configGroup.POST("/update/:id", configController.Update)
	configGroup.GET("/:id", configController.Detail)
	configGroup.POST("/delete/:id", configController.Delete)
	bGroup.POST("/admin/config/set-active/:id", configController.SetActive)
	bGroup.POST("/admin/config/set-inactive/:id", configController.SetInactive)

	// RoleController
	roleController := config.InjectRoleController(db)
	roleGroup := backendGroup.Group("/role", authorizationMiddleware.AuthorizationMiddleware(menus, "role"))
	roleGroup.GET("", roleController.Index)
	roleGroup.GET("/datatable", roleController.List)
	roleGroup.GET("/add", roleController.Add)
	roleGroup.POST("/store", roleController.Store)
	roleGroup.GET("/edit/:id", roleController.Edit)
	roleGroup.POST("/update/:id", roleController.Update)
	roleGroup.GET("/detail/:id", roleController.View)
	roleGroup.DELETE("/delete/:id", roleController.Delete)

	// divisiController
	divisiController := config.InjectDivisiController(db)
	divisiGroup := backendGroup.Group("/divisi", authorizationMiddleware.AuthorizationMiddleware(menus, "divisi"))
	divisiGroup.GET("", divisiController.Index)
	divisiGroup.GET("/datatable", divisiController.List)
	divisiGroup.GET("/add", divisiController.Add)
	divisiGroup.POST("/store", divisiController.Store)
	divisiGroup.GET("/edit/:id", divisiController.Edit)
	divisiGroup.POST("/update/:id", divisiController.Update)
	divisiGroup.GET("/detail/:id", divisiController.View)
	divisiGroup.DELETE("/delete/:id", divisiController.Delete)

	// JabatanController
	jabatanController := config.InjectJabatanController(db)
	jabatanGroup := backendGroup.Group("/jabatan", authorizationMiddleware.AuthorizationMiddleware(menus, "jabatan"))
	jabatanGroup.GET("", jabatanController.Index)
	jabatanGroup.GET("/datatable", jabatanController.List)
	jabatanGroup.GET("/add", jabatanController.Add)
	jabatanGroup.POST("/store", jabatanController.Store)
	jabatanGroup.GET("/edit/:id", jabatanController.Edit)
	jabatanGroup.POST("/update/:id", jabatanController.Update)
	jabatanGroup.GET("/detail/:id", jabatanController.View)
	jabatanGroup.DELETE("/delete/:id", jabatanController.Delete)
}
