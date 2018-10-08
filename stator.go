package main

type stator struct {
	cfg         *StatConfig
	statHandler *StatHandle
}

const (
	StatInvalidMethodReq = "Invalid_Http_Method"
	StatReadBody         = "Read_Body_Error"
	StatRedisGet         = "Redis_Get_Count"
	StatRedisResp        = "Redis_Resp"
)

func newStator(cfg *StatConfig) *stator {
	return &stator{cfg: cfg}
}

func (s *stator) start() {
	conf := StatConf{
		Interval:     int64(s.cfg.CycleSec),
		StatFile:     s.cfg.StatFile,
		RemoteMode:   s.cfg.RemoteMode == 1,
		MaxChanLen:   int(s.cfg.QueueCapacity),
		RoutineCount: int(s.cfg.RoutineCount)}

	isSync := conf.RoutineCount == 0
	s.statHandler = NewStatHandle(conf, isSync)
	s.statHandler.Run()
}

func (s *stator) getCmdStatTitle(cmdId uint16) string {
	return "UNKNOW"
}
