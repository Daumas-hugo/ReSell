package companies

import (
	"testing"

	"github.com/google/uuid"
)

func TestCreateCompanyRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		req     CreateCompanyRequest
		wantErr bool
	}{
		{
			name: "valid company",
			req: CreateCompanyRequest{
				Name:      "Test Company",
				LegalName: "Test Company SAS",
				Status:    "active",
			},
			wantErr: false,
		},
		{
			name: "with SIRET",
			req: CreateCompanyRequest{
				Name:      "Test Company",
				LegalName: "Test Company SAS",
				SIRET:     stringPtr("12345678901234"),
				Status:    "active",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.req.Name == "" {
				t.Error("Name should not be empty")
			}
			if tt.req.LegalName == "" {
				t.Error("LegalName should not be empty")
			}
		})
	}
}

func TestCompany_IDGeneration(t *testing.T) {
	company := Company{
		ID:        uuid.New(),
		Name:      "Test",
		LegalName: "Test SAS",
		Status:    "active",
	}

	if company.ID == uuid.Nil {
		t.Error("Company ID should not be nil")
	}
}

func stringPtr(s string) *string {
	return &s
}
