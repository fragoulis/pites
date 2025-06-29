export type MemberPaymentStatus = {
	total_paid: number;
	owed: number;
	ok: boolean;
	formatted: string;
	cutoff_date_formatted: string;
	registered_at_formatted: string;
	last_payment_until_formatted: string;
	is_payment_disabled: boolean;
};

export type User = {
	id: string;
	username: string;
};

export type Payment = {
	id: string;
	member_id: string;
	member_no?: string;
	member_name?: string;
	amount: number;
	months: number;
	receipt_block_no: number;
	receipt_no: number;
	receipt_id: string;
	created_by_user: User;
	issued_at: Date;
	issued_at_formatted: string;
	comments?: string;
	legacy_to?: Date;
	legacy_to_formatted?: string;
};

export type CreatePaymentForm = {
	id?: string;
	member_id?: string;
	amount?: number;
	months?: number;
	contains_registration_fee?: boolean;
	receipt_block_no?: number;
	receipt_no?: number;
	issued_at?: string;
	comments?: string;
	name?: string;
};

export type UpdatePaymentForm = {
	amount?: number;
	receipt_block_no?: number;
	receipt_no?: number;
	comments?: string;
	without_receipt?: boolean;
};

export type PaymentDetails = Payment & {
	status?: string;
};

export type Subscription = {
	id: string;
	member_id: string;
	active: boolean;
	fee_paid: boolean;
	start_date: Date | string;
	end_date: Date | string;
	start_date_formatted: string;
	end_date_formatted: string;
	status_formatted: string;
	months: number;
};

export type Employment = {
	id: string;
	member_id: string;
	company_id: string;
	start_date: Date | string;
	end_date: Date | string;
	start_date_formatted: string;
	end_date_formatted: string;
	company: Company;
};

export type CompanyBusinessType = {
	id: string;
	name: string;
};

export type Company = {
	id?: string;
	name?: string;
	email?: string;
	phone?: string;
	website?: string;
	branch?: string;
	name_formatted?: string;
	address_formatted?: string;
	address_street_id?: string;
	address_street_no?: string;
	business_type_id?: string;
	business_type?: CompanyBusinessType;
	comments?: string;
};

export type CompanyForm = {
	name?: string;
	business_type_id?: string;
	email?: string;
	phone?: string;
	website?: string;
	comments?: string;
	address_street_id?: string;
	address_street_no?: string;
};

export type Issue = {
	id: string;
	key: string;
	description: string;
	importance: string;
	member_id: string;
	company_id: string;
	resolved_at: string;
	created_at: string;
};

export type Member = {
	id: string;
	member_no: string;
	first_name: string;
	last_name: string;
	father_name: string;
	email: string;
	mobile: string;
	phone: string;
	birthdate: string;
	birthdate_formatted: string;
	id_card_number: string;
	social_security_num: string;
	other_union: boolean;
	name_formatted: string;
	address_formatted: string;
	address_street_name: string;
	address_city_name: string;
	company_formatted: string;
	company_name: string;
	company_branch_name: string;
	company_id: string;
	company_address: string;
	subscription_formatted: string;
	subscription_active: boolean;
	unemployed: boolean;
	payment_status: MemberPaymentStatus;
	payments: Payment[];
	issues: Issue[];
	subscriptions: Subscription[];
	employments: Employment[];
	address_street_id: string;
	address_city_id: string;
	address_street_no: string;
	comments: string;
	business_type_name: string;
	legacy_address: string;
	legacy_area: string;
	legacy_city: string;
	legacy_post_code: string;
	specialty: string;
	education: string;
};

export type Address = {
	id: string;
	name: string;
};

export type AssemblyForm = {
	date: string;
	comments: string;
};

export type UpdateMemberForm = {
	first_name?: string;
	last_name?: string;
	father_name?: string;
	email?: string;
	mobile?: string;
	phone?: string;
	birthdate?: string;
	id_card_number?: string;
	social_security_num?: string;
	other_union?: boolean;
	comments?: string;
	address_street_id?: string;
	address_city_id?: string;
	address_street_no?: string;
	company_id?: string;
	legacy_address?: string;
	legacy_area?: string;
	legacy_city?: string;
	legacy_post_code?: string;
	education?: string;
	specialty?: string;
	fixed_payment?: boolean;
};

export type MemberSubscriptionForm = {
	start_date?: string;
	fee_paid?: boolean;
};

export type NewMemberForm = {
	member_no?: string;
} & UpdateMemberForm &
	MemberSubscriptionForm;

export type MethodOptions = 'POST' | 'PATCH' | 'PUT';
export type DatatableSearchForm = {
	q?: string;
	active_only?: boolean;
	company_id?: string;
	address_city_ids?: string[];
	legacy_area?: string;
	business_type_ids?: string[];
	with_comments?: boolean;
	chapter_id?: string;
	with_fixed_monthly_payment?: boolean;
};

export type DatatableColumns = {
	[key: string]: string | object;
};

export type Chapter = {
	id?: string;
	name?: string;
	raw_city_query?: string;
	city_ids?: string[];
};

export type ImportMembersForm = {
	file?: string;
};
