package transformations

import (
	"testing"

	"github.com/capeprivacy/cape/connector/proto"
	gm "github.com/onsi/gomega"
)

var scrambSchema = &proto.Schema{
	Fields: []*proto.FieldInfo{
		{
			Field: proto.FieldType_TEXT,
			Name:  "text",
			Size:  500,
		},
	},
}

func TestWhiteLeakage(t *testing.T) {
	gm.RegisterTestingT(t)

	whiteSpacePattern := map[string]string{
		"\\ ": " ",
		"\n":  "\n",
		"\t":  "\t",
	}

	leakage := &whiteSpaceLeakage{pattern: whiteSpacePattern}
	originalString := "xxx\nx x\txx  x "
	scrambledString := "xxxxxxxx"
	actualOutputString := leakage.apply(originalString, scrambledString)
	gm.Expect(actualOutputString).To(gm.Equal(originalString))
}

func TestCapitalizationLeakage(t *testing.T) {
	gm.RegisterTestingT(t)

	leakage := &capitalizationLeakage{}
	originalString := "XxxxxXXxxX"
	scrambledString := "xxxxxxxxxx"
	actualOutputString := leakage.apply(originalString, scrambledString)
	gm.Expect(actualOutputString).To(gm.Equal(originalString))
}

func TestScrambler(t *testing.T) {
	gm.RegisterTestingT(t)

	transform, err := NewScramblerTransform("text")
	gm.Expect(err).To(gm.BeNil())
	var args Args = nil

	err = transform.Validate(args)
	gm.Expect(err).To(gm.BeNil())

	err = transform.Initialize(args)
	gm.Expect(err).To(gm.BeNil())

	transform.(*ScramblerTransform).key = []byte("secret")
	shakespeare := "First Citizen:\nBefore we proceed any further, hear me speak.\n\n" +
		"All:\nSpeak, speak.\n\nFirst Citizen:\nYou are all resolved rather to die than to " +
		"famish?\n\nAll:\nResolved.resolved.\n\nFirst Citizen:\nFirst, you know Caius " +
		"Marcius is chief enemy to the people.\n\nAll:\nWe know't, we know't."

	shakespeareScrambled := "Puxfa Poiyntv:\nNmppgq iv hmcihye rob teokoey, comu jf " +
		"yunma.\n\nIsb:\nYunma, yunma.\n\nPuxfa Poiyntv:\nNyb bqb isb escgzmzi pekyni " +
		"pj dye azbd pj sssyhi?\n\nIsb:\nSomtqxvtqiyhquhdye\n\nPuxfa Poiyntv:\nPuxfa, nyb " +
		"ecic Gaizg Qrccfqq tl uzmpo ivjzk pj lul edwqgh.\n\nIsb:\nIv zfigfz, iv zfigfz."

	inputField := &proto.Field{Value: &proto.Field_String_{String_: shakespeare}}
	expectedOutputField := &proto.Field{Value: &proto.Field_String_{String_: shakespeareScrambled}}
	record := &proto.Record{Fields: []*proto.Field{inputField}, Schema: scrambSchema}
	expectedRecord := &proto.Record{Fields: []*proto.Field{expectedOutputField}, Schema: scrambSchema}

	err = transform.Transform(scrambSchema, record)
	gm.Expect(err).To(gm.BeNil())

	gm.Expect(record).To(gm.Equal(expectedRecord))
}
