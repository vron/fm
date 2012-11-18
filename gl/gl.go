// This package works as a wrapper of a wrapper for openGL to
// allow for usage of go types directly in function calls, decreasing
// bookkeeping significantely
// Note that only funtion  thus have have needed are implemented!
package gl

import (
	"log"
	"math"
	"unsafe"

	// TODO: Update to use later version of GL...
	g "github.com/chsc/gogl/gl42"
)

type Vec3 [3]float32
type Vec4 [4]float32

// Stored in column major order
type Mat4 [4 * 4]float32

func NewVec4(v Vec3, f float32) Vec4 {
	return Vec4{
		v[0], v[1], v[2], f}
}

func Identity() Mat4 {
	return Mat4{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}
}

func Ind4(r, c int) int {
	return r + 4*c
}

func (m *Mat4) SetCol(c int, v Vec4) {
	m[Ind4(0, c)] = v[0]
	m[Ind4(1, c)] = v[1]
	m[Ind4(2, c)] = v[2]
	m[Ind4(3, c)] = v[3]
}

func (m1 *Mat4) Times(m2 Mat4) Mat4 {
	n := Identity()
	for r := 0; r < 4; r++ {
		for c := 0; c < 4; c++ {
			n[Ind4(r, c)] = 0
			for i := 0; i < 4; i++ {
				n[Ind4(r, c)] += m1[Ind4(r, i)] * m2[Ind4(i, c)]
			}
		}
	}
	m1 = &n
	return n
}

func (m *Mat4) Transpose() Mat4 {
	for r := 1; r < 4; r++ {
		for c := 0; c < r; c++ {
			t := m[Ind4(r, c)]
			m[Ind4(r, c)] = m[Ind4(c, r)]
			m[Ind4(c, r)] = t
		}
	}
	return *m
}

func Cross(a, b Vec3) Vec3 {
	return Vec3{
		a[1]*b[2] - a[2]*b[1],
		a[2]*b[0] - a[0]*b[2],
		a[0]*b[1] - a[1]*b[0],
	}
}

func (v *Vec3) Minus(v2 Vec3) {
	v[0] -= v2[0]
	v[1] -= v2[1]
	v[2] -= v2[2]
}

func (v *Vec3) Neg() Vec3 {
	v[0] = -v[0]
	v[1] = -v[1]
	v[2] = -v[2]
	return *v
}

func (v *Vec3) Normalize() {
	s := float32(math.Sqrt(float64(v[0]*v[0] + v[1]*v[1] + v[2]*v[2])))
	v[0] /= s
	v[1] /= s
	v[2] /= s
}

func (v *Vec3) Copy() Vec3 {
	pos := Vec3{v[0], v[1], v[2]}
	return pos
}

func GetUniformLocation(p g.Uint, s string) g.Int {
	cstr := g.GLString(s)
	defer g.GLStringFree(cstr)
	return g.GetUniformLocation(p, cstr)
}

// Loads a shader from a string
func LoadShader(typ g.Enum, s string) g.Uint {
	shader := g.CreateShader(typ)
	cst := g.GLString(s)

	defer g.GLStringFree(cst)
	g.ShaderSource(shader, 1, &cst, nil)

	g.CompileShader(shader)

	var status g.Int
	g.GetShaderiv(shader, g.COMPILE_STATUS, &status)

	if status == g.FALSE {
		var length g.Int
		g.GetShaderiv(shader, g.INFO_LOG_LENGTH, &length)

		log.Println(int(length))
		logStr := g.GLStringAlloc(g.Sizei(length + 1))
		// TODO: defer
		g.GetShaderInfoLog(shader, g.Sizei(length), nil, logStr)
		log.Println(g.GoString(logStr))
	}
	return shader
}

func CreateProgram(shaders []g.Uint) g.Uint {
	p := g.CreateProgram()

	for _, v := range shaders {
		g.AttachShader(p, v)
	}
	g.LinkProgram(p)

	var status g.Int
	g.GetProgramiv(p, g.LINK_STATUS, &status)
	if status == g.FALSE {
		var length g.Int
		g.GetProgramiv(p, g.INFO_LOG_LENGTH, &length)
		log.Println(length)
		logStr := g.GLStringAlloc(g.Sizei(length + 1))
		g.GetProgramInfoLog(p, g.Sizei(length), nil, logStr)
		log.Println(g.GoString(logStr))
	}
	return p
}

func UniformMatrix4fv(u g.Int, count int, t bool, v Mat4) {
	p := (unsafe.Pointer)(&v[0])
	g.UniformMatrix4fv(u, 1, g.FALSE, (*g.Float)(p))
}

func Perspective(fovy, aspect, znear, zfar float32) Mat4 {
	m := Identity()

	ymax := znear * float32(math.Tan(float64(fovy*math.Pi/360)))
	ymin := -ymax
	xmax := ymax * aspect
	xmin := ymin * aspect

	width := xmax - xmin
	height := ymax - ymin

	depth := zfar - znear
	q := -(zfar + znear) / depth
	qn := -2 * (zfar * znear) / depth

	w := 2 * znear / width
	//w = w / aspect;
	h := 2 * znear / height

	m[0] = w
	m[1] = 0
	m[2] = 0
	m[3] = 0

	m[4] = 0
	m[5] = h
	m[6] = 0
	m[7] = 0

	m[8] = 0
	m[9] = 0
	m[10] = q
	m[11] = -1

	m[12] = 0
	m[13] = 0
	m[14] = qn
	m[15] = 0
	return m
}
