package models

import (
	"errors"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" 
)

// Tenant models hold the attributes for a tenant
type Tenant struct {
	ID               int64  `json:"id" gorm:"column:id; primary_key:yes"`
	Guid             string `json:"guid" gorm:"column:guid"`
	TenantName       string `json:"tenant_name" gorm:"column:tenant_name"`
	TenantIdentifier string `json:"tenant_identifier" gorm:"column:tenant_identifier"`
}

// Error definitions
var (
	ErrTenantNameNotSpecified          = errors.New("Tenant name not specified")
	ErrTenantIdentifierNotSpecified    = errors.New("Tenant identifier not specified")
	ErrTenantNameOrIdentifierConflict  = errors.New("Tenant name or identifier already in use")
	ErrTenantNameLengthExceeded        = errors.New("Tenant name exceeds maximum length of 100 characters")
	ErrTenantIdentifierLengthExceeded  = errors.New("Tenant identifier exceeds maximum length of 100 characters")
)

// Validate checks the given tenant to make sure values are appropriate and complete
func (t *Tenant) Validate() error {
	if t.TenantName == "" {
		return ErrTenantNameNotSpecified
	}
	if t.TenantIdentifier == "" {
		return ErrTenantIdentifierNotSpecified
	}
	if len(t.TenantName) > 100 {
		return ErrTenantNameLengthExceeded
	}
	if len(t.TenantIdentifier) > 100 {
		return ErrTenantIdentifierLengthExceeded
	}
	return nil
}

// GetTenants returns all tenants in the database.
func GetTenants() ([]Tenant, error) {
	tenants := []Tenant{}
	err := db.Find(&tenants).Error
	if err != nil {
		log.Println("Error getting tenants:", err)
		return tenants, err
	}
	return tenants, nil
}

// GetTenant returns the tenant specified by the given id.
func GetTenant(id int64) (Tenant, error) {
	tenant := Tenant{}
	err := db.Where("id=?", id).First(&tenant).Error
	if err != nil {
		log.Printf("Error getting tenant by ID (%d): %v", id, err)
		return tenant, err
	}
	return tenant, nil
}

// GetTenantByIdentifier returns the tenant specified by the given identifier.
func GetTenantByIdentifier(identifier string) (Tenant, error) {
	tenant := Tenant{}
	err := db.Where("tenant_identifier=?", identifier).First(&tenant).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Printf("Error getting tenant by identifier (%s): %v", identifier, err)
	}
	return tenant, err
}

// PostTenant creates a new tenant in the database.
func PostTenant(t *Tenant) error {
	if err := t.Validate(); err != nil {
		return err
	}

	// Check if a tenant with the same identifier already exists
	_, err := GetTenantByIdentifier(t.TenantIdentifier)
	if err == nil {
		return ErrTenantNameOrIdentifierConflict
	}
	if err != gorm.ErrRecordNotFound {
		log.Println("Error checking tenant by identifier:", err)
		return err
	}

	// Insert into the DB
	err = db.Create(t).Error
	if err != nil {
		log.Println("Error saving tenant:", err)
		return err
	}
	return nil
}

// PutTenant edits an existing tenant in the database.
func PutTenant(t *Tenant) error {
	if err := t.Validate(); err != nil {
		return err
	}

	// Check if the tenant exists
	existingTenant, err := GetTenant(t.ID)
	if err != nil {
		return err
	}

	// Check if a tenant with the same identifier already exists and is not the current one
	if existingTenant.TenantIdentifier != t.TenantIdentifier {
		_, err = GetTenantByIdentifier(t.TenantIdentifier)
		if err == nil {
			return ErrTenantNameOrIdentifierConflict
		}
		if err != gorm.ErrRecordNotFound {
			log.Println("Error checking tenant by identifier during update:", err)
			return err
		}
	}

	// Save updated tenant
	err = db.Save(t).Error
	if err != nil {
		log.Println("Error updating tenant:", err)
		return err
	}
	return nil
}

// DeleteTenant deletes an existing tenant in the database.
func DeleteTenant(id int64) error {
	err := db.Where("id=?", id).Delete(Tenant{}).Error
	if err != nil {
		log.Printf("Error deleting tenant with ID %d: %v", id, err)
		return err
	}
	return nil
}