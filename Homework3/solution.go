package main

import (
	"github.com/fmi/go-homework/geom"
	"sync"
)

const epsilon = 1e-7

type Triangle struct {
	a, b, c geom.Vector
}

type Quad struct {
	a, b, c, d geom.Vector
}

type Sphere struct {
	origin geom.Vector
	r      float64
}

func NewTriangle(a, b, c geom.Vector) Triangle {
	return Triangle{
		a: a,
		b: b,
		c: c,
	}
}

func NewQuad(a, b, c, d geom.Vector) Quad {
	return Quad{
		a: a,
		b: b,
		c: c,
		d: d,
	}
}

func NewSphere(origin geom.Vector, r float64) Sphere {
	return Sphere{
		origin: origin,
		r:      r,
	}
}

func (triangle Triangle) Intersect(ray geom.Ray) bool {
	// Find vectors for two edges sharing the first vertex
	edge1 := geom.Sub(triangle.b, triangle.a)
	edge2 := geom.Sub(triangle.c, triangle.a)

	// Begin calculating determinant
	h := geom.Cross(ray.Direction, edge2)

	det := geom.Dot(edge1, h)
	if det > -epsilon && det < epsilon {
		return false // The ray is parallel to triangle plane, impossible that they intersect
	}

	f := 1 / det

	// Calculate vector from vertex to the ray origin
	s := geom.Sub(ray.Origin, triangle.a)

	// Calculating U parameter
	u := f * geom.Dot(s, h)
	if u < 0 || u > 1 {
		return false
	}

	// Prepare to test V parameter
	q := geom.Cross(s, edge1)

	v := f * geom.Dot(ray.Direction, q)
	if v < 0 || u+v > 1 {
		return false
	}

	// Calculating t - final check to see if ray intersects triangle
	t := f * geom.Dot(edge2, q)
	if t > epsilon {
		return true
	}

	return false
}

func (quad Quad) Intersect(ray geom.Ray) bool {
	var firstTriangle, secondTriangle Triangle
	if quad.isConvex() {
		firstTriangle = Triangle{a: quad.a, b: quad.c, c: quad.b}
		secondTriangle = Triangle{a: quad.a, b: quad.c, c: quad.d}
	} else {
		firstTriangle = Triangle{a: quad.b, b: quad.d, c: quad.a}
		secondTriangle = Triangle{a: quad.b, b: quad.d, c: quad.c}
	}

	var wg sync.WaitGroup
	var mutex sync.Mutex
	var foundIntersection bool

	checkIntersection := func(triangle *Triangle) {
		defer wg.Done()
		result := triangle.Intersect(ray)
		mutex.Lock()
		foundIntersection = foundIntersection || result
		mutex.Unlock()
	}

	wg.Add(2)
	go checkIntersection(&firstTriangle)
	go checkIntersection(&secondTriangle)
	wg.Wait()

	return foundIntersection
}

func (quad Quad) isConvex() bool {
	var sign bool
	vertices := []geom.Vector{quad.a, quad.b, quad.c, quad.d}
	n := len(vertices)

	for i := 0; i < n; i++ {
		dx1 := vertices[(i+2)%n].X - vertices[(i+1)%n].X
		dy1 := vertices[(i+2)%n].Y - vertices[(i+1)%n].Y

		dx2 := vertices[i].X - vertices[(i+1)%n].X
		dy2 := vertices[i].Y - vertices[(i+1)%n].Y

		zcross := dx1*dy2 - dy1*dx2

		if i == 0 {
			sign = zcross > 0
		} else if sign != (zcross > 0) {
			return false
		}
	}

	return true
}

func (sphere Sphere) Intersect(ray geom.Ray) bool {
	oc := geom.Sub(ray.Origin, sphere.origin)

	a := geom.Dot(ray.Direction, ray.Direction)
	b := 2.0 * geom.Dot(oc, ray.Direction)
	c := geom.Dot(oc, oc) - (sphere.r * sphere.r)

	discriminant := b*b - 4*a*c

	return discriminant > 0
}
