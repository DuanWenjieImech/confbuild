
package example

import(
	"sync"
	"errors"
	"encoding/json"
	"sync/atomic"
)

type ActivityDrop struct { 
	Comment	string   // optional 服务端本地化  
	TempID	uint32   // required 唯一ID  
	PeachRewardData	[]struct   { 
			 Type int32 	// optional 奖励类型    
			 Num int64 	// optional 数量    
			 Param1 int32 	// optional 参数1    
			 Param2 int32 	// optional 参数2    
			 Probability int32 	// optional 概率(千分比）    
	}  
	PVERewardData	[]struct   { 
			 Type int32 	// optional 奖励类型    
			 Num int64 	// optional 数量    
			 Param1 int32 	// optional 参数1    
			 Param2 int32 	// optional 参数2    
			 Probability int32 	// optional 概率(千分比）    
	}  
}


var(
	iActivityDropList = map[uint32]*ActivityDrop{}
	iActivityDropMutex 	sync.RWMutex
	iActivityDropSize  uint32
	iActivityDropHook	func(list map[uint32]*ActivityDrop)
)

//从文件读取数据到内存
func ActivityDrop_ListUpdate(){
	data, err := confRedis.HGet(GameConfDataKey, "ActivityDrop")
	if err != nil {
		panic(err)
	}

	list := []ActivityDrop{}

	err = json.Unmarshal(data, &list)
	if err != nil {
		panic(err)
	}

	
	iActivityDropMutex.Lock()
	for k, item := range list {
		iActivityDropList[item.TempID] = &list[k]
	}
	iActivityDropMutex.Unlock()

	if iActivityDropHook != nil {
		iActivityDropMutex.RLock()
		iActivityDropHook(iActivityDropList)
		iActivityDropMutex.RUnlock()
	}

	atomic.StoreUint32(&iActivityDropSize, uint32(len(iActivityDropList)))
}

//唯一主键查找
func ActivityDrop_FindByPk(ID uint32) (activityDrop *ActivityDrop, err error){
	iActivityDropMutex.RLock()
	defer iActivityDropMutex.RUnlock()

	var ok bool
	activityDrop, ok = iActivityDropList[ID]
	if ok == false {
		err = errors.New("Not Data Found")
		return
	}


	return
}

//map的数据量大小
func ActivityDrop_ListLen() uint32 {
	return atomic.LoadUint32(&iActivityDropSize)
}

//获取完整数据
func ActivityDrop_ListAll() map[uint32]*ActivityDrop{
	iActivityDropMutex.RLock()
	defer iActivityDropMutex.RUnlock()

	m := map[uint32]*ActivityDrop{}

	for k, _ := range iActivityDropList {
		m[k] = iActivityDropList[k]
	}

	return m
}

//自定义处理, 返回false, 终止遍历
func ActivityDrop_ListRange(f func(k uint32, v *ActivityDrop) bool) {
	iActivityDropMutex.RLock()
	defer iActivityDropMutex.RUnlock()


	for k, _ := range iActivityDropList {
		flag := f(k, iActivityDropList[k])
		if flag == false {
			return
		}
	}
}

//以下为兼容处理
func ActivityDropList() map[uint32]*ActivityDrop{
	return ActivityDrop_ListAll()
}

func FindByPkActivityDrop(ID uint32) (activityDrop *ActivityDrop, err error){
	return ActivityDrop_FindByPk(ID)
}

func ActivityDropLen() uint32 {
	return ActivityDrop_ListLen()
}
