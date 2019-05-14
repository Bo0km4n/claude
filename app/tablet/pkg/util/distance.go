package util

import "math"

func Distance(lat1 float64, lng1 float64, lat2 float64, lng2 float64) float64 {

	// 緯度経度をラジアンに変換
	rlat1 := lat1 * math.Pi / 180
	rlng1 := lng1 * math.Pi / 180
	rlat2 := lat2 * math.Pi / 180
	rlng2 := lng2 * math.Pi / 180

	// 2点の中心角(ラジアン)を求める
	a :=
		math.Sin(rlat1)*math.Sin(rlat2) +
			math.Cos(rlat1)*math.Cos(rlat2)*
				math.Cos(rlng1-rlng2)
	rr := math.Acos(a)

	earthRadius := 6378140. // 地球赤道半径(メートル)
	distance := earthRadius * rr
	return distance
}
