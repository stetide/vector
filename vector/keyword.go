package vector

type keyWord struct {
	name  string
	alias []string
}

func (k keyWord) getNameByAlias(al string) string {
	for _, a := range k.alias {
		if a == al {
			return a
		}
	}
	return k.name
}

var (
	kwVEC    = keyWord{name: "vec"}
	kwQUIT   = keyWord{name: "quit", alias: []string{"end", "exit", "close"}}
	kwCLEAR  = keyWord{name: "clear", alias: []string{"cls"}}
	kwHELP   = keyWord{name: "help"}
	kwANS    = keyWord{name: "ans"}
	kwEXPORT = keyWord{name: "export", alias: []string{"save"}}
)

var keywords = []keyWord{
	kwVEC,
	kwQUIT,
	kwCLEAR,
	kwHELP,
	kwANS,
	kwEXPORT,
}

func isKeyword(str string) bool {
	for _, kw := range keywords {
		if str == kw.name {
			return true
		}
		for _, kwa := range kw.alias {
			if str == kwa {
				return true
			}
		}
	}
	return false
}
