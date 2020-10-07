package account

import (
	"gitlab.com/michaelk99/birrdi/api-soa/services/profile"
)

// Creator is a public interface which creates a token for the provided
// subject. Implementations can decide what format this token can be (jwt, saml, etc...)
type TokenCreator interface {
	Create(acct *Account, prof *profile.Profile, ar []*AccountRole) (string, error)
}
