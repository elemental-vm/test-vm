package vm

import "bytes"
import "strconv"

type node struct {
	next *node
	val  int64
}

type stack struct {
	end *node
}

func (s *stack) push(val int64) {
	n := &node{
		val:  val,
		next: s.end,
	}
	s.end = n
}

func (s *stack) pop() int64 {
	v := s.end.val
	s.end = s.end.next
	return v
}

func (s *stack) tos() int64 {
	return s.end.val
}

func (s *stack) string() string {
	if s.end == nil {
		return "[]"
	}

	var out bytes.Buffer
	out.WriteByte('[')
	n := s.end

	for {
		out.WriteString(strconv.FormatInt(n.val, 10))
		if n.next != nil {
			out.WriteByte(',')
			n = n.next
		} else {
			break
		}
	}
	out.WriteByte(']')
	return out.String()
}
