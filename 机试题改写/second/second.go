package main

import "fmt"

func Share(situa *[100][100]int, row, col, i, j int) {
	if i+1 >= 0 && i+1 < row && j < col && j >= 0 {
		(*situa)[i+1][j] = 1 //situa[i+1][j]  或者 *(situa[i+1] + j) 或者 *(*(situa+i+1)+j)
	}
	if j >= 0 && j < col && i-1 < row && i-1 >= 0 {
		(*situa)[i-1][j] = 1
	}
	if j+1 < col && j+1 >= 0 && i >= 0 && i < row {
		(*situa)[i][j+1] = 1
	}
	if j-1 >= 0 && j-1 < col && i < row && i >= 0 {
		(*situa)[i][j-1] = 1
	}
}

var (
	row  int
	col  int
	i    int
	j    int
	flag int
)
var (
	situa [100][100]int
	tmp   [100][100]int
)

func main() {
	count := 1
	var a int
	fmt.Scanf("%d%d", &row, &col)
	for i = 0; i < row; i++ {
		for j = 0; j < col; j++ {
			fmt.Scanf("%d", &situa[i][j])
		}
		fmt.Scanln(&a)
	}

	for count = 1; ; count++ {
		flag = 1
		for i = 0; i < row; i++ {
			for j = 0; j < col; j++ {
				if count == 1 {
					tmp[i][j] = situa[i][j]
				} else {
					situa[i][j] = tmp[i][j]
				}

			}
		}
		for i = 0; i < row; i++ {
			for j = 0; j < col; j++ {
				if situa[i][j] == 1 {
					Share(&tmp, row, col, i, j)
				}
			}

		}
		for i = 0; i < row; i++ {
			for j = 0; j < col; j++ {
				if tmp[i][j] == 0 {
					flag = 0
				}
			}

		}
		if flag == 1 {
			break
		}
	}
	fmt.Printf("%d", count)
}
