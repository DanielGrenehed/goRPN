package main

// https://horman.net/avisynth/rpn.html

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type QueueNode struct {
	next *QueueNode
	data float64
}

type Stacker interface {
	Push(float64)
	Pop()
	Top()
	IsEmpty()
}

type Stack struct {
	top  *QueueNode
	size int
}

func CreateStack() Stack {
	return Stack{nil, 0}
}

func (s *Stack) Push(i float64) {
	push := &QueueNode{s.top, i}
	s.top = push
	s.size++
}

func (s *Stack) Pop() float64 {
	if s.size == 0 {
		return 0
	}
	result := s.top.data
	s.top = s.top.next
	s.size--
	return result
}

func (s Stack) Top() float64 {
	if s.IsEmpty() {
		return 0
	}
	return s.top.data
}

func (s Stack) IsEmpty() bool {
	return s.size == 0
}

type VariableListNoder interface {
	Add(string, float64)
	Get(string)
}

type VariableListNode struct {
	next  *VariableListNode
	name  string
	value float64
}

func (n *VariableListNode) Add(name string, value float64) bool {
	c := n
	for ; c.next != nil; c = c.next {
		if c.name == name {
			c.value = value
			return false
		}
	}
	c.next = &VariableListNode{nil, name, value}
	return true
}

func (n *VariableListNode) Get(name string) *VariableListNode {
	for c := n; c != nil; c = c.next {
		if c.name == name {
			return c
		}
	}
	return nil
}

type RPNCalculatorer interface {
	AddVariable(string, float64)
	GetVariable(string)
	StackPtr()
	VarCount()
}

type RPNCalculator struct {
	stack     Stack
	variables *VariableListNode
	vsize     int
}

func (calc *RPNCalculator) AddVariable(name string, value float64) {
	if calc.variables == nil {
		calc.variables = &VariableListNode{nil, name, value}
		calc.vsize++
	} else {
		if calc.variables.Add(name, value) {
			calc.vsize++
		}
	}
}

func (calc *RPNCalculator) GetVariable(name string) *VariableListNode {
	if calc.VarCount() == 0 {
		return nil
	} else {
		return calc.variables.Get(name)
	}
}

func ConstructCalculator() RPNCalculator {
	return RPNCalculator{CreateStack(), nil, 0}
}

func (calc *RPNCalculator) StackPtr() *Stack {
	return &calc.stack
}

func (calc RPNCalculator) VarCount() int {
	return calc.vsize
}

func Top(calc *RPNCalculator) float64 {
	return calc.StackPtr().Top()
}
func Pop(calc *RPNCalculator) float64 {
	return calc.StackPtr().Pop()
}
func Push(calc *RPNCalculator, p float64) {
	calc.StackPtr().Push(p)
}

/*
	Arithmetic operations
*/

func add(s *Stack) float64 {
	a := s.Pop()
	b := s.Pop()
	s.Push(a + b)
	return s.Top()
}

func subtract(s *Stack) float64 {
	a := s.Pop()
	b := s.Pop()
	s.Push(a - b)
	return s.Top()
}

func multiply(s *Stack) float64 {
	a := s.Pop()
	b := s.Pop()
	s.Push(a * b)
	return s.Top()
}

func divide(s *Stack) float64 {
	a := s.Pop()
	b := s.Pop()
	s.Push(a / b)
	return s.Top()
}

func revDivide(s *Stack) float64 {
	a := s.Pop()
	b := s.Pop()
	s.Push(b / a)
	return s.Top()
}

func power(s *Stack) float64 {
	a := s.Pop()
	b := s.Pop()
	s.Push(math.Pow(a, b))
	return s.Top()
}

func sqrt(s *Stack) float64 {
	a := s.Pop()
	s.Push(math.Sqrt(a))
	return s.Top()
}

func modulo(s *Stack) float64 {
	a := s.Pop()
	b := s.Pop()
	s.Push(math.Mod(a, b))
	return s.Top()
}

func abs(s *Stack) float64 {
	a := s.Pop()
	s.Push(math.Abs(a))
	return s.Top()
}

func negative(s *Stack) float64 {
	a := s.Pop()
	s.Push(-a)
	return s.Top()
}

/*
	Min/max
*/

func minimum(s *Stack) float64 {
	a := s.Pop()
	b := s.Pop()
	s.Push(math.Min(a, b))
	return s.Top()
}

func maximum(s *Stack) float64 {
	a := s.Pop()
	b := s.Pop()
	s.Push(math.Max(a, b))
	return s.Top()
}

func zmax(s *Stack) float64 {
	a := s.Pop()
	b := s.Pop()
	s.Push(math.Max(0, math.Min(a, b)))
	return s.Top()
}

/*
	Random
*/

func random(s *Stack) float64 {
	s.Push(rand.Float64())
	return s.Top()
}

func irandom(s *Stack) float64 {
	r := s.Pop()
	s.Push(float64(rand.Intn(int(r))))
	return s.Top()
}

/*
	Trigonometrical operations
*/

func sin(s *Stack) float64 {
	a := s.Pop()
	s.Push(math.Sin(a))
	return s.Top()
}

func cosine(s *Stack) float64 {
	a := s.Pop()
	s.Push(math.Cos(a))
	return s.Top()
}

func tan(s *Stack) float64 {
	a := s.Pop()
	s.Push(math.Tan(a))
	return s.Top()
}

func sincos(s *Stack) float64 {
	a := s.Pop()
	b, c := math.Sincos(a)
	s.Push(b)
	s.Push(c)
	return s.Top()
}

func arcsin(s *Stack) float64 {
	a := s.Pop()
	s.Push(math.Asin(a))
	return s.Top()
}

func arccosine(s *Stack) float64 {
	a := s.Pop()
	s.Push(math.Acos(a))
	return s.Top()
}

func atan(s *Stack) float64 {
	a := s.Pop()
	s.Push(math.Atan(a))
	return s.Top()
}

func atan2(s *Stack) float64 {
	a := s.Pop()
	b := s.Pop()
	s.Push(math.Atan2(a, b))
	return s.Top()
}

/*
	Comparison operations
*/

func lessThan(s *Stack) float64 {
	a := s.Pop()
	b := s.Pop()
	if a < b {
		s.Push(1)
	} else {
		s.Push(0)
	}
	return s.Top()
}

func lessThanOrEqual(s *Stack) float64 {
	a := s.Pop()
	b := s.Pop()
	if a <= b {
		s.Push(1)
	} else {
		s.Push(0)
	}
	return s.Top()
}

func greaterThan(s *Stack) float64 {
	a := s.Pop()
	b := s.Pop()
	if a > b {
		s.Push(1)
	} else {
		s.Push(0)
	}
	return s.Top()
}

func greaterThanOrEqual(s *Stack) float64 {
	a := s.Pop()
	b := s.Pop()
	if a >= b {
		s.Push(1)
	} else {
		s.Push(0)
	}
	return s.Top()
}

func equal(s *Stack) float64 {
	a := s.Pop()
	b := s.Pop()
	if a == b {
		s.Push(1)
	} else {
		s.Push(0)
	}
	return s.Top()
}

func notEqual(s *Stack) float64 {
	a := s.Pop()
	b := s.Pop()
	if a != b {
		s.Push(1)
	} else {
		s.Push(0)
	}
	return s.Top()
}

/*
	Bitwise operations
*/

func bAnd(s *Stack) float64 {
	a := int(s.Pop())
	b := int(s.Pop())
	c := a & b
	s.Push(float64(c))
	return s.Top()
}

func bOr(s *Stack) float64 {
	a := int(s.Pop())
	b := int(s.Pop())
	c := a | b
	s.Push(float64(c))
	return s.Top()
}

func bXor(s *Stack) float64 {
	a := int(s.Pop())
	b := int(s.Pop())
	c := a ^ b
	s.Push(float64(c))
	return s.Top()
}

func bNot(s *Stack) float64 {
	a := ^int(s.Pop())
	s.Push(float64(a))
	return s.Top()
}

/*
	Conditional operator
*/

func conditional(s *Stack) float64 {
	a := s.Pop()
	b := s.Pop()
	c := s.Pop()
	if c != 0 {
		s.Push(a)
	} else {
		s.Push(b)
	}
	return s.Top()
}

/*
	Stack operations
*/

func duplicate(s *Stack) float64 {
	s.Push(s.Top())
	return s.Top()
}

func swap(s *Stack) float64 {
	a := s.Pop()
	b := s.Pop()
	s.Push(a)
	s.Push(b)
	return s.Top()
}

/*
	Constants
*/

func pi(s *Stack) float64 {
	s.Push(math.Pi)
	return s.Top()
}

func tau(s *Stack) float64 {
	s.Push(math.Pi * 2)
	return s.Top()
}

func ipi(s *Stack) float64 {
	s.Push(1.0 / math.Pi)
	return s.Top()
}

func itau(s *Stack) float64 {
	s.Push(1.0 / (math.Pi * 2))
	return s.Top()
}

/*
	Function map
*/
type opFunc func(*Stack) float64
type opMap struct {
	name string
	fnc  opFunc
}

func oplist() []opMap {
	return []opMap{
		// Arithmetic operations
		{"+", add},
		{"-", subtract},
		{"/", divide},
		{"*", multiply},
		{"\\", revDivide},
		{"^", power},
		{"pow", power},
		{"sqrt", sqrt},
		{"%", modulo},
		{"abs", abs},
		{"neg", negative},
		// Min/max
		{"min", minimum},
		{"max", maximum},
		{"zmax", zmax},
		// Random
		{"rand", random},
		{"irand", irandom},
		// Trigonometrical operations
		{"sin", sin},
		{"cos", cosine},
		{"tan", tan},
		{"sincos", sincos},
		{"asin", arcsin},
		{"acos", arccosine},
		{"atan", atan},
		{"atan2", atan2},
		// Comparison operations
		{"<", lessThan},
		{"<=", lessThanOrEqual},
		{"=<", lessThanOrEqual},
		{">", greaterThan},
		{">=", greaterThanOrEqual},
		{"=>", greaterThanOrEqual},
		{"==", equal},
		{"!=", notEqual},
		{"=!", notEqual},
		// Bitwise operations
		{"and", bAnd},
		{"&", bAnd},
		{"or", bOr},
		{"|", bOr},
		{"xor", bXor},
		{"^", bXor},
		{"not", bNot},
		{"~", bNot},
		// Conditional operation
		{"?", conditional},
		// Stack operations
		{"dup", duplicate},
		{"swap", swap},
		// Constants
		{"pi", pi},
		{"tau", tau},
		{"ipi", ipi},
		{"itau", itau},
	}
}

func handleVariable(calc *RPNCalculator, op string) bool {
	if strings.HasPrefix(op, "@") {
		vname := strings.TrimPrefix(op, "@")
		//vname := op[1:]
		var value float64
		if strings.HasSuffix(op, "^") {
			vname = strings.TrimSuffix(vname, "^")
			value = Pop(calc)
		} else {
			value = Top(calc)
		}
		calc.AddVariable(vname, value)
		fmt.Println("Setting", vname, "to", value)
		return true
	} else {
		v := calc.GetVariable(op)
		if v != nil {
			Push(calc, v.value)
			//fmt.Println("Pushing", v.name, "value (", v.value, ") to stack")
			return true
		} else {
			fmt.Println("Failed to interpret", op, " VarCount:", calc.VarCount())
		}
	}
	return false
}

func handleFunction(calc *RPNCalculator, op string) bool {
	ops := oplist()
	nop := len(ops)
	for i := 0; i < nop; i++ {
		if ops[i].name == op {
			ops[i].fnc(calc.StackPtr())
			return true
		}
	}
	return false
}

func processOperation(calc *RPNCalculator, input string) {
	if len(input) == 0 {
		return
	}
	i, err := strconv.ParseFloat(input, 64)
	if err != nil {
		op := strings.ToLower(strings.TrimLeft(strings.TrimRight(input, "\n"), " "))
		found := handleFunction(calc, op)
		if !found {
			found = handleVariable(calc, op)
		}
	} else {
		Push(calc, i)
	}
	fmt.Println("Top of stack:", Top(calc))
}

func processFile(calc *RPNCalculator, filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ops := strings.Split(strings.TrimRight(scanner.Text(), "\n"), " ")
		for i := 0; i < len(ops); i++ {
			processOperation(calc, ops[i])
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func interactiveMode(calc *RPNCalculator) {
	reader := bufio.NewReader(os.Stdin)
	for {
		//fmt.Scanln(&line)
		line, _ := reader.ReadString('\n')
		line = strings.TrimRight(strings.ToLower(line), "\n")
		if line == "exit" || line == "quit" {
			break
		}
		ops := strings.Fields(line)
		//fmt.Println(ops)
		for i := 0; i < len(ops); i++ {
			//fmt.Println("Operation", ops[i])
			processOperation(calc, ops[i])
		}

	}
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	calc := ConstructCalculator()
	argc := len(os.Args)
	for i := 1; i < argc; i++ {
		arg := os.Args[i]
		if strings.HasSuffix(arg, ".rpn") {
			processFile(&calc, arg)
		} else {
			args := strings.Split(arg, " ")
			for j := 0; j < len(args); j++ {
				processOperation(&calc, args[j])
			}
		}
	}

	if argc == 1 {
		interactiveMode(&calc)
	}

	fmt.Println("Top of stack at end:", Top(&calc))
}
