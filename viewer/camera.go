package main

import (
	gls "github.com/vron/fm/gl"
//	"log"
	"math"
)

var (
	latpt gls.Vec3 = gls.Vec3{1,0,0}
	// TODO: Renormalise this hgut eacah time we change view!
	updir gls.Vec3 = gls.Vec3{0,0,1}
	theta, phi, r float64 = math.Pi/2, 0.0, 15
)

func getCamPos() gls.Vec3 {
	pos  := gls.Vec3{latpt[0], latpt[1], latpt[2]}
	pos[0] += float32(r*math.Cos(phi)*math.Sin(theta))
	pos[1] += float32(r*math.Sin(phi)*math.Sin(theta))
	pos[2] += float32(r*math.Cos(theta))
	return pos
}

func getLookAtMatrix() gls.Mat4 {
	lookDir := latpt.Copy()
	lookDir.Minus(getCamPos())
	lookDir.Normalize()
	updir.Normalize()
	
	rightDir := gls.Cross(lookDir, updir)
	rightDir.Normalize()
	perpUpDir := gls.Cross(rightDir, lookDir)

	rotMat := gls.Identity()
	rotMat.SetCol(0,gls.NewVec4(rightDir,0))
	rotMat.SetCol(1,gls.NewVec4(perpUpDir,0))
	rotMat.SetCol(2,gls.NewVec4(lookDir.Neg(),0))

	/*log.Println("updir: ", updir)
	log.Println("rightdir: ", rightDir)
	log.Println("perpupdir: ",perpUpDir)
	log.Println("lookdir: ", lookDir)
	log.Println(rotMat)
*/
	rotMat.Transpose()
	//log.Println(rotMat)


	transMat := gls.Identity()
	camPos := getCamPos()
	transMat.SetCol(3,gls.NewVec4(camPos.Neg(),1))
	//log.Println(transMat)
	if true {
		return rotMat.Times(transMat)
	} 
	//transMat = gls.Identity()
	return transMat

}