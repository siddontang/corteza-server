package gig

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"fmt"
	"github.com/spf13/cast"
)

type (
	decoderNoop struct {
		source uint64
	}
	decoderArchive struct {
		source uint64
	}
)

const (
	DecoderHandleNoop    = "noop"
	DecoderHandleArchive = "archive"
)

// ------------------------------------------------------------------------
// Constructors and utils

// DecoderNoopParams returns a new decoderNoop from the params
func DecoderNoopParams(params map[string]interface{}) (decoderNoop, error) {
	var (
		out = decoderNoop{}
		err error
	)

	// Param validation
	// - supported params
	index := map[string]bool{
		"source": true,
	}
	for p := range params {
		if !index[p] {
			return out, fmt.Errorf("unknown parameter provided to noop: %s", p)
		}
	}

	// Fill and check requirements
	out.source = cast.ToUint64(params["source"])
	return out, err
}

// DecoderNoopSource returns a new decoderNoop from the required fields and source
func DecoderNoopSource(source uint64) (decoderNoop, error) {
	var (
		err error
		out decoderNoop
	)
	out = decoderNoop{
		source: source,
	}

	return out, err
}

func (t decoderNoop) Ref() string {
	return DecoderHandleNoop
}

func (t decoderNoop) Params() map[string]interface{} {
	return map[string]interface{}{
		"source": t.source,
	}
}

// DecoderArchiveParams returns a new decoderArchive from the params
func DecoderArchiveParams(params map[string]interface{}) (decoderArchive, error) {
	var (
		out = decoderArchive{}
		err error
	)

	// Param validation
	// - supported params
	index := map[string]bool{
		"source": true,
	}
	for p := range params {
		if !index[p] {
			return out, fmt.Errorf("unknown parameter provided to archive: %s", p)
		}
	}

	// Fill and check requirements
	out.source = cast.ToUint64(params["source"])
	return out, err
}

// DecoderArchiveSource returns a new decoderArchive from the required fields and source
func DecoderArchiveSource(source uint64) (decoderArchive, error) {
	var (
		err error
		out decoderArchive
	)
	out = decoderArchive{
		source: source,
	}

	return out, err
}

func (t decoderArchive) Ref() string {
	return DecoderHandleArchive
}

func (t decoderArchive) Params() map[string]interface{} {
	return map[string]interface{}{
		"source": t.source,
	}
}

// ------------------------------------------------------------------------
// Task registry

func decoderDefinitions() TaskDefSet {
	return TaskDefSet{
		{
			Ref:         DecoderHandleNoop,
			Kind:        TaskDecoder,
			Description: "Noop does nothing.",
			Params: []taskDefParam{
				{
					Name:     "source",
					Kind:     "String",
					Required: false,
				},
			},
		},
		{
			Ref:         DecoderHandleArchive,
			Kind:        TaskDecoder,
			Description: "Extracts the contents of the archive into sepparate sources; extraction is not recursive.",
			Params: []taskDefParam{
				{
					Name:     "source",
					Kind:     "String",
					Required: false,
				},
			},
		},
	}
}