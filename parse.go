package main

import (
  "fmt"
)

type Op uint8

/*
  c     ->    { Op, Rune }
  cd    ->    { Op, Rune }
  c|d   ->    { Op, Rune }
  c*    ->    { Op, Rune }
  c?    ->    { Op, Rune }
  c+    ->    { Op, Rune }
*/

// push c
// pop c; make ?

const (
  Character     Op = iota + 1  // c
  Alternate                    // c|d
  Concat                       // cd
  Star                         // c*
  Question                     // c?
  Plus                         // c+
  Escape                       // \c
)

type Regex struct {
  Op Op
  Left *Regex
  Right *Regex
  Rune rune
}

func (re *Regex) Print() {
  if re.Left != nil {
    re.Left.Print()
  }
  fmt.Printf("%c", re.Rune)
  if re.Right != nil {
    re.Right.Print()
  }
}

type Stack []Regex

func (s *Stack) Pop() Regex {
  if len(*s) == 0 {
    return Regex{}
  }
  re := (*s)[len(*s)-1]
  *s = (*s)[:len(*s)-1]
  return re
}

func (s *Stack) Push(re Regex) {
  *s = append(*s, re)
}

type Parser struct {
  Stack Stack
  Len int
  Regex Regex
}

func (p *Parser) Init() {
  p.Stack = make(Stack, 1000)
  p.Len = 0
}

func (p *Parser) Push(r Regex) {
  p.Stack.Push(r)
  p.Len += 1
}

func (p *Parser) Pop() Regex {
  if p.Len == 0 {
    return Regex{}
  }
  p.Len -= 1
  return p.Stack.Pop()
}

func MakeRegex(r rune, op Op) Regex {
  return Regex{Op:op,Rune:r}
}


func GetOp(r rune) Op {
  switch(r) {
  case '|':
    return Alternate
  case '?':
    return Question
  case '*':
    return Star
  case '+':
    return Plus
  case '\\':
    return Escape
  default:
    return Character
  }
}

func (p *Parser) parse(text string) Regex {

  var last Op
  for _, c := range(text) {
    switch(c) {
    case '|':
      alt := Regex{Op:Alternate, Rune:c}
      e := p.Pop()
      alt.Left = &e
      p.Push(alt)
    case '?':
      qu := Regex{Op:Question,Rune:c}
      e := p.Pop()
      qu.Left = &e
      p.Push(qu)
    case '*':
      st := Regex{Op:Star,Rune:c}
      e := p.Pop()
      st.Left = &e
      p.Push(st)
    case '+':
      pl := Regex{Op:Plus,Rune:c}
      e := p.Pop()
      pl.Left = &e
      p.Push(pl)
    case '\\':
      p.Push(Regex{Op:Escape, Rune:c})
    default:
      ch := Regex{Op:Character, Rune:c}
      if last == Character || last == Concat {
        e := p.Pop()
        con := Regex{Op:Concat, Rune:'&'}
        con.Left = &e
        con.Right = &ch
        p.Push(con)
      } else if last == Alternate {
        alt := p.Pop()
        alt.Right = &ch
        p.Push(alt)
      } else if last == Escape {
        esc := p.Pop()
        esc.Right = &ch
        p.Push(esc)
      } else {
        p.Push(ch)
      }
    }
    last = GetOp(c)
  }
  return p.Pop()
}

func main() {
  var p Parser
  p.Init()
  p.parse("a")
  p.parse("bc")
  p.parse("d|e")
  p.parse("f*")
  reg := p.parse("g+")
  reg.Print()
 
  fmt.Println()

  p.parse("\\h")
  reg = p.parse("i?")
  reg.Print()
  reg = p.parse("i?d|eaa+abcc")
  reg.Print()
}
