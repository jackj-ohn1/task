package main

import "fmt"

func main() {
	var (
		row    int
		col    int
		i      int
		j      int
		re     [100][100]int
		scores [100]int
		answer [100]int
		sum    [100]int
	)
	fmt.Scanf("%d%d", &row, &col)
	fmt.Scanln()
	for i = 0; i < col; i++ {
		fmt.Scanf("%d", &scores[i])
	}
	fmt.Scanln()
	for i = 0; i < col; i++ {
		fmt.Scanf("%d", &answer[i])
	}
	fmt.Scanln()
	for i = 0; i < row; i++ {
		for j = 0; j < col; j++ {
			fmt.Scanf("%d", &re[i][j])
			if re[i][j] == answer[j] {
				sum[i] = sum[i] + scores[j]
			}
		}
		fmt.Scanln()
	}
	for i = 0; i < row; i++ {
		fmt.Printf("%d\n", sum[i])
	}
}
