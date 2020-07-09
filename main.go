package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
)

func getDigiID(seed, uri string, index uint32) (string, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, index)
	if err != nil {
		return "", err
	}
	err = binary.Write(buf, binary.BigEndian, []byte(uri))
	if err != nil {
		return "", err
	}
	h := sha256.New()
	h.Write(buf.Bytes())
	hash := h.Sum(nil)

	a := binary.LittleEndian.Uint32(hash[:4]) | 0x80000000
	b := binary.LittleEndian.Uint32(hash[4:8]) | 0x80000000
	c := binary.LittleEndian.Uint32(hash[8:12]) | 0x80000000
	d := binary.LittleEndian.Uint32(hash[12:16]) | 0x80000000

	fmt.Printf("%x %x %x %x\n", a, b, c, d)

	key, err := hdkeychain.NewKeyFromString(seed)
	if err != nil {
		return "", err
	}
	acct13, err := key.Child(hdkeychain.HardenedKeyStart + 13)
	if err != nil {
		return "", err
	}

	acctA, err := acct13.Child(a)
	if err != nil {
		return "", err
	}
	acctB, err := acctA.Child(b)
	if err != nil {
		return "", err
	}
	acctC, err := acctB.Child(c)
	if err != nil {
		return "", err
	}
	acctD, err := acctC.Child(d)
	if err != nil {
		return "", err
	}
	digiIDAddress, err := acctD.Address(&chaincfg.MainNetParams)
	if err != nil {
		return "", err
	}
	return digiIDAddress.EncodeAddress(), nil
}

func main() {
	const seed = "xprv9s21ZrQH143K3YjfLd4wnSrAowctp85Tp1BCL2EuwBVSSqVY4EPjFMTvY6DYxGbVkPp34gJYRxB9LwdJpJP62YxUby23WzvWQJebdG7bH1b"
	const index = 0
	const uri = "http://bitid.bitcoin.blue/callback"

	digiIDAddr, err := getDigiID(seed, uri, index)
	if err != nil {
		panic(err)
	}
	fmt.Println(digiIDAddr)
}
