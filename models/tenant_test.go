package models

import (
	"testing"

	"gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { check.TestingT(t) }

type ModelsSuite struct{}

var _ = check.Suite(&ModelsSuite{})

func (s *ModelsSuite) SetUpSuite(c *check.C) {
	// Any setup needed before tests are run, such as connecting to a test DB
}

// TestTenantPost tests the creation of a new tenant
func (s *ModelsSuite) TestTenantPost(c *check.C) {
	t := Tenant{
		TenantName:       "Test Tenant",
		TenantIdentifier: "test-tenant-identifier",
	}

	// Create a new tenant
	err := PostTenant(&t)
	c.Assert(err, check.Equals, nil)
	c.Assert(t.ID > 0, check.Equals, true) // Verify tenant was created and assigned an ID

	// Verify tenant was saved correctly
	storedTenant, err := GetTenant(t.ID)
	c.Assert(err, check.Equals, nil)
	c.Assert(storedTenant.TenantName, check.Equals, "Test Tenant")
	c.Assert(storedTenant.TenantIdentifier, check.Equals, "test-tenant-identifier")
}

// TestTenantUniqueIdentifier tests that tenant identifiers must be unique
func (s *ModelsSuite) TestTenantUniqueIdentifier(c *check.C) {
	t1 := Tenant{
		TenantName:       "Tenant One",
		TenantIdentifier: "unique-identifier",
	}

	// Create the first tenant
	err := PostTenant(&t1)
	c.Assert(err, check.Equals, nil)

	t2 := Tenant{
		TenantName:       "Tenant Two",
		TenantIdentifier: "unique-identifier",
	}

	// Attempt to create a second tenant with the same identifier
	err = PostTenant(&t2)
	c.Assert(err, check.Equals, ErrTenantNameOrIdentifierConflict)
}

// TestTenantValidation tests validation rules
func (s *ModelsSuite) TestTenantValidation(c *check.C) {
	t := Tenant{
		TenantName: "", // Missing name
	}

	err := t.Validate()
	c.Assert(err, check.Equals, ErrTenantNameNotSpecified)

	t.TenantName = "Valid Name"
	t.TenantIdentifier = "" // Missing identifier
	err = t.Validate()
	c.Assert(err, check.Equals, ErrTenantIdentifierNotSpecified)
}

// TestTenantGet tests fetching a tenant by ID
func (s *ModelsSuite) TestTenantGet(c *check.C) {
	t := Tenant{
		TenantName:       "Fetch Test Tenant",
		TenantIdentifier: "fetch-test-tenant",
	}

	err := PostTenant(&t)
	c.Assert(err, check.Equals, nil)

	// Fetch the tenant by ID
	fetchedTenant, err := GetTenant(t.ID)
	c.Assert(err, check.Equals, nil)
	c.Assert(fetchedTenant.TenantName, check.Equals, "Fetch Test Tenant")
	c.Assert(fetchedTenant.TenantIdentifier, check.Equals, "fetch-test-tenant")
}

// TestTenantUpdate tests updating a tenant's details
func (s *ModelsSuite) TestTenantUpdate(c *check.C) {
	t := Tenant{
		TenantName:       "Update Test Tenant",
		TenantIdentifier: "update-test-tenant",
	}

	err := PostTenant(&t)
	c.Assert(err, check.Equals, nil)

	// Update tenant's name
	t.TenantName = "Updated Tenant Name"
	err = PutTenant(&t)
	c.Assert(err, check.Equals, nil)

	// Verify the update
	updatedTenant, err := GetTenant(t.ID)
	c.Assert(err, check.Equals, nil)
	c.Assert(updatedTenant.TenantName, check.Equals, "Updated Tenant Name")
}

// TestTenantDelete tests deleting a tenant
func (s *ModelsSuite) TestTenantDelete(c *check.C) {
	t := Tenant{
		TenantName:       "Delete Test Tenant",
		TenantIdentifier: "delete-test-tenant",
	}

	err := PostTenant(&t)
	c.Assert(err, check.Equals, nil)

	// Delete the tenant
	err = DeleteTenant(t.ID)
	c.Assert(err, check.Equals, nil)

	// Verify tenant was deleted
	_, err = GetTenant(t.ID)
	c.Assert(err, check.NotNil)
}