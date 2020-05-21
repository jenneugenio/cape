package mage

import (
	"context"
	"fmt"
	"regexp"

	"github.com/Masterminds/semver"
	"github.com/magefile/mage/sh"
	"github.com/magefile/mage/target"
)

var protocVersionRegex = regexp.MustCompile(`libprotoc (([0-9]+\.?)*)`)

// Protoc is a dependency checker and generator for Protoc
type Protoc struct {
	Version    *semver.Version
	Configfile string
	Protofile  string
	Outputfile string
}

// NewProtoc returns a struct for managing and working with Protoc locally
func NewProtoc(required, cfgFile, protoFile, outputFile string) (*Protoc, error) {
	v, err := semver.NewVersion(required)
	if err != nil {
		return nil, err
	}

	return &Protoc{
		Version:    v,
		Configfile: cfgFile,
		Protofile:  protoFile,
		Outputfile: outputFile,
	}, nil
}

// MustProtoc returns a struct for managing and working with Protoc locally and panics on any error
func MustProtoc(required, cfgFile, protoFile, outputFile string) *Protoc {
	return &Protoc{
		Version:    semver.MustParse(required),
		Configfile: cfgFile,
		Protofile:  protoFile,
		Outputfile: outputFile,
	}
}

// Check returns an error if Protoc isn't available or the version is incorrect
func (p *Protoc) Check(_ context.Context) error {
	out, err := sh.Output("protoc", "--version")
	if err != nil {
		return err
	}

	matches := protocVersionRegex.FindStringSubmatch(out)
	if len(matches) != 3 {
		return fmt.Errorf("Could not parse output of `protoc --version`")
	}

	v, err := semver.NewVersion(matches[1])
	if err != nil {
		return fmt.Errorf("Could not parse output of `protoc --version`: %s", err.Error())
	}

	if v.LessThan(p.Version) {
		return fmt.Errorf("Please upgrade your version of Protoc from %s to %s or greater", v.String(), p.Version.String())
	}

	return nil
}

func (p *Protoc) Name() string {
	return "protoc"
}

func (p *Protoc) Generate(ctx context.Context) error {
	needsGeneration, err := target.Path(p.Outputfile, p.Protofile, p.Configfile)
	if err != nil {
		return err
	}

	if !needsGeneration {
		return nil
	}

	return sh.Run("go", "generate", p.Configfile)
}

func (p *Protoc) Setup(ctx context.Context) error {
	// Protoc is an external dependency, the user must have it installed with
	// the right version in their environment
	return p.Check(ctx)
}

func (p *Protoc) Clean(_ context.Context) error {
	return nil
}
