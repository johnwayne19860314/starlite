package controller

import (
	"net/http"

	"github.startlite.cn/itapp/startlite/pkg/lines"
	"github.startlite.cn/itapp/startlite/pkg/lines/appx"
)

// GetRoutes returns routes
func GetRoutes(appCtx appx.AppContext) []lines.Route {
	return []lines.Route{
		{Method: http.MethodPost, Path: "/v1/user/login", Handler: login},
		{Method: http.MethodPost, Path: "/v1/user", Handler: addUser},
		{Method: http.MethodPut, Path: "/v1/user", Handler: updateUser},
		{Method: http.MethodPost, Path: "/v1/user/list", Handler: listUsers},
		{Method: http.MethodDelete, Path: "/v1/user/:id", Handler: delUser},

		{Method: http.MethodGet, Path: "/v1/entry", Handler: getEntry},
		{Method: http.MethodPost, Path: "/v1/entry", Handler: addEntry},
		{Method: http.MethodPut, Path: "/v1/entry", Handler: updateEntry},
		{Method: http.MethodPost, Path: "/v1/entries/import", Handler: importEntries},
		{Method: http.MethodPost, Path: "/v1/entries/list", Handler: listEntries},
		{Method: http.MethodDelete, Path: "/v1/entry/:code", Handler: delEntry},
		{Method: http.MethodPost, Path: "/v1/entry/category/list", Handler: listEntryCategory},
		{Method: http.MethodPost, Path: "/v1/entry/category", Handler: addEntryCategory},
		{Method: http.MethodPut, Path: "/v1/entry/category", Handler: updateEntryCategory},
		{Method: http.MethodDelete, Path: "/v1/entry/category/:category", Handler: delEntryCategory},
	}
}
