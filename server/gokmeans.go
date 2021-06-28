package server

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"time"
)

// Node represents an observation of floating point values
type Node []float64

type Data struct {
	Component1 float64 `json:"x"`
	Component2 float64 `json:"y"`
}
type Cluster struct {
	Index int `json:"index"`
}
type Cent struct {
	Valor  float64 `json:"Claster 1"`
	Valor2 float64 `json:"Claster 2"`
}

func GetCentroids(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Centroids)
}

//Realiza la predicción
func PredictKmeans(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newdata Data
	_ = json.NewDecoder(r.Body).Decode(&newdata)
	var datita Node = Node{
		newdata.Component1,
		newdata.Component2,
	}
	cluster := Nearest(datita, Centroids)
	var ResponseCluster Cluster
	ResponseCluster.Index = cluster
	json.NewEncoder(w).Encode(ResponseCluster)
}
func wait(c chan int, values Node) {
	count := len(values)

	<-c
	for respCnt := 1; respCnt < count; respCnt++ {
		<-c
	}
}

func Train(Nodes []Node, clusterCount int, maxRounds int) (bool, []Node) {
	if int(len(Nodes)) < clusterCount {
		return false, nil
	}

	stdLen := 0
	for i, Node := range Nodes {
		curLen := len(Node)

		if i > 0 && len(Node) != stdLen {
			return false, nil
		}

		stdLen = curLen
	}

	centroids := make([]Node, clusterCount)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < clusterCount; i++ {
		srcIndex := r.Intn(len(Nodes))
		srcLen := len(Nodes[srcIndex])
		centroids[i] = make(Node, srcLen)
		copy(centroids[i], Nodes[r.Intn(len(Nodes))])
	}

	return Train2(Nodes, clusterCount, maxRounds, centroids)
}

func Train2(Nodes []Node, clusterCount int, maxRounds int, centroids []Node) (bool, []Node) {
	movement := true
	for i := 0; i < maxRounds && movement; i++ {
		movement = false

		groups := make(map[int][]Node)

		for _, Node := range Nodes {
			near := Nearest(Node, centroids)
			groups[near] = append(groups[near], Node)
		}

		for key, group := range groups {
			newNode := meanNode(group)

			if !equal(centroids[key], newNode) {
				centroids[key] = newNode
				movement = true
			}
		}
	}

	return true, centroids
}

func equal(node1, node2 Node) bool {
	if len(node1) != len(node2) {
		return false
	}

	for i, v := range node1 {
		if v != node2[i] {
			return false
		}
	}

	return true
}

func Nearest(in Node, nodes []Node) int {
	count := len(nodes)

	results := make(Node, count)
	cnt := make(chan int)
	for i, node := range nodes {
		go func(i int, node, cl Node) {
			results[i] = distance(in, node)
			cnt <- 1
		}(i, node, in)
	}

	wait(cnt, results)

	mindex := 0
	curdist := results[0]

	for i, dist := range results {
		if dist < curdist {
			curdist = dist
			mindex = i
		}
	}

	return mindex
}

func distance(node1 Node, node2 Node) float64 {
	length := len(node1)
	squares := make(Node, length)

	cnt := make(chan int)

	for i := range node1 {
		go func(i int) {
			diff := node1[i] - node2[i]
			squares[i] = diff * diff
			cnt <- 1
		}(i)
	}

	wait(cnt, squares)

	sum := 0.0
	for _, val := range squares {
		sum += val
	}

	return sum
}

func meanNode(values []Node) Node {
	newNode := make(Node, len(values[0]))

	for _, value := range values {
		for j := 0; j < len(newNode); j++ {
			newNode[j] += value[j]
		}
	}

	for i, value := range newNode {
		newNode[i] = value / float64(len(values))
	}

	return newNode
}
