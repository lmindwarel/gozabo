package gobudins

import "time"

const (
	RouteAccessToken   = "/auth/token/access"
	RouteUsers         = "/users"
	RouteAccounts      = "/accounts"
	RouteAuthTokenCode = "/auth/token/code"
)

type ErrorResponse struct {
	Code        string `json:"error"`
	Description string `json:"error_description"`
}

type Time struct {
	time.Time
}

func (mytime *Time) UnmarshalJSON(b []byte) (err error) {
	s := string(b)

	// Get rid of the quotes "" around the value.
	// A second option would be to include them
	// in the date format string instead, like so below:
	//   time.Parse(`"`+time.RFC3339Nano+`"`, s)
	s = s[1 : len(s)-1]

	t, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		t, err = time.Parse("2006-01-02 15:04:05", s)
	}
	mytime.Time = t
	return
}

type Date struct {
	time.Time
}

func (mydate *Date) UnmarshalJSON(b []byte) (err error) {
	s := string(b)

	// Get rid of the quotes "" around the value.
	// A second option would be to include them
	// in the date format string instead, like so below:
	//   time.Parse(`"`+time.RFC3339Nano+`"`, s)
	s = s[1 : len(s)-1]

	t, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		t, err = time.Parse("2006-01-02", s)
	}
	mydate.Time = t
	return
}

type ConnectCallbackData struct {
	Code         string `json:"code"`
	ConnectionID string `json:"connectionID"`
}

type APICredentials struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type AskForToken struct {
	APICredentials
	Code string `json:"code"`
}

type AskForTokenRenew struct {
	APICredentials
	UserID         int  `json:"id_user"`
	RevokePrevious bool `json:"revoke_previous"`
}

type Token struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

type TokenAccessType string

const (
	TokenAccessTypeStandard TokenAccessType = "standard"
)

type TemporaryCode struct {
	Code      string          `json:"code"`
	Type      string          `json:"type"`
	Access    TokenAccessType `json:"access"`
	ExpiresIn int             `json:"expires_in"`
}

const UserMe = "me"

type User struct {
	ID       int    `json:"id"`
	Signin   Time   `json:"signin"`
	Platform string `json:"platform"`
}

type AccountsResponse struct {
	Balance        float64            `json:"balance"`
	Balances       map[string]float64 `json:"balances"`
	ComingBalances map[string]float64 `json:"coming_balances"`
	Accounts       []Account          `json:"accounts"`
}

// Account as described at https://docs.budget-insight.com/reference/bank-accounts#response-bankaccount-object
type Account struct {
	ID           int              `json:"id"`
	ConnectionID *int             `json:"id_connection"`
	UserID       *int             `json:"id_user"`
	SourceID     *int             `json:"id_source"`
	ParentID     *int             `json:"id_parent"`
	Number       *string          `json:"number"`
	OriginalName string           `json:"original_name"`
	Balance      *float64         `json:"balance"`
	Coming       *float64         `json:"comming"`
	Display      bool             `json:"display"`
	LastUpdate   *Time            `json:"last_update"`
	Deleted      *Time            `json:"deleted"`
	Disabled     *Time            `json:"disabled"`
	IBAN         *string          `json:"iban"`
	BIC          *string          `json:"bic"`
	Currency     *Currency        `json:"currency"`
	Type         AccountTypeName  `json:"type"`
	TypeID       int              `json:"id_type"`
	Bookmarked   int              `json:"bookmarked"`
	Name         string           `json:"name"`
	Error        *string          `json:"error"`
	Usage        BankAccountUsage `json:"usage"`
	Ownsership   string           `json:"ownership"`
	CompanyName  *string          `json:"company_name"`
	Loan         *Loan            `json:"loan"`
}

type SyncedAccount struct {
	Account
	// TODO: Transaction Transaction `json:"transactions"`
}

type UpdateAccount struct {
	Display  bool `json:"display"`
	Disabled bool `json:"disabled"`
}

// BankAccountUsage as described at https://docs.budget-insight.com/reference/bank-accounts#bankaccountusage-values
type BankAccountUsage string

const (
	BankAccountUsagePriv BankAccountUsage = "PRIV"
	BankAccountUsageOrga BankAccountUsage = "ORGA"
	BankAccountUsageAsso BankAccountUsage = "ASSO"
)

// Loan as described at https://docs.budget-insight.com/reference/bank-accounts#loan-object
type Loan struct {
	TotalAmount       *float64 `json:"total_amount"`
	AvailableAmount   *float64 `json:"available_amount"`
	UsedAcmount       *float64 `json:"used_amount"`
	SubscriptionDate  *Time    `json:"subscription_date"`
	MaturityDate      *Time    `json:"maturity_date"`
	NextPaymentAmount *float64 `json:"next_payment_amount"`
	NextPatmentAmount *Time    `json:"next_payment_date"`
	Rate              *float64 `json:"rate"`
	NbPaymentsLeft    *int     `json:"nb_payments_left"`
	NbPaymentsDone    *int     `json:"nb_payments_done"`
	NbPaymentsTotal   *int     `json:"nb_payments_total"`
	LastPaymentAmount *float64 `json:"last_payment_amount"`
	LastPaymentDate   *Time    `json:"last_payment_date"`
	AccountLabel      *string  `json:"account_label"`
	InsuranceLabel    *string  `json:"insurance_label"`
	Duration          *int     `json:"duration"`
}

type Currency struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Symbol    string `json:"symbol"`
	Precision int    `json:"precision"`
}

type AccountTypeName string

const (
	AccountTypeNameCheckings      AccountTypeName = "checking"
	AccountTypeNameSavings        AccountTypeName = "savings"
	AccountTypeNameDeposit        AccountTypeName = "deposit"
	AccountTypeNameLoan           AccountTypeName = "loan"
	AccountTypeNameMarket         AccountTypeName = "market"
	AccountTypeNameJoint          AccountTypeName = "joint"
	AccountTypeNameCard           AccountTypeName = "card"
	AccountTypeNameLifeInsurance  AccountTypeName = "lifeinsurance"
	AccountTypeNamePEE            AccountTypeName = "pee"
	AccountTypeNamePERCO          AccountTypeName = "perco"
	AccountTypeNameArticle83      AccountTypeName = "article83"
	AccountTypeNameRSP            AccountTypeName = "rsp"
	AccountTypeNamePEA            AccountTypeName = "pea"
	AccountTypeNameCapitalisation AccountTypeName = "capitalisation"
	AccountTypeNamePERP           AccountTypeName = "perp"
	AccountTypeNameMadelin        AccountTypeName = "madelin"
	AccountTypeNameUnknow         AccountTypeName = "unknown"
)

// AccountType as described at https://docs.budget-insight.com/reference/bank-account-types#response-accounttype-object
type AccountType struct {
	ID           int             `json:"id"`
	Name         AccountTypeName `json:"name"`
	ParentID     *int            `json:"id_parent"`
	IsInvest     bool            `json:"is_invest"`
	DisplayName  string          `json:"display_name"`
	DisplayNameP string          `json:"display_name_p"`
}

type FinanceSecurityType string

const (
	FinanceSecurityTypeOPCVM   FinanceSecurityType = "OPCVM"
	FinanceSecurityTypeETF     FinanceSecurityType = "Trackers - ETF"
	FinanceSecurityTypeActions FinanceSecurityType = "Actions"
)

type InvestmentsResponse struct {
	Diff            float64      `json:"diff"`
	DiffPercent     float64      `json:"diff_percent"`
	PrevDiff        *float64     `json:"prev_diff"`
	PrevDiffPercent *float64     `json:"prev_diff_percent"`
	Valuation       float64      `json:"valuation"`
	Investments     []Investment `json:"investments"`
}

type Investment struct {
	ID                int                  `json:"id"`
	AccountID         int                  `json:"id_account"`
	SecurityID        int                  `json:"id_security"`
	TypeID            *FinanceSecurityType `json:"id_type"`
	Label             string               `json:"label"`
	Code              *string              `json:"code"`
	CodeType          string               `json:"code_type"`
	Source            string               `json:"source"`
	Description       *string              `json:"description"`
	Quantity          float64              `json:"quantity"`
	UnitPrice         float64              `json:"unitprice"`
	UnitValue         float64              `json:"unitvalue"`
	Valuation         float64              `json:"valuation"`
	Diff              float64              `json:"diff"`
	DiffPercent       float64              `json:"diff_percent"`
	PrevDiff          *float64             `json:"prev_diff"`
	PrevDiffPercent   *float64             `json:"prev_diff_percent"`
	VDate             Date                 `json:"vdate"`
	PrevVDate         *Date                `json:"prev_vdate"`
	PortfolioShare    float64              `json:"portfolio_share"`
	Calculated        []string             `json:"calculated"`
	Deleted           *Time                `json:"deleted"`
	LastUpdate        *Time                `json:"last_update"`
	OriginalCurrency  *Currency            `json:"original_currency"`
	OriginalValuation *float64             `json:"original_valuation"`
	OriginalUnitvalue *float64             `json:"original_unitvalue"`
	OriginalUnitprice *float64             `json:"original_unitprice"`
	OriginalDiff      int                  `json:"original_diff"`
	// Details           int             `json:"details"`
}

// Connection as described at https://docs.budget-insight.com/reference/connections#response-connection-object
type Connection struct {
	ID           int              `json:"id"`
	UserID       *int             `json:"id_user"`
	ConnectorID  int              `json:"id_connector"`
	State        *ConnectionState `json:"state"`
	ErrorMessage *string          `json:"error_message"`
	// Fields           []FormFields              `json:"fields"`
	LastUpdate *Time `json:"last_update"`
	Created    *Time `json:"created"`
	Active     bool  `json:"active"`
	LastPush   *Time `json:"last_push"`
	NextTry    *Time `json:"next_try"`
}

type ConnectionState struct {
	SCARequired                 bool `json:"SCARequired"`
	WebauthRequired             bool `json:"webauthRequired"`
	AdditionalInformationNeeded bool `json:"additionalInformationNeeded"`
	Decoupled                   bool `json:"decoupled"`
	Validating                  bool `json:"validating"`
	ActionNeeded                bool `json:"actionNeeded"`
	PasswordExpired             bool `json:"passwordExpired"`
	Wrongpass                   bool `json:"SCAwrongpassRequired"`
	RateLimiting                bool `json:"rateLimiting"`
	WebsiteUnavailable          bool `json:"websiteUnavailable"`
	Bug                         bool `json:"bug"`
}
