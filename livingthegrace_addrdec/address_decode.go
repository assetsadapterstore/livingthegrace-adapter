package livingthegrace_addrdec

import (
	"fmt"
	"github.com/assetsadapterstore/livingthegrace-adapter/extends/account"
	"github.com/assetsadapterstore/livingthegrace-adapter/extends/crypto"
	"github.com/assetsadapterstore/livingthegrace-adapter/extends/crypto/rsencoding"
	"github.com/blocktree/go-owcrypt"
	"github.com/blocktree/openwallet/v2/openwallet"
)

var (
	Default = AddressDecoderV2{}
)

//AddressDecoderV2
type AddressDecoderV2 struct {
	*openwallet.AddressDecoderV2Base
}

//AddressEncode 地址编码
func (dec *AddressDecoderV2) AddressEncode(hash []byte, opts ...interface{}) (string, error) {

	edpub, _ := owcrypt.CURVE25519_convert_Ed_to_X(hash)
	_, id := crypto.BytesToHashAndID(edpub)
	accounts := account.NewAccount(id)

	return fmt.Sprintf("%s%s", "LTG-", accounts.Address), nil
}

// AddressVerify 地址校验
func (dec *AddressDecoderV2) AddressVerify(address string, opts ...interface{}) bool {

	if address[:4] != "LTG-" {
		return false
	}

	_, err := rsencoding.Decode(address[4:])
	if err != nil {
		return false
	}

	return true
}
