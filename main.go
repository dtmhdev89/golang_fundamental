package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/example/tester/sample"
	_ "github.com/go-sql-driver/mysql"
)

var variableInitialization int = 5

// Golang runtime will actually have some sort of type inference to detect the
// type of variable.
// variableInitialization :=5

var wg sync.WaitGroup
var messages = make(chan string)

func main() {
	// fmt.Println("Hello, 世界")
	// fmt.Println(rand.Int())
	// output, _ := temp.Marshal(map[string]string{"a": "b", "c": "d"})
	// fmt.Println(string(output))
	fmt.Println("Program Begin")
	e := sample.ExposedStruct{
		ExposedSample: "sample",
		// unexposedSample: "sample",
	}

	fmt.Println(e)

	f := sample.NewUnexposedFuction()

	fmt.Println(f.Sample)
	f.ChangeSample("hello")

	fmt.Println(f)

	f.Sample = "new value sample"

	fmt.Println(f)

	fmt.Println("Start main")
	wg.Add(3)
	go side()
	// fmt.Println(msg)
	go side()
	// msg = <-messages
	// fmt.Println(msg)
	go nodeRunner()
	wg.Wait()

	wg.Add(1)
	DoubleNodeRunner()
	wg.Wait()

	wg.Add(1)
	StackRunner()
	wg.Wait()

	wg.Add(1)
	QueueRunner()
	wg.Wait()

	wg.Add(1)
	BinaryTreeRunner()
	wg.Wait()

	wg.Add(1)
	HashedMapsRunner()
	wg.Wait()

	wg.Add(1)
	HashNodeRunner()
	wg.Wait()

	wg.Add(1)
	SortRunner()
	wg.Wait()

	go createPing()
	fmt.Println("Return to main")
	msg := <-messages
	fmt.Println("MSG from ping through channel messages: ", msg)
	fmt.Println("Return to main")
	// time.Sleep(5 * time.Second)

	log.Print("Hello world sample started")
	// http.HandleFunc("/", sample.Handler)
	http.Handle("/", &sample.HandleViaStruct{})
	http.ListenAndServe(":8080", nil)

	fmt.Println("End main")
}

func side() {
	fmt.Println("Start side process")
	time.Sleep(1 * time.Second)
	// messages <- "ping from side"
	fmt.Println("End side process")
	wg.Done()
}

func createPing() {
	time.Sleep(1 * time.Second)
	messages <- "ping"
}

func nodeRunner() {
	fmt.Println("Start node runner")
	aa := Node{Value: 1}
	bb := Node{Value: 2}
	cc := Node{Value: 3}
	aa.Next = &bb
	bb.Next = &cc
	fmt.Println(aa)
	fmt.Println(bb)
	fmt.Println(cc)

	PrintNodes(&aa)
	fmt.Println("Len of Nodes:", LenNodes(aa))

	fmt.Println("Node zz at first creation")
	zz := Node{Value: 0}
	PrintNodes(&zz)
	insertNodeAtFront(&aa, &zz)
	fmt.Println("Node zz after append in front of a single linked list")
	PrintNodes(&zz)
	fmt.Println("Len of Nodes:", LenNodes(zz))
	fmt.Println(&aa)
	fmt.Println("End node runner")
	wg.Done()
}

type Node struct {
	Value int
	Next  *Node
}

func PrintNodes(root *Node) {
	nodewalk := root

	for nodewalk.Next != nil {
		fmt.Println(nodewalk.Value)
		nodewalk = nodewalk.Next
	}

	if nodewalk.Next == nil {
		fmt.Println(nodewalk.Value)
	}
}

func LenNodes(root Node) int {
	nodewalk := root
	counter := 1
	for nodewalk.Next != nil {
		counter = counter + 1
		nodewalk = *nodewalk.Next
	}

	return counter
}

func insertNodeAtFront(root *Node, newNode *Node) {
	newNode.Next = root
	// newNode.Next.Value = 44
}

func Append(root *Node, newNode *Node) {
	nodewalk := root

	for nodewalk.Next != nil {
		nodewalk = nodewalk.Next
	}

	nodewalk.Next = newNode
}

func Insert(root *Node, loc int, newNode *Node) error {
	nodeWalk := root
	counter := 0

	for nodeWalk.Next != nil {
		if counter == loc-1 {
			temp := nodeWalk.Next
			newNode.Next = temp
			nodeWalk.Next = newNode

			return nil
		}
		nodeWalk = nodeWalk.Next
		counter = counter + 1
	}

	return fmt.Errorf("went past no of elements in list")
}

func Delete(root *Node, loc int) error {
	nodeWalk := root
	previousNode := nodeWalk
	counter := int(0)

	for nodeWalk.Next != nil {
		if counter == loc {
			previousNode.Next = nodeWalk.Next

			return nil
		}

		previousNode = nodeWalk
		nodeWalk = nodeWalk.Next
		counter = counter + 1
	}

	return fmt.Errorf("wen past the expected list")
}

func Search(root *Node, searchVal int) *Node {
	nodeWalk := root

	for nodeWalk.Next != nil {
		if nodeWalk.Value == searchVal {
			return nodeWalk
		}
		nodeWalk = nodeWalk.Next
	}

	return nil
}

func DoubleNodeRunner() {
	aa := DoubleNode{Value: 1}
	bb := DoubleNode{Value: 2}
	cc := DoubleNode{Value: 3}
	dd := DoubleNode{Value: 4}

	aa.Next = &bb
	bb.Next = &cc
	bb.Previous = &aa
	cc.Previous = &bb

	root := InsertDoubleNode(&aa, &dd, 0)
	fmt.Printf("root: %v and %p\n", *root, root)
	fmt.Printf("dd at root - 0: %v and %p\n", dd, &dd)
	fmt.Printf("aa at 1: %v and %p\n", aa, &aa)
	fmt.Printf("bb at 2: %v and %p\n", bb, &bb)
	fmt.Printf("cc at 3: %v and %p\n", cc, &cc)

	fmt.Println("Delete first node")
	DeleteDoubleNode(&dd, 2)
	fmt.Printf("aa : %v and %p\n", aa, &aa)
	fmt.Printf("bb : %v and %p\n", bb, &bb)
	fmt.Printf("cc : %v and %p\n", cc, &cc)
	fmt.Printf("dd : %v and %p\n", dd, &dd)

	wg.Done()
}

type DoubleNode struct {
	Value    int
	Next     *DoubleNode
	Previous *DoubleNode
}

func InsertDoubleNode(root *DoubleNode, newNode *DoubleNode, loc int) *DoubleNode {
	n := root
	var p *DoubleNode
	fmt.Println("var p *DoubleNode: ", p)
	counter := 0

	for n != nil {
		if counter == loc {
			newNode.Next = n
			newNode.Previous = p
			if n.Next != nil {
				n.Previous = newNode
			}

			if p != nil {
				p.Next = newNode
				return root
			}

			return newNode
		}
		p = n
		n = n.Next
		counter = counter + 1
	}

	if counter == loc {
		newNode.Previous = p
		p.Next = newNode

		return root
	}

	return nil
}

func DeleteDoubleNode(root *DoubleNode, loc int) *DoubleNode {
	counter := 0
	n := root
	var p *DoubleNode

	for n != nil {
		if counter == loc {
			if p == nil {
				temp := n.Next
				n.Next.Previous = nil
				n.Next = nil
				n.Previous = nil

				return temp
			}

			p.Next = n.Next
			n.Next.Previous = p
			n.Next = nil
			n.Previous = nil

			return root
		}

		p = n
		n = n.Next
		counter = counter + 1
	}

	return nil
}

type Stack struct {
	stack []int
}

func NewStack() Stack {
	return Stack{stack: []int{}}
}

func StackRunner() {
	fmt.Println("Stack Runner")

	aa := NewStack()
	aa.AddToStack(1)
	aa.AddToStack(2)
	fmt.Println(aa.stack)
	fmt.Println(aa.RemoveItemFromStack())
	fmt.Println(aa.RemoveItemFromStack())
	fmt.Println(aa.RemoveItemFromStack())

	wg.Done()
}

func (s *Stack) AddToStack(item int) {
	s.stack = append(s.stack, item)
}

func (s *Stack) RemoveItemFromStack() (int, error) {
	if len(s.stack) == 0 {
		return 0, fmt.Errorf("stack is empty")
	}

	temp := s.stack[len(s.stack)-1]
	s.stack = s.stack[:len(s.stack)-1]

	return temp, nil
}

func QueueRunner() {
	fmt.Println("Queue Runner")
	aa := NewQueue()
	aa.Enqueue(1)
	aa.Enqueue(2)

	fmt.Println(aa)
	aa.Dequeue()
	aa.Dequeue()

	fmt.Println(aa)

	aa.Enqueue(4)
	fmt.Println(aa)

	aa.Dequeue()
	aa.Dequeue()

	fmt.Println(aa)

	wg.Done()
}

type Queue struct {
	queue []int
}

func NewQueue() Queue {
	return Queue{queue: []int{}}
}

func (q *Queue) Enqueue(item int) {
	q.queue = append(q.queue, item)
}

func (q *Queue) Dequeue() (int, error) {
	if len(q.queue) == 0 {
		return 0, fmt.Errorf("queue is empty")
	}

	temp := q.queue[0]
	q.queue = q.queue[1:]

	return temp, nil
}

func BinaryTreeRunner() {
	A := TreeNode{data: "A"}
	B := TreeNode{data: "B"}
	C := TreeNode{data: "C"}
	D := TreeNode{data: "D"}
	E := TreeNode{data: "E"}
	F := TreeNode{data: "F"}
	G := TreeNode{data: "G"}

	A.left = &B
	A.right = &C
	B.left = &D
	B.right = &E
	C.left = &F
	C.right = &G

	InorderPrint(&A)

	wg.Done()
}

type TreeNode struct {
	data  string
	left  *TreeNode
	right *TreeNode
}

func InorderPrint(root *TreeNode) {
	if root == nil {
		return
	}

	if root.left != nil {
		InorderPrint(root.left)
	}

	fmt.Println(root.data)

	if root.right != nil {
		InorderPrint(root.right)
	}
}

func HashedMapsRunner() {
	fmt.Println("HashedMaps Runner")

	aa := NewHashedMaps()
	// aa.Set("aa", "sample value")
	// fmt.Println(aa.Get("aa"))

	aa.Set("a", "sample value")
	aa.Set("ABB", "unexpected sample value")
	fmt.Println(aa.Get("a"))

	wg.Done()
}

type HashedMaps struct {
	items [100]string
}

func NewHashedMaps() HashedMaps {
	return HashedMaps{items: [100]string{}}
}

func (h HashedMaps) GetHash(key string) int {
	totalSum := 0
	for _, v := range key {
		totalSum = totalSum + int(v)
	}

	hashKey := totalSum % 100

	return hashKey
}

func (h *HashedMaps) Set(key, val string) {
	hashedKey := h.GetHash(key)
	h.items[hashedKey] = val
}

func (h *HashedMaps) Get(key string) string {
	hashedKey := h.GetHash(key)
	return h.items[hashedKey]
}

func HashNodeRunner() {
	fmt.Println("HashNode Runner")

	aa := NewHashedNodeMaps()
	aa.Set("a", "sample value")
	aa.Set("ABB", "unknown sample value")
	fmt.Println(aa.Get("a"))
	fmt.Println(aa.Get("ABB"))
	wg.Done()
}

type HashNode struct {
	Key   string
	Value string
}

type HashedNodeMaps struct {
	items [100][]HashNode
}

func NewHashedNodeMaps() HashedNodeMaps {
	return HashedNodeMaps{items: [100][]HashNode{}}
}

func (h HashedNodeMaps) GetHash2(key string) int {
	totalSum := 0
	for _, v := range key {
		totalSum = totalSum + int(v)
	}

	hashKey := totalSum % 100
	return hashKey
}

func (h *HashedNodeMaps) Set(key, val string) {
	hashedKey := h.GetHash2(key)
	if len(h.items[hashedKey]) == 0 {
		h.items[hashedKey] = []HashNode{HashNode{Key: key, Value: val}}
	}

	for _, i := range h.items[hashedKey] {
		if i.Key == key {
			return
		}
	}

	h.items[hashedKey] = append(h.items[hashedKey], HashNode{Key: key, Value: val})
}

func (h *HashedNodeMaps) Get(key string) string {
	hashedKey := h.GetHash2(key)

	if len(h.items[hashedKey]) == 0 {
		return ""
	}

	for _, item := range h.items[hashedKey] {
		if item.Key == key {
			return item.Value
		}
	}

	return ""
}

func SortRunner() {
	fmt.Println("Sort Runner")
	values := []int{4, 3, 5, 2, 1}
	fmt.Println(values)
	BubleSort(values)
	fmt.Println(values)

	fmt.Println("Merge Sort")
	bb := []int{4, 3, 5, 2, 1}
	out := MergeSort(bb)
	fmt.Println("Sorted: ", out)

	fmt.Println("Quick Sort")
	cc := []int{3, 5, 4, 2, 1}
	out = QuickSort(cc)
	fmt.Println("Sorted", out)

	fmt.Println("Binary Search")
	dd := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	found := BinarySearch(8, dd)
	fmt.Printf("Was 8 found: %v\n", found)

	fmt.Println("Fibonacci Number")
	fmt.Println(fibonacci(10))

	fmt.Println("Fibonacci Number with memoization (dynamic programming)")
	fmt.Println(Fibonacci(100))

	fmt.Println("Fibonacci Number using tabulate approach")
	fmt.Println(fibonacciTabulate(100))
	wg.Done()
}

// slices is mutable type, but array (fixed length) is immutable

// O(n^2)
func BubleSort(items []int) {
	if len(items) <= 1 {
		return
	}

	for {
		sortHappended := false
		for i := 0; i < len(items)-1; i++ {
			if items[i] > items[i+1] {
				temp := items[i]
				items[i] = items[i+1]
				items[i+1] = temp
				sortHappended = true
			}
		}
		if !sortHappended {
			break
		}
	}
}

// Time complexity O(n log(n))
func MergeSort(items []int) []int {
	// Divide part
	if len(items) <= 1 {
		return items
	}

	leftSide := MergeSort(items[0 : len(items)/2])
	rightSide := MergeSort(items[len(items)/2:])

	// Conquer part
	i := 0
	j := 0
	combined := []int{}
	for i < len(leftSide) || j < len(rightSide) {
		if i >= len(leftSide) {
			combined = append(combined, rightSide[j:]...)
			j = len(rightSide)
			continue
		}

		if j >= len(rightSide) {
			combined = append(combined, leftSide[i:]...)
			i = len(leftSide)
			continue
		}

		if leftSide[i] < rightSide[j] {
			combined = append(combined, leftSide[i])
			i = i + 1
			continue
		}

		combined = append(combined, rightSide[j])
		j = j + 1
	}

	return combined
}

// Time complexity O(n log(n))
func QuickSort(values []int) []int {
	if len(values) <= 1 {
		return values
	}

	leftSide := []int{}
	rightSide := []int{}

	pivot := values[len(values)-1]

	for _, v := range values[:len(values)-1] {
		if v < pivot {
			leftSide = append(leftSide, v)
			continue
		}

		rightSide = append(rightSide, v)
	}

	sortedLeftSide := QuickSort(leftSide)
	sortedRightSide := QuickSort(rightSide)

	sorted := append(sortedLeftSide, pivot)
	sorted = append(sorted, sortedRightSide...)

	return sorted

}

// Time Complexity O(log(n))
func BinarySearch(finding int, values []int) bool {
	if len(values) == 0 {
		return false
	}

	if len(values) == 1 {
		return values[0] == finding
	}

	found := false
	leftHalf := values[:len(values)/2]
	rightHalf := values[len(values)/2:]

	if finding >= rightHalf[0] {
		found = BinarySearch(finding, rightHalf)
	} else {
		found = BinarySearch(finding, leftHalf)
	}

	return found
}

// Fibonacci number. Time complexity O(2^n)
func fibonacci(n int) int {
	if n < 0 {
		return 0
	}

	if n == 1 || n == 2 {
		return 1
	}

	return fibonacci(n-1) + fibonacci(n-2)
}

var store = map[int]int{
	1: 1,
	2: 1,
}

// Fibonacci number with memoization. Time complexity: O(n). But need large memory to store the values storage
func Fibonacci(n int) int {
	if n <= 0 {
		return 0
	}

	if store[n] != 0 {
		return store[n]
	}
	val := Fibonacci(n-1) + Fibonacci(n-2)
	store[n] = val

	return val
}

// Tabulate way for fibinacci. Time complexity O(n). More efficient memory usage.
func fibonacciTabulate(n int) int {
	if n <= 0 {
		return 0
	}

	if n <= 2 {
		return 1
	}

	previous1 := 1
	previous2 := 1
	currentVal := 0

	for i := 3; i <= n; i++ {
		currentVal = previous1 + previous2
		previous1 = previous2
		previous2 = currentVal
	}

	return currentVal
}
