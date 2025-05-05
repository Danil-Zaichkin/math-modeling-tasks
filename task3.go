package main

import (
	"fmt"
	"math"
	"strings"
)

type Edge struct {
	to, rev, cap, cost, flow int
}

type Graph struct {
	adj [][]*Edge
}

func NewGraph(n int) *Graph {
	adj := make([][]*Edge, n)
	for i := range adj {
		adj[i] = make([]*Edge, 0)
	}
	return &Graph{adj: adj}
}

func (g *Graph) AddEdge(u, v, cap, cost int) {
	fwd := &Edge{to: v, rev: len(g.adj[v]), cap: cap, cost: cost}
	bwd := &Edge{to: u, rev: len(g.adj[u]), cap: 0, cost: -cost}
	g.adj[u] = append(g.adj[u], fwd)
	g.adj[v] = append(g.adj[v], bwd)
}

func minCostMaxFlow(g *Graph, s, t int, idToName []string) (int, int) {
	n := len(g.adj)
	const INF = math.MaxInt32 / 2
	totalFlow, totalCost := 0, 0
	iteration := 0

	for {
		dist := make([]int, n)
		inQ := make([]bool, n)
		prevNode := make([]int, n)
		prevEdge := make([]int, n)
		for i := range dist {
			dist[i] = INF
			prevNode[i] = -1
			prevEdge[i] = -1
		}
		dist[s] = 0
		queue := []int{s}
		inQ[s] = true

		for len(queue) > 0 {
			u := queue[0]
			queue = queue[1:]
			inQ[u] = false
			for i, e := range g.adj[u] {
				if e.cap > e.flow && dist[e.to] > dist[u]+e.cost {
					dist[e.to] = dist[u] + e.cost
					prevNode[e.to] = u
					prevEdge[e.to] = i
					if !inQ[e.to] {
						queue = append(queue, e.to)
						inQ[e.to] = true
					}
				}
			}
		}

		if dist[t] == INF {
			break
		}

		augFlow := INF
		for v := t; v != s; v = prevNode[v] {
			u := prevNode[v]
			e := g.adj[u][prevEdge[v]]
			if rem := e.cap - e.flow; rem < augFlow {
				augFlow = rem
			}
		}

		if iteration < 5 {
			fmt.Printf("\nИтерация %d\n", iteration+1)
			pathNodes := []int{}
			for v := t; v != -1; v = prevNode[v] {
				pathNodes = append(pathNodes, v)
			}
			for i, j := 0, len(pathNodes)-1; i < j; i, j = i+1, j-1 {
				pathNodes[i], pathNodes[j] = pathNodes[j], pathNodes[i]
			}
			names := make([]string, len(pathNodes))
			for i, u := range pathNodes {
				names[i] = idToName[u]
			}
			pathStr := strings.Join(names, " → ")
			pathCost := 0
			for v := t; v != s; v = prevNode[v] {
				u := prevNode[v]
				e := g.adj[u][prevEdge[v]]
				pathCost += e.cost
			}
			fmt.Printf("Путь: %s\n", pathStr)
			fmt.Printf("Пропускаем поток: %d\n", augFlow)
			fmt.Printf("Стоимость на данной итерации: %d\n", augFlow*pathCost)
			fmt.Printf("Общий поток после итерации: %d\n", totalFlow+augFlow)
			fmt.Printf("Общая стоимость после итерации: %d\n", totalCost+augFlow*pathCost)
		}

		for v := t; v != s; v = prevNode[v] {
			u := prevNode[v]
			e := g.adj[u][prevEdge[v]]
			e.flow += augFlow
			g.adj[e.to][e.rev].flow -= augFlow
			totalCost += augFlow * e.cost
		}
		totalFlow += augFlow
		iteration++
	}

	fmt.Printf("\nВсего итераций: %d\n", iteration)
	return totalFlow, totalCost
}

func main() {
	n := 7
	edgeCaps := [][]int{
		{-1, 30, 45, 25, 30, 20, 40},
		{30, -1, 55, 25, 35, 40, 25},
		{25, 30, -1, 45, 75, 30, 40},
		{15, 10, 25, -1, 40, 30, 80},
		{10, 45, 15, 60, -1, 60, 75},
		{10, 30, 45, 30, 55, -1, 40},
		{15, 25, 45, 30, 40, 50, -1},
	}
	edgeCosts := [][]int{
		{-1, 5, 10, 4, 5, 6, 10},
		{1, -1, 7, 10, 15, 5, 5},
		{1, 10, -1, 5, 4, 7, 12},
		{2, 6, 4, -1, 5, 10, 8},
		{1, 7, 4, 4, -1, 9, 2},
		{1, 4, 2, 3, 8, -1, 12},
		{1, 10, 5, 6, 8, 16, -1},
	}
	nodeCaps := []int{1000, 55, 35, 40, 50, 45, 1000}
	totalNodes := 2 * n
	g := NewGraph(totalNodes)

	idToName := make([]string, totalNodes)
	for i := 0; i < n; i++ {
		idToName[i] = fmt.Sprintf("n%d-in", i+1)
		idToName[i+n] = fmt.Sprintf("n%d-out", i+1)
	}

	for i := 0; i < n; i++ {
		inID, outID := i, i+n
		g.AddEdge(inID, outID, nodeCaps[i], 0)
	}

	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if i != j && edgeCaps[i][j] >= 0 {
				fromID, toID := i+n, j
				g.AddEdge(fromID, toID, edgeCaps[i][j], edgeCosts[i][j])
			}
		}
	}

	source, sink := 0, totalNodes-1
	flow, cost := minCostMaxFlow(g, source, sink, idToName)

	fmt.Printf("\nМаксимальное кол-во поездов: %d\n", flow)
	fmt.Printf("Минимальная стоимость проезда максимальным кол-вом поездов: %d\n\n", cost)

	fmt.Println("Итоговые потоки:")
	for u := 0; u < totalNodes; u++ {
		for _, e := range g.adj[u] {
			if e.flow > 0 && e.cost >= 0 {
				fmt.Printf("%s -> %s: %d\n", idToName[u], idToName[e.to], e.flow)
			}
		}
	}
}
