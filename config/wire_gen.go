// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package config

import (
	"go-checkin/controllers"
	"go-checkin/repository"
	"go-checkin/service"
	"gorm.io/gorm"
)

// Injectors from di.go:

func InjectUserController(db *gorm.DB) controllers.UserController {
	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	userController := controllers.NewUserController(userService)
	return userController
}

func InjectAuthController(db *gorm.DB) controllers.AuthController {
	userRepository := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepository)
	authController := controllers.NewAuthController(authService)
	return authController
}

func InjectMenuController(db *gorm.DB) controllers.MenuController {
	menuRepository := repository.NewMenuRepository(db)
	menuService := service.NewMenuService(menuRepository)
	menuController := controllers.NewMenuController(menuService)
	return menuController
}

func InjectConfigController(db *gorm.DB) controllers.ConfigController {
	repository := repository.NewConfigRepository(db)
	service := service.NewConfigService(repository)
	controller := controllers.NewConfigController(service)
	return controller
}

func InjectRoleController(db *gorm.DB) controllers.RoleController {
	repository := repository.NewRoleRepository(db)
	service := service.NewRoleService(repository)
	controller := controllers.NewRoleController(service)
	return controller
}
func InjectDevisiController(db *gorm.DB) controllers.DevisiController {
	repository := repository.NewDevisiRepository(db)
	service := service.NewDevisiService(repository)
	controller := controllers.NewDevisiController(service)
	return controller
}
func InjectJabatanController(db *gorm.DB) controllers.JabatanController {
	repository := repository.NewJabatanRepository(db)
	service := service.NewJabatanService(repository)
	controller := controllers.NewJabatanController(service)
	return controller
}
