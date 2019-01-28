package topological

// Graph represents a dependency graph
type Graph struct {
	nodes   []string
	outputs map[string]map[string]int
	inputs  map[string]int
}

// NewGraph creates new Graph object
func NewGraph(cap int) *Graph {
	return &Graph{
		nodes:   make([]string, 0, cap),
		inputs:  make(map[string]int),
		outputs: make(map[string]map[string]int),
	}
}

// AddNode adds new node to Graph
func (g *Graph) AddNode(node string) bool {
	g.nodes = append(g.nodes, node)

	if _, ok := g.outputs[node]; ok {
		return false
	}
	g.outputs[node] = make(map[string]int)
	g.inputs[node] = 0
	return true
}

// AddNodes adds one or more new node to Graph
func (g *Graph) AddNodes(nodes ...string) bool {
	for _, node := range nodes {
		if ok := g.AddNode(node); !ok {
			return false
		}
	}
	return true
}

// AddEdge adds new source to the destination
func (g *Graph) AddEdge(from, to string) bool {
	m, ok := g.outputs[from]
	if !ok {
		return false
	}

	m[to] = len(m) + 1
	g.inputs[to]++

	return true
}

func (g *Graph) unsafeRemoveEdge(from, to string) {
	delete(g.outputs[from], to)
	g.inputs[to]--
}

// RemoveEdge removes the description from the source
func (g *Graph) RemoveEdge(from, to string) bool {
	if _, ok := g.outputs[from]; !ok {
		return false
	}
	g.unsafeRemoveEdge(from, to)
	return true
}

// Sort sorts nodes based on topological sort algorithm
func (g *Graph) Sort() ([]string, bool) {
	L := make([]string, 0, len(g.nodes))
	S := make([]string, 0, len(g.nodes))

	for _, node := range g.nodes {
		if g.inputs[node] == 0 {
			S = append(S, node)
		}
	}

	for len(S) > 0 {
		var n string
		n, S = S[0], S[1:]
		L = append(L, n)

		ms := make([]string, len(g.outputs[n]))
		for m, i := range g.outputs[n] {
			ms[i-1] = m
		}

		for _, m := range ms {
			g.unsafeRemoveEdge(n, m)

			if g.inputs[m] == 0 {
				S = append(S, m)
			}
		}
	}

	N := 0
	for _, v := range g.inputs {
		N += v
	}

	if N > 0 {
		return L, false
	}

	return L, true
}
