package password

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/skip2/go-qrcode"

	"github.com/makiuchi-d/gozxing"
	qrread "github.com/makiuchi-d/gozxing/qrcode"
)

var (
	Digits6 [2]int = [2]int{6, 1000000}
	Digits8 [2]int = [2]int{8, 100000000}

	Digits = Digits6

	OTP otp
)

type otp struct{}

func (otp) GenSecret(timeInterval int64) (string, error) {
	buf := bytes.Buffer{}
	err := binary.Write(&buf, binary.BigEndian, timeInterval)
	if err != nil {
		return "", err
	}
	hasher := hmac.New(sha1.New, buf.Bytes())
	secret := base32.StdEncoding.EncodeToString(hasher.Sum(nil))
	return secret, nil
}

func (otp) HOTP(secret string, timeInterval int64) string {
	key, err := base32.StdEncoding.DecodeString(strings.ToUpper(secret))
	if err != nil {
		fmt.Println(err)
	}
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(timeInterval))
	hasher := hmac.New(sha1.New, key)
	hasher.Write(buf)
	h := hasher.Sum(nil)
	offset := h[len(h)-1] & 0xf
	r := bytes.NewReader(h[offset : offset+4])

	var data uint32
	err = binary.Read(r, binary.BigEndian, &data)
	if err != nil {
		fmt.Println(err)
	}
	h12 := (int(data) & 0x7fffffff) % Digits[1]
	passcode := strconv.Itoa(h12)

	length := len(passcode)
	if length == Digits[0] {
		return passcode
	}
	for i := (Digits[0] - length); i > 0; i-- {
		passcode = "0" + passcode
	}
	return passcode
}

func (o *otp) TOTP(secret string) string {
	t := time.Now().Local().Unix() / 30
	return o.HOTP(secret, t)
}

func (o *otp) Verify(secret string, input string) bool {
	return o.TOTP(secret) == input
}

func NewOTP(account, issuer string) (string, error) {
	const uri string = "otpauth://totp/%s:%s?secret=%s&issuer=%s"
	t := time.Now().Local().Unix() / 30
	secret, err := OTP.GenSecret(t)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(uri, issuer, account, secret, issuer), nil
}

func GenQRCode(content string, size int) {
	dest := "/tmp/qrcode.tmp.png"
	err := qrcode.WriteColorFile(content, qrcode.Medium, size, color.White, color.Black, dest)
	if err != nil {
		fmt.Print(err)
	}
}

func ReadQRCode(filename string) string {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return err.Error()
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		fmt.Println(err)
		return err.Error()
	}
	png, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		fmt.Println(err)
		return err.Error()
	}
	qrReader := qrread.NewQRCodeReader()
	result, err := qrReader.Decode(png, nil)
	if err != nil {
		fmt.Println(err)
		return err.Error()
	}
	fmt.Println(result.String())
	return result.String()
}
