package main

var fragmentShader = `
#version 420

out vec4 outputColor;

void main()
{
	outputColor = vec4(1.0f, 1.0f, 1.0f, 1.0f);
}

`

var vertexShader = `
#version 420

layout(location = 0) in vec4 position;

uniform mat4 ctcm;
uniform mat4 wtcm;
uniform mat4 mtwm;


void main()
{
	vec4 temp = mtwm * position;
	temp = wtcm * temp;
	gl_Position =ctcm * temp;
}
`
