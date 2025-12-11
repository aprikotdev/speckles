package cfg

import (
	pb "github.com/aprikotdev/speckles/cfg/gen/specs/v1"
)

var Default = &pb.Namespaces{
	Namespaces: []*pb.Namespace{
		HTML,
		SVG,
		MathML,
	},
}
