package pgm

import (
	"math"
)

// Ported from https://github.com/bkiers/RotatingCalipers.

const sigma = 0.00000000001

type Caliper struct {
	hull         []point
	pointIndex   int
	currentAngle float64
}

func newCaliper(hull []point, pointIndex int, currentAngle float64) *Caliper {
	return &Caliper{
		hull:         hull,
		pointIndex:   pointIndex,
		currentAngle: currentAngle,
	}
}

func (c *Caliper) getAngleNextPoint() float64 {
	p1 := c.hull[c.pointIndex]
	p2 := c.hull[(c.pointIndex+1)%len(c.hull)]

	dx := p2.x - p1.x
	dy := p2.y - p1.y

	angle := math.Atan2(dy, dx) * 180 / math.Pi

	if angle < 0 {
		return 360 + angle
	}
	return angle
}

func (c *Caliper) getConstant() float64 {
	p := c.hull[c.pointIndex]
	return p.y - (c.getSlope() * p.x)
}

func (c *Caliper) getDeltaAngleNextPoint() float64 {
	angle := c.getAngleNextPoint()

	if angle < 0 {
		angle = 360 + angle - c.currentAngle
	} else {
		angle = angle - c.currentAngle
	}

	if angle < 0 {
		return 360
	}
	return angle
}

func (c *Caliper) getIntersection(d *Caliper) point {
	var p point
	switch {
	case c.isVertical():
		p.x = c.hull[c.pointIndex].x
		p.y = d.getConstant()
	case c.isHorizontal():
		p.x = d.hull[d.pointIndex].x
		p.y = c.getConstant()
	default:
		p.x = (d.getConstant() - c.getConstant()) / (c.getSlope() - d.getSlope())
		p.y = (c.getSlope() * p.x) + c.getConstant()
	}

	return p
}

func (c *Caliper) getSlope() float64 {
	return math.Tan(c.currentAngle * math.Pi / 180)
}

func (c *Caliper) isHorizontal() bool {
	return (math.Abs(c.currentAngle) < sigma) || (math.Abs(c.currentAngle-180) < sigma)
}

func (c *Caliper) isVertical() bool {
	return (math.Abs(c.currentAngle-90) < sigma) || (math.Abs(c.currentAngle-270) < sigma)
}

func (c *Caliper) rotateBy(angle float64) {
	if c.getDeltaAngleNextPoint() == angle {
		c.pointIndex = (c.pointIndex + 1) % len(c.hull)
	}
	c.currentAngle = math.Mod(c.currentAngle+angle, 360)
}
