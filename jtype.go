package json2

const (
	ARR JTYPE = 1 << iota
	STR
	BOOL
	NUM
	NULL
	OBJ
)

var (
	JT = map[JTYPE]string{
		ARR:        "ARR",
		STR:        "STR",
		BOOL:       "BOOL",
		NUM:        "NUM",
		NULL:       "NULL",
		OBJ:        "OBJ",
		ARR | STR:  "ARR_STR",
		ARR | BOOL: "ARR_BOOL",
		ARR | NUM:  "ARR_NUM",
		ARR | NULL: "ARR_NULL",
		ARR | OBJ:  "ARR_OBJ",
	}
)

// Str : JSON Type string
func (jt JTYPE) Str() string {
	return JT[jt]
}

// IsArr : is json array type
func (jt JTYPE) IsArr() bool {
	return jt&ARR == ARR
}

// IsObj : is json object type
func (jt JTYPE) IsObj() bool {
	return jt&OBJ == OBJ && jt&ARR != ARR
}

// IsObjArr : is json object array type
func (jt JTYPE) IsObjArr() bool {
	return jt&OBJ == OBJ && jt&ARR == ARR
}

// IsPrimitive : is json primitive type
func (jt JTYPE) IsPrimitive() bool {
	return jt&OBJ != OBJ && jt&ARR != ARR
}

// IsLeafValue : is json Primitive OR Primitive array
func (jt JTYPE) IsLeafValue() bool {
	return jt&OBJ != OBJ
}
