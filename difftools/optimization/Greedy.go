package optimization

import (
	"os"
	"fmt"
	diff "m/difftools/diffusion"
	"math"
	"bufio"
	"strings"
	"strconv"
	"math/rand"

)

func Greedy(seed int64, sample_size int, adj [][]int, Seed_set []int, prob_map [2][2][2][2]float64, pop [2]int, interest_list [][]int, assum_list [][]int, ans_len int, Count_true bool, sample_size2 int) ([]int, float64, []float64) {
	//sample_size2はグリーディで求めた解をより詳しくやる
	var n int = len(adj)
	var max float64 = 0
	var result float64
	var index int
	var ans []int
	var ans_v []float64

	ans = make([]int, 0, ans_len)
	S := make([]int, len(Seed_set))
	_ = copy(S, Seed_set)
	S_test := make([]int, len(Seed_set))
	_ = copy(S_test, Seed_set)

	var info_num int

	if Count_true {
		info_num = 2
	} else {
		info_num = 1
	}

	for i := 0; i < ans_len; i++ {
		fmt.Println(i)
		max = 0
		for j := 0; j < n; j++ {
			if (j+1)%100 == 0 {
				fmt.Println(i, "-", (j+1)/100)
			}
			_ = copy(S_test, S)
			if S_test[j] != 0 { //すでに発信源のユーザだったら
				continue
			}
			S_test[j] = info_num

			dist := Infl_prop_exp(seed, sample_size, adj, S_test, prob_map, pop, interest_list, assum_list)
			if Count_true {
				result = dist[diff.InfoType_T]
			} else {
				result = dist[diff.InfoType_F]
			}

			if result > max {
				max = result
				index = j
			}
		} //subloop end

		ans = append(ans, index)
		ans_v = append(ans_v, max)
		S[index] = info_num

	} //mainloop end

	// var max_2 float64
	// dist2 := Infl_prop_exp(seed, sample_size2, adj, S, prob_map, pop, interest_list, assum_list)
	// if Count_true {
	// 	max_2 = dist2[diff.InfoType_T]
	// } else {
	// 	max_2 = dist2[diff.InfoType_F]
	// }
	return ans, max, ans_v
}

type User_Infl struct {
    users []int
    infl  float64
}

type ByInfl []User_Infl

func (a ByInfl) Len() int           { return len(a) }
func (a ByInfl) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByInfl) Less(i, j int) bool { return a[i].infl < a[j].infl }

func Greedy_exp(seed int64, sample_size int, adj [][]int, Seed_set []int, prob_map [2][2][2][2]float64, pop [2]int, interest_list [][]int, assum_list [][]int, ans_len int, Count_true bool, capacity float64, max_user int, OnlyInfler bool, user_weight float64, use_kaiki bool)([]int, float64) {

	var costcal func(float64, float64,[][]int,int,int) float64
	if use_kaiki{
		costcal = Cal_cost_kaiki
	}else{
		costcal = Cal_cost
	}
	var n int = len(adj)
	var max float64 = 0
	var result float64
	var index int
	var ans []int

	ans = make([]int, 0, ans_len)
	S := make([]int, len(Seed_set))
	_ = copy(S, Seed_set)
	S_test := make([]int, len(Seed_set))
	_ = copy(S_test, Seed_set)
	cap_use := capacity

	var info_num int

	if Count_true {
		info_num = 2
	} else {
		info_num = 1
	}

	for{
		max = -1
		for j := 0; j < n; j++ {

			_ = copy(S_test, S)//初期化
			for i:=0; i< len(ans); i++{//初期設定
				S_test[ans[i]] = info_num
			}
			if S_test[j] != 0 { //すでに発信源のユーザだったら
				continue
			}
			if(OnlyInfler){
				if(FolowerSize(adj,j)==10000000000000){
					continue
				}
			}
			if costcal(user_weight, 1-user_weight, adj, j, max_user)>cap_use{//コストが大きすぎるユーザなら
				continue
			}
			S_test[j] = info_num
			rand.Seed(100)//おそらく後で消す　重要
			dist := Infl_prop_exp(seed, sample_size, adj, S_test, prob_map, pop, interest_list, assum_list)
			if Count_true {
				result = dist[diff.InfoType_T]/costcal(user_weight,1-user_weight,adj,j,max_user)
			} else {
				result = dist[diff.InfoType_F]/costcal(user_weight,1-user_weight,adj,j,max_user)
			}

			if result > max {
				max = result
				index = j
			}
		} //subloop end

		if max == -1{
			break;
		}
		ans = append(ans, index)
		cap_use -= costcal(user_weight,1-user_weight,adj,index,max_user)

		S[index] = info_num
	}
	return ans, max
}

func Cal_cost(u_weight float64, f_wight float64,adj [][]int,node int, max_user int)float64{
	f := FolowerSize(adj,node)
	u_max := len(adj)
	f_max := 0
	for i:=0; i<u_max; i++{
		if(adj[max_user][i] == 1){
			f_max ++
		}
	}
	// fmt.Println("folowersize",f)
	// fmt.Println("f:",f,"f_wight:",f_wight,"f_max:",f_max,"u_weight:",u_weight)
	return float64(f)*f_wight/float64(f_max) + u_weight/2
}

func Cal_cost_kaiki(u_weight float64, f_wight float64,adj [][]int,node int, max_user int)float64{
	return 100.0
	f := FolowerSize(adj,node)
	file, err := os.Open("kaiki.txt")
	if err != nil {
		fmt.Println("ファイルを開く際にエラーが発生しました:", err)
		return -1
	}
	defer file.Close()

	// ファイルを読み込む
	scanner := bufio.NewScanner(file)
	var line string
	if scanner.Scan() {
		line = scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("ファイルを読み取る際にエラーが発生しました:", err)
		return -1
	}

	// コンマで分割して整数に変換する
	parts := strings.Split(line, ",")
	if len(parts) != 2 {
		fmt.Println("予期しないデータ形式です")
		return -1
	}

	intercept, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		fmt.Println("最初の部分を整数に変換する際にエラーが発生しました:", err)
		return -1
	}

	slope, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		fmt.Println("二番目の部分を整数に変換する際にエラーが発生しました:", err)
		return -1
	}
	// fmt.Println(math.Log(float64(f))*slope +intercept)
	return math.Log(float64(f))*slope +intercept
}
