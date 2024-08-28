package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	"gophish/logger"
	"gophish/models"
)

// Helper function for handling errors
func handleError(w http.ResponseWriter, err error, message string, status int, context string) {
	logger.Errorf("Error in %s: %v", context, err)
	JSONResponse(w, models.Response{Success: false, Message: message}, status)
}

// Tenants handles the functionality for the /api/tenants endpoint.
func (as *Server) Tenants(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// Get all tenants
		ts, err := models.GetTenants()
		if err != nil {
			handleError(w, err, "Error retrieving tenants", http.StatusInternalServerError, "Tenants handler")
			return
		}
		JSONResponse(w, ts, http.StatusOK)

	case http.MethodPost:
		// Create a new tenant
		t := models.Tenant{}
		if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
			handleError(w, err, "Invalid JSON structure", http.StatusBadRequest, "Tenants handler")
			return
		}

		// Validate tenant identifier
		if _, err := models.GetTenantByIdentifier(t.TenantIdentifier); err == nil {
			JSONResponse(w, models.Response{Success: false, Message: "Tenant identifier already in use"}, http.StatusConflict)
			return
		} else if err != gorm.ErrRecordNotFound {
			handleError(w, err, "Error checking tenant identifier", http.StatusInternalServerError, "Tenants handler")
			return
		}

		if err := models.PostTenant(&t); err != nil {
			handleError(w, err, "Error inserting tenant into database", http.StatusInternalServerError, "Tenants handler")
			return
		}
		JSONResponse(w, t, http.StatusCreated)
	}
}

// Tenant handles the functions for the /api/tenants/:id endpoint.
func (as *Server) Tenant(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		handleError(w, err, "Invalid tenant ID", http.StatusBadRequest, "Tenant handler")
		return
	}

	t, err := models.GetTenant(id)
	if err != nil {
		handleError(w, err, "Tenant not found", http.StatusNotFound, "Tenant handler")
		return
	}

	switch r.Method {
	case http.MethodGet:
		JSONResponse(w, t, http.StatusOK)

	case http.MethodDelete:
		if err := models.DeleteTenant(id); err != nil {
			handleError(w, err, "Error deleting tenant", http.StatusInternalServerError, "Tenant handler")
			return
		}
		JSONResponse(w, models.Response{Success: true, Message: "Tenant deleted successfully!"}, http.StatusOK)

	case http.MethodPut:
		// Update tenant details
		tUpdate := models.Tenant{}
		if err := json.NewDecoder(r.Body).Decode(&tUpdate); err != nil {
			handleError(w, err, "Invalid JSON structure", http.StatusBadRequest, "Tenant handler")
			return
		}
		if tUpdate.ID != id {
			JSONResponse(w, models.Response{Success: false, Message: "Error: /:id and tenant_id mismatch"}, http.StatusBadRequest)
			return
		}
		if err := models.PutTenant(&tUpdate); err != nil {
			handleError(w, err, err.Error(), http.StatusBadRequest, "Tenant handler")
			return
		}
		JSONResponse(w, tUpdate, http.StatusOK)
	}
}
