package main

import (
	"flag"
	"log"
	"runtime"
	"time"

	gl "github.com/chsc/gogl/gl42"
	"github.com/jteeuwen/glfw"
	gls "github.com/vron/fm/gl"
)

var (
	theProgram gl.Uint
	mtwm_unif  gl.Int
	wtcm_unif  gl.Int
	ctcm_unif  gl.Int
)

var (
	fScript string
)

func init() {
	flag.StringVar(&fScript, "l", "", "File to run as we start")
}

func main() {
	log.Println("Starting")
	flag.Parse()

	runtime.GOMAXPROCS(4)

	initGLFW()
	initGL()

	initScene()

	initInput()

	glfw.SetWindowSizeCallback(resize)

	initLua()

	// Run the script file if we had one
	// TODO: Do this before accepting input from user...
	if fScript != "" {
		runFile(fScript)
	}

	for glfw.WindowParam(glfw.Opened) == 1 {
		display()
		glfw.SwapBuffers()
		time.Sleep(10 * time.Millisecond)
	}

}

func initGLFW() {
	log.SetFlags(log.Lshortfile)
	if err := glfw.Init(); err != nil {
		log.Fatalf("glfw: %s\n", err)
		return
	}
	//defer glfw.Terminate()

	//glfw.OpenWindowHint(glfw.Windowed, 1)

	Width := 800
	Height := 600
	if err := glfw.OpenWindow(Width, Height, 0, 0, 0, 0, 16, 0, glfw.Windowed); err != nil {
		log.Fatalf("glfw: %s\n", err)
		return
	}
	//defer glfw.CloseWindow()

	glfw.SetSwapInterval(2)
	glfw.SetWindowTitle("This is my title")
}

func initGL() {
	if err := gl.Init(); err != nil {
		log.Fatalf("gl: %s\n", err)
	}
}

func initScene() {
	initProgram()

	//gl.Enable(gl.CULL_FACE)
	//gl.CullFace(gl.BACK)
	gl.FrontFace(gl.CW)
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthMask(gl.TRUE)
	gl.DepthFunc(gl.LEQUAL)
	gl.DepthRange(0.0, 1.0)
	gl.Enable(gl.DEPTH_CLAMP)
}

func initProgram() {
	vs := gls.LoadShader(gl.VERTEX_SHADER, vertexShader)
	fs := gls.LoadShader(gl.FRAGMENT_SHADER, fragmentShader)
	theProgram = gls.CreateProgram([]gl.Uint{vs, fs})
	mtwm_unif = gls.GetUniformLocation(theProgram, "mtwm")
	wtcm_unif = gls.GetUniformLocation(theProgram, "wtcm")
	ctcm_unif = gls.GetUniformLocation(theProgram, "ctcm")
}

func resize(w, h int) {

	// TODO: Fix the near and far limits
	persMat := gls.Perspective(45, float32(w)/float32(h), 1, 1000)
	//persMat := gls.Identity()
	gl.UseProgram(theProgram)
	gls.UniformMatrix4fv(ctcm_unif, 1, false, persMat)

	//log.Println("Resizing ", persMat)

	gl.Viewport(0, 0, gl.Sizei(w), gl.Sizei(h))
}

type mS struct {
	Start, Length int
}

func display() {
	gl.ClearColor(0.19, 0.19, 0.21, 0)
	gl.ClearDepth(1)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	gl.UseProgram(theProgram)

	// Temporarily set the cam t clip
	id := gls.Identity()

	trans := getLookAtMatrix()
	//	log.Println(trans)
	gls.UniformMatrix4fv(wtcm_unif, 1, false, trans)

	gls.UniformMatrix4fv(mtwm_unif, 1, false, id)

	// And here we do the actual rendering

	// Generate a name for vufferobj and store it
	gl.GenBuffers(1, &positionBufferObject)
	// Bind it to a target
	gl.BindBuffer(gl.ARRAY_BUFFER, positionBufferObject)

	// Conver all nodes to a long array of vertexes!
	i := 0
	// TODO: Fix this, based on 4 per eleemnts
	vpos := make([]float32, len(elements)*4*4)
	elms := make(map[uint32]mS, len(elements))
	// Note that the order here is undefined!
	ii := 0
	for k, v := range elements {
		elms[k] = mS{ii, len(v.Nodes())}
		for _, vv := range v.Nodes() {
			ii++
			vpos[i] = float32(nodes[vv][0])
			i++
			vpos[i] = float32(nodes[vv][1])
			i++
			vpos[i] = float32(nodes[vv][2])
			i++
			vpos[i] = 1.0
			i++
		}
	}
	/*
		for _, v := range nodes {
			vpos[i] = float32(v[0])
			i++
			vpos[i] = float32(v[1])
			i++
			vpos[i] = float32(v[2])
			i++
		}
	*/
	gl.BufferData(gl.ARRAY_BUFFER, gl.Sizeiptr(32/8*len(vpos)), (gl.Pointer)(&vpos[0]), gl.STATIC_DRAW)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	gl.BindBuffer(gl.ARRAY_BUFFER, positionBufferObject)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 4, gl.FLOAT, gl.FALSE, 0, (gl.Pointer)(nil))

	/*	for _,v := range elms {

		}
	*/
	for _, v := range elms {
		gl.DrawArrays(gl.LINE_LOOP, gl.Int(v.Start), gl.Sizei(v.Length))
	}

	// Draw Coordinate acis
	/*for i := 9; i < 9+2*5; i += 2 {
		gl.DrawArrays(gl.LINE_STRIP, gl.Int(i), 2)
	}*/
}

var positionBufferObject gl.Uint
var vertexPositions = []float32{
	1, -1, -2, 1.0,
	-1, 1, -2, 1.0,
	-1, -1, -2, 1.0,
	1, -1, -2, 1.0,
	-1, 1, -2, 1.0,
	-1, 1, -3, 1.0,
	1, -1, -2, 1.0,
	-1, -1, -2, 1.0,
	-1, -1, -3, 1.0,
	0, 0, -1, 1,
	0, 0, 1, 1,
	0, 0, 0, 1,
	0, 1, 0, 1,
	0, 0, 0, 1,
	1, 0, 0, 1,
	1, 0, 0, 1,
	0.9, 0.1, 0, 1,
	1, 0, 0, 1,
	0.9, -0.1, 0, 1,
}
