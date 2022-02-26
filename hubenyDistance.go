package gomap

import (
	"math"
)

func HubenyDistance(src, dst Node) float64 {
	dx := degree2radian(dst.Lon - src.Lon)
	dy := degree2radian(dst.Lat - src.Lat)
	my := degree2radian((src.Lat + dst.Lat) / 2)

	W := math.Sqrt(1 - (power2(ECCENTRICITY) * power2(math.Sin(my)))) // 卯酉線曲率半径の分母
	m_numer := EQUATORIAL_RADIUS * (1 - power2(ECCENTRICITY))         // 子午線曲率半径の分子

	M := m_numer / math.Pow(W, 3) // 子午線曲率半径
	N := EQUATORIAL_RADIUS / W    // 卯酉線曲率半径

	d := math.Sqrt(power2(dy*M) + power2(dx*N*math.Cos(my)))

	return d
}

func degree2radian(x float64) float64 {
	return x * math.Pi / 180
}

func power2(x float64) float64 {
	return math.Pow(x, 2)
}

const (
	EQUATORIAL_RADIUS = 6378137.0            // 赤道半径 GRS80
	POLAR_RADIUS      = 6356752.314          // 極半径 GRS80
	ECCENTRICITY      = 0.081819191042815790 // 第一離心率 GRS80
)

type Point struct {
	Lat float64
	Lon float64
}
