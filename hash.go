package main

import (
  "encoding/hex"
  "io"
  "fmt"
  "hash"
  "crypto/md5"
  "crypto/sha1"
  "crypto/sha256"
  "crypto/sha512"
  flow "github.com/trustmaster/goflow"
)

type HashResult struct {
	Word, Result string
}

// format HashResult
func (hr HashResult) String() string {
	return fmt.Sprintf("%s %s", hr.Result, hr.Word)
}

// make a hashresult from hashtype, word
func HashString(hashtype, word string) HashResult {
  var hsh hash.Hash

  switch hashtype {
  case "MD5":
    hsh = md5.New()
  case "SHA1":
    hsh = sha1.New()
  case "SHA256":
    hsh = sha256.New()
  case "SHA512":
    hsh = sha512.New()
  default:
    panic("invalid hash type")
    return HashResult{}
  }

  io.WriteString(hsh, word)
  return HashResult{word, hex.EncodeToString(hsh.Sum(nil))}
}

// a flow component that computes hashes
type Hasher struct {
	flow.Component

	Word   <-chan string     // input for strings to hash
	Result chan<- HashResult // output for the computed hash

	hashtype string // hash function we are using
}

func NewHasher(hashtype string) *Hasher {
  h := new(Hasher)

  h.hashtype = hashtype

  return h
}

func (h *Hasher) OnWord(word string) {
  h.Result <- HashString(h.hashtype, word)
}

