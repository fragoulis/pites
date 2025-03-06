package model

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/pocketbase/pocketbase/models"

	companyModel "github.com/fragoulis/setip_v2/internal/app/company/model"
	"github.com/fragoulis/setip_v2/internal/app/employment"
	issueModel "github.com/fragoulis/setip_v2/internal/app/issue/model"
	paymentModel "github.com/fragoulis/setip_v2/internal/app/payment/model"
	"github.com/fragoulis/setip_v2/internal/utils"
)

const (
	Unemployed                 = "Άνεργη/ος"
	ActiveSubscriptionNotFound = "Ανενεργή/ός"
)

var ErrUnableToDeterminePaymentStatus = errors.New("missing subscription: unable to determine payment status")

type Member struct {
	ID                    string                   `json:"id"`
	MemberNo              string                   `json:"member_no"`
	Email                 string                   `json:"email"`
	Mobile                string                   `json:"mobile"`
	Phone                 string                   `json:"phone"`
	Birthdate             string                   `json:"birthdate"`
	BirthdateFormatted    string                   `json:"birthdate_formatted"`
	IDCardNumber          string                   `json:"id_card_number"`
	SocialSecurityNum     string                   `json:"social_security_num"`
	OtherUnion            bool                     `json:"other_union"`
	NameFormatted         string                   `json:"name_formatted"`
	FirstName             string                   `json:"first_name"`
	LastName              string                   `json:"last_name"`
	FatherName            string                   `json:"father_name"`
	FullName              string                   `json:"full_name"`
	AddressFormatted      string                   `json:"address_formatted"`
	AddressStreetName     string                   `json:"address_street_name"`
	AddressCityName       string                   `json:"address_city_name"`
	AddressPostCode       string                   `json:"address_post_code"`
	CompanyFormatted      string                   `json:"company_formatted"`
	CompanyName           string                   `json:"company_name"`
	CompanyBranchName     string                   `json:"company_branch_name"`
	CompanyID             string                   `json:"company_id"`
	CompanyAddress        string                   `json:"company_address"`
	SubscriptionFormatted string                   `json:"subscription_formatted"`
	SubscriptionActive    bool                     `json:"subscription_active"`
	Employments           []*employment.Employment `json:"employments"`
	Unemployed            bool                     `json:"unemployed"`
	Subscriptions         []*Subscription          `json:"subscriptions"`
	Payments              []*paymentModel.Payment  `json:"payments"`
	PaymentStatus         *PaymentStatus           `json:"payment_status"`
	AddressStreetID       string                   `json:"address_street_id"`
	AddressCityID         string                   `json:"address_city_id"`
	AddressStreetNo       string                   `json:"address_street_no"`
	Comments              string                   `json:"comments"`
	Issues                []*issueModel.Issue      `json:"issues"`
	LegacyAddress         string                   `json:"legacy_address"`
	LegacyArea            string                   `json:"legacy_area"`
	LegacyCity            string                   `json:"legacy_city"`
	LegacyPostCode        string                   `json:"legacy_post_code"`
	BusinessTypeName      string                   `json:"business_type_name"`
	Specialty             string                   `json:"specialty"`
	Education             string                   `json:"education"`
}

func NewFromRecord(
	rec *models.Record,
	employments []*employment.Employment,
	subscriptions []*Subscription,
	payments []*paymentModel.Payment,
	issues []*issueModel.Issue,
) *Member {
	member := &Member{
		ID:                rec.GetId(),
		MemberNo:          fmt.Sprintf("%06d", rec.GetInt("member_no")),
		Email:             rec.GetString("email"),
		Mobile:            rec.GetString("mobile"),
		Phone:             rec.GetString("phone"),
		IDCardNumber:      rec.GetString("id_card_number"),
		SocialSecurityNum: rec.GetString("social_security_num"),
		OtherUnion:        rec.GetBool("other_union"),
		NameFormatted:     newNameFormattedFromRecord(rec).String(),
		FullName:          rec.GetString("full_name"),
		FirstName:         rec.GetString("first_name"),
		LastName:          rec.GetString("last_name"),
		FatherName:        rec.GetString("father_name"),
		AddressFormatted:  newAddressFormattedFromRecord(rec),
		AddressStreetName: newAddressStreetNameFromRecord(rec),
		AddressCityName:   newAddressCityNameFromRecord(rec),
		AddressPostCode:   newAddressPostCodeFromRecord(rec),
		Employments:       employments,
		Subscriptions:     subscriptions,
		Payments:          payments,
		AddressStreetID:   rec.GetString("address_street_id"),
		AddressCityID:     rec.GetString("address_city_id"),
		AddressStreetNo:   rec.GetString("address_street_no"),
		Issues:            issues,
		LegacyAddress:     rec.GetString("legacy_address"),
		LegacyArea:        rec.GetString("legacy_area"),
		LegacyCity:        rec.GetString("legacy_city"),
		LegacyPostCode:    rec.GetString("legacy_post_code"),
		CompanyID:         rec.GetString("company_id"),
		Specialty:         rec.GetString("specialty"),
		Education:         rec.GetString("education"),
	}

	member.setEmploymentStatus(rec)
	member.setPaymentStatus()
	member.setSubscriptionStatus()

	birhtdate := rec.GetDateTime("birthdate").Time()
	member.Birthdate = utils.ForInput(birhtdate)
	member.BirthdateFormatted = utils.Year(birhtdate)

	var err error

	member.Comments, err = strconv.Unquote(rec.GetString("comments"))
	if err != nil {
		member.Comments = err.Error()
	}

	return member
}

func (m *Member) setEmploymentStatus(rec *models.Record) {
	company := companyModel.NewFromRecord(rec.ExpandedOne("company_id"))

	if company == nil {
		m.Unemployed = true
		m.CompanyFormatted = Unemployed

		return
	}

	// We do not want to override these two with the
	// parent's (if present).
	m.CompanyID = company.ID
	m.CompanyAddress = company.AddressFormatted

	var branch *companyModel.Company

	// If parent is present, that means that this company record is
	// a branch, in which case we only want the name, then we update
	// the company pointer to point to the parent.
	if company.Parent != nil {
		branch = company

		company = company.Parent
	}

	m.CompanyName = company.Name
	m.CompanyFormatted = company.Name

	if branch != nil {
		m.CompanyBranchName = branch.Name
		m.CompanyFormatted = fmt.Sprintf("%s | %s", company.Name, branch.Name)

		if branch.BusinessType != nil {
			m.BusinessTypeName = branch.BusinessType.Name
		}
	}

	if company.BusinessType != nil {
		m.BusinessTypeName = company.BusinessType.Name
	}
}

func (m *Member) setPaymentStatus() {
	m.PaymentStatus = newPaymentStatusFromMember(m)
}

func (m *Member) setSubscriptionStatus() {
	subscription := m.activeSubscription()

	if subscription == nil {
		m.SubscriptionFormatted = ActiveSubscriptionNotFound

		return
	}

	m.SubscriptionFormatted = subscription.StatusFormatted
	m.SubscriptionActive = true
}

func (m *Member) activeSubscription() *Subscription {
	for _, sub := range m.Subscriptions {
		if sub.Active {
			return sub
		}
	}

	return nil
}

func (m *Member) ActiveSubscriptionStartedAt() time.Time {
	subscription := m.activeSubscription()

	if subscription == nil {
		return time.Time{}
	}

	return subscription.StartDate
}

func (m *Member) LastPaymentUtil() time.Time {
	if len(m.Payments) == 0 {
		return time.Time{}
	}

	return m.Payments[0].LegacyTo
}

// Temporary method until we get our data right.
// Takes the latest from:
// 1. cutoff date
// 2. active subscription's start date (exclusive)
// 3. last payment legacy_to column (inclusive).
func (m *Member) MemberHasPaidUntil() (time.Time, error) {
	subscription := m.activeSubscription()

	if subscription == nil {
		// If there is no active subscription, the member needs to
		// re-register or we are mising the registration date.
		// In that case, we need to have at least one payment,
		// otherwise we fail.
		if len(m.Payments) == 0 {
			return time.Time{}, ErrUnableToDeterminePaymentStatus
		}
	}

	// Assume the cutoff point is the latest.
	paidUntil := cutoffDate

	// Compare with the subscription's start date (exclusive).
	if subscription != nil && subscription.StartDate.After(paidUntil) {
		paidUntil = utils.EndOfMonthAhead(subscription.StartDate, -1)
	}

	// Compare the latest payment (inclusive).
	for _, payment := range m.Payments {
		lpt := payment.LegacyTo
		if lpt.IsZero() {
			continue
		}

		if lpt.After(paidUntil) {
			paidUntil = lpt
		}
	}

	// Transform to end of month
	return utils.EndOfMonth(paidUntil), nil
}
