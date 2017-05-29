package skypeapi

import (
	"strings"
	"fmt"
	"encoding/base64"
	"encoding/json"
	"bytes"
	"time"
	"crypto/x509"
	"net/http"
	"encoding/pem"
)

const (
	openIdRequestPath                   string = "https://login.botframework.com/v1/.well-known/openidconfiguration"
	authorizationHeaderValuePrefix      string = "Bearer "
	wrongAuthorizationHeaderFormatError string = "The provided authorization header is in the wrong format: %v"
	wrongSplitLengthError               string = "The authorize value split length with character \"%v\" is not valid: %v (%v)"
	splitCharacter                      string = "."
	issuerUrl                           string = "https://api.botframework.com"
)

type OpenIdDocument struct {
	Issuer                            string `json:"issuer"`
	AuthorizationEndpoint             string `json:"authorization_endpoint"`
	JwksURI                           string `json:"jwks_uri"`
	IDTokenSigningAlgValuesSupported  []string `json:"id_token_signing_alg_values_supported"`
	TokenEndpointAuthMethodsSupported []string `json:"token_endpoint_auth_methods_supported"`
}

type SigningKeys struct {
	Keys []struct {
		Kty          string `json:"kty"`
		Use          string `json:"use"`
		KeyId        string `json:"kid"`
		X5T          string `json:"x5t"`
		N            string `json:"n"`
		E            string `json:"e"`
		X5C          []string `json:"x5c"`
		Endorsements []string `json:"endorsements,omitempty"`
	} `json:"keys"`
}

type JwtHeader struct {
	Type            string `json:"typ"`
	Algorithm       string `json:"alg"`
	SigningKeyId    string `json:"kid"`
	SigningKeyIdX5T string `json:"x5t"`
}

type JwtPayload struct {
	ServiceUrl   string `json:"serviceurl"`
	Issuer       string `json:"iss"`
	Audience     string `json:"aud"`
	Expires      int `json:"exp"`
	CreatedOnNbf int `json:"nbf"`
}

type MicrosoftJsonWebToken struct {
	HeaderBase64, PayloadBase64 string
	Header                      JwtHeader
	Payload                     JwtPayload
	VerifySignature             []byte
}

func (microSoftJsonWebToken MicrosoftJsonWebToken) Verify(microsoftAppId string, signingKeys SigningKeys) bool {
	if microSoftJsonWebToken.Payload.Issuer != issuerUrl {
		return false
	} else if microSoftJsonWebToken.Payload.Audience != microsoftAppId {
		return false
	} else if int64(microSoftJsonWebToken.Payload.Expires) <= time.Now().Unix() {
		return false
	} else {
		return microSoftJsonWebToken.verifyCertificate(signingKeys)
	}
}

func GetSigningKeys() (SigningKeys, error) {
	openIdDocument := &OpenIdDocument{}
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, openIdRequestPath, nil)
	if err != nil {
		return SigningKeys{}, err
	} else {
		resp, err := client.Do(req)
		if err != nil {
			return SigningKeys{}, err
		} else {
			defer resp.Body.Close()
			err := json.NewDecoder(resp.Body).Decode(openIdDocument)
			if err != nil {
				return SigningKeys{}, err
			} else {
				return GetSigningKeysByUrl(openIdDocument.JwksURI)
			}
		}
	}
}

func GetSigningKeysByUrl(url string) (SigningKeys, error) {
	signingKeys := &SigningKeys{}
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return SigningKeys{}, err
	} else {
		resp, err := client.Do(req)
		if err != nil {
			return SigningKeys{}, err
		} else {
			defer resp.Body.Close()
			err := json.NewDecoder(resp.Body).Decode(signingKeys)
			if err != nil {
				return SigningKeys{}, err
			} else {
				return *signingKeys, nil
			}
		}
	}
}

func (microSoftJsonWebToken MicrosoftJsonWebToken) verifyCertificate(signingKeys SigningKeys) bool {
	for _, key := range signingKeys.Keys {
		if key.KeyId == microSoftJsonWebToken.Header.SigningKeyId {
			certificateParsed := parseCertificateString(key.X5C[0])
			block, _ := pem.Decode([]byte(certificateParsed))
			if certificate, err := x509.ParseCertificate(block.Bytes); err != nil {
				return false
			} else {
				hashed := []byte(microSoftJsonWebToken.HeaderBase64 + splitCharacter + microSoftJsonWebToken.PayloadBase64)
				return certificate.CheckSignature(x509.SHA256WithRSA, hashed, microSoftJsonWebToken.VerifySignature) == nil
			}
		}
	}
	return false
}

func parseCertificateString(rawCertificate string) string {
	parsedCertificate := "-----BEGIN CERTIFICATE-----\n"
	buffer := bytes.NewBuffer(make([]byte, 64))
	buffer.Reset()
	for index, charByte := range []byte(rawCertificate) {
		buffer.WriteByte(charByte)
		if (index+1)%64 == 0 {
			parsedCertificate += string(buffer.Bytes()) + "\n"
			buffer.Reset()
		}
	}
	parsedCertificate += string(buffer.Bytes()) + "\n-----END CERTIFICATE-----"
	buffer.Reset()
	return parsedCertificate
}

func ParseMicrosoftJsonWebToken(headerValue string) (MicrosoftJsonWebToken, error) {
	microsoftJsonWebToken := &MicrosoftJsonWebToken{}
	parsedHeaderValue := parseHeaderValue(headerValue)
	if len(parsedHeaderValue) == 0 {
		return *microsoftJsonWebToken, fmt.Errorf(wrongAuthorizationHeaderFormatError, parsedHeaderValue)
	} else {
		var split []string = strings.Split(parsedHeaderValue, splitCharacter)
		if len(split) == 3 {
			jwtHeader := &JwtHeader{}
			jwtPayload := &JwtPayload{}
			if err := decodeBase64JsonPart(split[0], jwtHeader); err != nil {
				return *microsoftJsonWebToken, err
			}
			if err := decodeBase64JsonPart(split[1], jwtPayload); err != nil {
				return *microsoftJsonWebToken, err
			}
			jwtVerifySignature, err := base64.RawURLEncoding.DecodeString(split[2])
			if err != nil {
				return *microsoftJsonWebToken, err
			}
			microsoftJsonWebToken.HeaderBase64 = split[0]
			microsoftJsonWebToken.PayloadBase64 = split[1]
			microsoftJsonWebToken.Header = *jwtHeader
			microsoftJsonWebToken.Payload = *jwtPayload
			microsoftJsonWebToken.VerifySignature = jwtVerifySignature
			return *microsoftJsonWebToken, nil
		} else {
			return *microsoftJsonWebToken, fmt.Errorf(wrongSplitLengthError, splitCharacter, len(split), parseHeaderValue)
		}
	}
}

func decodeBase64JsonPart(rawPart string, partObj interface{}) error {
	if partBytes, err := base64.RawURLEncoding.DecodeString(rawPart); err != nil {
		return err
	} else {
		return json.NewDecoder(bytes.NewReader(partBytes)).Decode(&partObj)
	}
}

func parseHeaderValue(headerValue string) string {
	if index := strings.Index(headerValue, authorizationHeaderValuePrefix); index == 0 &&
		len(headerValue) > len(authorizationHeaderValuePrefix) {
		return headerValue[len(authorizationHeaderValuePrefix):]
	} else {
		return ""
	}
}
