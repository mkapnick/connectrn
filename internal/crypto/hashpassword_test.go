package crypto_test

import (
	"testing"

	"gitlab.com/michaelk99/connectrn/internal/crypto"
	"golang.org/x/crypto/bcrypt"
)

var testHashPasswordTT = []struct {
	password string
}{
	{
		`pass1234!@#$`,
	},
	{
		`d#VYG5h}5{q^C\x*`,
	},
	{
		`5nh^}+>vk&=Qc{y'8+_wYn==a8)BVA9hwv2~%9`,
	},
	{
		`5#P5eQG#m#[h9G@k6Np9&{-ADF`,
	},
	{
		`L4*K*mwBm&&S&T?2`,
	},
	{
		`Epd$E9+c@7XR9TZ5pB^xRCx5V24^EGX@d6AjX6NrPuLDJY^J`,
	},
	{
		`v23UGR!Bafdr`,
	},
}

func TestHashPassword(t *testing.T) {
	for _, tt := range testHashPasswordTT {
		hash, err := crypto.HashPassword(tt.password)
		if err != nil {
			t.Fatalf("failed to hash password %s: %v", tt.password, err)
		}

		err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(tt.password))
		if err != nil {
			t.Fatalf("HashPassword() output does not match provided password: %v", err)
		}
	}
}

var testValidatePasswordTT = []struct {
	password string
	hash     string
	expect   bool
}{
	{
		password: "pass123!@#",
		hash:     "$2a$07$o479KsDT3QOmLcyQi7WlPujmu7mX899eP0BTdpebkfTHIYdlUCGhe",
		expect:   true,
	},
	{
		password: "pass123!@#",
		hash:     "$2a$07$o479KsDT3QOmLcyQi7WlPuj",
		expect:   false,
	},
	{
		password: "/9?d=U&N*{*H$2@p",
		hash:     "$2a$07$qsgnu3dhUG8zEGsX8BKKCug.2.VbCWCHXaxOVIp78IEAzQ9bfz8i.",
		expect:   true,
	},
	{
		password: "/9?d=U&N*{*H$2@p",
		hash:     "ju!nk$2a$07$qsgnu3dhUG8zEGsX8BKKCug.2.VbCWCHXaxOVIp78IEAzQ9bfz8i.",
		expect:   false,
	},
}

func TestValidatePassword(t *testing.T) {
	for _, tt := range testValidatePasswordTT {
		valid := crypto.ValidatePassword(tt.password, tt.hash)

		if valid != tt.expect {
			t.Fatalf("expected validation to be %v got %v", tt.expect, valid)
		}
	}
}
