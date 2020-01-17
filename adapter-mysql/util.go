package mysql

// snake string, XxYy to xx_yy , XxYY to xx_yy
func SnakeString(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
			d = d + 32
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	return string(data[:])
}

// camel string, xx_yy to xxYy
func CamelString(s string) string {
	data := make([]byte, 0, len(s))
	j := false

	num := len(s) - 1
	for i := 0; i <= num; i++ {
		d := s[i]
		if j {
			d = d - 32
			j = false
		}
		if d == '_' && num > i && s[i+1] >= 'a' && s[i+1] <= 'z' {
			j = true
			continue
		}
		data = append(data, d)
	}
	return string(data[:])
}
