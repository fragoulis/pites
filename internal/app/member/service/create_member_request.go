package service

import validation "github.com/go-ozzo/ozzo-validation/v4"

// This is here simply to have greek validation error messages.
//
//nolint:gochecknoglobals
var requiredErr = validation.NewError("validation_required", "Το πεδίο είναι υποχρεωτικό.")

type CreateMemberRequest struct {
	MemberNo int `json:"member_no"`
	UpdateMemberRequest
	CreateSubscriptionRequest
}

func (r *CreateMemberRequest) ToUpdateRequest() *UpdateMemberRequest {
	return &UpdateMemberRequest{
		FirstName:         r.FirstName,
		LastName:          r.LastName,
		FatherName:        r.FatherName,
		Email:             r.Email,
		Mobile:            r.Mobile,
		Phone:             r.Phone,
		Birthdate:         r.Birthdate,
		IDCardNumber:      r.IDCardNumber,
		SocialSecurityNum: r.SocialSecurityNum,
		OtherUnion:        r.OtherUnion,
		Comments:          r.Comments,
		Education:         r.Education,
		AddressStreetID:   r.AddressStreetID,
		AddressCityID:     r.AddressCityID,
		AddressStreetNo:   r.AddressStreetNo,
		CompanyID:         r.CompanyID,
		Specialty:         r.Specialty,
		LegacyAddress:     r.LegacyAddress,
		LegacyArea:        r.LegacyArea,
		LegacyCity:        r.LegacyCity,
		LegacyPostCode:    r.LegacyPostCode,
	}
}

func (r *CreateMemberRequest) Validate() validation.Errors {
	errs := validation.Errors{}

	if r.FirstName == "" {
		errs["first_name"] = requiredErr
	}

	if r.LastName == "" {
		errs["last_name"] = requiredErr
	}

	if r.StartDate == "" {
		errs["start_date"] = requiredErr
	}

	return errs
}
