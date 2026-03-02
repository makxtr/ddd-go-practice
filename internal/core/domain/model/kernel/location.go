/**
Location - это координата на доске, она состоит из X (горизонталь) и Y (вертикаль)
Минимально возможная для установки координата 1,1
Максимально возможная для установки координата 10,10
2 координаты равны, если их X и Y равны, обеспечьте функционал проверки на эквивалентность
Нельзя изменять объект Location после создания
Должна быть возможность рассчитать расстояние между двумя Location. Расстояние между Location - это совокупное количество шагов по X и Y,
которое необходимо сделать курьеру, чтобы достигнуть точки. Курьер может двигаться только по вертикали и горизонтали.
На картинке ниже: расстояние между курьером и заказом - 2 шага по X и 3 шага по Y, суммарно 5 шагов.
Должна быть возможность создать рандомную координату. В будущем эта функциональность будет использована в целях тестирования

*/

package kernel

import (
	"delivery/internal/pkg/errs"
	"math/rand/v2"
)

const minCoordinate = 1
const maxCoordinate = 10

type Location struct {
	x     int
	y     int
	isSet bool
}

func NewLocation(x, y int) (Location, error) {
	if x < minCoordinate || x > maxCoordinate {
		return Location{}, errs.NewValueIsOutOfRangeError("x", x, minCoordinate, maxCoordinate)
	}
	if y < minCoordinate || y > maxCoordinate {
		return Location{}, errs.NewValueIsOutOfRangeError("y", y, minCoordinate, maxCoordinate)
	}

	return Location{x, y, true}, nil
}

func NewRandomLocation() Location {
	return Location{
		x:     rand.IntN(maxCoordinate-minCoordinate+1) + minCoordinate,
		y:     rand.IntN(maxCoordinate-minCoordinate+1) + minCoordinate,
		isSet: true,
	}
}

func (l Location) DistanceTo(target Location) (int, error) {
	if target.IsEmpty() {
		return 0, errs.NewValueIsRequiredError("target")
	}

	xDist := l.x - target.x
	if xDist < 0 {
		xDist = -xDist
	}

	yDist := l.y - target.y
	if yDist < 0 {
		yDist = -yDist
	}

	return xDist + yDist, nil
}

func (l Location) X() int {
	return l.x
}

func (l Location) Y() int {
	return l.y
}

func (l Location) Equals(other Location) bool {
	return l == other
}

func (l Location) IsEmpty() bool {
	return !l.isSet
}

func (l Location) IsValid() bool {
	return l.isSet
}
