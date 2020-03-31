package livingthegrace_addrdec

import (
	"encoding/hex"
	"fmt"
	"github.com/blocktree/go-owcrypt"
	"testing"
)

func TestAddressDecoder_AddressEncode(t *testing.T) {

	p2pk, _ := hex.DecodeString("3")
	p2pkAddr, _ := Default.AddressEncode(p2pk)
	t.Logf("p2pkAddr: %s", p2pkAddr)
}

func TestAddressDecoder_PubKeyEncode(t *testing.T) {

	p2pk, _ := hex.DecodeString("1")
	p2pkAddr, _ := Default.AddressEncode(p2pk)
	t.Logf("p2pkAddr: %s", p2pkAddr)

	p2pk2, _ := hex.DecodeString("2")
	p2pkAddr2, _ := Default.AddressEncode(p2pk2)
	t.Logf("p2pkAddr2: %s", p2pkAddr2)
}


func TestAddressDecoder_AddressDecode(t *testing.T) {

	p2pkAddr := "4"
	p2pkHash, _ := Default.AddressDecode(p2pkAddr)
	t.Logf("p2pkHash: %s", hex.EncodeToString(p2pkHash))
}

func TestDecompressPubKey(t *testing.T) {
	pub, _ := hex.DecodeString("5")
	uncompessedPublicKey := owcrypt.PointDecompress(pub, owcrypt.ECC_CURVE_SECP256K1)
	t.Logf("pub: %s", hex.EncodeToString(uncompessedPublicKey))
}


func TestAddressVerify(t *testing.T) {

	fmt.Println( Default.AddressVerify("LTG-73XT-BS88-QN5R-9P2KQ"))
	fmt.Println( Default.AddressVerify("LTG1-NKUN-SFN9-4D4U-7X8UE"))
	fmt.Println( Default.AddressVerify("LTG-NKUN-SFN9-4D4U-7X8U1"))
}


