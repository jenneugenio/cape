package transformations

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"index/suffixarray"
	"math/rand"
	"regexp"
	"sort"
	"strings"
	"unicode"

	"github.com/capeprivacy/cape/connector/proto"
	errors "github.com/capeprivacy/cape/partyerrors"
	"github.com/jdkato/prose/tokenize"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type ScramblerTransform struct {
	field                         string
	key                           []byte
	whiteSpaceLeakageInstance     *whiteSpaceLeakage
	capitalizationLeakageInstance *capitalizationLeakage
}

func (s *ScramblerTransform) scrambleText(text string) (string, error) {
	textLowerCased := strings.ToLower(text)

	tokens := tokenize.TextToWords(textLowerCased)
	hashedTokens := make([]string, len(tokens))

	for i, token := range tokens {
		if isPunct(token) {
			hashedTokens[i] = token
		} else {
			var err error
			hashedTokens[i], err = s.hash(token)
			if err != nil {
				return "", err
			}
		}
	}

	tokensStitched := strings.Join(hashedTokens[:], "")
	output := s.whiteSpaceLeakageInstance.apply(text, tokensStitched)
	output = s.capitalizationLeakageInstance.apply(text, output)

	return output, nil
}

func (s *ScramblerTransform) Transform(schema *proto.Schema, input *proto.Record) error {
	field, err := GetField(schema, input, s.field)
	if err != nil {
		return err
	}

	output := &proto.Field{}
	switch val := field.GetValue().(type) {
	case *proto.Field_String_:
		res, err := s.scrambleText(val.String_)
		if err != nil {
			return errors.New(UnsupportedType, "Attempted to call %s transform on an unsupported type %T", s.Function(), val)
		}
		output.Value = &proto.Field_String_{String_: res}
	}

	return SetField(schema, input, output, s.field)
}

func (s *ScramblerTransform) Initialize(args Args) error {
	key := make([]byte, keyLen)
	_, err := rand.Read(key)
	if err != nil {
		return err
	}
	s.key = key

	whiteSpacePattern := map[string]string{
		"\\ ": " ",
		"\n":  "\n",
		"\t":  "\t",
		"\r":  "\r",
	}

	s.whiteSpaceLeakageInstance = &whiteSpaceLeakage{pattern: whiteSpacePattern}
	s.capitalizationLeakageInstance = &capitalizationLeakage{}

	return nil
}

func (s *ScramblerTransform) Validate(args Args) error {
	return nil
}

func (s *ScramblerTransform) SupportedTypes() []proto.FieldType {
	return []proto.FieldType{
		proto.FieldType_CHAR,
		proto.FieldType_TEXT,
		proto.FieldType_VARCHAR,
	}
}

func (s *ScramblerTransform) Function() string {
	return "scrambler"
}

func (s *ScramblerTransform) Field() string {
	return s.field
}

func NewScramblerTransform(field string) (Transformation, error) {
	s := &ScramblerTransform{field: field}
	return s, nil
}

func isPunct(token string) bool {
	unicodeToken := []rune(token)
	return unicode.IsPunct(unicodeToken[0])
}

func (s *ScramblerTransform) hash(token string) (string, error) {
	h := hmac.New(sha256.New, s.key)
	_, err := h.Write([]byte(token))
	if err != nil {
		return "", err
	}

	seed := int(binary.BigEndian.Uint64(h.Sum(nil)))
	tokenLen := len(token)
	hashedToken := sampleString(tokenLen, int64(seed))
	return hashedToken, nil
}

func sampleString(stringLen int, seed int64) string {
	seededRand := rand.New(rand.NewSource(seed))
	randomString := make([]byte, stringLen)
	for i := range randomString {
		randomString[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(randomString)
}

type whiteSpaceLeakage struct {
	pattern map[string]string
}

func (w *whiteSpaceLeakage) apply(originalText string, scrambledText string) string {
	idx := suffixarray.New([]byte(originalText))
	registerSpaces := make(map[int]string)

	for spacePattern, replacement := range w.pattern {
		regex := regexp.MustCompile(spacePattern)
		spaceIdx := idx.FindAllIndex(regex, -1)

		for _, i := range spaceIdx {
			registerSpaces[i[0]] = replacement
		}
	}

	// Order Dict by key
	var whitePositions []int
	for k := range registerSpaces {
		whitePositions = append(whitePositions, k)
	}
	sort.Ints(whitePositions)

	for _, position := range whitePositions {
		scrambledText = scrambledText[:position] + registerSpaces[position] + scrambledText[position:]
	}

	return scrambledText
}

type capitalizationLeakage struct {
}

func (c *capitalizationLeakage) apply(originalText string, scrambledText string) string {
	scrambledText = strings.ToLower(scrambledText)

	for position, c := range originalText {
		if unicode.IsUpper(c) {
			char := string(scrambledText[position])
			scrambledText = scrambledText[:position] + strings.ToUpper(char) + scrambledText[position+1:]
		}
	}
	return scrambledText
}
