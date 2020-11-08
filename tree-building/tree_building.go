package tree

import (
	"errors"
	"sort"
)

// req: Record is a struct containing int fields ID and Parent
type Record struct {
	ID     int
	Parent int
}

// req: Node is a struct containing int field ID and []*Node field Children.
type Node struct {
	ID       int
	Children []*Node
}

func addChild(node *Node, child *Node) {
	node.Children = append(node.Children, child)
}

func Build(records []Record) (*Node, error) {
	if len(records) == 0 {
		return nil, nil
	}

	index := map[int]*Node{}

	sort.Slice(records[:], func(i, j int) bool {
		return records[i].ID < records[j].ID
	})

	//backlog := []*Record{}

	for i, rec := range records {
		if _, ok := index[rec.ID]; ok {
			return nil, errors.New("duplicate node")
		}
		if i != rec.ID {
			return nil, errors.New("non-continious")
		}
		node := Node{ID: rec.ID, Children: nil}
		index[rec.ID] = &node
		if rec.ID > 0 {
			if rec.ID == rec.Parent {
				return nil, errors.New("only root can self-parent")
			}
			if parent, ok := index[rec.Parent]; ok {
				addChild(parent, &node)
			} else {
				return nil, errors.New("orphaned node")
			}
		} else if rec.Parent > 0 {
			return nil, errors.New("root has a parent")
		}
	}

	return index[0], nil
}
