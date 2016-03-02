package main

import (
  "fmt"
)

// TODO: 
//    - construct parse tree from regex(?)
//    - construct DFA from parse tree

// TODO: define something a bit more sophisticated than this
var END rune = '0'

type State struct {
  Op rune
  Next *State
}

func (s *State) Matches(op rune) bool {
  return s.Op == op
}

type List struct {
  N int
  States []State
}

func (l *List) Init() {
  // TODO: decide reasonable limit
  l.N = 0
  l.States = make([]State, 100)
}

func (l *List) Add(state State) {
  l.States[l.N] = state
  l.N += 1
}

type Regex struct {
  Start State
}

func (re *Regex) Match(text string) bool {

  var current, next List
  current.Init()
  next.Init()
  current.Add(re.Start)
  for _, l := range(text) {
    for i := 0; i < current.N; i++ {
      if s := current.States[i]; s.Matches(l) {
        next.Add(*s.Next)
      }
    }
    current, next = next, current
    next.N = 0
  }

  for i := 0; i < current.N; i++ {
    if current.States[i].Matches(END) {
      return true
    }
  }
  return false
}

func main() {

  /* manually construct nfa until parser */
  end := State{Op:END, Next:nil}
  i := State{Op:'i', Next:&end}
  h := State{Op:'h', Next:&i}
  re := Regex{Start: h}
  fmt.Println("true: ", re.Match("hi"))
  fmt.Println("false: ", re.Match("ho"))
}
