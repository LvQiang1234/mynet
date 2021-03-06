package base

import (
	"fmt"
	"io/ioutil"
	"mynet/base/vector"
	"os"
)

//datatype
const (
	DType_none        = iota
	DType_String      = iota
	DType_Enum        = iota
	DType_S8          = iota
	DType_S16         = iota
	DType_S32         = iota
	DType_F32         = iota
	DType_F64         = iota
	DType_S64         = iota
	DType_StringArray = iota
	DType_S8Array     = iota
	DType_S16Array    = iota
	DType_S32Array    = iota
	DType_F32Array    = iota
	DType_F64Array    = iota
	DType_S64Array    = iota
)

type (
	RData struct {
		m_FileName string
		m_Type     int

		m_String      string
		m_Enum        int
		m_S8          int8
		m_S16         int16
		m_S32         int
		m_F32         float32
		m_F64         float64
		m_S64         int64
		m_StringArray []string
		m_S8Array     []int8
		m_S16Array    []int16
		m_S32Array    []int
		m_F32Array    []float32
		m_F64Array    []float64
		m_S64Array    []int64
	}

	CDataFile struct {
		RecordNum int //记录数量
		ColumNum  int //列数量

		fileName           string //文件名
		fstream            *BitStream
		readstep           int //控制读的总数量
		dataTypes          vector.Vector
		currentColumnIndex int
	}

	IDateFile interface {
		ReadDataFile(string) bool
		GetData(*RData) bool
		ReadDataInit()
	}
)

func (this *CDataFile) ReadDataInit() {
	this.ColumNum = 0
	this.RecordNum = 0
	this.readstep = 0
	this.fstream = nil
}

func (this *CDataFile) ReadDataFile(fileName string) bool {
	this.dataTypes.Clear()
	this.currentColumnIndex = 0
	this.fileName = fileName

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("[%s] open failed", fileName)
		return false
	}
	defer file.Close()
	buf, err := ioutil.ReadAll(file)
	if err != nil {
		return false
	}
	this.fstream = NewBitStream(buf, len(buf))

	for {
		tchr := this.fstream.ReadInt(8)
		if tchr == '@' { //找到数据文件的开头
			tchr = this.fstream.ReadInt(8) //这个是换行字符
			//fmt.Println(tchr)
			break
		}
	}
	//得到记录总数
	this.RecordNum = this.fstream.ReadInt(32)
	//得到列的总数
	this.ColumNum = this.fstream.ReadInt(32)
	//sheet name
	this.fstream.ReadString()

	this.readstep = this.RecordNum * this.ColumNum
	for nColumnIndex := 0; nColumnIndex < this.ColumNum; nColumnIndex++ {
		//col name
		this.fstream.ReadString()
		nDataType := this.fstream.ReadInt(8)
		this.dataTypes.PushBack(int(nDataType))
	}
	return true
}

/****************************
	格式:
	头文件:
	1、总记录数(int)
	2、总字段数(int)
	字段格式:
	1、字段长度(int)
	2、字读数据类型(int->2,string->1,enum->3,float->4)
	3、字段内容(int,string)
*************************/
func (this *CDataFile) GetData(pData *RData) bool {
	if this.readstep == 0 || this.fstream == nil {
		return false
	}

	pData.m_FileName = this.fileName
	switch this.dataTypes.Get(this.currentColumnIndex).(int) {
	case DType_String:
		pData.m_String = this.fstream.ReadString()
	case DType_S8:
		pData.m_S8 = int8(this.fstream.ReadInt(8))
	case DType_S16:
		pData.m_S16 = int16(this.fstream.ReadInt(16))
	case DType_S32:
		pData.m_S32 = this.fstream.ReadInt(32)
	case DType_Enum:
		pData.m_Enum = this.fstream.ReadInt(16)
	case DType_F32:
		pData.m_F32 = this.fstream.ReadFloat()
	case DType_F64:
		pData.m_F64 = this.fstream.ReadFloat64()
	case DType_S64:
		pData.m_S64 = this.fstream.ReadInt64(64)

	case DType_StringArray:
		nLen := this.fstream.ReadInt(8)
		pData.m_StringArray = make([]string, nLen)
		for i := 0; i < nLen; i++ {
			pData.m_StringArray[i] = this.fstream.ReadString()
		}
	case DType_S8Array:
		nLen := this.fstream.ReadInt(8)
		pData.m_S8Array = make([]int8, nLen)
		for i := 0; i < nLen; i++ {
			pData.m_S8Array[i] = int8(this.fstream.ReadInt(8))
		}
	case DType_S16Array:
		nLen := this.fstream.ReadInt(8)
		pData.m_S16Array = make([]int16, nLen)
		for i := 0; i < nLen; i++ {
			pData.m_S16Array[i] = int16(this.fstream.ReadInt(16))
		}
	case DType_S32Array:
		nLen := this.fstream.ReadInt(8)
		pData.m_S32Array = make([]int, nLen)
		for i := 0; i < nLen; i++ {
			pData.m_S32Array[i] = this.fstream.ReadInt(32)
		}
	case DType_F32Array:
		nLen := this.fstream.ReadInt(8)
		pData.m_F32Array = make([]float32, nLen)
		for i := 0; i < nLen; i++ {
			pData.m_F32Array[i] = this.fstream.ReadFloat()
		}
	case DType_F64Array:
		nLen := this.fstream.ReadInt(8)
		pData.m_F64Array = make([]float64, nLen)
		for i := 0; i < nLen; i++ {
			pData.m_F64Array[i] = this.fstream.ReadFloat64()
		}
	case DType_S64Array:
		nLen := this.fstream.ReadInt(8)
		pData.m_S64Array = make([]int64, nLen)
		for i := 0; i < nLen; i++ {
			pData.m_S64Array[i] = this.fstream.ReadInt64(64)
		}
	}

	pData.m_Type = this.dataTypes.Get(this.currentColumnIndex).(int)
	this.currentColumnIndex = (this.currentColumnIndex + 1) % this.ColumNum
	this.readstep--
	return true
}

/****************************
	RData funciton
****************************/
func (this *RData) String(datacol string) string {
	IFAssert(this.m_Type == DType_String, fmt.Sprintf("read [%s] col[%s] error", this.m_FileName, datacol))
	return this.m_String
}

func (this *RData) Enum(datacol string) int {
	IFAssert(this.m_Type == DType_Enum, fmt.Sprintf("read [%s] col[%s] error", this.m_FileName, datacol))
	return this.m_Enum
}

func (this *RData) Int8(datacol string) int8 {
	IFAssert(this.m_Type == DType_S8, fmt.Sprintf("read [%s] col[%s] error", this.m_FileName, datacol))
	return this.m_S8
}

func (this *RData) Int16(datacol string) int16 {
	IFAssert(this.m_Type == DType_S16, fmt.Sprintf("read [%s] col[%s] error", this.m_FileName, datacol))
	return this.m_S16
}

func (this *RData) Int(datacol string) int {
	IFAssert(this.m_Type == DType_S32, fmt.Sprintf("read [%s] col[%s] error", this.m_FileName, datacol))
	return this.m_S32
}

func (this *RData) Float32(datacol string) float32 {
	IFAssert(this.m_Type == DType_F32, fmt.Sprintf("read [%s] col[%s] error", this.m_FileName, datacol))
	return this.m_F32
}

func (this *RData) Float64(datacol string) float64 {
	IFAssert(this.m_Type == DType_F64, fmt.Sprintf("read [%s] col[%s] error", this.m_FileName, datacol))
	return this.m_F64
}

func (this *RData) Int64(datacol string) int64 {
	IFAssert(this.m_Type == DType_S64, fmt.Sprintf("read [%s] col[%s] error", this.m_FileName, datacol))
	return this.m_S64
}

func (this *RData) StringArray(datacol string) []string {
	IFAssert(this.m_Type == DType_StringArray, fmt.Sprintf("read [%s] col[%s] error", this.m_FileName, datacol))
	return this.m_StringArray
}

func (this *RData) Int8Array(datacol string) []int8 {
	IFAssert(this.m_Type == DType_S8Array, fmt.Sprintf("read [%s] col[%s] error", this.m_FileName, datacol))
	return this.m_S8Array
}

func (this *RData) Int16Array(datacol string) []int16 {
	IFAssert(this.m_Type == DType_S16Array, fmt.Sprintf("read [%s] col[%s] error", this.m_FileName, datacol))
	return this.m_S16Array
}

func (this *RData) IntArray(datacol string) []int {
	IFAssert(this.m_Type == DType_S32Array, fmt.Sprintf("read [%s] col[%s] error", this.m_FileName, datacol))
	return this.m_S32Array
}

func (this *RData) Float32Array(datacol string) []float32 {
	IFAssert(this.m_Type == DType_F32Array, fmt.Sprintf("read [%s] col[%s] error", this.m_FileName, datacol))
	return this.m_F32Array
}

func (this *RData) Float64Array(datacol string) []float64 {
	IFAssert(this.m_Type == DType_F64Array, fmt.Sprintf("read [%s] col[%s] error", this.m_FileName, datacol))
	return this.m_F64Array
}

func (this *RData) Int64Array(datacol string) []int64 {
	IFAssert(this.m_Type == DType_S64Array, fmt.Sprintf("read [%s] col[%s] error", this.m_FileName, datacol))
	return this.m_S64Array
}
