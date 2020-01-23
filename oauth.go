package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"io"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	log "github.com/sirupsen/logrus"
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

func authorizationHeader(url string) (string, error) {
	var oauthParams = map[string]string{
		"oauth_consumer_key":     apiKey,
		"oauth_signature_method": "HMAC-SHA1",
		"oauth_version":          "1.0",
		"oauth_token":            token,
		"oauth_timestamp":        strconv.FormatInt(time.Now().Unix(), 10),
		"oauth_nonce":            nonce(),
	}

	signature := getHMACSignature(url, oauthParams)
	oauthParams["oauth_signature"] = signature
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

func getHMACSignature(urlStr string, oauthParams map[string]string) string {
	key := encode(apiSecret, false)
	key = append(key, '&')
	u, err := url.Parse(urlStr)
	if err != nil {
		log.Fatal(err)
	}
	// if r.credentials != nil {
	key = append(key, encode(secret, false)...)
	// }
	h := hmac.New(sha1.New, key)
	writeBaseString(h, "GET", u, url.Values{}, oauthParams)
	signature := base64.StdEncoding.EncodeToString(h.Sum(key[:0]))
	// return string(rawSignature)
	return signature
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

type keyValue struct{ key, value []byte }
type byKeyValue []keyValue

func (p byKeyValue) Len() int      { return len(p) }
func (p byKeyValue) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p byKeyValue) Less(i, j int) bool {
	sgn := bytes.Compare(p[i].key, p[j].key)
	if sgn == 0 {
		sgn = bytes.Compare(p[i].value, p[j].value)
	}
	return sgn < 0
}

func (p byKeyValue) appendValues(values url.Values) byKeyValue {
	for k, vs := range values {
		k := encode(k, true)
		for _, v := range vs {
			v := encode(v, true)
			p = append(p, keyValue{k, v})
		}
	}
	return p
}

func writeBaseString(w io.Writer, method string, u *url.URL, form url.Values, oauthParams map[string]string) {
	// Method
	w.Write(encode(strings.ToUpper(method), false))
	w.Write([]byte{'&'})

	// URL
	scheme := strings.ToLower(u.Scheme)
	host := strings.ToLower(u.Host)

	uNoQuery := *u
	uNoQuery.RawQuery = ""
	path := uNoQuery.RequestURI()

	switch {
	case scheme == "http" && strings.HasSuffix(host, ":80"):
		host = host[:len(host)-len(":80")]
	case scheme == "https" && strings.HasSuffix(host, ":443"):
		host = host[:len(host)-len(":443")]
	}

	w.Write(encode(scheme, false))
	w.Write(encode("://", false))
	w.Write(encode(host, false))
	w.Write(encode(path, false))
	w.Write([]byte{'&'})

	// Create sorted slice of encoded parameters. Parameter keys and values are
	// double encoded in a single step. This is safe because double encoding
	// does not change the sort order.
	queryParams := u.Query()
	p := make(byKeyValue, 0, len(form)+len(queryParams)+len(oauthParams))
	p = p.appendValues(form)
	p = p.appendValues(queryParams)
	for k, v := range oauthParams {
		p = append(p, keyValue{encode(k, true), encode(v, true)})
	}
	sort.Sort(p)

	// Write the parameters.
	encodedAmp := encode("&", false)
	encodedEqual := encode("=", false)
	sep := false
	for _, kv := range p {
		if sep {
			w.Write(encodedAmp)
		} else {
			sep = true
		}
		w.Write(kv.key)
		w.Write(encodedEqual)
		w.Write(kv.value)
	}
}
