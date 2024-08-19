package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gophish/gophish/logger"
	"github.com/gorilla/mux"
	"gophish/models"
)

// Tenants handles the functionality for the /api/tenants endpointapi_key=d42ceb1d38b7f95ed4fabbb357b054c8c1f4cc42849
func (as *Server) Tenants(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// Get all tenants
		ts, err := models.GetTenants()
		if err != nil {
			logger.Error(err)
			JSONResponse(w, models.Response{Success: false, Message: "Error retrieving tenants"}, http.StatusInternalServerError)
			return
		}
		JSONResponse(w, ts, http.StatusOK)

	case http.MethodPost:
		// Create a new tenant
		t := models.Tenant{}
		err := json.NewDecoder(r.Body).Decode(&t)
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: "Invalid JSON structure"}, http.StatusBadRequest)
			return
		}

		// Validate tenant identifier
		_, err = models.GetTenantByIdentifier(t.TenantIdentifier)
		if err == nil {
			JSONResponse(w, models.Response{Success: false, Message: "Tenant identifier already in use"}, http.StatusConflict)
			return
		}
		if err != gorm.ErrRecordNotFound {
			JSONResponse(w, models.Response{Success: false, Message: "Error checking tenant identifier"}, http.StatusInternalServerError)
			logger.Error(err)
			return
		}

		err = models.PostTenant(&t)
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: "Error inserting tenant into database"}, http.StatusInternalServerError)
			logger.Error(err)
			return
		}
		JSONResponse(w, t, http.StatusCreated)
	}
}

// Tenant handles the functions for the /api/tenants/:id endpoint
func (as *Server) Tenant(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		JSONResponse(w, models.Response{Success: false, Message: "Invalid tenant ID"}, http.StatusBadRequest)
		return
	}

	t, err := models.GetTenant(id)
	if err != nil {
		JSONResponse(w, models.Response{Success: false, Message: "Tenant not found"}, http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		JSONResponse(w, t, http.StatusOK)

	case http.MethodDelete:
		err = models.DeleteTenant(id)
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: "Error deleting tenant"}, http.StatusInternalServerError)
			return
		}
		JSONResponse(w, models.Response{Success: true, Message: "Tenant deleted successfully!"}, http.StatusOK)

	case http.MethodPut:
		// Update tenant details
		tUpdate := models.Tenant{}
		err = json.NewDecoder(r.Body).Decode(&tUpdate)
		if err != nil {
			logger.Error(err)
			JSONResponse(w, models.Response{Success: false, Message: "Invalid JSON structure"}, http.StatusBadRequest)
			return
		}
		if tUpdate.ID != id {
			JSONResponse(w, models.Response{Success: false, Message: "Error: /:id and tenant_id mismatch"}, http.StatusBadRequest)
			return
		}
		err = models.PutTenant(&tUpdate)
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusBadRequest)
			return
		}
		JSONResponse(w, tUpdate, http.StatusOK)
	}
}