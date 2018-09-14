package daemon

import (
	"encoding/json"
	"fmt"
	"httpschedule2/engine"
	"io/ioutil"
	"math"
	"sort"
)

type ReceiveData struct {
	DistanceWeight    float64    `json:"distanceweight,omitempty"`
	DiffWeight        float64    `json:"diffweight,omitempty"`
	ResourceWeight    float64    `json:"resourceweight,omitempty"`
	NearWeight        float64    `json:"nearweight,omitempty"`
	SameWeight        float64    `json:"sameweight,omitempty"`
	ReceiveContainers Containers `json:"containers,omitempty"`
	ReceiveNode       Nodes      `json:"nodes,omitempty"`
	ReceiveContainer  Container  `json:"container,omitempty"`
}
type Container struct {
	Cpu      float64           `json:"cpu,omitempty"`   //容器cpu占用量
	Mem      float64           `json:"mem,omitempty"`   //容器内存占用量
	Image    string            `json:"image,omitempty"` //容器镜像名称
	Label    map[string]string `json:"label,omitempty"` //容器标签
	Hostnode string            `json:"hostnode"`        //容器所属node编号
}
type Node struct {
	Index        int       `json:"index,omitempty"`
	Name         string    `json:"nodename,omitempty"`
	Datadistance float64   `json:"datadistance,omitempty"` //与数据中心的距离
	Cpucap       float64   `json:"cpucap,omitempty"`       //cpu容量
	Memcap       float64   `json:"memcap,omitempty"`       //内存容量
	Distance     []float64 `json:"distance,omitempty"`     //与各节点的距离
}
type Containers struct {
	AllContainer []*Container `json:"allcontainer,omitempty"` //容器集合
}
type Nodes struct {
	AllNodes []*Node `json:"allnodes,omitempty"` //节点集合
}
type ReturnNode struct {
	NodeIndex    int       `json:"index,omitempty"`
	Nodename     string    `json:"name,omitempty"`
	Datadistance float64   `json:"datadistance,omitempty"` //与数据中心的距离
	Cpucap       float64   `json:"cpucap,omitempty"`       //cpu容量
	Memcap       float64   `json:"memcap,omitempty"`       //内存容量
	Distance     []float64 `json:"distance,omitempty"`
}

var faraway map[string]string //互斥标签对
var near map[string]string    //相近标签对
var same map[string]string    //亲和标签对

var resourcescore []float64 //资源得分
var dispersescore []float64 //互斥得分
var nearscore []float64     //相近得分
var samescore []float64     // 亲和得分
var nametonum map[string]int

func calculatedistancescore(allnode *Nodes, number int) []float64 { //计算与数据源的距离得分
	intscore := make([]float64, 0, number)
	for _, v := range allnode.AllNodes {
		intscore = append(intscore, v.Datadistance)
		//intscore[i] = v.Datadistance
		//println(v.Datadistance)
		//println(intscore[i])
	}
	//var datadistancescore []float64
	datadistancescore := make([]float64, number)
	//print(intscore)
	datadistancescore = normalizationfloat(intscore) //归一化操作
	//print(datadistancescore)
	for i, _ := range datadistancescore {
		datadistancescore[i] = 1 - datadistancescore[i]
	}
	//for _, v := range datadistancescore {
	//	println(v)
	//}
	//print(datadistancescore)
	printSlice(datadistancescore)
	return datadistancescore
}

func calculatediffscore(allcontainer *Containers, number int, faraway map[string]string, newcontainer *Container) []float64 { //计算分散性得分

	diff := make([]float64, number) //分散性变量中间得分
	var hostnode []int
	for _, v := range allcontainer.AllContainer {
		if v.Hostnode == "" {
			//fmt.Println("findnull")
			continue
		}
		if diff[nametonum[v.Hostnode]-1] == 0 {
			diff[nametonum[v.Hostnode]-1] = float64(comparestring(v.Image, newcontainer.Image) + comparemap(v.Label, newcontainer.Label))
		} else {
			if float64(comparestring(v.Image, newcontainer.Image)+comparemap(v.Label, newcontainer.Label)) > diff[nametonum[v.Hostnode]-1] {
				diff[nametonum[v.Hostnode]-1] = float64(comparestring(v.Image, newcontainer.Image) + comparemap(v.Label, newcontainer.Label))
			}
		}
		for key, _ := range faraway { //检测是否在互斥标签对内
			if (key == v.Label["selector"] && faraway[key] == newcontainer.Label["selector"]) || (key == newcontainer.Label["selector"] && faraway[key] == v.Label["selector"]) {
				//fmt.Println(key, faraway[key])
				//	fmt.Println(v.Hostnode)
				diff[nametonum[v.Hostnode]-1] = 3
				hostnode = append(hostnode, nametonum[v.Hostnode])
			}
		}
	}
	//fmt.Println(diff)
	nordiff := normalizationfloat(diff)
	//fmt.Println(nordiff)
	for i, _ := range nordiff {
		nordiff[i] = 1 - nordiff[i]
		if nordiff[i] == 0 {
			nordiff[i] += 0.1
		}
		for _, v := range hostnode {
			if i == v-1 {
				nordiff[i] = 0
			}
		}
	}
	printSlice(nordiff)
	//check farawaymap
	return nordiff
}

func calculateresourcescore(allcontainer *Containers, number int, newcontainer *Container) []float64 { //计算资源得分

	//calculate used resource
	usedresource := make([][2]float64, number)

	for _, v := range allcontainer.AllContainer {
		if v.Hostnode == "" {
			//	fmt.Println("findnull")
			continue
		}
		usedresource[nametonum[v.Hostnode]-1][0] += v.Cpu
		usedresource[nametonum[v.Hostnode]-1][1] += v.Mem
		//usedresource[v.Hostnode-1][2] += v.Store
		//	fmt.Println(nametonum[v.Hostnode] - 1)
	}
	//fmt.Println("1")
	//fmt.Println(usedresource)
	score := make([]float64, 0, number)
	for i := 0; i < number; i++ {
		score = append(score, math.Sqrt(math.Pow(math.Abs(float64(newcontainer.Cpu-usedresource[i][0])), 2)+math.Pow(math.Abs(float64(newcontainer.Mem-usedresource[i][1])), 2))) //求欧氏距离
	}
	printSlice(normalizationfloat(score))
	return normalizationfloat(score)
}

func calculatenearscore(allnode *Nodes, allcontainer *Containers, number int, near map[string]string, newcontainer *Container) []float64 { //计算相近得分

	neardata := make([]float64, number)
	var findnear string //查找的相近标签
	var find bool       //是否找到相近标签
	var max float64
	for k, v := range allnode.AllNodes {
		if k == 0 {
			max = findmaxfloat(v.Distance)
		} else {
			if findmaxfloat(v.Distance) > max {
				max = findmaxfloat(v.Distance)
			}
		}

	}
	for key, value := range near {
		find = false
		if key == newcontainer.Label["selector"] {
			findnear = value
			//	fmt.Println(findnear)
			find = true
		} else if value == newcontainer.Label["selector"] {
			findnear = key
			//fmt.Println(findnear)
			find = true
		}
		if find == true { //如果找到，则计算相近的最小值
			var hostnode []int
			findpipei := false
			// find the node has this container
			for _, v := range allcontainer.AllContainer {
				if v.Hostnode == "" {
					//fmt.Println("findnull")
					continue
				}
				if v.Label["selector"] == findnear {
					//	fmt.Println(v.Label, v.Hostnode)
					hostnode = append(hostnode, nametonum[v.Hostnode])
					findpipei = true
				}
			}
			//fmt.Println(hostnode)
			if findpipei == false {
				continue
			}
			//calculate the near score
			for k, v := range allnode.AllNodes {
				var min float64
				for i := 0; i < len(hostnode); i++ {
					if i == 0 {
						min = v.Distance[hostnode[i]-1]
					}
					if i > 0 {
						if v.Distance[hostnode[i]-1] < min {
							min = v.Distance[hostnode[i]-1]
						}
					}
				}
				neardata[k] += float64(min)
				//fmt.Println(k)
				if neardata[k] == 0 { //if this node has near containers
					neardata[k] = max
				}
			}
		}
	}
	//fmt.Println(neardata)
	normnear := normalizationfloat(neardata)
	for i, _ := range normnear {
		normnear[i] = 1 - normnear[i]
	}
	printSlice(normnear)
	return normnear
}
func calculatesamesocre(allcontainer *Containers, number int, same map[string]string, newcontainer *Container) ([]float64, bool) { //计算亲和性得分

	samedata := make([]float64, number)
	var findsame string //找到的亲和性标签
	var find bool       //是否找到
	finalfind := false
	for key, value := range same {
		//	fmt.Println(key, value)
		find = false
		if key == newcontainer.Label["selector"] {
			findsame = value
			//fmt.Println(findsame)
			find = true
		} else if value == newcontainer.Label["selector"] {
			findsame = key
			//fmt.Println(findsame)
			find = true
		}
		if find == true { //如果找到，则结果置1
			finalfind = true
			var hostnode []int
			// find the node has this container
			for _, v := range allcontainer.AllContainer {
				if v.Hostnode == "" {
					//fmt.Println("findnull")
					continue
				}
				if v.Label["selector"] == findsame {
					//fmt.Println(v.Label, v.Hostnode)
					hostnode = append(hostnode, nametonum[v.Hostnode])
				}
			}
			//	fmt.Println(hostnode)
			//calculate the near score
			for i := 0; i < len(hostnode); i++ {
				samedata[hostnode[i]-1] = 1
			}
		}
	}
	printSlice(samedata)
	return samedata, finalfind
}
func calculatetotal(a, b, c, d, e float64, number int, allnode *Nodes, allcontainer *Containers, container *Container, faraway, near, same map[string]string) (totalscore float64, choice int, err error, bl bool) {
	distancescore := calculatedistancescore(allnode, number)                              //datadistance
	resourcescore := calculateresourcescore(allcontainer, number, container)              //resource
	diffscore := calculatediffscore(allcontainer, number, faraway, container)             //different
	nearscore := calculatenearscore(allnode, allcontainer, number, near, container)       //near
	samescore, findsamelabel := calculatesamesocre(allcontainer, number, same, container) //same
	if (e != 0 || b != 0) && findsamelabel == true {
		a = 0
		c = 0
		d = 0
	}
	//var total [num]float64
	total := make([]float64, 0, number)
	for i := 0; i < number; i++ {
		total = append(total, a*distancescore[i]+b*diffscore[i]+c*resourcescore[i]+d*nearscore[i]+e*samescore[i])
	}
	fmt.Println(total)
	score, index := findmaxnode(total)
	return score, index, nil, findsamelabel
}

func (d *Daemon) selectnode(job *engine.Job) engine.Status {
	var myfaraway []byte
	var mynear []byte
	var mysame []byte

	if contents3, err3 := ioutil.ReadFile("different.txt"); err3 == nil {
		//因为contents是[]byte类型，直接转换成string类型后会多一行空格,需要使用strings.Replace替换换行符
		myfaraway = contents3
		fmt.Println("readfaraway")
	} else {
		return job.Errorf("error: %v", err3)
	}

	if contents4, err4 := ioutil.ReadFile("near.txt"); err4 == nil {
		//因为contents是[]byte类型，直接转换成string类型后会多一行空格,需要使用strings.Replace替换换行符
		mynear = contents4
		fmt.Println("readnear")
	} else {
		return job.Errorf("error: %v", err4)
	}

	if contents5, err5 := ioutil.ReadFile("same.txt"); err5 == nil {
		//因为contents是[]byte类型，直接转换成string类型后会多一行空格,需要使用strings.Replace替换换行符
		mysame = contents5
		fmt.Println("readsame")
	} else {
		return job.Errorf("error: %v", err5)
	}
	rd := &ReceiveData{}
	receivedata := []byte(job.GetEnv("data"))
	if err := json.Unmarshal(receivedata, rd); err != nil {
		fmt.Println("nodes error")
		fmt.Println(err)
		return job.Errorf("receive data error: %v", err)
	}
	containers := &Containers{
		AllContainer: rd.ReceiveContainers.AllContainer,
	}
	nodes := &Nodes{
		AllNodes: rd.ReceiveNode.AllNodes,
	}
	mycontainer := &Container{
		Cpu:      rd.ReceiveContainer.Cpu,
		Mem:      rd.ReceiveContainer.Mem,
		Image:    rd.ReceiveContainer.Image,
		Label:    rd.ReceiveContainer.Label,
		Hostnode: rd.ReceiveContainer.Hostnode,
	}
	aa := rd.DistanceWeight
	bb := rd.DiffWeight
	cc := rd.ResourceWeight
	dd := rd.NearWeight
	ee := rd.SameWeight

	mapfaraway := make(map[string]string)
	mapnear := make(map[string]string)
	mapsame := make(map[string]string)
	nametonum = make(map[string]int)
	for _, v := range nodes.AllNodes {
		nametonum[v.Name] = v.Index
	}
	//fmt.Println(nametonum)
	if err1 := json.Unmarshal(myfaraway, &mapfaraway); err1 != nil {
		fmt.Println("faraway error")
		return job.Errorf("error: %v", err1)
	}
	if err2 := json.Unmarshal(mynear, &mapnear); err2 != nil {
		fmt.Println("near error")
		return job.Errorf("error:%v", err2)
	}
	if err6 := json.Unmarshal(mysame, &mapsame); err6 != nil {
		fmt.Println("same error")
		return job.Errorf("error:%v", err6)
	}
	number := len(nodes.AllNodes)
	//fmt.Println(number)
	sort.Sort(nodes)
	// for _, v := range nodes.AllNodes {
	// 	fmt.Println(v.Index)
	// }
	score, index, err, findsamelabel := calculatetotal(aa, bb, cc, dd, ee, number, nodes, containers, mycontainer, mapfaraway, mapnear, mapsame)
	if err != nil {
		return job.Errorf("calculateerror:%v", err)
	}
	fmt.Println("score:", score, "choose:", index)
	var nodename string
	var datadistance float64
	var cpucap float64
	var memcap float64
	var distance []float64
	for k, v := range nodes.AllNodes {
		if k == index-1 {
			nodename = v.Name
			datadistance = v.Datadistance
			cpucap = v.Cpucap
			memcap = v.Memcap
			distance = v.Distance
			fmt.Println("find it ")
			//fmt.Println(v.Name)
		}
	}
	var returnnode *ReturnNode
	if ee != 0 && bb != 0 && findsamelabel == true && score <= ee {
		returnnode = &ReturnNode{
			NodeIndex: -1,
		}
	} else if ee != 0 && bb == 0 && findsamelabel == true && score == 0 {
		returnnode = &ReturnNode{
			NodeIndex: -1,
		}
	} else {
		returnnode = &ReturnNode{
			NodeIndex:    index,
			Nodename:     nodename,
			Datadistance: datadistance,
			Cpucap:       cpucap,
			Memcap:       memcap,
			Distance:     distance,
		}
	}

	str, err8 := json.Marshal(returnnode)
	if err8 != nil {
		return job.Errorf("%v", err8)
	}

	//fmt.Println(str)
	if _, err9 := job.Stdout.Write(str); err9 != nil {
		return job.Errorf("%v", err9)
	}
	return engine.StatusOK
}
