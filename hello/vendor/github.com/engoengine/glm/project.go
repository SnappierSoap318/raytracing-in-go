package glm

import (
	"github.com/EngoEngine/math"
)

// Ortho returns a Mat4 that represents a orthographic projection from the given
// arguments.
func Ortho(left, right, bottom, top, near, far float32) Mat4 {
	rml, tmb, fmn := 1/(right-left), 1/(top-bottom), 1/(far-near)

	return Mat4{
		2 * rml, 0, 0, 0,
		0, 2 * tmb, 0, 0,
		0, 0, -2 * fmn, 0,
		-(right + left) * rml, -(top + bottom) * tmb, -(far + near) * fmn, 1,
	}
}

// Ortho2D is equivalent to Ortho with the near and far planes being -1 and 1,
// respectively.
func Ortho2D(left, right, bottom, top float32) Mat4 {
	return Ortho(left, right, bottom, top, -1, 1)
}

// Perspective returns a Mat4 representing a perspective projection of the given
// arguments.
func Perspective(fovy, aspect, near, far float32) Mat4 {
	nmf, f := 1/(near-far), 1./math.Tan(fovy/2.0)

	return Mat4{
		f / aspect, 0, 0, 0,
		0, f, 0, 0,
		0, 0, (near + far) * nmf, -1,
		0, 0, (2. * far * near) * nmf, 0,
	}
}

// Frustum returns a Mat4 representing a frustrum transform (squared pyramid with the top cut off)
func Frustum(left, right, bottom, top, near, far float32) Mat4 {
	rml, tmb, fmn := 1/(right-left), 1/(top-bottom), 1/(far-near)
	A, B, C, D := (right+left)*rml, (top+bottom)*tmb, -(far+near)*fmn, -(2*far*near)*fmn

	return Mat4{
		(2 * near) * rml, 0, 0, 0,
		0, (2 * near) * tmb, 0, 0,
		A, B, C, -1,
		0, 0, D, 0,
	}
}

// LookAt returns a Mat4 that represents a camera transform from the given
// arguments.
func LookAt(eyeX, eyeY, eyeZ, centerX, centerY, centerZ, upX, upY, upZ float32) Mat4 {
	return LookAtV(
		&Vec3{eyeX, eyeY, eyeZ},
		&Vec3{centerX, centerY, centerZ},
		&Vec3{upX, upY, upZ},
	)
}

// LookAtV generates a transform matrix from world space into the specific eye
// space.
func LookAtV(eye, center, up *Vec3) Mat4 {
	var f Vec3
	f.SubOf(center, eye)
	f.Normalize()
	var nup Vec3
	nup.SetNormalizeOf(up)
	var s Vec3
	s.CrossOf(&f, &nup)
	s.Normalize()
	var u Vec3
	u.CrossOf(&s, &f)

	M := Mat4{
		s[0], u[0], -f[0], 0,
		s[1], u[1], -f[1], 0,
		s[2], u[2], -f[2], 0,
		0, 0, 0, 1,
	}

	t := Translate3D(-eye[0], -eye[1], -eye[2])
	return M.Mul4(&t)
}

// Project transforms a set of coordinates from object space (in obj) to window
// coordinates (with depth)
//
// Window coordinates are continuous, not discrete, so you won't get exact pixel
// locations without rounding.
func Project(obj *Vec3, modelview, projection *Mat4, initialX, initialY, width, height int) Vec3 {
	obj4 := obj.Vec4(1)

	pm := projection.Mul4(modelview)
	vpp := pm.Mul4x1(&obj4)
	return Vec3{
		float32(initialX) + (float32(width)*(vpp[0]+1))*0.5,
		float32(initialY) + (float32(height)*(vpp[1]+1))*0.5,
		(vpp[2] + 1) * 0.5,
	}
}

// UnProject transforms a set of window coordinates to object space. If your MVP
// matrix is not invertible this will return garbage.
//
// Note that the projection may not be perfect if you use strict pixel locations
// rather than the exact values given by Project.
func UnProject(win *Vec3, modelview, projection *Mat4, initialX, initialY, width, height int) Vec3 {
	pm := projection.Mul4(modelview)
	inv := pm.Inverse()

	obj4 := inv.Mul4x1(&Vec4{
		(2 * (win[0] - float32(initialX)) / float32(width)) - 1,
		(2 * (win[1] - float32(initialY)) / float32(height)) - 1,
		2*win[2] - 1,
		1.0,
	})
	obj := obj4.Vec3()

	//if obj4[3] > MinValue {}
	over := 1 / obj4[3]
	obj[0] *= over
	obj[1] *= over
	obj[2] *= over

	return obj
}