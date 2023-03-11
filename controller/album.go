package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/pulse227/server-recruit-challenge-sample/model"
	"github.com/pulse227/server-recruit-challenge-sample/service"
)

type albumController struct {
	service       service.AlbumService
	singerService service.SingerService
}

// アルバムに歌手の情報を取得API
type albumResponse struct {
	AlbumID int           `json:"id"`
	Title   string        `json:"title"`
	Singer  *model.Singer `json:"singer"`
}

func NewAlbumController(s service.AlbumService, ss service.SingerService) *albumController {
	return &albumController{service: s, singerService: ss}
}

// GET /albums のハンドラー
func (c *albumController) GetAlbumListHandler(w http.ResponseWriter, r *http.Request) {
	albums, err := c.service.GetAlbumListService(r.Context())
	if err != nil {
		errorHandler(w, r, 500, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(albums)
}

// GET /albums/{id} のハンドラー
func (c *albumController) GetAlbumDetailHandler(w http.ResponseWriter, r *http.Request) {
	albumID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		err = fmt.Errorf("invalid path param: %w", err)
		errorHandler(w, r, 400, err.Error())
		return
	}

	album, err := c.service.GetAlbumService(r.Context(), model.AlbumID(albumID))
	if err != nil {
		errorHandler(w, r, 500, err.Error())
		return
	}

	//アルバム情報を元にSingerIDを取得し、それを元にGetSingerServiceで情報を取得
	singerID := album.SingerID
	singer, err := c.singerService.GetSingerService(r.Context(), singerID)
	if err != nil {
		errorHandler(w, r, 500, err.Error())
		return
	}

	//album情報とsinger情報をmerge
	response := albumResponse{AlbumID: int(album.ID), Title: album.Title, Singer: singer}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(response)
}

// POST /albums のハンドラー
func (c *albumController) PostAlbumHandler(w http.ResponseWriter, r *http.Request) {
	var album *model.Album
	if err := json.NewDecoder(r.Body).Decode(&album); err != nil {
		err = fmt.Errorf("invalid body param: %w", err)
		errorHandler(w, r, 400, err.Error())
		return
	}

	if err := c.service.PostAlbumService(r.Context(), album); err != nil {
		errorHandler(w, r, 500, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(album)
}

// DELETE /albums/{id} のハンドラー
func (c *albumController) DeleteAlbumHandler(w http.ResponseWriter, r *http.Request) {
	albumID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		err = fmt.Errorf("invalid path param: %w", err)
		errorHandler(w, r, 400, err.Error())
		return
	}

	if err := c.service.DeleteAlbumService(r.Context(), model.AlbumID(albumID)); err != nil {
		errorHandler(w, r, 500, err.Error())
		return
	}
	w.WriteHeader(204)
}
