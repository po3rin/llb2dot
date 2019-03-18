package llb2dot

import (
	"io"
	"io/ioutil"

	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/client/llb/imagemetaresolver"
	"github.com/moby/buildkit/frontend/dockerfile/dockerfile2llb"
	"github.com/moby/buildkit/solver/pb"
	"github.com/moby/buildkit/util/appcontext"
	"github.com/opencontainers/go-digest"
	"github.com/pkg/errors"
)

// LLBOps is LLBOp slice.
type LLBOps []LLBOp

// LLBOp has simple BuildKit LLB Op.
type LLBOp struct {
	Op         pb.Op
	Digest     digest.Digest
	OpMetadata pb.OpMetadata
}

func (llbOp LLBOp) getDesc() string {
	desc, ok := llbOp.OpMetadata.Description["com.docker.dockerfile.v1.command"]
	if !ok {
		desc, ok = llbOp.OpMetadata.Description["llb.customname"]
		if !ok {
			return string(llbOp.Digest[:10]) + "..."
		}
	}
	return desc
}

func def2Op(def *llb.Definition) ([]LLBOp, error) {
	var ops []LLBOp
	for _, dt := range def.Def {
		var op pb.Op
		if err := (&op).Unmarshal(dt); err != nil {
			return nil, errors.Wrap(err, "failed to parse op")
		}
		dgst := digest.FromBytes(dt)
		ent := LLBOp{Op: op, Digest: dgst, OpMetadata: def.Metadata[dgst]}
		ops = append(ops, ent)
	}
	return ops, nil
}

// LoadLLB received LLB from io.Reader and convert llbOp struct.
func LoadLLB(r io.Reader) ([]LLBOp, error) {
	def, err := llb.ReadFrom(r)
	if err != nil {
		return nil, err
	}
	return def2Op(def)
}

// LoadDockerfile received Dockerfile from io.Reader and convert llbOp struct.
func LoadDockerfile(r io.Reader) ([]LLBOp, error) {
	df, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	state, _, err := dockerfile2llb.Dockerfile2LLB(appcontext.Context(), df, dockerfile2llb.ConvertOpt{
		MetaResolver: imagemetaresolver.Default(),
	})
	if err != nil {
		return nil, err
	}

	def, err := state.Marshal()
	return def2Op(def)
}
