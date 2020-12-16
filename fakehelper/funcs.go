package fakehelper

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image/color"
	"image/png"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unsafe"

	"github.com/brianvoe/gofakeit/v5"
	"github.com/issue9/identicon"

	"github.com/winwisely268/sanfood-faker/utilities"
)

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	roles         = "admin acrfixed acrleader acrmobile anonymous asm exec staff"
)

type RefCount struct {
	Sequence      int
	Previous      int
	LastReference int
}

func NewRefCount() *RefCount {
	return &RefCount{
		Sequence:      0,
		Previous:      0,
		LastReference: 0,
	}
}

// Copies the refcount and resets the last reference
func (r *RefCount) ResetLastReference() *RefCount {
	rc := r
	rc.LastReference = 0
	return rc
}

func FakeGenRole() (string, gofakeit.Info) {
	return "rolegen", gofakeit.Info{
		Category:    "rolegen",
		Description: "generate random user role",
		Example:     "rolegen",
		Output:      "*string",
		Params:      []gofakeit.Param{},
		Call: func(m *map[string][]string, info *gofakeit.Info) (interface{}, error) {
			splitted := strings.Split(roles, " ")
			lenSplit := len(splitted)
			randIdx := rand.Intn(lenSplit)
			return splitted[randIdx], nil
		},
	}
}

func FakeMailGen() (string, gofakeit.Info) {
	return "mailgen", gofakeit.Info{
		Category:    "mailgen",
		Description: "generate random user email",
		Example:     "mailgen",
		Output:      "string",
		Params:      []gofakeit.Param{},
		Call: func(m *map[string][]string, info *gofakeit.Info) (interface{}, error) {
			mail := gofakeit.Email()
			return mail, nil
		},
	}
}

func FakeNameSequence(callFunc func(prefix, referral string, isRef, isUniqueRef, reset bool) (interface{}, error)) (string, gofakeit.Info) {
	return "nameseq", gofakeit.Info{
		Category:    "nameseq",
		Description: "generate sequential string with number, i.e.: sys_account_project_1",
		Example:     "sys_account",
		Output:      "string",
		Params: []gofakeit.Param{
			{Field: "prefix", Type: "string", Description: "module name or prefix"},
			{Field: "isRef", Type: "bool", Description: "the current field is a reference to another field in another struct", Default: "false"},
			{Field: "referral", Type: "string", Description: "the current field refers to which field"},
			{Field: "isUnique", Type: "bool", Description: "unique flag for reference"},
			{Field: "reset", Type: "bool", Description: "reset reference"},
		},
		Call: func(m *map[string][]string, info *gofakeit.Info) (interface{}, error) {
			prefix, err := info.GetString(m, "prefix")
			if err != nil {
				return nil, err
			}
			isRef, err := info.GetBool(m, "isRef")
			if err != nil {
				return nil, err
			}
			referral, err := info.GetString(m, "referral")
			if err != nil {
				return nil, err
			}
			isUniqueRef, err := info.GetBool(m, "isUnique")
			if err != nil {
				return nil, err
			}
			reset, err := info.GetBool(m, "reset")
			if err != nil {
				return nil, err
			}
			return callFunc(prefix, referral, isRef, isUniqueRef, reset)
		},
	}
}

func FakeMailSequence(callFunc func(prefix, referral string, isRef, isUniqueRef bool) (interface{}, error)) (string, gofakeit.Info) {
	return "mailseq", gofakeit.Info{
		Category:    "mailseq",
		Description: "generate sequential email string with number, i.e.: account1@domain.com",
		Example:     "sys_account",
		Output:      "string",
		Params: []gofakeit.Param{
			{Field: "prefix", Type: "string", Description: "module name or prefix"},
			{Field: "isRef", Type: "bool", Description: "the current field is a reference to another field in another struct", Default: "false"},
			{Field: "referral", Type: "string", Description: "the current field refers to which field"},
			{Field: "isUnique", Type: "bool", Description: "unique flag for reference"},
		},
		Call: func(m *map[string][]string, info *gofakeit.Info) (interface{}, error) {
			prefix, err := info.GetString(m, "prefix")
			if err != nil {
				return nil, err
			}
			isRef, err := info.GetBool(m, "isRef")
			if err != nil {
				return nil, err
			}
			referral, err := info.GetString(m, "referral")
			if err != nil {
				return nil, err
			}
			isUniqueRef, err := info.GetBool(m, "isUnique")
			if err != nil {
				return nil, err
			}
			return callFunc(prefix, referral, isRef, isUniqueRef)
		},
	}
}

func FakeRandomTs() (string, gofakeit.Info) {
	return "randomts", gofakeit.Info{
		Category:    "randomts",
		Description: "generate random protobuf timestamp",
		Example:     "randomts",
		Output:      "int64",
		Params:      []gofakeit.Param{},
		Call: func(m *map[string][]string, info *gofakeit.Info) (interface{}, error) {
			min := time.Date(2021, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
			max := time.Date(2025, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
			delta := max - min

			sec := rand.Int63n(delta) + min
			randDate := time.Unix(sec, 0).UnixNano()
			return randDate, nil
		},
	}
}

// FakeAvatarGen generates and writes random user / project / org letter avatar
// Outputs the filepath to the generated image
func FakeAvatarGen() (string, gofakeit.Info) {
	return "avatargen", gofakeit.Info{
		Category:    "avatargen",
		Description: "generate random user / project / org avatar",
		Example:     "avatargen:<OUTPUT_DIR>",
		Output:      "string",
		Params: []gofakeit.Param{
			{Field: "outputdir", Type: "string", Description: "output directory for the generated avatars"},
			{Field: "size", Type: "int", Description: "generated avatar size", Default: "128"},
		},
		Call: func(m *map[string][]string, info *gofakeit.Info) (interface{}, error) {
			outDir, err := info.GetString(m, "outputdir")
			if err != nil {
				return nil, err
			}
			size, err := info.GetInt(m, "size")
			if err != nil {
				return nil, err
			}
			return GenFakeLogo(outDir, size)
		},
	}
}

// FakeAvatarGenBytes generates and writes random user / project / org letter avatar
// Outputs the filepath to the generated image
func FakeAvatarGenBytes() (string, gofakeit.Info) {
	return "logogen", gofakeit.Info{
		Category:    "logogen",
		Description: "generate random user / project / org avatar",
		Example:     "logogen:<SIZE>",
		Output:      "string",
		Params: []gofakeit.Param{
			{Field: "size", Type: "int", Description: "generated logo size", Default: "100"},
		},
		Call: func(m *map[string][]string, info *gofakeit.Info) (interface{}, error) {
			size, err := info.GetInt(m, "size")
			if err != nil {
				return nil, err
			}
			return genFakeLogoBase64(size)
		},
	}
}

func GenFakeLogo(outDir string, size int) (string, error) {
	if ex, _ := utilities.PathExists(outDir); !ex {
		_ = os.MkdirAll(outDir, 0755)
	}
	imgId := utilities.NewID()
	b, err := GenFakeLogoBytes(size)
	if err != nil {
		return "", err
	}
	filename := filepath.Join(outDir, fmt.Sprintf("%s.png", imgId))
	if err = ioutil.WriteFile(filename, b, 0644); err != nil {
		return "", err
	}
	return filename, nil
}

func GenFakeLogoBytes(size int) ([]byte, error) {
	img, err := identicon.Make(size, randRGB(100), randRGB(255), []byte(randString(18)))
	if err != nil {
		return nil, err
	}
	b := bytes.Buffer{}
	if err = png.Encode(&b, img); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func genFakeLogoBase64(size int) (string, error) {
	b, err := GenFakeLogoBytes(size)
	if err != nil {
		return "", err
	}
	return base64.RawStdEncoding.EncodeToString(b), nil
}

func UnmarshalFromFilepath(path string, any interface{}) error {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	ext := filepath.Ext(path)
	switch ext {
	case ".json":
		if err = utilities.UnmarshalJson(f, any); err != nil {
			return err
		}
	case ".yml", ".yaml":
		if err = utilities.UnmarshalYAML(f, any); err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid format specified, cannot load bootstrap")
	}
	return nil
}

// getRandomColorInRgb Returns a random RGBColor
func randRGB(colorRange int) color.NRGBA {
	rand.Seed(utilities.CurrentTimestamp())
	red := uint8(rand.Intn(colorRange))
	green := uint8(rand.Intn(colorRange))
	blue := uint8(rand.Intn(colorRange))
	return color.NRGBA{R: red, G: green, B: blue, A: 255}
}

func randString(n int) string {
	var src = rand.NewSource(utilities.CurrentTimestamp())
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}
