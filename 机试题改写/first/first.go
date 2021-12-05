package main
import "fmt"

var (
	seconds int
	minutes int
	hours int
	addsec int
)
func main(){
	fmt.Printf("请输入小时：")
	fmt.Scanf("%d\n", &hours)
	fmt.Printf("请输入分钟：")
	fmt.Scanf("%d\n", &minutes)
	fmt.Printf("请输入秒数：")
	fmt.Scanf("%d\n", &seconds)
	fmt.Printf("请输入增加的秒数：")
	fmt.Scanf("%d\n", &addsec)
	var tmp = seconds + addsec;
	if tmp < 60 {
		seconds = tmp;
	} else {
		minutes = minutes + tmp / 60
		hours = hours + minutes/60
	}
	if hours >= 24{
		hours -=24
	}
	if minutes >= 60{
		minutes -=60
	}
	if tmp >= 60{
		seconds = tmp - 60
	}
	fmt.Printf("%d %d %d",hours, minutes, seconds)
}
