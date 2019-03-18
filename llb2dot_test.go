package llb2dot_test

import (
	"bytes"
	"testing"

	"github.com/moby/buildkit/solver/pb"
	"github.com/po3rin/llb2dot"
)

func TestLLB2Graph(t *testing.T) {
	tests := []struct {
		input llb2dot.LLBOps
		want  []byte
	}{
		{
			input: llb2dot.LLBOps{
				llb2dot.LLBOp{
					Digest: "aaaaaaaaaaaa",
				},
				llb2dot.LLBOp{
					Op: pb.Op{
						Inputs: []*pb.Input{
							&pb.Input{
								Digest: "aaaaaaaaaaaa",
							},
						},
					},
					Digest: "bbbbbbbbbbbb",
				},
			},
			want: []byte(`strict digraph llb {
// Node definitions.
"aaaaaaaaaa..." [digest=aaaaaaaaaaaa];
"bbbbbbbbbb..." [digest=bbbbbbbbbbbb];

// Edge definitions.
"aaaaaaaaaa..." -> "bbbbbbbbbb...";
}`),
		},
	}
	for _, tt := range tests {
		g, err := llb2dot.LLB2Graph(tt.input)
		if err != nil {
			t.Fatalf("Unexpected error : %+v", err)
		}

		var got bytes.Buffer
		err := llb2dot.WriteDOT(&got, g)
		if err != nil {
			t.Fatalf("Unexpected error : %+v", err)
		}

		if string(got.String()) != string(tt.want) {
			t.Fatalf("failed test. want: %+v, got: %+v", string(tt.want), got.String())
		}
	}
}
