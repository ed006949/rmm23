package mod_strings

var (
	FVEnclosure = map[string][2]string{
		RedisearchTagTypeText:    {enclosureEmpty0, enclosureEmpty1},
		RedisearchTagTypeTag:     {enclosureCurly0, enclosureCurly1},
		RedisearchTagTypeNumeric: {enclosureSquare0, enclosureSquare1},
		RedisearchTagTypeGeo:     {enclosureSquare0, enclosureSquare1},
	}
)
