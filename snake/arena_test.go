package snake

import "testing"

var pointsDouble = make(chan int)

func newDoubleArenaWithFoodFinder(h, w int, f func(*Arena, []int) bool) *Arena {
	a := newDoubleArena(h, w)
	a.hasFood = f
	return a
}

func newDoubleArena(h, w int) *Arena {
	s := newSnake(RIGHT, [][]int{
		{1, 0},
		{1, 1},
		{1, 2},
		{1, 3},
		{1, 4},
	})

	return newArena(s, pointsDouble, h, w)
}

func TestArenaHaveFoodPlaced(t *testing.T) {
	if a := newDoubleArena(20, 20); a.Food == nil {
		t.Fatal("Arena expected to have food placed")
	}
}

func TestMoveSnakeOutOfArenaHeightLimit(t *testing.T) {
	a := newDoubleArena(4, 10)
	a.Snake.changeDirection(UP)

	if err := a.moveSnake(); err == nil || err.Error() != "Died" {
		t.Fatal("Expected Snake to die when moving outside the Arena height limits")
	}
}

func TestMoveSnakeOutOfArenaWidthLimit(t *testing.T) {
	a := newDoubleArena(10, 1)
	a.Snake.changeDirection(LEFT)

	if err := a.moveSnake(); err == nil || err.Error() != "Died" {
		t.Fatal("Expected Snake to die when moving outside the Arena height limits")
	}
}

func TestPlaceNewFoodWhenEatFood(t *testing.T) {
	a := newDoubleArenaWithFoodFinder(10, 10, func(*Arena, []int) bool {
		return true
	})

	f := a.Food

	a.moveSnake()

	if a.Food.X == f.X && a.Food.Y == f.Y {
		t.Fatal("Expected new food to have been placed on Arena")
	}
}

func TestIncreaseSnakeLengthWhenEatFood(t *testing.T) {
	a := newDoubleArenaWithFoodFinder(10, 10, func(*Arena, []int) bool {
		return true
	})

	l := a.Snake.Length

	a.moveSnake()

	if a.Snake.Length != l+1 {
		t.Fatal("Expected Snake to have grown")
	}
}

func TestAddPointsWhenEatFood(t *testing.T) {
	a := newDoubleArenaWithFoodFinder(10, 10, func(*Arena, []int) bool {
		return true
	})

	if p, ok := <-pointsDouble; ok && p != a.Food.Points {
		t.Fatalf("Value %d was expected but got %d", a.Food.Points, p)
	}

	a.moveSnake()
}

func TestDoesNotAddPointsWhenFoodNotFound(t *testing.T) {
	a := newDoubleArenaWithFoodFinder(10, 10, func(*Arena, []int) bool {
		return false
	})

	select {
	case p, _ := <-points:
		t.Fatalf("No point was expected to be received but received %d", p)
	default:
		close(points)
	}

	a.moveSnake()
}

func TestDoesNotPlaceNewFoodWhenFoodNotFound(t *testing.T) {
	a := newDoubleArenaWithFoodFinder(10, 10, func(*Arena, []int) bool {
		return false
	})

	f := a.Food

	a.moveSnake()

	if a.Food.X != f.X || a.Food.Y != f.Y {
		t.Fatal("Food in Arena expected not to have changed")
	}
}

func TestDoesNotIncreaseSnakeLengthWhenFoodNotFound(t *testing.T) {
	a := newDoubleArenaWithFoodFinder(10, 10, func(*Arena, []int) bool {
		return false
	})

	l := a.Snake.Length

	a.moveSnake()

	if a.Snake.Length != l {
		t.Fatal("Expected Snake not to have grown")
	}
}