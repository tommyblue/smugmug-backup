package main

import (
	"crypto/rand"
	"encoding/binary"
	"os"
	"strconv"
	"sync/atomic"
	"time"
)

var apiKey string
var apiSecret string
var token string
var secret string

var oauthKeys = []string{
	"oauth_consumer_key",
	"oauth_nonce",
	"oauth_signature",
	"oauth_signature_method",
	"oauth_timestamp",
	"oauth_token",
	"oauth_version",
	"oauth_callback",
	"oauth_verifier",
	"oauth_session_handle",
}

var nonceCounter uint64

func init() {
	if err := binary.Read(rand.Reader, binary.BigEndian, &nonceCounter); err != nil {
		// fallback to time if rand reader is broken
		nonceCounter = uint64(time.Now().UnixNano())
	}

	apiKey = os.Getenv("API_KEY")
	apiSecret = os.Getenv("API_SECRET")
	token = os.Getenv("USER_TOKEN")
	secret = os.Getenv("USER_SECRET")
}

// nonce returns a unique string.
func nonce() string {
	return strconv.FormatUint(atomic.AddUint64(&nonceCounter, 1), 16)
}

func authorizationHeader() (string, error) {
	var oauthParams = map[string]string{
		"oauth_consumer_key":     apiKey,
		"oauth_signature_method": "PLAINTEXT",
		"oauth_version":          "1.0",
		"oauth_token":            token,
		"oauth_signature":        getSignature(),
		"oauth_timestamp":        strconv.FormatInt(time.Now().Unix(), 10),
		"oauth_nonce":            nonce(),
	}
	var h []byte
	// Append parameters in a fixed order to support testing.
	for _, k := range oauthKeys {
		if v, ok := oauthParams[k]; ok {
			if h == nil {
				h = []byte(`OAuth `)
			} else {
				h = append(h, ", "...)
			}
			h = append(h, k...)
			h = append(h, `="`...)
			h = append(h, encode(v, false)...)
			h = append(h, '"')
		}
	}

	return string(h), nil
}

func getSignature() string {
	rawSignature := encode(apiSecret, false)
	rawSignature = append(rawSignature, '&')
	// if r.credentials != nil {
	rawSignature = append(rawSignature, encode(secret, false)...)
	// }
	return string(rawSignature)
}

// noscape[b] is true if b should not be escaped per section 3.6 of the RFC.
var noEscape = [256]bool{
	'A': true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true,
	'a': true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true,
	'0': true, true, true, true, true, true, true, true, true, true,
	'-': true,
	'.': true,
	'_': true,
	'~': true,
}

// encode encodes string per section 3.6 of the RFC. If double is true, then
// the encoding is applied twice.
func encode(s string, double bool) []byte {
	// Compute size of result.
	m := 3
	if double {
		m = 5
	}
	n := 0
	for i := 0; i < len(s); i++ {
		if noEscape[s[i]] {
			n++
		} else {
			n += m
		}
	}

	p := make([]byte, n)

	// Encode it.
	j := 0
	for i := 0; i < len(s); i++ {
		b := s[i]
		if noEscape[b] {
			p[j] = b
			j++
		} else if double {
			p[j] = '%'
			p[j+1] = '2'
			p[j+2] = '5'
			p[j+3] = "0123456789ABCDEF"[b>>4]
			p[j+4] = "0123456789ABCDEF"[b&15]
			j += 5
		} else {
			p[j] = '%'
			p[j+1] = "0123456789ABCDEF"[b>>4]
			p[j+2] = "0123456789ABCDEF"[b&15]
			j += 3
		}
	}
	return p
}
