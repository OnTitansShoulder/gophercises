package parser

import (
	"encoding/json"

	"github.com/ontitansshoulder/textgame/story"
	"github.com/pkg/errors"
)

func Parse(jsonBytes []byte) (map[string]story.Arc, error) {
	var arcs map[string]story.Arc
	err := json.Unmarshal(jsonBytes, &arcs)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal json")
	}

	for k, v := range arcs {
		v.Name = k
		arcs[k] = v
	}

	return arcs, nil
}
