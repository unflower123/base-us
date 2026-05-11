package balance

import (
	"fmt"
	"support_platform/pkg/bcrypt"
	"testing"
)

func Test_ReidsDeductBalance(t *testing.T) {
	amount := int64(403)
	rate := int64(6000)
	fee := int64(0)
	fmt.Println(GetCost(amount, rate, fee))

	amount = int64(100)
	rate = int64(0)
	fee = int64(300)
	fmt.Println(GetCost(amount, rate, fee))

	password := bcrypt.BcryptHash("12345678")
	fmt.Println(password)
	fmt.Println(password == "$2a$10$X1kRZ2ZacInshv2Z87GxI.y5bs/Ki42spYAVlz2pwhQImxcwgMPrG")
}
