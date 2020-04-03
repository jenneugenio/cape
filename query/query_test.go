package query

import (
	"bufio"
	"fmt"
	"github.com/dropoutlabs/cape/primitives"
	gm "github.com/onsi/gomega"
	"io/ioutil"
	"os"
	"testing"
)

type QueryFixture struct {
	input    string
	expected string
}

func loadPolicySpec(file string) (*primitives.PolicySpec, error) {
	f, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return primitives.ParsePolicySpec(f)
}

func loadQueries(file string) ([]*QueryFixture, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var inputs []string
	var outputs []string

	writingTo := &inputs

	for scanner.Scan() {
		fmt.Println("Reading", scanner.Text())

		if len(scanner.Text()) == 0 {
			continue
		}

		if scanner.Text() == "--- input above, output below ---" {
			writingTo = &outputs
			continue
		}

		*writingTo = append(*writingTo, scanner.Text())
	}

	fixtures := make([]*QueryFixture, len(inputs))
	for i := 0; i < len(inputs); i++ {
		fixtures[i] = &QueryFixture{
			input:    inputs[i],
			expected: outputs[i],
		}
	}

	return fixtures, nil
}

type Fixture struct {
	ps      *primitives.PolicySpec
	queries []*QueryFixture
}

func loadFixtures() ([]*Fixture, error) {
	files, err := ioutil.ReadDir("./testdata")
	if err != nil {
		return nil, err
	}

	if len(files)%2 != 0 {
		panic("There must be an even number of files!")
	}

	var fixtures []*Fixture

	for i := 0; i < len(files)/2; i++ {
		var fileIdent string
		if i >= 99 {
			fileIdent = fmt.Sprintf("%d", i+1)
		} else if i >= 9 {
			fileIdent = fmt.Sprintf("0%d", i+1)
		} else {
			fileIdent = fmt.Sprintf("00%d", i+1)
		}

		policyFile := fmt.Sprintf("./testdata/%s_policy.yaml", fileIdent)
		sqlFile := fmt.Sprintf("./testdata/%s_query.sql", fileIdent)

		queries, err := loadQueries(sqlFile)
		if err != nil {
			return nil, err
		}

		ps, err := loadPolicySpec(policyFile)
		if err != nil {
			return nil, err
		}

		fixtures = append(fixtures, &Fixture{
			ps:      ps,
			queries: queries,
		})
	}

	return fixtures, nil
}

func TestQuery(t *testing.T) {
	gm.RegisterTestingT(t)

	label, err := primitives.NewLabel("my-data")
	gm.Expect(err).To(gm.BeNil())

	t.Run("Can parse a valid query", func(t *testing.T) {
		gm.RegisterTestingT(t)

		_, err := New(label, "SELECT * from transactions")
		gm.Expect(err).To(gm.BeNil())
	})

	t.Run("Errors on an invalid query", func(t *testing.T) {
		gm.RegisterTestingT(t)

		_, err := New(label, "jdksajdksajdkldklasj")
		gm.Expect(err).ToNot(gm.BeNil())
	})

	t.Run("You cannot pass a join", func(t *testing.T) {
		gm.RegisterTestingT(t)

		_, err := New(label, "select * from transactions join people on transactions.person_id = people.id")
		gm.Expect(err).ToNot(gm.BeNil())
	})

	t.Run("Only accepts selects", func(t *testing.T) {
		gm.RegisterTestingT(t)

		_, err := New(label, "delete from transactions")
		gm.Expect(err).ToNot(gm.BeNil())
	})

	t.Run("Tracks the collection", func(t *testing.T) {
		gm.RegisterTestingT(t)

		q, err := New(label, "SELECT * from transactions")
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(q.Collection()).To(gm.Equal("transactions"))
	})
}

func TestRewriting(t *testing.T) {
	gm.RegisterTestingT(t)
	type TestCase struct {
		// test name
		name string

		// rule/policy info
		target string
		effect primitives.Effect
		fields []string
		where  map[string]string

		// input & expected query
		input    string
		expected string
	}

	testCases := []*TestCase{
		{
			"It redacts a field you cannot access",
			"records:mycollection.transactions",
			primitives.Deny,
			[]string{"processor"},
			nil,
			"SELECT processor, card_number, value FROM transactions",
			"SELECT card_number, value FROM transactions",
		},

		{
			"It can give you access to only things you can have",
			"records:mycollection.transactions",
			primitives.Allow,
			[]string{"processor"},
			nil,
			"SELECT processor, card_number, value FROM transactions",
			"SELECT processor FROM transactions",
		},

		{
			"It can rewrite a star command",
			"records:mycollection.transactions",
			primitives.Allow,
			[]string{"processor", "card_number", "processor"},
			nil,
			"SELECT * FROM transactions",
			"SELECT processor, card_number, processor FROM transactions",
		},

		{
			"It can filter based on row",
			"records:mycollection.transactions",
			primitives.Allow,
			[]string{"card_number"},
			map[string]string{
				"processor": "visa",
			},
			"SELECT * FROM transactions",
			"SELECT card_number FROM transactions WHERE processor = 'visa'",
		},

		{
			"It can filter multiple conditions",
			"records:mycollection.transactions",
			primitives.Allow,
			[]string{"card_number"},
			map[string]string{
				"processor": "visa",
				"vendor":    "Cool Shirts Inc.",
			},
			"SELECT * FROM transactions",
			"SELECT card_number FROM transactions WHERE processor = 'visa' AND vendor = 'Cool Shirts Inc.'",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gm.RegisterTestingT(t)

			target, err := primitives.NewTarget(tc.target)
			gm.Expect(err).To(gm.BeNil())

			fields := make([]primitives.Field, len(tc.fields))
			for i := 0; i < len(fields); i++ {
				f, err := primitives.NewField(tc.fields[i])
				gm.Expect(err).To(gm.BeNil())

				fields[i] = f
			}

			r := primitives.Rule{
				Target: target,
				Action: primitives.Read,
				Effect: tc.effect,
				Fields: fields,
				Where: []primitives.Where{
					tc.where,
				},
			}

			label := primitives.Label("cool-rule")
			spec := &primitives.PolicySpec{
				Version: 1,
				Label:   label,
				Rules:   []primitives.Rule{r},
			}

			p, err := primitives.NewPolicy(label, spec)
			gm.Expect(err).To(gm.BeNil())

			q, err := New(label, tc.input)
			gm.Expect(err).To(gm.BeNil())

			q, err = q.Rewrite(p)
			gm.Expect(err).To(gm.BeNil())

			gm.Expect(q.Raw()).To(gm.Equal(tc.expected))
		})
	}

	t.Run("Errors when you can't access anything you've asked for", func(t *testing.T) {
		gm.RegisterTestingT(t)

		label := primitives.Label("bad-policy")
		r := primitives.Rule{
			Target: "records:mycollection.transactions",
			Action: "read",
			Effect: primitives.Allow,
			Fields: []primitives.Field{"processor"},
		}

		spec := &primitives.PolicySpec{
			Version: 1,
			Label:   label,
			Rules:   []primitives.Rule{r},
		}

		p, err := primitives.NewPolicy(label, spec)
		gm.Expect(err).To(gm.BeNil())

		// I've only asked for things I can't see!
		q, err := New("cool-query", "SELECT card_number, value FROM transactions")
		gm.Expect(err).To(gm.BeNil())

		_, err = q.Rewrite(p)
		gm.Expect(err.Error()).To(gm.Equal("no_possible_fields: Cannot access any requested fields"))
	})
}

func TestQueryAgainstPolicyFile(t *testing.T) {
	gm.RegisterTestingT(t)

	fixtures, err := loadFixtures()
	gm.Expect(err).To(gm.BeNil())

	for _, f := range fixtures {
		for _, q := range f.queries {
			testName := fmt.Sprintf("Checks query rewriting for policy %s & query %s", f.ps.Label, q.input)
			t.Run(testName, func(t *testing.T) {
				query, err := New(f.ps.Label, q.input)
				gm.Expect(err).To(gm.BeNil())

				policy, err := primitives.NewPolicy(f.ps.Label, f.ps)
				gm.Expect(err).To(gm.BeNil())

				query, err = query.Rewrite(policy)
				gm.Expect(err).To(gm.BeNil())

				raw := query.Raw()
				gm.Expect(raw).To(gm.Equal(q.expected))
			})
		}
	}
}
