module github.com/aperturerobotics/cli/cmd/urfave-cli-genflags

go 1.24

replace github.com/aperturerobotics/cli => ../../

require (
	github.com/aperturerobotics/cli v1.0.0
	golang.org/x/text v0.23.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/aperturerobotics/common v0.21.2 // indirect
	github.com/xrash/smetrics v0.0.0-20240521201337-686a1a2994c1 // indirect
)
