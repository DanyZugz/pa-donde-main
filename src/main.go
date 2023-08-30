package main

import (
	"net/http"
	"pa-donde/src/utils"
	api_v1 "pa-donde/src/api/v1"
)

func main() {
	s := utils.CreateNewServer(utils.TEST_DB_NAME, utils.SQLITE)
	s.PlacesHandler = &api_v1.PlaceHandler{DB: s.DB}
	s.MountHandlers()
	http.ListenAndServe(":3000", s.Router)
}
