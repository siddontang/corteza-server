package gig

import (
	"testing"

	"github.com/cortezaproject/corteza-server/pkg/gig"
)

func Test_gig_create(t *testing.T) {
	var (
		ctx, svc, h, s = setup(t)
		err            error
		g              gig.Gig
	)
	_ = s
	_ = svc
	_ = g

	t.Run("no worker defined", func(_ *testing.T) {
		_, err = svc.Create(ctx, gig.UpdatePayload{})
		h.a.Error(err)
	})

	t.Run("complete payload", func(_ *testing.T) {
		g, err = svc.Create(ctx, gig.UpdatePayload{
			Worker: gig.WorkerNoop(),
			Decode: gig.DecoderSet{
				decoderSafe(gig.DecoderNoopSource(0)),
			},
			Preprocess: gig.PreprocessorSet{
				gig.PreprocessorNoop(),
			},
			Postprocess: gig.PostprocessorSet{
				gig.PostprocessorNoop(),
			},
			Sources: []gig.SourceWrap{
				{
					Uri:  testSource(t, "secret-key.txt"),
					Name: "secret-key.txt",
				},
			},
		})
		h.a.NoError(err)

		h.a.NotEqual(0, g.ID)
		h.a.False(g.CreatedAt.IsZero())

		h.a.Len(g.Sources, 1)
		h.a.Len(g.Sources[0].Decoders(), 1)
		h.a.IsType(decoderSafe(gig.DecoderNoopSource(0)), g.Sources[0].Decoders()[0])

		h.a.Len(g.Preprocess, 1)
		pre := gig.PreprocessorNoop()
		h.a.IsType(pre, g.Preprocess[0])

		h.a.Len(g.Postprocess, 1)
		post := gig.PostprocessorNoop()
		h.a.IsType(post, g.Postprocess[0])
	})
}