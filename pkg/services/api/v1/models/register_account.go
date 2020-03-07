package models

import (
	"crypto/aes"
	"crypto/cipher"
	crand "crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"math/rand"
	"protoChatServices/pkg/services/api/v1/global"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/mergermarket/go-pkcs7"
)

// Set global environment variable
var conf *global.Configuration
var level, cases, fatal string

// Function initialization
func init() {
	conf = global.New()
}

// function to encrypt using AES 256
func Encrypt(data string) (string, error) {
	key := []byte(conf.KeyAes)
	plainText := []byte(data)
	plainText, err := pkcs7.Pad(plainText, aes.BlockSize)

	if err != nil {
		err := fmt.Errorf(`plainText: "%s" has the wrong block size`, plainText)
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(crand.Reader, iv); err != nil {
		return "", err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText[aes.BlockSize:], plainText)

	return fmt.Sprintf("%x", cipherText), nil
}

// Function to decrypt
func Decrypt(encrypted string) (string, error) {
	key := []byte(conf.KeyAes)
	cipherText, _ := hex.DecodeString(encrypted)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	if len(cipherText) < aes.BlockSize {
		panic("cipherText too short")
	}
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]
	if len(cipherText)%aes.BlockSize != 0 {
		panic("cipherText is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(cipherText, cipherText)

	cipherText, _ = pkcs7.Unpad(cipherText, aes.BlockSize)
	return fmt.Sprintf("%s", cipherText), nil
}

// Function Generate random number
func getRandomString() string {
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

	const charset = "1234567890" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	b := make([]byte, 6)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}

	return string(b)
}

// Function to create a JWT Token
func JwtTokenCreate(token *global.Tokenization) (string, error) {
	tokenAuth := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), token)
	tokenNew, _ := tokenAuth.SignedString([]byte(conf.Token))

	return tokenNew, nil
}

// Function for register account
func RegisterAccount(username string) (code, status, message, token, registeredDate string) {

	// Get data UserId
	userId := getRandomString()
	timeDate := time.Now().Format("2006-01-02 15:04:05")

	tokenJwt := &global.Tokenization{
		UserId:   userId,
		Time:     timeDate,
		Username: username,
	}

	tokenAuth, _ := JwtTokenCreate(tokenJwt)
	encrypted, _ := Encrypt(tokenAuth)

	token = encrypted
	registeredDate = timeDate
	code = "00"
	message = "Registered Account Success"
	status = "Register Account Success"

	return code, status, message, token, registeredDate
}
