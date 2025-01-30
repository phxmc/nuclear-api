package utils

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func MustNewId() string {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		panic(err)
	}

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	hash := crypto.Keccak256Hash(publicKeyBytes[1:])
	address := hash[12:]

	return hexutil.Encode(address)
}
