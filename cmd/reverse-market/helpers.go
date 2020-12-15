package main

import (
	"fmt"
	"github.com/reverse-market/backend/pkg/database/models"
	"io"
	"io/ioutil"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"
)

func (app *Application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.loggers.error.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *Application) clientError(w http.ResponseWriter, err error, status int) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.loggers.info.Output(2, trace)

	http.Error(w, http.StatusText(status), status)
}

func savePhoto(source io.Reader) (string, error) {
	file, err := ioutil.TempFile("images", "*.jpg")
	if err != nil {
		return "", err
	}
	defer file.Close()

	if err := file.Chmod(0777); err != nil {
		return "", err
	}

	if _, err := io.Copy(file, source); err != nil {
		return "", err
	}

	return fmt.Sprintf("/%s", file.Name()), nil
}

func getTagFilters(r *http.Request) *models.TagFilters {
	var category *int
	if categoryParam, err := strconv.Atoi(r.URL.Query().Get("category")); err != nil {
		category = nil
	} else {
		category = &categoryParam
	}

	search := r.URL.Query().Get("search")

	return &models.TagFilters{
		CategoryID: category,
		Search:     strings.ToLower(strings.TrimSpace(search)),
	}
}

func getRequestsFilters(r *http.Request) *models.RequestFilters {
	filters := &models.RequestFilters{
		Page:          0,
		Size:          10,
		Tags:          []int{},
		SortColumn:    "date",
		SortDirection: "desc",
	}

	query := r.URL.Query()

	if page, err := strconv.Atoi(query.Get("page")); err == nil {
		filters.Page = page
	}

	if size, err := strconv.Atoi(query.Get("size")); err == nil {
		filters.Size = size
	}

	if category, err := strconv.Atoi(query.Get("category")); err == nil {
		filters.CategoryID = &category
	}

	if tags, ok := query["tag"]; ok {
		for _, tagStr := range tags {
			if tag, err := strconv.Atoi(tagStr); err == nil {
				filters.Tags = append(filters.Tags, tag)
			}
		}
	}

	if priceFrom, err := strconv.Atoi(query.Get("price_from")); err == nil {
		filters.PriceFrom = &priceFrom
	}

	if priceTo, err := strconv.Atoi(query.Get("price_to")); err == nil {
		filters.PriceTo = &priceTo
	}

	if sort := query.Get("sort"); sort != "" {
		params := strings.Split(strings.TrimSpace(sort), "_")
		if len(params) == 2 {
			if params[0] == "date" || params[0] == "price" {
				filters.SortColumn = params[0]
			}
			if params[1] == "asc" || params[1] == "desc" {
				filters.SortDirection = params[1]
			}
		}
	}

	if search := query.Get("search"); search != "" {
		filters.Search = strings.ToLower(strings.TrimSpace(search))
	}

	return filters
}
