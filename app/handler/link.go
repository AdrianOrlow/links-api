package handler

import (
	"encoding/json"
	"errors"
	"github.com/AdrianOrlow/links-api/app/model"
	"github.com/AdrianOrlow/links-api/app/utils"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"net/http"
)

func HandleRedirect(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var link *model.Link
	vars := mux.Vars(r)

	hashIdOrLink := vars["hashIdOrLink"]
	id, err := utils.DecodeId(hashIdOrLink)

	if err != nil {
		link, err = getLinkByPermalink(hashIdOrLink, db)
	} else {
		link, err = getLinkById(uint(id), db)
	}

	if err != nil {
		respondError(w, http.StatusNotFound, err)
		return
	}

	if !link.IsValid() {
		err = errors.New("link expired")
		respondError(w, http.StatusUnauthorized, err)
		return
	}

	err = increaseLinksVisitsByOne(link, db)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	http.Redirect(w, r, link.Url, http.StatusMovedPermanently)
}

func GetAllLinks(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	links, err := getAllLinks(db)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	for i := 0; i < len(links); i++ {
		links[i].WithHashId()
	}
	respondJSON(w, http.StatusOK, links)
}

func CreateLink(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var link model.Link

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&link)
	if err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}
	defer r.Body.Close()

	err = db.Save(&link).Error
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	respondJSON(w, http.StatusCreated, link.WithHashId())
}

func getAllLinks (db *gorm.DB) ([]model.Link, error) {
	var links []model.Link
	err := db.
		Order("visits DESC, valid_until DESC").
		Find(&links).
		Error
	if err != nil {
		return nil, err
	}

	return links, nil
}

func getLinkById (id uint, db *gorm.DB) (*model.Link, error) {
	var link model.Link
	err := db.
		First(&link, id).
		Error
	if err != nil {
		return nil, err
	}

	return &link, nil
}

func getLinkByPermalink (permalink string, db *gorm.DB) (*model.Link, error) {
	var link model.Link
	err := db.
		First(&link, model.Link{Permalink: permalink}).
		Error
	if err != nil {
		return nil, err
	}

	return &link, nil
}

func increaseLinksVisitsByOne (link *model.Link, db *gorm.DB) error {
	err := db.
		Model(&link).
		Update("visits", gorm.Expr("visits + ?", 1)).
		Error

	return err
}