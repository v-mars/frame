package utils

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"hash"
	"io"
)

// Nil empty UUID, all zeros
var Nil UUID

// A UUID is a 128 bit (16 byte) Universal Unique IDentifier as defined in RFC
// 4122.
type UUID [16]byte

// Bytes returns bytes slice representation of UUID.
func (uuid UUID) Bytes() []byte {
	return uuid[:]
}

// String returns the string form of uuid, xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
// , or "" if uuid is invalid.
func (uuid UUID) String() string {
	var buf [36]byte
	encodeHex(buf[:], uuid)
	return string(buf[:])
}

func encodeHex(dst []byte, uuid UUID) {
	hex.Encode(dst, uuid[:4])
	dst[8] = '-'
	hex.Encode(dst[9:13], uuid[4:6])
	dst[13] = '-'
	hex.Encode(dst[14:18], uuid[6:8])
	dst[18] = '-'
	hex.Encode(dst[19:23], uuid[8:10])
	dst[23] = '-'
	hex.Encode(dst[24:], uuid[10:])
}

// Must returns uuid if err is nil and panics otherwise.
func Must(uuid UUID, err error) UUID {
	if err != nil {
		panic(err)
	}
	return uuid
}

// NewRandom returns a Random (Version 4) UUID.
//
// The strength of the UUIDs is based on the strength of the crypto/rand
// package.
//
// A note about uniqueness derived from the UUID Wikipedia entry:
//
//  Randomly generated UUIDs have 122 random bits.  One's annual risk of being
//  hit by a meteorite is estimated to be one chance in 17 billion, that
//  means the probability is about 0.00000000006 (6 × 10−11),
//  equivalent to the odds of creating a few tens of trillions of UUIDs in a
//  year and having one duplicate.
func NewRandom() (UUID, error) {
	var uuid UUID
	_, err := io.ReadFull(rand.Reader, uuid[:])
	if err != nil {
		return Nil, err
	}
	uuid[6] = (uuid[6] & 0x0f) | 0x40 // Version 4
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // Variant is 10
	return uuid, nil
}

// NewHash returns a new UUID derived from the hash of space concatenated with
// data generated by h.  The hash should be at least 16 byte in length.  The
// first 16 bytes of the hash are used to form the UUID.  The version of the
// UUID will be the lower 4 bits of version.
func NewHash(h hash.Hash, space UUID, data []byte, version int) UUID {
	h.Reset()
	h.Write(space[:])
	h.Write(data)
	s := h.Sum(nil)
	var uuid UUID
	copy(uuid[:], s)
	uuid[6] = (uuid[6] & 0x0f) | uint8((version&0xf)<<4)
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // RFC 4122 variant
	return uuid
}

// NewMD5 returns a new MD5 (Version 3) UUID based on the
// supplied name space and data.  It is the same as calling:
//
//  NewHash(md5.New(), space, data, 3)
func NewMD5(space UUID, data []byte) UUID {
	return NewHash(md5.New(), space, data, 3)
}

// NewSHA1 returns a new SHA1 (Version 5) UUID based on the
// supplied name space and data.  It is the same as calling:
//
//  NewHash(sha1.New(), space, data, 5)
func NewSHA1(space UUID, data []byte) UUID {
	return NewHash(sha1.New(), space, data, 5)
}
