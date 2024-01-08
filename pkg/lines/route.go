package lines

import (
	"strings"

	"github.com/gin-gonic/gin"

	"github.startlite.cn/itapp/startlite/pkg/lines/appx"
)

type (
	Route struct {
		Name              string
		Method            string
		Path              string
		Handler           EntryPoint
		CustomHandlerFunc gin.HandlerFunc
		LoginRequired     bool
		RolesRequired     []string
		FeaturesRequired  []string
		Group             string
	}
)

func (route Route) MakeRelativePath() string {
	path := route.Path
	if route.Group != "" {
		path = strings.Join([]string{route.Group, route.Path}, "/")
	}
	return path
}

type methodPath struct {
	Method string
	Path   string
}

func (route Route) makeUniqueMethodPath() methodPath {
	return methodPath{
		Method: route.Method,
		Path:   route.MakeRelativePath(),
	}
}

func CombinedRoutes(routesArr ...[]Route) []Route {
	var combinedRoutes []Route
	var uniqueMethodPathRouteMap = make(map[methodPath]Route)
	for i := 0; i < len(routesArr); i++ {
		for _, v := range routesArr[i] {
			// remove the old one when duplicated
			uniqueMethodPathRouteMap[v.makeUniqueMethodPath()] = v
		}
	}
	for _, v := range uniqueMethodPathRouteMap {
		combinedRoutes = append(combinedRoutes, v)
	}

	return combinedRoutes
}

func ExtendRoutes(appCtx appx.AppContext, routes []Route) []Route {
	// combine routes together
	//var routesInCtx []Route
	// if err == nil {
	// 	routes = CombinedRoutes(routesInCtx, routes)
	// }

	appCtx.Provide(&routes)
	return routes
}
