package numpygo

import (
	"github.com/praveenpenumaka/numpygo/domain"
	"github.com/praveenpenumaka/numpygo/utils"
)

// TODO: Implement and add test cases
func Argmax(a NDArray, axis ...int) NDArray {
	if a.Size == 0 {
		return NDArray{}
	}
	if len(axis) == 0 {
		d := Ones("FLOAT64", 1)
		_, maxIndex := a.Elements.Max()
		d.Elements.Values[0] = float64(maxIndex)
		return d
	}
	tAxis := utils.GetAxis(axis)
	if tAxis > a.Dims {
		return NDArray{}
	}
	return DimensionReduction(func(v domain.Vector) float64 {
		_, maxIndex := v.Max()
		return float64(maxIndex)
	}, a, tAxis)
}

// return maximum along given axis
func Amax(a NDArray, axis ...int) NDArray {
	if a.Size == 0 {
		return NDArray{}
	}
	if len(axis) == 0 {
		d := Ones("FLOAT64", 1)
		d.Elements.Values[0], _ = a.Elements.Max()
		return d
	}
	tAxis := utils.GetAxis(axis)
	if tAxis > a.Dims {
		return NDArray{}
	}
	return DimensionReduction(func(v domain.Vector) float64 {
		max, _ := v.Max()
		return max
	}, a, tAxis)
}

func Argmin(a NDArray, axis ...int) NDArray {
	if a.Size == 0 {
		return NDArray{}
	}
	if len(axis) == 0 {
		d := Ones("FLOAT64", 1)
		_, minIndex := a.Elements.Max()
		d.Elements.Values[0] = float64(minIndex)
		return d
	}
	tAxis := utils.GetAxis(axis)
	if tAxis > a.Dims {
		return NDArray{}
	}
	return DimensionReduction(func(v domain.Vector) float64 {
		min, _ := v.Min()
		return min
	}, a, tAxis)
}

func Amin(a NDArray, axis ...int) NDArray {
	if a.Size == 0 {
		return NDArray{}
	}
	if len(axis) == 0 {
		d := Ones("FLOAT64", 1)
		d.Elements.Values[0], _ = a.Elements.Max()
		return d
	}
	tAxis := utils.GetAxis(axis)
	if tAxis > a.Dims {
		return NDArray{}
	}
	return DimensionReduction(func(v domain.Vector) float64 {
		min, _ := v.Min()
		return min
	}, a, tAxis)
}

// TODO: Verify this
func Unique(a NDArray, axis ...int) NDArray {
	if a.Size == 0 {
		return NDArray{}
	}
	if len(axis) == 0 {
		uniqueValues := *a.Elements.Unique()
		d := Ones("FLOAT64", 0)
		d.Elements.Values = append(d.Elements.Values, uniqueValues.Values...)
		d.Size = len(d.Elements.Values)
		d.Shape = domain.IVector{Values: []int{1}}
		d.Dims = 1
		d.Shape = domain.IVector{Values: []int{1}}
		return d
	}
	tAxis := utils.GetAxis(axis)
	if tAxis > a.Dims {
		return NDArray{}
	}
	return DimensionReduction(func(v domain.Vector) float64 {
		return v.Unique().Values[0]
	}, a, tAxis)
}

func Sum(a NDArray, axis ...int) NDArray {
	if a.Size == 0 {
		return NDArray{}
	}
	if len(axis) == 0 {
		d := Ones("FLOAT64", 1)
		d.Elements.Values[0] = a.Elements.Sum()
		return d
	}
	tAxis := utils.GetAxis(axis)
	if tAxis > a.Dims {
		return NDArray{}
	}
	return DimensionReduction(func(v domain.Vector) float64 {
		return v.Sum()
	}, a, tAxis)
}

func DimensionReduction(lambda func(v domain.Vector) float64, a NDArray, axis int) NDArray {
	newShape := a.Shape.Remove(axis)
	ndIndex := NewNDIndex(a.Shape.Values)
	newArray := Zeros(a.DType, newShape.Values...)
	if a.Dims == 1 {
		d := Ones("FLOAT64", 1)
		d.Elements.Values[0] = lambda(a.Elements)
		return d
	}
	for vector := ndIndex.Next(); vector != nil; vector = ndIndex.Next() {
		oldIndex, err := utils.GetIndexFromVector(vector, &a.Strides, &a.Shape)
		newVector := vector.Remove(axis)
		newIndex, err := utils.GetIndexFromVector(newVector, &newArray.Strides, &newArray.Shape)
		if err != nil {
			return NDArray{}
		}
		newArray.Elements.Values[newIndex] = lambda(domain.Vector{Values: []float64{newArray.Elements.Values[newIndex], a.Elements.Values[oldIndex]}})
	}
	return newArray
}