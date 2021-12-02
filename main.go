package main

// reference: https://horman.net/avisynth/rpn.html

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

/*
	Variable list
*/

type VariableListNoder interface {
	Set(string, float64)
	Get(string)
}

type VariableListNode struct {
	next  *VariableListNode
	name  string
	value float64
}

// returns true if call adds a new element to list
func (n *VariableListNode) Set(name string, value float64) bool {
	c := n
	for ; c != nil; c = c.next {
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

/*
	Calculator, containing the stack and all variables
*/

type RPNCalculatorer interface {
	SetVariable(string, float64)
	GetVariable(string)
	StackPtr()
	VarCount()
	PrintList()
}

type RPNCalculator struct {
	stack     Stack
	variables *VariableListNode
	vsize     int
}

func ConstructCalculator() *RPNCalculator {
	return &RPNCalculator{CreateStack(), nil, 0}
}

func (calc *RPNCalculator) SetVariable(name string, value float64) {
	if calc.variables == nil {
		calc.variables = &VariableListNode{nil, name, value}
		calc.vsize++
	} else {
		if calc.variables.Set(name, value) {
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

func (calc RPNCalculator) PrintVariables() {
	c := calc.variables
	fmt.Println("Variables: (", calc.VarCount(), ")")
	for ; c != nil; c = c.next {
		fmt.Println("{'"+c.name+"' ,", c.value, "}")
	}
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

/*

	Stack operations

*/

/*
	Arithmetic operations
*/

func add(s *Stack) {
	a := s.Pop()
	b := s.Pop()
	s.Push(a + b)
}

func subtract(s *Stack) {
	a := s.Pop()
	b := s.Pop()
	s.Push(a - b)
}

func multiply(s *Stack) {
	a := s.Pop()
	b := s.Pop()
	s.Push(a * b)
}

func divide(s *Stack) {
	a := s.Pop()
	b := s.Pop()
	s.Push(a / b)
}

func revDivide(s *Stack) {
	a := s.Pop()
	b := s.Pop()
	s.Push(b / a)
}

func power(s *Stack) {
	a := s.Pop()
	b := s.Pop()
	s.Push(math.Pow(a, b))
}

func sqrt(s *Stack) {
	a := s.Pop()
	s.Push(math.Sqrt(a))
}

func modulo(s *Stack) {
	a := s.Pop()
	b := s.Pop()
	s.Push(math.Mod(a, b))
}

func abs(s *Stack) {
	a := s.Pop()
	s.Push(math.Abs(a))
}

func negative(s *Stack) {
	a := s.Pop()
	s.Push(-a)
}

/*
	Min/max
*/

func minimum(s *Stack) {
	a := s.Pop()
	b := s.Pop()
	s.Push(math.Min(a, b))
}

func maximum(s *Stack) {
	a := s.Pop()
	b := s.Pop()
	s.Push(math.Max(a, b))
}

func zmax(s *Stack) {
	a := s.Pop()
	b := s.Pop()
	s.Push(math.Max(0, math.Min(a, b)))
}

/*
	Random
*/

func random(s *Stack) {
	s.Push(rand.Float64())
}

func irandom(s *Stack) {
	r := s.Pop()
	s.Push(float64(rand.Intn(int(r))))
}

/*
	Trigonometrical operations
*/

func sin(s *Stack) {
	a := s.Pop()
	s.Push(math.Sin(a))
}

func cosine(s *Stack) {
	a := s.Pop()
	s.Push(math.Cos(a))
}

func tan(s *Stack) {
	a := s.Pop()
	s.Push(math.Tan(a))
}

func sincos(s *Stack) {
	a := s.Pop()
	b, c := math.Sincos(a)
	s.Push(b)
	s.Push(c)
}

func arcsin(s *Stack) {
	a := s.Pop()
	s.Push(math.Asin(a))
}

func arccosine(s *Stack) {
	a := s.Pop()
	s.Push(math.Acos(a))
}

func atan(s *Stack) {
	a := s.Pop()
	s.Push(math.Atan(a))
}

func atan2(s *Stack) {
	a := s.Pop()
	b := s.Pop()
	s.Push(math.Atan2(a, b))
}

/*
	Comparison operations
*/

func lessThan(s *Stack) {
	a := s.Pop()
	b := s.Pop()
	if a < b {
		s.Push(1)
	} else {
		s.Push(0)
	}
}

func lessThanOrEqual(s *Stack) {
	a := s.Pop()
	b := s.Pop()
	if a <= b {
		s.Push(1)
	} else {
		s.Push(0)
	}
}

func greaterThan(s *Stack) {
	a := s.Pop()
	b := s.Pop()
	if a > b {
		s.Push(1)
	} else {
		s.Push(0)
	}
}

func greaterThanOrEqual(s *Stack) {
	a := s.Pop()
	b := s.Pop()
	if a >= b {
		s.Push(1)
	} else {
		s.Push(0)
	}
}

func equal(s *Stack) {
	a := s.Pop()
	b := s.Pop()
	if a == b {
		s.Push(1)
	} else {
		s.Push(0)
	}
}

func notEqual(s *Stack) {
	a := s.Pop()
	b := s.Pop()
	if a != b {
		s.Push(1)
	} else {
		s.Push(0)
	}
}

/*
	Bitwise operations
*/

func bAnd(s *Stack) {
	a := int(s.Pop())
	b := int(s.Pop())
	c := a & b
	s.Push(float64(c))
}

func bOr(s *Stack) {
	a := int(s.Pop())
	b := int(s.Pop())
	c := a | b
	s.Push(float64(c))
}

func bXor(s *Stack) {
	a := int(s.Pop())
	b := int(s.Pop())
	c := a ^ b
	s.Push(float64(c))
}

func bNot(s *Stack) {
	a := ^int(s.Pop())
	s.Push(float64(a))
}

/*
	Conditional operator
*/

func conditional(s *Stack) {
	a := s.Pop()
	b := s.Pop()
	c := s.Pop()
	if c != 0 {
		s.Push(a)
	} else {
		s.Push(b)
	}
}

/*
	Stack operations
*/

func duplicate(s *Stack) {
	s.Push(s.Top())
}

func swap(s *Stack) {
	a := s.Pop()
	b := s.Pop()
	s.Push(a)
	s.Push(b)
}

func height(s *Stack) {
	s.Push(float64(s.size))
}

func pop(s *Stack) {
	s.Pop()
}

/*
	Constants
*/

func pi(s *Stack) {
	s.Push(math.Pi)
}

func tau(s *Stack) {
	s.Push(math.Pi * 2)
}

func ipi(s *Stack) {
	s.Push(1.0 / math.Pi)
}

func itau(s *Stack) {
	s.Push(1.0 / (math.Pi * 2))
}

/*
	Function map
*/
type opFunc func(*Stack)
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
		{"pop", pop},
		{"hgt", height},
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
		var value float64
		if strings.HasSuffix(op, "^") {
			vname = strings.TrimSuffix(vname, "^")
			value = Pop(calc)
		} else {
			value = Top(calc)
		}
		calc.SetVariable(vname, value)
		fmt.Println("Setting", vname, "to", value)
		return true
	} else {
		v := calc.GetVariable(op)
		if v != nil {
			calc.StackPtr().Push(v.value)
			return true
		} else {
			fmt.Println("Failed to interpret '" + op + "'")
			calc.PrintVariables()
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

func processFile(calc *RPNCalculator, filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		processLine(calc, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func processOperation(calc *RPNCalculator, op string) {
	if len(op) == 0 {
		return
	} else if strings.HasSuffix(op, ".rpn") {
		fmt.Println("Loading file", op)
		processFile(calc, op)
		fmt.Println("End of file", op)
		fmt.Println("Top of stack:", Top(calc))
		return
	}
	i, err := strconv.ParseFloat(op, 64)
	if err != nil {
		found := handleFunction(calc, op)
		if !found {
			found = handleVariable(calc, op)
		}
	} else {
		calc.StackPtr().Push(i)
	}
	fmt.Println("Top of stack:", Top(calc), "op:", op)
}

func processLine(calc *RPNCalculator, line string) {
	line = strings.TrimLeft(strings.TrimRight(strings.ToLower(line), "\n"), " ")
	ops := strings.Split(line, " ")
	for i := 0; i < len(ops); i++ {
		processOperation(calc, ops[i])
	}
}

func interactiveMode(calc *RPNCalculator) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Entering interactive mode")
	for {
		fmt.Print(":")
		line, _ := reader.ReadString('\n')
		line = strings.TrimRight(strings.ToLower(line), "\n")
		if line == "exit" || line == "quit" {
			break
		}
		processLine(calc, line)
	}
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	calc := ConstructCalculator()
	argc := len(os.Args)
	for i := 1; i < argc; i++ {
		arg := os.Args[i]
		processLine(calc, arg)
	}

	if argc == 1 {
		interactiveMode(calc)
	}

	fmt.Println("Top of stack at end:", Top(calc))
}
