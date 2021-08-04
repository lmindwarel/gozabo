package gozabo

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
		t, err = time.Parse(time.RFC1123, s)
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

// Account as described at https://docs.budget-insight.com/reference/bank-accounts#response-bankaccount-object
type Account struct {
	ID             int       `json:"id"`
	Token          string    `json:"token"`
	ExpirationTime time.Time `json:"exp_time"`
	Provider       Provider  `json:"id_source"`
	// Balances []Balance `json:"balances"`
	Blockchain *string
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type AuthType string

const (
	AuthTypeReadOnly          AuthType = "read_only"
	AuthTypePassword          AuthType = "password"
	AuthTypeToken             AuthType = "token"
	AuthTypeExtendedPublicKey AuthType = "xpub"
	AuthTypeOAuth             AuthType = "oauth"
	AuthTypeWeb3Wallet        AuthType = "web3"
)

type Scope string

const (
	ScopeGetDepositAddress      Scope = "get_deposit_address"
	ScopeCreateDepositAddress   Scope = "create__deposit_address"
	ScopeReadTransactionHistory Scope = "read_transaction_history"
	ScopeReadBalances           Scope = "read_balances"
)

type Provider struct {
	Name        string     `json:"name"`
	DisplayName string     `json:"display_name"`
	Logo        string     `json:"logo"`
	AuthType    AuthType   `json:"auth_type"`
	Scopes      []Scope    `json:"scopes"`
	Currencies  []Currency `json:"currencies"`
}

type AssetType string

const (
	AssetTypeUTXO    AssetType = "utxo"
	AssetTypeAccount AssetType = "account"
	AssetTypeERC20   AssetType = "erc20"
)

type Currency struct {
	Type AssetType `json:"type"`
	List []string  `json:"list"`
}
