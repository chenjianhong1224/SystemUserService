package main

import (
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

// Item type(attr)
const (
	ITEM_TYPE_UNKNOWN = iota
	ITEM_TYPE_COUNTER
	ITEM_TYPE_AMOUNT
	ITEM_TYPE_TIME
	ITEM_TYPE_AMOUNT_TIME
)

// stat option
const (
	OPT_UNKNOWN = iota
	OPT_FAIL_COUNT
	OPT_SUCC_COUNT
	OPT_TIME_COUNT
)

// Item data members
type Item struct {
	attr      int
	flag      bool
	count     int64
	total     int64
	succ      int64
	fail      int64
	totalTime int64
	avgTime   int64
	maxTime   int64
	minTime   int64
}

var index uint64
var itemMapIndex uint64

func (it *Item) Reset() {
	it.flag = false
	it.count = 0
	it.total = 0
	it.succ = 0
	it.fail = 0
	it.totalTime = 0
	it.avgTime = 0
	it.maxTime = 0
	it.minTime = 0
}

func (it *Item) merge(item *Item) {
	it.attr = item.attr
	it.flag = item.flag
	oldCount := it.count
	it.count += item.count
	it.total += item.total
	it.succ += item.succ
	it.fail += item.fail
	it.totalTime += item.totalTime
	if it.count > 0 {
		it.avgTime = (it.avgTime*oldCount + item.avgTime*item.count) / it.count
	}
	if it.maxTime < item.maxTime {
		it.maxTime = item.maxTime
	}

	if it.minTime > item.minTime {
		it.minTime = item.minTime
	}
}

func (it *Item) Dump(name string) string {
	switch it.attr {
	case ITEM_TYPE_COUNTER:
		return fmt.Sprintf("[%s:%d]", name, it.total)
	case ITEM_TYPE_AMOUNT:
		return fmt.Sprintf("[%s:%d,%d,%d]", name, it.total, it.succ, it.fail)
	case ITEM_TYPE_TIME:
		return fmt.Sprintf("[%s:%d,%d,%d,%d]", name, it.total, it.avgTime, it.maxTime, it.minTime)
	case ITEM_TYPE_AMOUNT_TIME:
		return fmt.Sprintf("[%s:%d,%d,%d,%d,%d,%d,%d]", name, it.total, it.succ, it.fail, it.totalTime, it.avgTime, it.maxTime, it.minTime)
	}
	return ""
}

func (it *Item) Add(opt int, value int64) {
	switch it.attr {
	case ITEM_TYPE_COUNTER:
		it.flag = true
		it.total += value
	case ITEM_TYPE_AMOUNT:
		it.flag = true
		switch opt {
		case OPT_FAIL_COUNT:
			it.total += value
			it.fail += value
		case OPT_SUCC_COUNT:
			it.total += value
			it.succ += value
		}
	case ITEM_TYPE_TIME:
		it.flag = true
		it.count++
		it.total += value
		it.avgTime = (int64)(it.total / it.count)
		if it.count == 1 {
			it.minTime = value
			it.maxTime = value
		} else {
			if it.maxTime < value {
				it.maxTime = value
			}
			if it.minTime > value {
				it.minTime = value
			}
		}
	case ITEM_TYPE_AMOUNT_TIME:
		it.flag = true
		switch opt {
		case OPT_FAIL_COUNT:
			it.fail += value
			it.total += value
		case OPT_SUCC_COUNT:
			it.succ += value
			it.total += value
		case OPT_TIME_COUNT:
			it.count++
			it.totalTime += value
			it.avgTime = (int64)(it.totalTime / it.count)
			if it.count == 1 {
				it.maxTime = value
				it.minTime = value
			} else {
				if it.maxTime < value {
					it.maxTime = value
				}
				if it.minTime > value {
					it.minTime = value
				}
			}
		}
	}
}

func timeTransform(t int64) string {
	tm := time.Unix(t, 0)
	return fmt.Sprintf("%s", tm.Format("20060102150405"))
}

func appendToFile(fileName string, content string) error {
	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "open stat.log file failed. err: "+err.Error())
	} else {
		n, _ := f.Seek(0, os.SEEK_END)
		_, err = f.WriteAt([]byte(content), n)
	}
	return err
}

const (
	RWMapCount      = 2
	DefaultMapCount = 16
)

type statItemMap map[string]*Item

type statItemEntry [RWMapCount]statItemMap

type ItemMutex [RWMapCount]sync.Mutex

type StatHandle struct {
	endpoint        string
	interval        int64
	statFile        string
	startTime       int64
	dropCount       uint64
	maxChanLen      int
	routineCount    int
	chanCount       int
	itemMapCount    int
	currentMapIndex uint64
	currentRWIndex  int32
	isSync          bool
	statChans       []chan *statItem
	statItemEntrys  []statItemEntry
	statItemMus     []ItemMutex
}

type amountTimeItem struct {
	name   string
	isSucc bool
	time   int64
}

type countItem struct {
	name string
}

type amountItem struct {
	name   string
	isSucc bool
}

type timeItem struct {
	name string
	time int64
}

type statItem interface{}

type StatInfo struct {
	ChanMaxLen     int
	ChanCount      int
	RoutineCount   int
	DropStatCount  uint64
	ChanLenDetails []int
}

func (stat *StatHandle) StatAmountTime(name string, isSucc bool, time int64) bool {
	if stat.isSync {
		return stat.syncStatAmountTime(name, isSucc, time)
	}

	item := amountTimeItem{name: name, isSucc: isSucc, time: time}
	statItem := statItem(&item)
	for i := 0; i < stat.chanCount; i++ {
		id := atomic.AddUint64(&index, 1)
		id = id % uint64(stat.chanCount)
		select {
		case stat.statChans[id] <- &statItem:
			return true
		default:
		}
	}

	atomic.AddUint64(&stat.dropCount, 1)
	return false
}

func (stat *StatHandle) syncStatAmountTime(name string, isSucc bool, time int64) bool {
	item := amountTimeItem{name: name, isSucc: isSucc, time: time}
	index := atomic.AddUint64(&stat.currentMapIndex, 1) % uint64(stat.itemMapCount)
	return stat.DoStatAmountTime(&item, &stat.statItemEntrys[index], &stat.statItemMus[index])
}

func (stat *StatHandle) DoStatAmountTime(amt *amountTimeItem, m *statItemEntry, mus *ItemMutex) bool {
	index := atomic.LoadInt32(&stat.currentRWIndex)
	if mus != nil {
		(*mus)[index].Lock()
		defer (*mus)[index].Unlock()
	}

	if item, found := (*m)[index][amt.name]; found {
		if amt.isSucc {
			item.Add(OPT_SUCC_COUNT, 1)
		} else {
			item.Add(OPT_FAIL_COUNT, 1)
		}
		item.Add(OPT_TIME_COUNT, amt.time)
	} else {
		it := new(Item)
		it.attr = ITEM_TYPE_AMOUNT_TIME
		if amt.isSucc {
			it.Add(OPT_SUCC_COUNT, 1)
		} else {
			it.Add(OPT_FAIL_COUNT, 1)
		}
		it.Add(OPT_TIME_COUNT, amt.time)
		(*m)[index][amt.name] = it
	}

	return true
}

func (stat *StatHandle) StatCount(name string) bool {
	if stat.isSync {
		return stat.syncStatCount(name)
	}

	item := countItem{name: name}
	statItem := statItem(&item)
	for i := 0; i < stat.chanCount; i++ {
		id := atomic.AddUint64(&index, 1)
		id = id % uint64(stat.chanCount)
		select {
		case stat.statChans[id] <- &statItem:
			return true
		default:
		}
	}

	atomic.AddUint64(&stat.dropCount, 1)
	return false
}

func (stat *StatHandle) syncStatCount(name string) bool {
	item := countItem{name: name}
	index := atomic.AddUint64(&stat.currentMapIndex, 1) % uint64(stat.itemMapCount)
	return stat.DoStatCount(&item, &stat.statItemEntrys[index], &stat.statItemMus[index])
}

func (stat *StatHandle) DoStatCount(count *countItem, m *statItemEntry, mus *ItemMutex) bool {
	index := atomic.LoadInt32(&stat.currentRWIndex)
	if mus != nil {
		(*mus)[index].Lock()
		defer (*mus)[index].Unlock()
	}

	if item, found := (*m)[index][count.name]; found {
		item.Add(OPT_SUCC_COUNT, 1)
	} else {
		it := new(Item)
		it.attr = ITEM_TYPE_COUNTER
		it.Add(OPT_SUCC_COUNT, 1)
		(*m)[index][count.name] = it
	}

	return true
}

func (stat *StatHandle) StatAmount(name string, isSucc bool) bool {
	if stat.isSync {
		return stat.syncStatAmount(name, isSucc)
	}

	item := amountItem{name: name, isSucc: isSucc}
	statItem := statItem(&item)
	for i := 0; i < stat.chanCount; i++ {
		id := atomic.AddUint64(&index, 1)
		id = id % uint64(stat.chanCount)
		select {
		case stat.statChans[id] <- &statItem:
			return true
		default:
		}
	}

	atomic.AddUint64(&stat.dropCount, 1)
	return false
}

func (stat *StatHandle) syncStatAmount(name string, isSucc bool) bool {
	item := amountItem{name: name, isSucc: isSucc}
	index := atomic.AddUint64(&stat.currentMapIndex, 1) % uint64(stat.itemMapCount)
	return stat.DoStatAmount(&item, &stat.statItemEntrys[index], &stat.statItemMus[index])
}

func (stat *StatHandle) DoStatAmount(amount *amountItem, m *statItemEntry, mus *ItemMutex) bool {
	index := atomic.LoadInt32(&stat.currentRWIndex)
	if mus != nil {
		(*mus)[index].Lock()
		defer (*mus)[index].Unlock()
	}

	if item, found := (*m)[index][amount.name]; found {
		if amount.isSucc {
			item.Add(OPT_SUCC_COUNT, 1)
		} else {
			item.Add(OPT_FAIL_COUNT, 1)
		}
	} else {
		it := new(Item)
		it.attr = ITEM_TYPE_AMOUNT
		if amount.isSucc {
			it.Add(OPT_SUCC_COUNT, 1)
		} else {
			it.Add(OPT_FAIL_COUNT, 1)
		}
		(*m)[index][amount.name] = it
	}

	return true
}

func (stat *StatHandle) StatTime(name string, time int64) bool {
	if stat.isSync {
		return stat.syncStatTime(name, time)
	}

	item := timeItem{name: name, time: time}
	statItem := statItem(&item)
	for i := 0; i < stat.chanCount; i++ {
		id := atomic.AddUint64(&index, 1)
		id = id % uint64(stat.chanCount)
		select {
		case stat.statChans[id] <- &statItem:
			return true
		default:
		}
	}

	atomic.AddUint64(&stat.dropCount, 1)
	return false
}

func (stat *StatHandle) syncStatTime(name string, time int64) bool {
	item := timeItem{name: name, time: time}
	index := atomic.AddUint64(&stat.currentMapIndex, 1) % uint64(stat.itemMapCount)
	return stat.DoStatTime(&item, &stat.statItemEntrys[index], &stat.statItemMus[index])
}

func (stat *StatHandle) DoStatTime(t *timeItem, m *statItemEntry, mus *ItemMutex) bool {
	index := atomic.LoadInt32(&stat.currentRWIndex)
	if mus != nil {
		(*mus)[index].Lock()
		defer (*mus)[index].Unlock()
	}

	if item, found := (*m)[index][t.name]; found {
		item.Add(OPT_SUCC_COUNT, t.time)
	} else {
		it := new(Item)
		it.attr = ITEM_TYPE_TIME
		it.Add(OPT_SUCC_COUNT, t.time)
		(*m)[index][t.name] = it
	}

	return true
}

func (stat *StatHandle) processStat(ch chan *statItem, m *statItemEntry, mus *ItemMutex) {
	for {
		item := <-ch
		switch item := (*item).(type) {
		case *amountTimeItem:
			stat.DoStatAmountTime(item, m, mus)
		case *countItem:
			stat.DoStatCount(item, m, mus)
		case *amountItem:
			stat.DoStatAmount(item, m, mus)
		case *timeItem:
			stat.DoStatTime(item, m, mus)
		default:
		}
	}
}

func (stat *StatHandle) mergeStatMap(index int32) *statItemMap {
	// not concurrency
	result := make(statItemMap)
	for i := 0; i < stat.itemMapCount; i++ {
		for k, v := range stat.statItemEntrys[i][index] {
			item, ok := result[k]
			if !ok {
				item = new(Item)
			}

			if v.flag {
				item.merge(v)
				result[k] = item
				v.Reset()
			}
		}
	}

	return &result
}

func (stat *StatHandle) pushStatMsg(statFile string) bool {
	now := time.Now().Unix()
	s := fmt.Sprintf("\n[interval:%s,%s]", timeTransform(stat.startTime), timeTransform(now))
	ifRealStat := false
	index := atomic.LoadInt32(&stat.currentRWIndex)
	atomic.StoreInt32(&stat.currentRWIndex, (index+1)%RWMapCount)
	result := stat.mergeStatMap(index)
	for k, v := range *result {
		if v.flag {
			s = fmt.Sprintf("%s%s", s, v.Dump(k))
			ifRealStat = true
		}
	}

	stat.startTime = now
	if ifRealStat {
		if appendToFile(statFile, s) != nil {
			fmt.Fprintf(os.Stderr, "appendToFile failed!")
			return false
		}
	}
	return true
}

func (stat *StatHandle) getStatFileSuffix() string {
	now := time.Now()
	return fmt.Sprintf("%d%02d%02d", now.Year(), now.Month(), now.Day())
}

func (stat *StatHandle) dumpStatPool() {
	now := time.Now().Unix()
	if (now - stat.startTime) >= stat.interval {
		stat.pushStatMsg(fmt.Sprintf("%s.%s", stat.statFile, stat.getStatFileSuffix()))
	}
}

func (stat *StatHandle) Run() {
	for i := 0; i < stat.routineCount; i++ {
		go func(ch chan *statItem, m *statItemEntry, mus *ItemMutex) {
			stat.processStat(ch, m, mus)
		}(stat.statChans[i], &stat.statItemEntrys[i], nil)
	}

	go func() {
		for {
			stat.dumpStatPool()
			time.Sleep(time.Second)
		}
	}()
}

type StatConf struct {
	Interval     int64
	StatFile     string
	RemoteMode   bool
	RoutineCount int
	MaxChanLen   int
	WaitMs       time.Duration
}

func NewStatHandle(conf StatConf, isSync bool) *StatHandle {
	statHandle := new(StatHandle)
	statHandle.startTime = time.Now().Unix()
	statHandle.interval = conf.Interval
	statHandle.statFile = conf.StatFile
	statHandle.isSync = isSync
	if !isSync {
		statHandle.routineCount = conf.RoutineCount
		if conf.RoutineCount <= 0 {
			panic("async stat goroutine should be at least 1")
		}

		statHandle.chanCount = conf.RoutineCount
		statHandle.statChans = make([]chan *statItem, statHandle.chanCount)
		statHandle.maxChanLen = conf.MaxChanLen
		for i := 0; i < statHandle.chanCount; i++ {
			statHandle.statChans[i] = make(chan *statItem, statHandle.maxChanLen)
		}

		statHandle.itemMapCount = conf.RoutineCount
	} else {
		statHandle.itemMapCount = DefaultMapCount
		statHandle.statItemMus = make([]ItemMutex, statHandle.itemMapCount)
	}

	statHandle.statItemEntrys = make([]statItemEntry, statHandle.itemMapCount)
	for i := 0; i < statHandle.itemMapCount; i++ {
		for j := 0; j < RWMapCount; j++ {
			statHandle.statItemEntrys[i][j] = make(statItemMap)
		}
	}

	return statHandle
}

func (stat *StatHandle) GetStatChansInfo() *StatInfo {
	statInfo := StatInfo{}
	statInfo.ChanMaxLen = stat.maxChanLen
	statInfo.ChanCount = stat.chanCount
	statInfo.RoutineCount = stat.routineCount
	statInfo.DropStatCount = atomic.LoadUint64(&stat.dropCount)
	statInfo.ChanLenDetails = make([]int, stat.chanCount)
	for i := 0; i < stat.chanCount; i++ {
		// not exact
		statInfo.ChanLenDetails[i] = len(stat.statChans[i])
	}

	return &statInfo
}
