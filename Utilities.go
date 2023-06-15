package FlowX

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

func ReplaceString(input, old, new string) string {
	if input == "" {
		return ""
	}
	return strings.Replace(input, old, new, -1)
}

func CountStringOccurrences(input, substr string) (int, error) {
	count := strings.Count(input, substr)
	if count == -1 {
		return 0, fmt.Errorf("error counting occurrences of substring")
	}
	return count, nil
}

func CurrentUnixTime() int64 {
	return time.Now().Unix()
}

func getFieldValue(jsonData, fieldName string) (interface{}, error) {
	var data map[string]interface{}
	err := json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON: %v", err)
	}

	fieldValue, ok := data[fieldName]
	if !ok {
		return nil, fmt.Errorf("field not found: %s", fieldName)
	}

	return fieldValue, nil
}

func Base64Encode(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

func Base64Decode(str string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return "", fmt.Errorf("error decoding base64: %v", err)
	}
	return string(data), nil
}

func ParseLR(str, before, after string) (string, error) {
	idx := strings.Index(str, before)
	if idx == -1 {
		return "", fmt.Errorf("substring '%s' not found in input", before)
	}

	start := idx + len(before)
	end := strings.Index(str[start:], after)
	if end == -1 {
		return "", fmt.Errorf("substring '%s' not found after '%s'", after, before)
	}

	return str[start : start+end], nil
}

func URLEncodeString(input string) string {
	result := url.QueryEscape(input)
	return result
}

func URLDecodeString(input string) (string, error) {
	result, err := url.QueryUnescape(input)
	if err != nil {
		return "", err
	}
	return result, nil
}

func StringSliceContains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func RandomString(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}

func MapKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	return keys
}

func ReverseString(str string) string {
	runes := []rune(str)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func SerializeJSON(data interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("error serializing JSON: %v", err)
	}
	return jsonData, nil
}

func DeserializeJSON(jsonData []byte, target interface{}) error {
	err := json.Unmarshal(jsonData, target)
	if err != nil {
		return fmt.Errorf("error deserializing JSON: %v", err)
	}
	return nil
}

func ReadFile(filename string) ([]byte, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}
	return data, nil
}

func WriteFile(filename string, data []byte) error {
	err := ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("error writing file: %v", err)
	}
	return nil
}

func CreateDirectory(dirName string) error {
	err := os.Mkdir(dirName, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error creating directory: %v", err)
	}
	return nil
}

func DeleteFile(filename string) error {
	err := os.Remove(filename)
	if err != nil {
		return fmt.Errorf("error deleting file: %v", err)
	}
	return nil
}

func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func IsEmailValid(email string) bool {
	emailRegex := `^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9-]+(?:\.[a-zA-Z0-9-]+)*$`
	match, _ := regexp.MatchString(emailRegex, email)
	return match
}

func IsURLValid(urlString string) bool {
	_, err := url.ParseRequestURI(urlString)
	return err == nil
}

func IsPhoneNumberValid(phoneNumber string) bool {
	phoneRegex := `^[0-9]{10}$`
	match, _ := regexp.MatchString(phoneRegex, phoneNumber)
	return match
}

func EncryptAES(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("error creating AES cipher: %v", err)
	}

	ciphertext := make([]byte, aes.BlockSize+len(data))
	iv := ciphertext[:aes.BlockSize]
	if _, err := rand.Read(iv); err != nil {
		return nil, fmt.Errorf("error generating IV: %v", err)
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], data)

	return ciphertext, nil
}

func DecryptAES(ciphertext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("error creating AES cipher: %v", err)
	}

	if len(ciphertext) < aes.BlockSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	data := ciphertext[aes.BlockSize:]

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(data, data)

	return data, nil
}

func HashSHA256(data []byte) ([]byte, error) {
	hash := sha256.Sum256(data)
	return hash[:], nil
}
