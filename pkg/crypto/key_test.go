package crypto

import (
	"testing"
)

func TestGenerateRSAKeyAndEncodeAndLoad(t *testing.T) {

	k, err := GenerateRSAKey()
	if err != nil {
		t.Fatal(err)
	}

	raw, err := EncodePrivateKeyToPEM(k)
	if err != nil {
		t.Fatal(err)
	}

	_, err = LoadRSAPrivateKeyPEM(raw)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGenerateECDSAKeyAndEncodeAndLoad(t *testing.T) {

	k, err := GenerateECDSAKey()
	if err != nil {
		t.Fatal(err)
	}

	raw, err := EncodePrivateKeyToPEM(k)
	if err != nil {
		t.Fatal(err)
	}

	_, err = LoadECDSAPrivateKeyPEM(raw)
	if err != nil {
		t.Fatal(err)
	}
}
