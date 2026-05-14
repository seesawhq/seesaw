package models

type Target struct {
	FlagID int64
	Flag   Flag
	Type   string
	data   string
}

// { segmentID: "1", condition: "IS_IN", "serve_type": "percentage",  "server_data": {} }
// {"1": "30", "2":"40", "3": 20}
// { condition: "IS_IN", values: "aaa,aaa,aaaa", "serve_type": "vairent",  "server_data": {} }
// {"value": "1"}

type SegmentTarget struct {
	TagetID   int64
	Condition string
}

type IndiviualTarget struct {
	TagetID   int64
	Condition string
	Values    string
}
