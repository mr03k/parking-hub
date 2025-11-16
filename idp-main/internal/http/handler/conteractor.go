package handler

import (
	"log/slog"
	"net/http"

	"application/internal/entity/device"

	"application/internal/http/dto"
	"application/internal/http/response"

	"github.com/google/uuid"
	"github.com/swaggest/openapi-go"
)

var _ Handler = (*ConteractorHandler)(nil)

type ConteractorHandler struct {
	logger *slog.Logger
}

// NewConteractorHandler creates a new ConteractorHandler instance
func NewConteractorHandler(logger *slog.Logger) *ConteractorHandler {
	return &ConteractorHandler{
		logger: logger.With("layer", "ConteractorHandler"),
	}
}

// Get Contractor
func (h *ConteractorHandler) GetContractor(w http.ResponseWriter, r *http.Request) {
	// logger := h.logger.With("method", "GetContractor")
	// ctx := r.Context()

	id := r.PathValue("id")
	// get data form request

	contractor := new(dto.ContractorResponse)

	contractor.ID = id
	contractor.Name = "name"
	contractor.Code = "code"
	contractor.RegisterNumber = "register_number"
	contractor.ContactPerson = "contact_person"
	contractor.CeoName = "ceo_name"
	contractor.AutorizationSignatories = "autorization_signatories"
	contractor.PhoneNumbers = "phone_numbers"
	contractor.Email = "email"
	contractor.Address = "address"
	contractor.ContractType = "contract_type"
	contractor.BankAccountNumber = "bank_account_number"
	contractor.Description = "description"

	response.Pure(w, http.StatusOK, contractor)
}

// Create Contractor
func (h *ConteractorHandler) CreateContractor(w http.ResponseWriter, r *http.Request) {
	// logger := h.logger.With("method", "CreateContractor")
	// ctx := r.Context()

	contractor, err := dto.NewConttactorCreateRequestFromRequest(r)
	if err != nil {
		// logger.Error("error parse request", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	contractorResponse := new(dto.ConttactorCreateResponse)
	contractorResponse.ID = uuid.NewString()
	contractorResponse.ContractorRequest = contractor.ContractorRequest

	response.Pure(w, http.StatusOK, contractorResponse)
}

// Update Contractor
func (h *ConteractorHandler) UpdateContractor(w http.ResponseWriter, r *http.Request) {
	// logger := h.logger.With("method", "UpdateContractor")
	// ctx := r.Context()

	id := r.PathValue("id")
	contractor, err := dto.NewConttactorCreateRequestFromRequest(r)
	if err != nil {
		// logger.Error("error parse request", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	contractorResponse := new(dto.ConttactorCreateResponse)
	contractorResponse.ID = id
	contractorResponse.ContractorRequest = contractor.ContractorRequest

	response.Pure(w, http.StatusOK, contractorResponse)
}

// Delete Contractor
func (h *ConteractorHandler) DeleteContractor(w http.ResponseWriter, r *http.Request) {
	// logger := h.logger.With("method", "DeleteContractor")
	// ctx := r.Context()

	id := r.PathValue("id")

	response.Pure(w, http.StatusOK, id)
}

// List Contractor
// List Contractor
func (h *ConteractorHandler) ListContractor(w http.ResponseWriter, r *http.Request) {
	// Create mock contractor data
	mockContractors := []device.Contractor{
		{
			ID:                    uuid.New(),
			ContractorName:        "Contractor 1",
			CodeContractor:        "C001",
			NumberRegistration:    "RN12345",
			PersonContact:         "John Doe",
			CEOName:               "Alice Smith",
			SignatoriesAuthorized: "Jane Roe",
			PhoneNumber:           "123456789",
			Email:                 "contractor1@example.com",
			Address:               "123 Elm Street",
			TypeContract:          "Type A",
			NumberAccountBank:     "987654321",
			Description:           "First mock contractor",
		},
		{
			ID:                    uuid.New(),
			ContractorName:        "Contractor 2",
			CodeContractor:        "C002",
			NumberRegistration:    "RN67890",
			PersonContact:         "Bob Brown",
			CEOName:               "Charlie Johnson",
			SignatoriesAuthorized: "Diana White",
			PhoneNumber:           "987654321",
			Email:                 "contractor2@example.com",
			Address:               "456 Maple Avenue",
			TypeContract:          "Type B",
			NumberAccountBank:     "123456789",
			Description:           "Second mock contractor",
		},
	}

	// Prepare the response
	contractorsResponse := struct {
		Count       int                 `json:"count"`
		Contractors []device.Contractor `json:"contractors"`
	}{
		Count:       len(mockContractors),
		Contractors: mockContractors,
	}

	// Send the response
	response.Pure(w, http.StatusOK, contractorsResponse)
}

// RegisterMuxRouter registers the handler to the given mux
func (h *ConteractorHandler) RegisterMuxRouter(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/contractor", h.CreateContractor)
	mux.HandleFunc("GET /api/v1/contractor", h.ListContractor)
	mux.HandleFunc("GET /api/v1/contractor/{id}", h.GetContractor)
}

type ContractorID struct {
	ContractorID string `json:"contractor_id" path:"contractor_id"`
}

// Create Vehicle
func (s *ConteractorHandler) CreateContractorOAPI(op openapi.OperationContext) {
	op.SetSummary("Create Vehicle")

	op.AddReqStructure(new(dto.ContractorRequest), openapi.WithContentType("application/json"))
	op.AddRespStructure(new(dto.ContractorResponse), openapi.WithHTTPStatus(http.StatusOK))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusInternalServerError))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusBadRequest))
}

// List Vehicles
func (s *ConteractorHandler) ListContractorOAPI(op openapi.OperationContext) {
	op.SetSummary("List Vehicles")

	op.AddRespStructure(new(dto.ConttactorListResponse), openapi.WithHTTPStatus(http.StatusOK))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusInternalServerError))
}

// Delete Vehicle
func (s *ConteractorHandler) DeleteContractorOAPI(op openapi.OperationContext) {
	op.SetSummary("Delete Vehicle")
	op.AddReqStructure(new(ContractorID), openapi.WithContentType("application/json"))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusOK))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusInternalServerError))
}

// Get Vehicle By ID
func (s *ConteractorHandler) GetContractorOAPI(op openapi.OperationContext) {
	op.SetSummary("Get Vehicle By ID")
	op.AddReqStructure(new(ContractorID), openapi.WithContentType("application/json"))
	op.AddRespStructure(new(dto.VehicleResponse), openapi.WithHTTPStatus(http.StatusOK))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusInternalServerError))
}

func (s *ConteractorHandler) OpenApiSpec(api OAPI) {
	vehicleTags := WithTags("Contractor")
	api.Register("POST", "/api/v1/contractor", s.CreateContractorOAPI, vehicleTags)
	api.Register("GET", "/api/v1/contractor", s.ListContractorOAPI, vehicleTags)
	api.Register("DELETE", "/api/v1/contractor/{contractor_id}", s.DeleteContractorOAPI, vehicleTags)
	api.Register("GET", "/api/v1/contractor/{contractor_id}", s.GetContractorOAPI, vehicleTags)
}
