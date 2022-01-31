/* SPDX-License-Identifier: MIT
 *
 * Copyright (C) 2017-2021 WireGuard LLC. All Rights Reserved.
 */

package device

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/subtle"
	"fmt"
	"hash"

	"golang.org/x/crypto/blake2s"
	"golang.org/x/crypto/curve25519"
)

/* KDF related functions.
 * HMAC-based Key Derivation Function (HKDF)
 * https://tools.ietf.org/html/rfc5869
 */

func HMAC1(sum *[blake2s.Size]byte, key, in0 []byte) {
	mac := hmac.New(func() hash.Hash {
		h, _ := blake2s.New256(nil)
		return h
	}, key)
	mac.Write(in0)
	mac.Sum(sum[:0])
}

func HMAC2(sum *[blake2s.Size]byte, key, in0, in1 []byte) {
	mac := hmac.New(func() hash.Hash {
		h, _ := blake2s.New256(nil)
		return h
	}, key)
	mac.Write(in0)
	mac.Write(in1)
	mac.Sum(sum[:0])
}

func KDF1(t0 *[blake2s.Size]byte, key, input []byte) {
	HMAC1(t0, key, input)
	HMAC1(t0, t0[:], []byte{0x1})
}

func KDF2(t0, t1 *[blake2s.Size]byte, key, input []byte) {
	println("KDF2 called")
	var prk [blake2s.Size]byte
	HMAC1(&prk, key, input)
	HMAC1(t0, prk[:], []byte{0x1})
	HMAC2(t1, prk[:], t0[:], []byte{0x2})
	setZero(prk[:])
}

func KDF3(t0, t1, t2 *[blake2s.Size]byte, key, input []byte) {
	println("KDF3 called")
	var prk [blake2s.Size]byte
	HMAC1(&prk, key, input)
	HMAC1(t0, prk[:], []byte{0x1})
	HMAC2(t1, prk[:], t0[:], []byte{0x2})
	HMAC2(t2, prk[:], t1[:], []byte{0x3})
	setZero(prk[:])
}

func isZero(val []byte) bool {
	acc := 1
	for _, b := range val {
		acc &= subtle.ConstantTimeByteEq(b, 0)
	}
	return acc == 1
}

/* This function is not used as pervasively as it should because this is mostly impossible in Go at the moment */
func setZero(arr []byte) {
	for i := range arr {
		arr[i] = 0
	}
}

/* I think the confusion comes from thinking that what you supply to
Curve25519 as the private key is already a scalar value.
It’s not, it’s a uniformly random key. Curve25519 then applies a surjective mapping
function (clamping) to this key to derive a suitable scalar value.
(Essentially, 5 bits of the key are ignored, leaving 2^251 distinct keys).
It’s not pretty, but it gets the job done and is fast.  */
func (sk *NoisePrivateKey) clamp() {
	sk[0] &= 248
	sk[31] = (sk[31] & 127) | 64
}

// If the private key is held on the device, can a clamp be performed since it looks
// like wireguard modifies the private key?
func newPrivateKey() (sk NoisePrivateKey, err error) {
	_, err = rand.Read(sk[:])
	sk.clamp()
	fmt.Printf("newPrivateKey() - called\n")
	return
}

// generate the public key from the private key
// we can handle this in pkclient
func (sk *NoisePrivateKey) publicKey() (pk NoisePublicKey) {
	apk := (*[NoisePublicKeySize]byte)(&pk)
	fmt.Printf("apk: %x\n", *apk)
	ask := (*[NoisePrivateKeySize]byte)(sk)
	fmt.Printf("ask: %x\n", *ask)
	curve25519.ScalarBaseMult(apk, ask)
	fmt.Printf("()privateKey.publicKey()\napk: %x\nask: %x\n----------end priv.publicKey()-----------------\n", *apk, *ask)
	return
}

// return the private key from the hsm instead
func (sk *NoisePrivateKey) publicKeyHSM(dev *Device) (pk NoisePublicKey) {
	if dev.staticIdentity.hsmEnabled {
		pk, _ := dev.staticIdentity.hsm.PublicKeyNoise()
		return pk
	}
	fmt.Println("THIS SHOULD NOT HAVE BEEN CALLED\n")
	return pk
}

func (sk *NoisePrivateKey) sharedSecret(pk NoisePublicKey) (ss [NoisePublicKeySize]byte) {
	apk := (*[NoisePublicKeySize]byte)(&pk)
	ask := (*[NoisePrivateKeySize]byte)(sk)
	curve25519.ScalarMult(&ss, ask, apk)
	fmt.Printf("priv.sharedSecret() - \n ss: %x\nask: %x\napk: %x\n----------end priv.sharedSecret()-----------------\n", ss, *ask, *apk)
	return ss
}

func ByteToNoisePublicKey(input []byte) (pk NoisePublicKey) {
	if len(input) > NoisePublicKeySize {
		fmt.Printf("Warning, bad key length input: len of input was: %d\n", len(input))
	}
	copy(pk[0:], input[0:])
	return pk
}
