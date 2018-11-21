package main

import (
	"github.com/fmi/go-homework/geom"
)

const Episilon = 1e-7

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
	// Find vectors for two edges sharing vertex 0
	edge1 := geom.Sub(triangle.b, triangle.a)
	edge2 := geom.Sub(triangle.c, triangle.a)

	// Begin calculating determinant
	h := geom.Cross(ray.Direction, edge2)

	det := geom.Dot(edge1, h)
	if det > -Episilon && det < Episilon {
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
	if t > Episilon {
		return true
	}

	return false
}

func (quad Quad) Intersect(ray geom.Ray) bool {
	return true
}

func (sphere Sphere) Intersect(ray geom.Ray) bool {
	return true
}
