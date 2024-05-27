package snowflake

import "github.com/bwmarrin/snowflake"

type SnowflakeIDGenerator interface {
	Generate() snowflake.ID
}
