
package example

import(
	"sync"
	"errors"
	"encoding/json"
	"sync/atomic"
)

type ActivityAward struct { 
	Comment	string   // optional 服务端本地化  
	TempID	uint32   // required 模板ID  
	ActivtyId	uint32   // optional 活动ID  
	Condition	int32   // optional 活动条件  
	Num	int32   // optional 条件  
	UnLock	int32   // optional 永久累充解锁条件  
	RewardData	[]struct   { 
			 Type int32 	// optional 奖励类型    
			 Num int64 	// optional 数量    
			 Param1 int32 	// optional 参数1    
			 Param2 int32 	// optional 参数2    
	}  
}


var(
	iActivityAwardList = map[uint32]*ActivityAward{}
	iActivityAwardMutex 	sync.RWMutex
	iActivityAwardSize  uint32
	iActivityAwardHook	func(list map[uint32]*ActivityAward)
)

//从文件读取数据到内存
func ActivityAward_ListUpdate(){
	data, err := confRedis.HGet(GameConfDataKey, "ActivityAward")
	if err != nil {
		panic(err)
	}

	list := []ActivityAward{}

	err = json.Unmarshal(data, &list)
	if err != nil {
		panic(err)
	}

	
	iActivityAwardMutex.Lock()
	for k, item := range list {
		iActivityAwardList[item.TempID] = &list[k]
	}
	iActivityAwardMutex.Unlock()

	if iActivityAwardHook != nil {
		iActivityAwardMutex.RLock()
		iActivityAwardHook(iActivityAwardList)
		iActivityAwardMutex.RUnlock()
	}

	atomic.StoreUint32(&iActivityAwardSize, uint32(len(iActivityAwardList)))
}

//唯一主键查找
func ActivityAward_FindByPk(ID uint32) (activityAward *ActivityAward, err error){
	iActivityAwardMutex.RLock()
	defer iActivityAwardMutex.RUnlock()

	var ok bool
	activityAward, ok = iActivityAwardList[ID]
	if ok == false {
		err = errors.New("Not Data Found")
		return
	}


	return
}

//map的数据量大小
func ActivityAward_ListLen() uint32 {
	return atomic.LoadUint32(&iActivityAwardSize)
}

//获取完整数据
func ActivityAward_ListAll() map[uint32]*ActivityAward{
	iActivityAwardMutex.RLock()
	defer iActivityAwardMutex.RUnlock()

	m := map[uint32]*ActivityAward{}

	for k, _ := range iActivityAwardList {
		m[k] = iActivityAwardList[k]
	}

	return m
}

//自定义处理, 返回false, 终止遍历
func ActivityAward_ListRange(f func(k uint32, v *ActivityAward) bool) {
	iActivityAwardMutex.RLock()
	defer iActivityAwardMutex.RUnlock()


	for k, _ := range iActivityAwardList {
		flag := f(k, iActivityAwardList[k])
		if flag == false {
			return
		}
	}
}

//以下为兼容处理
func ActivityAwardList() map[uint32]*ActivityAward{
	return ActivityAward_ListAll()
}

func FindByPkActivityAward(ID uint32) (activityAward *ActivityAward, err error){
	return ActivityAward_FindByPk(ID)
}

func ActivityAwardLen() uint32 {
	return ActivityAward_ListLen()
}
