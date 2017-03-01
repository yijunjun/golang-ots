package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	ots2 "github.com/GiterLab/goots"
	. "github.com/GiterLab/goots/otstype"
	"github.com/tealeg/xlsx"
)

// modify it to yours
var (
	ACCESSID  = "LTAIqjq4OplpzZRS"
	ACCESSKEY = "GbtpMcQHxTBElHLVwC4UDl1lSsmFK4"
)

type TFile struct {
	Name       string
	Sheets     []string
	RowHeaders []string
	Rows       [][]string
}

func Save(tf *TFile) error {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var err error

	if _, err := os.Stat(tf.Name); err == nil {
		file, err = xlsx.OpenFile(tf.Name)
		if err != nil {
			return err
		}
	} else {
		file = xlsx.NewFile()
	}
	for _, item := range tf.Sheets {
		sheet, err = file.AddSheet(item)
		if err != nil {
			return err
		}

		sheet.SheetFormat.DefaultColWidth = 150

		sheet.AddRow().WriteSlice(&tf.RowHeaders, len(tf.RowHeaders))

		for _, tr := range tf.Rows {
			sheet.AddRow().WriteSlice(&tr, len(tr))
		}
	}

	return file.Save(tf.Name)
}

type TInsert struct {
	//"成功次数", "失败次数", "成功比例", "总成本", "平均成本", "总时间", "平均时间", "总请求时间", "平均请求时间"
	Ok         int32
	Fail       int32
	OkRate     float32
	CostAll    int32
	CostAvg    int32
	TimeAll    float64
	TimeAvg    float64
	ReqTimeAll float64
	ReqTimeAvg float64
}

var ID_LIMIT = 10000

// 网址,实列,表名,执行次数,列数
func do_insert(ots *ots2.OTSClient, table string, num, cols int) (*TInsert, error) {
	isSingle := table == "single"

	var ti TInsert

	var attr OTSAttribute

	errList := make([]error, num)

	// insert a row
	for j := 0; j < num; j++ {
		begin := time.Now()
		if isSingle {
			abs := make([]string, cols)
			for i := 0; i < cols; i++ {
				abs[i] = fmt.Sprintf("%020d", ID_LIMIT+i)
			}
			attr = OTSAttribute{
				"attr": strings.Join(abs, ","),
			}
		} else {
			abs := make(map[string]interface{}, cols)
			for i := 0; i < cols; i++ {
				abs[fmt.Sprintf("%020d", ID_LIMIT+i)] = true
			}
			attr = OTSAttribute(abs)
		}
		mid := time.Now()
		res, ots_err := ots.PutRow(
			table,
			OTSCondition_EXPECT_NOT_EXIST,
			&OTSPrimaryKey{
				"id": fmt.Sprintf("%020d", ID_LIMIT+j),
			},
			&attr,
		)
		if ots_err != nil {
			errList[ti.Fail] = ots_err
			ti.Fail += 1
		} else {
			ti.Ok += 1
			ti.CostAll += res.GetWriteConsumed()
			now := time.Now()
			ti.TimeAll += now.Sub(begin).Seconds()
			ti.ReqTimeAll += now.Sub(mid).Seconds()
		}
	}

	for i := int32(0); i < ti.Fail; i++ {
		fmt.Println(i, errList[i].Error())
	}

	if ti.Ok > 0 {
		ti.OkRate = float32(ti.Ok) / float32(ti.Ok+ti.Fail)

		ti.CostAvg = ti.CostAll / ti.Ok

		ti.TimeAvg = ti.TimeAll / float64(ti.Ok)

		ti.ReqTimeAvg = ti.ReqTimeAll / float64(ti.Ok)
	}
	return &ti, nil
}

func Insert(num, cols int) {
	tf := &TFile{
		Name:       "inesrt.xlsx",
		Sheets:     []string{fmt.Sprintf("%v列-%v次", cols, num)},
		RowHeaders: []string{"节点表", "成功次数", "失败次数", "成功比例", "总成本", "平均成本", "总时间", "平均时间", "总请求时间", "平均请求时间"},
		Rows:       make([][]string, len(configs)),
	}
	var ots_client *ots2.OTSClient
	var err error
	var ti *TInsert
	for i, item := range configs {
		// 删除表并且重建表,保证数据不冲突
		ots_client, err = ots2.New(item[0], ACCESSID, ACCESSKEY, item[1])
		if err != nil {
			fmt.Println("New:", item[0], item[1], err)
			return
		}

		if errDT := ots_client.DeleteTable(item[2]); errDT != nil {
			fmt.Println("DeleteTable:", item[0], item[1], errDT)
		}

		errCT := ots_client.CreateTable(
			&OTSTableMeta{
				TableName: item[2],
				SchemaOfPrimaryKey: OTSSchemaOfPrimaryKey{
					{K: "id", V: "STRING"},
				},
			},
			&OTSReservedThroughput{
				OTSCapacityUnit{0, 0},
			},
		)
		if errCT != nil {
			fmt.Println("CreateTable:", item[0], item[1], errCT)
			return
		}

		ti, err = do_insert(ots_client, item[2], num, cols)
		if err != nil {
			fmt.Println("do_insert:", err)
			continue
		}
		tf.Rows[i] = []string{
			item[1] + "-" + item[2],
			fmt.Sprint(ti.Ok),
			fmt.Sprint(ti.Fail),
			fmt.Sprintf("%2.2f%%", ti.OkRate*100),
			fmt.Sprint(ti.CostAll),
			fmt.Sprint(ti.CostAvg),
			fmt.Sprintf("%3.3f", ti.TimeAll),
			fmt.Sprintf("%3.3f", ti.TimeAvg),
			fmt.Sprintf("%3.3f", ti.ReqTimeAll),
			fmt.Sprintf("%3.3f", ti.ReqTimeAvg),
		}
	}

	if err := Save(tf); err != nil {
		fmt.Println("save", err)
	}
}

var configs = [][]string{
	[]string{
		"http://rongliang-test.cn-shenzhen.ots.aliyuncs.com",
		"rongliang-test",
		"single",
	},
	[]string{
		"http://rongliang-test.cn-shenzhen.ots.aliyuncs.com",
		"rongliang-test",
		"multi",
	},
	[]string{
		"http://xingneng-test.cn-shenzhen.ots.aliyuncs.com",
		"xingneng-test",
		"single",
	},
	[]string{
		"http://xingneng-test.cn-shenzhen.ots.aliyuncs.com",
		"xingneng-test",
		"multi",
	},
}

func main() {
	insideEnv := flag.String("inside", "", "使用内网地址")
	insertEnv := flag.String("insert", "", "测试插入")
	num := flag.Int("num", 10, "次数")
	cols := flag.Int("cols", 10, "列数")

	flag.Parse()

	if *insideEnv != "" {
		configs[0][0] = "http://rongliang-test.cn-shenzhen.ots-internal.aliyuncs.com"
		configs[1][0] = "http://rongliang-test.cn-shenzhen.ots-internal.aliyuncs.com"

		configs[2][0] = "http://xingneng-test.cn-shenzhen.ots-internal.aliyuncs.com"
		configs[3][0] = "http://xingneng-test.cn-shenzhen.ots-internal.aliyuncs.com"
	}

	// set running environment
	// ots2.OTSDebugEnable = true
	// ots2.OTSLoggerEnable = true
	ots2.OTSErrorPanicMode = true // 默认为开启，如果不喜欢panic则设置此为false

	if *insertEnv != "" {
		Insert(*num, *cols)
	}
}
