module github.com/aperturerobotics/cli/cmd/urfave-cli-genflags

go 1.25

replace github.com/aperturerobotics/cli => ../../

require (
	github.com/aperturerobotics/cli v1.0.2-0.20260131035933-6db6a670406d
	golang.org/x/text v0.33.0
	gopkg.in/yaml.v3 v3.0.1
)

require github.com/xrash/smetrics v0.0.0-20250705151800-55b8f293f342 // indirect
