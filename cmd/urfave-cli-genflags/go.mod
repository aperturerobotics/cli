module github.com/aperturerobotics/cli/cmd/urfave-cli-genflags

go 1.25

replace github.com/aperturerobotics/cli => ../../

require (
	github.com/aperturerobotics/cli v1.0.0
	golang.org/x/text v0.30.0
	gopkg.in/yaml.v3 v3.0.1
)

require github.com/xrash/smetrics v0.0.0-20250705151800-55b8f293f342 // indirect
