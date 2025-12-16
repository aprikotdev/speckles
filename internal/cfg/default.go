package cfg

import (
	pb "github.com/aprikotdev/speckles/internal/pb/gen/specs/v1"
)

var Default = &pb.Namespaces{
	Namespaces: []*pb.Namespace{HTML, SVG, MathML},
}
