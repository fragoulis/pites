package service

import (
	"fmt"
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/daos"

	"github.com/fragoulis/setip_v2/internal/utils"
)

const (
	subscriptionFee = 2
)

type UpdateMemberRequest struct {
	FirstName         string `json:"first_name"`
	LastName          string `json:"last_name"`
	FatherName        string `json:"father_name"`
	Email             string `json:"email"`
	Mobile            string `json:"mobile"`
	Phone             string `json:"phone"`
	Birthdate         string `json:"birthdate"`
	IDCardNumber      string `json:"id_card_number"`
	SocialSecurityNum string `json:"social_security_num"`
	OtherUnion        bool   `json:"other_union"`
	Comments          string `json:"comments"`
	Education         string `json:"education"`
	AddressStreetID   string `json:"address_street_id"`
	AddressCityID     string `json:"address_city_id"`
	AddressStreetNo   string `json:"address_street_no"`
	CompanyID         string `json:"company_id"`
	Specialty         string `json:"specialty"`
	FixedPayment      bool   `json:"fixed_payment"`

	// Support unstructured addresses for the rare occussions
	// someone lives outside of Attica.
	// Temporarily use the legacy fields.
	LegacyAddress  string `json:"legacy_address"`
	LegacyArea     string `json:"legacy_area"`
	LegacyCity     string `json:"legacy_city"`
	LegacyPostCode string `json:"legacy_post_code"`
}

func (r *UpdateMemberRequest) ToFormData(dao *daos.Dao) (map[string]any, error) {
	firstName := utils.Normalize(r.FirstName)
	lastName := utils.Normalize(r.LastName)
	fullName := fmt.Sprintf("%s %s", lastName, firstName)

	// If street id is present, set the city id from that record.
	// Otherwise keep things as sent.
	if r.AddressStreetID != "" {
		addressStreet, err := dao.FindRecordById("address_streets", r.AddressStreetID)
		if addressStreet == nil || err != nil {
			return nil, fmt.Errorf("failed to find street: %w", err)
		}

		r.AddressCityID = addressStreet.GetString("city_id")
	}

	fixedMonthlyAmount := 0
	if r.FixedPayment {
		fixedMonthlyAmount = subscriptionFee
	}

	return map[string]any{
		"first_name":                    firstName,
		"last_name":                     lastName,
		"full_name":                     fullName,
		"father_name":                   utils.Normalize(r.FatherName),
		"email":                         utils.Normalize(r.Email),
		"mobile":                        r.Mobile,
		"phone":                         r.Phone,
		"birthdate":                     r.Birthdate,
		"id_card_number":                utils.Normalize(r.IDCardNumber),
		"social_security_num":           utils.Normalize(r.SocialSecurityNum),
		"other_union":                   r.OtherUnion,
		"education":                     utils.Normalize(r.Education),
		"comments":                      strconv.Quote(r.Comments),
		"company_id":                    r.CompanyID,
		"specialty":                     utils.Normalize(r.Specialty),
		"address_city_id":               r.AddressCityID,
		"address_street_id":             r.AddressStreetID,
		"address_street_no":             r.AddressStreetNo,
		"legacy_address":                utils.Normalize(r.LegacyAddress),
		"legacy_area":                   utils.Normalize(r.LegacyArea),
		"legacy_city":                   utils.Normalize(r.LegacyCity),
		"legacy_post_code":              utils.Normalize(r.LegacyPostCode),
		"fixed_monthly_amount_in_euros": fixedMonthlyAmount,
	}, nil
}

func (r *UpdateMemberRequest) Validate() validation.Errors {
	errs := validation.Errors{}

	if r.FirstName == "" {
		errs["first_name"] = requiredErr
	}

	if r.LastName == "" {
		errs["last_name"] = requiredErr
	}

	return errs
}
