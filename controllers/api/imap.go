package api

import (
	"encoding/json"
	"net/http"
	"time"
	"strconv"

	ctx "gophish/context"
	"gophish/imap"
	"gophish/models"
	"github.com/gorilla/mux"
)

// IMAPServerValidate handles requests for the /api/imapserver/validate endpoint
func (as *Server) IMAPServerValidate(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == "GET":
		JSONResponse(w, models.Response{Success: false, Message: "Only POSTs allowed"}, http.StatusBadRequest)
	case r.Method == "POST":
		im := models.IMAP{}
		err := json.NewDecoder(r.Body).Decode(&im)
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: "Invalid request"}, http.StatusBadRequest)
			return
		}
		err = imap.Validate(&im)
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusOK)
			return
		}
		JSONResponse(w, models.Response{Success: true, Message: "Successful login."}, http.StatusCreated)
	}
}

// IMAPServer handles requests for the /api/imapserver/ endpoint
func (as *Server) IMAPServer(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == "GET":
		ss, err := models.GetIMAP(ctx.Get(r, "user_id").(int64))
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
			return
		}
		JSONResponse(w, ss, http.StatusOK)

	// POST: Update database
	case r.Method == "POST":
		im := models.IMAP{}
		err := json.NewDecoder(r.Body).Decode(&im)
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: "Invalid data. Please check your IMAP settings."}, http.StatusBadRequest)
			return
		}
		im.ModifiedDate = time.Now().UTC()
		im.UserId = ctx.Get(r, "user_id").(int64)
		err = models.PostIMAP(&im, ctx.Get(r, "user_id").(int64))
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
			return
		}
		JSONResponse(w, models.Response{Success: true, Message: "Successfully saved IMAP settings."}, http.StatusCreated)
	}
}

// ImapsByTenant handles requests for the /imap/{tenant_id:[0-9]+} endpoint
func (as *Server) ImapsByTenant(w http.ResponseWriter, r *http.Request) {
    // Extract tenant_id from the URL path
    vars := mux.Vars(r)
    tenantIDStr := vars["tenant_id"]
    tenantID, err := strconv.ParseInt(tenantIDStr, 10, 64)
    if err != nil {
        JSONResponse(w, models.Response{Success: false, Message: "Invalid tenant ID."}, http.StatusBadRequest)
        return
    }

    switch r.Method {
    case "GET":
        // Retrieve IMAP settings for the given tenant ID
        imaps, err := models.GetIMAPByTenantID(tenantID)
        if err != nil {
            JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
            return
        }
        JSONResponse(w, imaps, http.StatusOK)
    default:
        JSONResponse(w, models.Response{Success: false, Message: "Only GET requests are allowed."}, http.StatusMethodNotAllowed)
    }
}

