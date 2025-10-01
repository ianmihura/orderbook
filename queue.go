package main

import (
	"slices"
)

/*
A Stack with a unique field v: array of Orders.

The Top(), head, next, is the last element.

Example of a normal queue, 3 items deep: the ask queue followed by the bid queue (id:val):

0: 53
1: 52
2: 51
ASK TOP
-- here we would match --
BID TOP
2: 49
1: 48
0: 47
*/
type Queue struct {
	v []Order
}

/* Pushes new value to the stack, this new value becomes the Top */
func (q *Queue) Push(val Order) {
	q.v = append(q.v, val)
}

/* Removes the Top of the stack, and returns it */
func (q *Queue) Pop() *Order {
	if q.IsEmpty() {
		return nil
	} else {
		val := q.Top()
		q.v = q.v[:len(q.v)-1]
		return val
	}
}

func (q *Queue) IsEmpty() bool {
	return len(q.v) == 0
}

/* Returns pointer to the head element */
func (q *Queue) Top() *Order {
	if q.IsEmpty() {
		return nil
	} else {
		return &q.v[len(q.v)-1]
	}
}

func (q *Queue) Len() int {
	return len(q.v)
}

// Inserts the Order o into the Queue at index i
func (q *Queue) Insert(i int, o Order) {
	q.v = slices.Insert(q.v, i, o)
}

// Finds all elements (Orders) in the stack with the same id.
// Returns the list of indexes of these elements.
func (q *Queue) FindAll(o Order) []int {
	var indexes []int
	// TODO faster bisect search (is sorted by price)
	for i := range q.v {
		if q.v[i].id == o.id {
			indexes = append(indexes, i)
		}
	}
	return indexes
}

// Random Access Remove: removes element (Order) at random index i.
// Returns the removed element.
func (q *Queue) Remove(i int) Order {
	order := q.v[i]
	q.v = slices.Delete(q.v, i, i+1)
	return order
}

func (q *Queue) CopyFromSlice(in []Order) {
	q.v = in
}
