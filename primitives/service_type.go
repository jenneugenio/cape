package primitives

import (
	"fmt"
	"io"
	"strconv"

	errors "github.com/capeprivacy/cape/partyerrors"
)

// ServiceType enum holding the supported service types
type ServiceType string

var (
	// UserServiceType is the user service type
	UserServiceType ServiceType = "user"

	// DataConnectorServiceType is the data connector service type
	DataConnectorServiceType ServiceType = "data-connector"

	WorkerServiceType ServiceType = "worker"
)

var typeRegistry map[ServiceType]string

func init() {
	typeRegistry = map[ServiceType]string{
		UserServiceType:          UserServiceType.String(),
		DataConnectorServiceType: DataConnectorServiceType.String(),
		WorkerServiceType:        WorkerServiceType.String(),
	}
}

// ServiceTypes returns a map of a type to string representation
func ServiceTypes() map[ServiceType]string {
	return typeRegistry
}

// NewServiceType validates the input and returns a new ServiceType
func NewServiceType(typ string) (ServiceType, error) {
	s := ServiceType(typ)
	err := s.Validate()

	return s, err
}

// String returns the string represented by the enum value
func (s *ServiceType) String() string {
	return string(*s)
}

// Validate checks to see if the service type is valid
func (s *ServiceType) Validate() error {
	switch *s {
	case UserServiceType:
		return nil
	case DataConnectorServiceType:
		return nil
	case WorkerServiceType:
		return nil
	default:
		return errors.New(InvalidServiceType, "%s is not a valid ServiceType", *s)
	}
}

// UnmarshalGQL unmarshals a string in the CredentialsAlgType enum
func (s *ServiceType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return errors.New(InvalidServiceType, "Cannot unmarshal ServiceType")
	}

	*s = ServiceType(str)

	return s.Validate()
}

// MarshalGQL marshals a CredentailsAlgType enum to string
func (s ServiceType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(s.String()))
}
