package main

import (
	"testing"

	"github.com/fmi/go-homework/geom"
)

func TestSampleSimpleTriangle(t *testing.T) {
	var prim geom.Intersectable

	a, b, c := geom.NewVector(-1, -1, 0), geom.NewVector(1, -1, 0), geom.NewVector(0, 1, 0)
	prim = NewTriangle(a, b, c)
	ray := geom.NewRay(geom.NewVector(0, 0, -1), geom.NewVector(0, 0, 1))

	if !prim.Intersect(ray) {
		t.Errorf("Expected ray %#v to intersect triangle %#v but it did not.", ray, prim)
	}
}

func TestSampleSimpleSphere(t *testing.T) {
	var prim geom.Intersectable

	origin, r := geom.NewVector(0, 0, 0), 2.0
	prim = NewSphere(origin, r)
	ray := geom.NewRay(geom.NewVector(0, 0, 2.5), geom.NewVector(0, 0, 0))

	if !prim.Intersect(ray) {
		t.Errorf("Expected ray %#v to intersect sphere %#v but it did not.", ray, prim)
	}
}

func TestSampleSimpleQuad(t *testing.T) {
	var prim geom.Intersectable

	a, b, c, d := geom.NewVector(0, -2, 0), geom.NewVector(3, 0, 0), geom.NewVector(0, 1, 0), geom.Vector{-1, 0, 0}
	prim = NewQuad(a, b, c, d)
	ray := geom.Ray{geom.NewVector(0, 0, 2), geom.NewVector(0, 0, -1)}

	if !prim.Intersect(ray) {
		t.Errorf("Expected ray %#v to intersect quad %#v but it did not.", ray, prim)
	}
}

func TestSampleIntersectableImplementations(t *testing.T) {
	var prim geom.Intersectable

	a, b, c, d := geom.NewVector(-1, -1, 0),
		geom.NewVector(1, -1, 0),
		geom.NewVector(0, 1, 0),
		geom.NewVector(-1, 1, 0)

	prim = NewTriangle(a, b, c)
	prim = NewQuad(a, b, c, d)
	prim = NewSphere(a, 5)

	_ = prim
}