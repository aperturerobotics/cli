# urfave-cli-genflags

This is a tool that is used internally to generate flag types and methods from a
YAML input. It intentionally pins usage of `github.com/aperturerobotics/cli` to
a *release* rather than using the adjacent code so that changes don't result in
*this* tool refusing to compile. It's almost like dogfooding?
