package filter

import (
	"encoding/json"
	"github.com/pa001024/reflex/source"
	"github.com/pa001024/reflex/util"
	"text/template"
)

type IFilter interface {
	Process(src []*source.FeedInfo) (dst []*source.FeedInfo)
}

type Filter struct {
	Type string `json:"type"` // 类型 filter工厂ID 如[moegirlwiki,简繁转换]等
	Name string `json:"name"` // 名字 不能跟别的target或source名字相同
}

func New(name string, b []byte) (rst IFilter) {
	obj := &Filter{}
	err := json.Unmarshal(b, obj)
	if err != nil {
		util.ERROR.Err("JSON Parse Error", err)
		return
	}
	switch obj.Type {
	default:
	case "moegirlwiki":
		dst := &FilterMoegirlwiki{}
		json.Unmarshal(b, dst)
		dst.Name = name
		rst = dst
		util.INFO.Logf("filter.moegirlwiki \"%s\" Loaded.", name)
	case "word":
		dst := &FilterWords{}
		json.Unmarshal(b, dst)
		dst.Name = name
		rst = dst
		util.INFO.Logf("filter.word \"%s\" Loaded.", name)
	case "zhconv":
		dst := &FilterZhConv{}
		json.Unmarshal(b, dst)
		dst.Name = name
		rst = dst
		util.INFO.Logf("filter.zhconv \"%s\" Loaded.", name)
	case "basic":
		dst := &FilterBasic{}
		json.Unmarshal(b, dst)
		dst.Name = name
		dst.compFormat = template.Must(template.New(name).Parse(dst.Format))
		if dst.MaxLength == 0 {
			dst.MaxLength = 120
		}
		rst = dst
		util.INFO.Logf("filter.basic \"%s\" Loaded.", name)
	}
	return
}
