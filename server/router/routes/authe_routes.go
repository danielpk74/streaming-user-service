package routes

import (
	"api/controllers"

	"net/http"
)

var autheRoutes = []Route{
	Route{
		Uri:     "/login",
		Method:  http.MethodPost,
		Handler: controllers.Login,
	},
}
