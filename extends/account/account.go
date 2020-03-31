package account

import (
	"errors"
	"fmt"

	"github.com/assetsadapterstore/livingthegrace-adapter/extends/account/pb"
	"github.com/assetsadapterstore/livingthegrace-adapter/extends/crypto"
	"github.com/assetsadapterstore/livingthegrace-adapter/extends/crypto/rsencoding"

	"github.com/golang/protobuf/proto"
)

var (
	ErrPublicKeyInvalidLen = errors.New("public key has invalid length")
)

type Account struct {
	*pb.Account
}

func NewAccount(id uint64) *Account {
	return &Account{
		Account: &pb.Account{
			Id:              id,
			RewardRecipient: id,
			Address:         rsencoding.Encode(id),
		},
	}
}

func (a *Account) ToBytes() []byte {
	if bs, err := proto.Marshal(a.Account); err == nil {
		return bs
	} else {
		panic(err)
	}
}

func FromBytes(bs []byte) *Account {
	var a pb.Account
	if err := proto.Unmarshal(bs, &a); err == nil {
		return &Account{Account: &a}
	} else {
		panic(err)
	}
}

func PublicKeyToID(publicKey []byte) uint64 {
	_, id := crypto.BytesToHashAndID(publicKey)
	return id
}

func IdToStringID(id uint64) string {
	return fmt.Sprintf("%020d", id)
}
