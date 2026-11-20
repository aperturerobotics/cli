module github.com/aperturerobotics/cli/cmd/urfave-cli-genflags

go 1.26.0

replace github.com/aperturerobotics/cli => ../../

require (
	github.com/aperturerobotics/cli v1.1.0
	golang.org/x/text v0.44.0
	gopkg.in/yaml.v3 v3.0.1
)

require github.com/xrash/smetrics v0.0.0-20261013231157-39c47d93efba // indirect
