package snowflake

import "github.com/bwmarrin/snowflake"

// IDGeneratorImpl is a struct that represents a snowflake ID generator
type IDGeneratorImpl struct {
	node *snowflake.Node
}

// New creates a new snowflake ID generator
func New(nodeID int64) (*IDGeneratorImpl, error) {
	node, err := snowflake.NewNode(nodeID)
	if err != nil {
		return nil, err
	}

	return &IDGeneratorImpl{
		node: node,
	}, nil
}

// Generate generates a new snowflake ID
func (s *IDGeneratorImpl) Generate() snowflake.ID {
	return s.node.Generate()
}
