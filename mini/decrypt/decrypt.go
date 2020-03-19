package decrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/xiya-team/gowechat/mini/base"
	"github.com/xiya-team/gowechat/wxcontext"
)

var (
	// ErrAppIDNotMatch appid不匹配
	ErrAppIDNotMatch = errors.New("app id not match")
	// ErrInvalidBlockSize block size不合法
	ErrInvalidBlockSize = errors.New("invalid block size")
	// ErrInvalidPKCS7Data PKCS7数据不合法
	ErrInvalidPKCS7Data = errors.New("invalid PKCS7 data")
	// ErrInvalidPKCS7Padding 输入padding失败
	ErrInvalidPKCS7Padding = errors.New("invalid padding on input")
)

type watermark struct {
	AppID     string `json:"appid"`
	Timestamp int64  `json:"timestamp"`
}

// Mobile 解密后的用户手机号码信息
type Mobile struct {
	PhoneNumber     string    `json:"phoneNumber"`
	PurePhoneNumber string    `json:"purePhoneNumber"`
	CountryCode     string    `json:"countryCode"`
	Watermark       watermark `json:"watermark"`
}

//Pay pay
type Decrypt struct {
	base.MiniBase
}

//NewOauth 实例化授权信息
func NewDecrypt(context *wxcontext.Context) (decrypt *Decrypt) {
	decrypt = new(Decrypt)
	decrypt.Context = context
	return
}

// DecryptMobile 解密手机号码
//
// sessionKey 通过 Login 向微信服务端请求得到的 session_key
// encryptedData 小程序通过 api 得到的加密数据(encryptedData)
// iv 小程序通过 api 得到的初始向量(iv)
func (c *Decrypt) DecryptMobile(sessionKey, encryptedData, iv string) (*Mobile, error) {
	raw, err := decryptUserData(sessionKey, encryptedData, iv)
	if err != nil {
		return nil, err
	}

	mobile := new(Mobile)
	if err := json.Unmarshal(raw, mobile); err != nil {
		return nil, err
	}

	return mobile, nil
}

// ShareInfo 解密后的分享信息
type ShareInfo struct {
	GID string `json:"openGId"`
}

// DecryptShareInfo 解密转发信息的加密数据
//
// sessionKey 通过 Login 向微信服务端请求得到的 session_key
// encryptedData 小程序通过 api 得到的加密数据(encryptedData)
// iv 小程序通过 api 得到的初始向量(iv)
//
// gid 小程序唯一群号
func (c *Decrypt) DecryptShareInfo(sessionKey, encryptedData, iv string) (*ShareInfo, error) {

	raw, err := decryptUserData(sessionKey, encryptedData, iv)
	if err != nil {
		return nil, err
	}

	info := new(ShareInfo)
	if err = json.Unmarshal(raw, info); err != nil {
		return nil, err
	}

	return info, nil
}

// UserInfo 解密后的用户信息
type UserInfo struct {
	OpenID    string    `json:"openId"`
	Nickname  string    `json:"nickName"`
	Gender    int       `json:"gender"`
	Province  string    `json:"province"`
	Language  string    `json:"language"`
	Country   string    `json:"country"`
	City      string    `json:"city"`
	Avatar    string    `json:"avatarUrl"`
	UnionID   string    `json:"unionId"`
	Watermark watermark `json:"watermark"`
}

// DecryptUserInfo 解密用户信息
//
// sessionKey 微信 session_key
// rawData 不包括敏感信息的原始数据字符串，用于计算签名。
// encryptedData 包括敏感数据在内的完整用户信息的加密数据
// signature 使用 sha1( rawData + session_key ) 得到字符串，用于校验用户信息
// iv 加密算法的初始向量
func (c *Decrypt) DecryptUserInfo(sessionKey, rawData, encryptedData, signature, iv string) (*UserInfo, error) {

	if ok := validateUserInfo(signature, rawData, sessionKey); !ok {
		return nil, errors.New("failed to validate signature")
	}

	raw, err := decryptUserData(sessionKey, encryptedData, iv)
	if err != nil {
		return nil, err
	}

	info := new(UserInfo)
	if err := json.Unmarshal(raw, info); err != nil {
		return nil, err
	}

	return info, nil
}

// Decrypt 解密数据
func (wxa *Decrypt) Decrypt(sessionKey, encryptedData, iv string) (*UserInfo, error) {
	cipherText, err := getCipherText(sessionKey, encryptedData, iv)
	if err != nil {
		return nil, err
	}
	var userInfo UserInfo
	err = json.Unmarshal(cipherText, &userInfo)
	if err != nil {
		return nil, err
	}
	if userInfo.Watermark.AppID != wxa.AppID {
		return nil, ErrAppIDNotMatch
	}
	return &userInfo, nil
}

// getCipherText returns slice of the cipher text
func getCipherText(sessionKey, encryptedData, iv string) ([]byte, error) {
	aesKey, err := base64.StdEncoding.DecodeString(sessionKey)
	if err != nil {
		return nil, err
	}
	cipherText, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return nil, err
	}
	ivBytes, err := base64.StdEncoding.DecodeString(iv)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}
	mode := cipher.NewCBCDecrypter(block, ivBytes)
	mode.CryptBlocks(cipherText, cipherText)
	cipherText, err = pkcs7Unpad(cipherText, block.BlockSize())
	if err != nil {
		return nil, err
	}
	return cipherText, nil
}

// pkcs7Unpad returns slice of the original data without padding
func pkcs7Unpad(data []byte, blockSize int) ([]byte, error) {
	if blockSize <= 0 {
		return nil, ErrInvalidBlockSize
	}
	if len(data)%blockSize != 0 || len(data) == 0 {
		return nil, ErrInvalidPKCS7Data
	}
	c := data[len(data)-1]
	n := int(c)
	if n == 0 || n > len(data) {
		return nil, ErrInvalidPKCS7Padding
	}
	for i := 0; i < n; i++ {
		if data[len(data)-n+i] != c {
			return nil, ErrInvalidPKCS7Padding
		}
	}
	return data[:len(data)-n], nil
}