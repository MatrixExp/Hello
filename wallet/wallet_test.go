package wallet

import "testing"

const (
	testSecretPhrase = "little orbit comfort eyebrow talk pink flame ridge bring milk equip blood"
	testSecretSeed   = "0x167d9a020688544ea246b056799d6a771e97c9da057e4d0b87024537f99177bc"
	testPubKey       = "0xdc64bef918ddda3126a39a11113767741ddfdf91399f055e1d963f2ae1ec2535"
	address44        = "5UczqUVGsoQpZnBCZkDtxvLxJ42KnUfaGTzPkQmZeAAug4s9"
)

func TestNewWallet(t *testing.T) {
	wallet, err := NewWallet(testSecretPhrase)
	if err != nil {
		t.Fatal(err)
	}
	sign, err := wallet.Sign([]byte{}, "")
	if err != nil {
		t.Fatal(sign)
	}
	t.Log(sign)
}

func TestGetPrivateKey(t *testing.T) {
	w, err := NewWallet(testSecretPhrase)
	if err != nil {
		t.Fatal(err)
	}
	privateKey, err := w.GetPrivateKeyHex()
	if err != nil {
		t.Fatal(err)
	}
	if testSecretSeed != privateKey {
		t.Fatal(privateKey)
	}
}
