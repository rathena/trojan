package integration

// A specific type to restrict invalid gender suffix values.
type GenderSuffix string

// The gender suffx values used when creating a new account.
const (
	CreateAccountMale   GenderSuffix = "_M"
	CreateAccountFemale GenderSuffix = "_F"
)

// The core account structure consisting of common values used by all
// dependent systems.
type Account struct {
	AccountID  uint32
	AuthKeyOne uint32
	AuthKeyTwo uint32
	Gender     uint8
}
