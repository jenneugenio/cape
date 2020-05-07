module github.com/capeprivacy/cape

go 1.14

require (
	cloud.google.com/go v0.44.3 // indirect
	github.com/99designs/gqlgen v0.11.2
	github.com/Masterminds/semver v1.5.0
	github.com/Masterminds/sprig v2.22.0+incompatible
	github.com/NYTimes/gziphandler v1.1.1
	github.com/badoux/checkmail v0.0.0-20181210160741-9661bd69e9ad
	github.com/bombsimon/wsl/v2 v2.2.0 // indirect
	github.com/chzyer/readline v0.0.0-20180603132655-2972be24d48e
	github.com/fatih/color v1.9.0 // indirect
	github.com/felixge/httpsnoop v1.0.1
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/gofrs/uuid v3.2.0+incompatible
	github.com/golang/protobuf v1.4.0
	github.com/golangci/gocyclo v0.0.0-20180528144436-0a533e8fa43d // indirect
	github.com/golangci/golangci-lint v1.24.0
	github.com/golangci/revgrep v0.0.0-20180812185044-276a5c0a1039 // indirect
	github.com/gorilla/websocket v1.4.1 // indirect
	github.com/gostaticanalysis/analysisutil v0.0.3 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.1-0.20190118093823-f849b5445de4
	github.com/imdario/mergo v0.3.9 // indirect
	github.com/jackc/pgconn v1.3.2
	github.com/jackc/pgproto3/v2 v2.0.1
	github.com/jackc/pgtype v1.2.0
	github.com/jackc/pgx/v4 v4.4.1
	github.com/jackc/tern v1.9.1
	github.com/jirfag/go-printf-func-name v0.0.0-20200119135958-7558a9eaa5af // indirect
	github.com/juju/ansiterm v0.0.0-20180109212912-720a0952cc2a
	github.com/justinas/alice v1.2.0
	github.com/leekchan/gtf v0.0.0-20190214083521-5fba33c5b00b
	github.com/lunixbochs/vtclean v1.0.0 // indirect
	github.com/machinebox/graphql v0.2.2
	github.com/magefile/mage v1.9.0
	github.com/manifoldco/go-base32 v1.0.4
	github.com/manifoldco/go-base64 v1.0.3
	github.com/manifoldco/healthz v1.2.0
	github.com/manifoldco/promptui v0.7.0
	github.com/marianogappa/sqlparser v0.0.0-20190512194142-a82c3f44d2fc
	github.com/markbates/pkger v0.15.1
	github.com/matryer/is v1.3.0 // indirect
	github.com/mattn/go-colorable v0.1.6 // indirect
	github.com/mitchellh/mapstructure v1.2.2
	github.com/mitchellh/reflectwalk v1.0.1 // indirect
	github.com/onsi/gomega v1.9.0
	github.com/rs/cors v1.7.0
	github.com/rs/zerolog v1.18.0
	github.com/securego/gosec v0.0.0-20200316084457-7da9f46445fd // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/viper v1.6.2 // indirect
	github.com/urfave/cli/v2 v2.1.1
	github.com/vektah/gqlparser/v2 v2.0.1
	go.opencensus.io v0.22.2 // indirect
	go.uber.org/multierr v1.1.0
	golang.org/x/crypto v0.0.0-20200414173820-0848c9571904
	golang.org/x/net v0.0.0-20200301022130-244492dfa37a // indirect
	golang.org/x/sys v0.0.0-20200323222414-85ca7c5b95cd // indirect
	google.golang.org/grpc v1.28.1
	google.golang.org/protobuf v1.21.0
	gopkg.in/ini.v1 v1.55.0 // indirect
	gopkg.in/square/go-jose.v2 v2.4.1
	helm.sh/helm/v3 v3.2.0
	mvdan.cc/unparam v0.0.0-20200314162735-0ac8026f7d06 // indirect
	rsc.io/letsencrypt v0.0.3 // indirect
	sigs.k8s.io/kind v0.8.0
	sigs.k8s.io/yaml v1.2.0
	sourcegraph.com/sqs/pbtypes v1.0.0 // indirect
)

replace github.com/grpc-ecosystem/go-grpc-middleware => github.com/capeprivacy/go-grpc-middleware v1.0.1-0.20200421173811-abd58a9536e9

// The following can be removed once helm/helm and all of its dependencies support go modules.
replace github.com/docker/docker => github.com/moby/moby v0.7.3-0.20190826074503-38ab9da00309

replace github.com/Azure/go-autorest => github.com/Azure/go-autorest v13.3.3+incompatible
