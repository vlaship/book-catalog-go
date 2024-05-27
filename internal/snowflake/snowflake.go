package snowflake

import "github.com/bwmarrin/snowflake"

type IDGenerator interface {
	Generate() snowflake.ID
}
