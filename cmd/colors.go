package cmd

import "strconv"

type Hex string
type Ansi string

const AnsiReset = "\033[0m"

// The RGB type holds three values: one for red (R), green, (G) and
// blue (B). Each of these colors are on the domain of [0, 255].
type RGB struct {
	R int `json:"R"`
	G int `json:"G"`
	B int `json:"B"`
}

// HextoRGB converts a hexadecimal string to RGB values
func HextoRGB(hex Hex) RGB {
	if hex[0:1] == "#" {
		hex = hex[1:]
	}
	r := string(hex)[0:2]
	g := string(hex)[2:4]
	b := string(hex)[4:6]
	R, _ := strconv.ParseInt(r, 16, 0)
	G, _ := strconv.ParseInt(g, 16, 0)
	B, _ := strconv.ParseInt(b, 16, 0)

	return RGB{int(R), int(G), int(B)}
}

// HextoAnsi converts a hexadecimal string to an Ansi escape code
func HextoAnsi(hex Hex) Ansi {
	rgb := HextoRGB(hex)
	str := "\x1b[38;2;" + strconv.FormatInt(int64(rgb.R), 10) + ";" + strconv.FormatInt(int64(rgb.G), 10) + ";" + strconv.FormatInt(int64(rgb.B), 10) + "m"
	return Ansi(str)
}
