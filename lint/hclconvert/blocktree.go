package hclconvert

type blockTree struct {
	name     string
	subTrees []*blockTree
	values   []interface{}
}

func NewBlockTree(name string) *blockTree {
	bt := &blockTree{}
	bt.name = name
	bt.subTrees = []*blockTree{}
	bt.values = make([]interface{}, 0)
	return bt
}

func (bt *blockTree) appendSubTree(s *blockTree) {
	bt.subTrees = append(bt.subTrees, s)
}

func (bt *blockTree) appendValue(v interface{}) {
	bt.values = append(bt.values, v)
}

func (bt *blockTree) out() map[string][]interface{} {
	out := make(map[string][]interface{})
	if len(bt.subTrees) > 0 {
		out[bt.name] = make([]interface{}, len(bt.subTrees))
		for i, s := range bt.subTrees {
			out[bt.name][i] = s.out()
		}
		return out

	} else if len(bt.values) > 0 {
		out[bt.name] = make([]interface{}, len(bt.values))
		for i, v := range bt.values {
			out[bt.name][i] = v
		}
		return out
	}
	return nil
}

func mergeSliceMap(m1 map[string][]interface{}, m2 map[string][]interface{}) map[string][]interface{} {
	result := map[string][]interface{}{}
	for k, v := range m1 {
		result[k] = v
	}
	for k, v := range m2 {
		if s, ok := result[k]; ok {
			result[k] = append(s, v...)
		} else {
			result[k] = v
		}
	}
	return result
}
