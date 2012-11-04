package main

import (
	"fmt"
	gl "github.com/chsc/gogl/gl21"
	"github.com/jteeuwen/glfw"
	"log"
	"os"
)

const (
	Title  = "Spinning Gopher"
	Width  = 640
	Height = 480
)

func main() {
	if err := glfw.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "glfw: %s\n", err)
		return
	}
	defer glfw.Terminate()

	//glfw.OpenWindowHint(glfw.Windowed, 1)

	if err := glfw.OpenWindow(Width, Height, 0, 0, 0, 0, 16, 0, glfw.Windowed); err != nil {
		fmt.Fprintf(os.Stderr, "glfw: %s\n", err)
		return
	}
	defer glfw.CloseWindow()

	glfw.SetSwapInterval(1)
	glfw.SetWindowTitle(Title)

	if err := gl.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "gl: %s\n", err)
	}

	if err := initScene(); err != nil {
		fmt.Fprintf(os.Stderr, "init: %s\n", err)
		return
	}
	defer destroyScene()

	for glfw.WindowParam(glfw.Opened) == 1 {
		drawScene()
		glfw.SwapBuffers()
	}
}

func initScene() error {
	// Vertex buffers
	gl.GenBuffers(1, &positionBufferObject)
	
	gl.BindBuffer(gl.ARRAY_BUFFER, positionBufferObject)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(vertexPositions), vertexPositions, gl.STATIC_DRAW)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)


	// InitializeProgram
	shaderList := make([]gl.Uint,0)
	shaderList = append(shaderList, CreateShader(gl.VERTEX_SHADER, vertexShader))
	shaderList = append(shaderList, CreateShader(gl.FRAGMENT_SHADER, fragmentShader))

	theProg := CreateProgram(shaderList)
	for _,v := range shaderList {
		gl.DeleteShader(v)
	}
}

func destroyScene() {

}

// Shaders
var (
	vertexShader string = `
#version 330

layout(location = 0) in vec4 position;
void main()
{
	gl_Position = position;
}
`
	fragmentShader string = `
#version 330

out vec4 outputColor;
void main()
{
	outputColor = vec4(1.0f,1.0f,1.0f,1.0f)
}
`
)

var vertexPositions = []float32{
	0.75, 0.75, 0, 1,
	0.75, -0.75, 0.0, 1.0,
	-0.75, -0.75, 0.0, 1.0,
}

func drawScene() {
	gl.ClearColor(0,0,0,0)
	gl.Clear(gl.COLOR_BUFFER_BIT);

	gl.UseProgram(myProgram)

	gl.BindBuffer(gl.ARRAY_BUFFER, positionBufferObject)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0,4,gl.FLOAT, gl.FALSE, 0, nil)

	gl.DrawArrays(gl.TRIANGLES, 0, 3)

	gl.DisableVertexAttribArray(0)
	gl.UseProgram(0)

}

func CreateShader(typ gl.Enum, code string) gl.Uint {
	shader := gl.CreateShader(typ)
	gl.ShadeModel(shader, 1, code, 0)

	gl.CompileShader(shader)

	var status int
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	
	if (status == gl.FALSE) {
		var length gl.Int
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &length)
		
		logStr := new(gl.Char[length+1])
		gl.GetShaderInfoLog(shader, length, 0, logStr)
		log.Println(logStr)
	}
	return shader
}

func CreateProgram(shaderList []gl.Uint) gl.Uint {
	program := gl.CreateProgram()

	for _, v := range shaderList {
		gl.AttachShader(program, v)
	}

	gl.LinkProgram(program)

	var status gl.Int
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if (status == gl.FALSE) {
		var length gl.Int
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &length)
		
		logStr := new(gl.Char[length+1])
		gl.GetProgramInfoLog(program, length, 0, logStr)
		log.Println(logStr)
	}
	return program
}
