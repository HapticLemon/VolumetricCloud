package Ruido

import "../Vectores"

/*vec3 snoiseVec3( vec3 x ){

float s  = snoise(vec3( x ));
float s1 = snoise(vec3( x.y - 19.1 , x.z + 33.4 , x.x + 47.2 ));
float s2 = snoise(vec3( x.z + 74.2 , x.x - 124.5 , x.y + 99.4 ));
vec3 c = vec3( s , s1 , s2 );
return c;

}*/

// Implementación sacada de :
//https://github.com/cabbibo/glsl-curl-noise/blob/master/curl.glsl
//
func snoiseVec3(punto Vectores.Vector) Vectores.Vector {
	var s float64 = Noise3(punto.X, punto.Y, punto.Z)
	var s1 float64 = Noise3(punto.Y-19.1, punto.Z+33.4, punto.X+47.2)
	var s2 float64 = Noise3(punto.Z+74.2, punto.X-124.5, punto.Y+99.4)

	return Vectores.Vector{X: s, Y: s1, Z: s2}
}

/*vec3 curlNoise( vec3 p ){

const float e = .1;
vec3 dx = vec3( e   , 0.0 , 0.0 );
vec3 dy = vec3( 0.0 , e   , 0.0 );
vec3 dz = vec3( 0.0 , 0.0 , e   );

vec3 p_x0 = snoiseVec3( p - dx );
vec3 p_x1 = snoiseVec3( p + dx );
vec3 p_y0 = snoiseVec3( p - dy );
vec3 p_y1 = snoiseVec3( p + dy );
vec3 p_z0 = snoiseVec3( p - dz );
vec3 p_z1 = snoiseVec3( p + dz );

float x = p_y1.z - p_y0.z - p_z1.y + p_z0.y;
float y = p_z1.x - p_z0.x - p_x1.z + p_x0.z;
float z = p_x1.y - p_x0.y - p_y1.x + p_y0.x;

const float divisor = 1.0 / ( 2.0 * e );
return normalize( vec3( x , y , z ) * divisor );

}*/

func CurlNoise(punto Vectores.Vector) Vectores.Vector {

	const e float64 = .01
	var dx = Vectores.Vector{e, 0.0, 0.0}
	var dy = Vectores.Vector{0.0, e, 0.0}
	var dz = Vectores.Vector{0.0, 0.0, e}

	var p_x0 Vectores.Vector = snoiseVec3(punto.Sub(dx))
	var p_x1 Vectores.Vector = snoiseVec3(punto.Add(dx))

	var p_y0 Vectores.Vector = snoiseVec3(punto.Sub(dy))
	var p_y1 Vectores.Vector = snoiseVec3(punto.Add(dy))

	var p_z0 Vectores.Vector = snoiseVec3(punto.Sub(dz))
	var p_z1 Vectores.Vector = snoiseVec3(punto.Add(dz))

	var x float64 = p_y1.Z - p_y0.Z - p_z1.Y + p_z0.Y
	var y float64 = p_z1.X - p_z0.X - p_x1.Z + p_x0.Z
	var z float64 = p_x1.Y - p_x0.Y - p_y1.X + p_y0.X

	const divisor float64 = 1.0 / (2.0 * e)

	return Vectores.Vector{x, y, z}.MultiplyByScalar(divisor).Normalize()

}
