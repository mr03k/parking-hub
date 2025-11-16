package contractorentity

import "github.com/google/uuid"

type Contractor struct {
	UUID                    uuid.UUID `json:"uuid"`
	Name                    string    `json:"contractor_name"`
	Code                    string    `json:"contractor_code"`
	RegisterNumber          string    `json:"register_number"`
	ContactPerson           string    `json:"contact_person"`
	CeoName                 string    `json:"ceo_name"`
	AutorizationSignatories string    `json:"autorization_signatories"`
	PhoneNumbers            string    `json:"phone_numbers"`
	Email                   string    `json:"email"`
	Address                 string    `json:"address"`
	ContractType            string    `json:"contract_type"`
	BankAccountNumber       string    `json:"bank_account_number"`
	Description             string    `json:"description"`
}
