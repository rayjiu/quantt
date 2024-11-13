package wogoo

import (
	"encoding/binary"
	"io"
	"net"
	"strconv"
	"time"

	"github.com/dablelv/cyan/encoding"
	"github.com/rayjiu/quantt/watchdog/internal/chinaA/wogoo/constants"
	"github.com/rayjiu/quantt/watchdog/internal/chinaA/wogoo/model"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

var WogooHQServer wogoo = wogoo{
	snapshotMap: make(map[string]func(snp *model.StockSnapshot)),
}

// wogoo 实时行情订阅
type wogoo struct {
	commonRq    *model.CommonReq
	conn        net.Conn
	remainBuf   []byte
	snapshotMap map[string]func(snp *model.StockSnapshot)
}

func (w *wogoo) SetupClient(host, userId, deviceNo, token string) {
	log.Infof("Host:%v", host)
	conn, err := net.Dial("tcp", host)

	conn.SetDeadline(time.Time{})

	if err != nil {
		panic(err)
	}
	defer func(conn net.Conn) {
		//conn.Close()
		//buffer := make([]byte, 1000000)
		//_, err1 := conn.Read(buffer)
		//if err1 != nil {
		//
		//}
	}(conn)

	w.conn = conn

	w.createCommonRq(userId, deviceNo, token)

	go w.receiver()

}

func (w *wogoo) createCommonRq(userId string, deviceNo string, token string) {
	w.commonRq = &model.CommonReq{
		UserId:   userId,
		Token:    token,
		TermType: 1,
		DeviceNo: deviceNo,
	}
}

func (w *wogoo) receiver() {
	timeSend := time.Now()
	for {
		resByte, err := w.readMsg()
		// log.Infof("Receive data.%+v, elapse:%v", resByte, time.Since(timeSend))
		if err != nil {
			break
		}
		for _, v := range resByte {
			v := v
			go w.decodeMsg(v)
		}

		if time.Since(timeSend) >= 20*time.Second {
			//c.sendUnSub()
			//if time.Now().Second() %10 ==0 {
			w.sendKeepLiving()
			timeSend = time.Now()
		}
	}
}

func (w *wogoo) readMsg() ([][]byte, error) {
	buffer := make([]byte, 1000000)
	receiveLen, err := w.conn.Read(buffer)
	w.conn.SetReadDeadline(time.Now().Add(time.Millisecond * 100))

	if err == io.EOF {
		log.Infof("server closed: %v", err.Error())
		return nil, err
	}
	remain := append(w.remainBuf, buffer[:receiveLen]...)
	res := make([][]byte, 0)
	for len(remain) >= 11 && binary.BigEndian.Uint32(remain[2:6]) <= uint32(len(remain)) {
		res = append(res, remain[:binary.BigEndian.Uint32(remain[2:6])])
		remain = remain[binary.BigEndian.Uint32(remain[2:6]):]
	}
	w.remainBuf = remain
	return res, nil
}

func (w *wogoo) decodeMsg(buff []byte) {
	length := len(buff)
	if length < 14 {
		return
	}
	commandNo := int(binary.BigEndian.Uint16(buff[6:8]))
	switch commandNo {
	case 1002:
		pb := &model.StockSnapshot{}
		proto.Unmarshal(buff[14:length], pb)
		// s, _ := encoding.ToIndentJSON(&pb)
		if len(buff) > 0 {
			w.snapshotMap[pb.StockId+"."+strconv.Itoa(int(pb.MarketType))](pb)
		}
	case 2011:
		pb := &model.KeepLivingRs{}
		proto.Unmarshal(buff[14:length], pb)
		s, _ := encoding.ToIndentJSON(&pb)
		if len(buff) > 0 {
			log.Infof("心跳包: %v", s)
		}
	}
}

func (w *wogoo) sendKeepLiving() {
	kpReq := &model.MarketOverviewReq{
		CommonReq:  w.commonRq,
		MarketType: 1,
	}
	body, _ := proto.Marshal(kpReq)
	d := []byte{'W', 'G', 0, 0, 0, 18, 7, 219, 0, 1, 1, 4, 5, 1}
	for _, v := range body {
		d = append(d, v)
	}
	d[5] = byte(len(d))
	_, err := w.conn.Write(d)
	if err != nil {
		log.Infof("发送错误%v", err)
	}
}

func (w *wogoo) sendSub(stocks []*model.SubStockInfo) error {
	req := &model.SubReq{
		CommonReq: w.commonRq,
		SubStocks: stocks,
		SubMarkets: []*model.SubMarketInfo{

			{
				MarketType: 1,
			},
			{
				MarketType: 2,
			},
		},
	}

	body, _ := proto.Marshal(req)

	d := []byte{'W', 'G', 0, 0, 0, 18, 0, 101, 0, 2, 1, 5, 15, 1}

	d = append(d, body...)
	Len := []byte{1, 2, 3, 4}
	binary.BigEndian.PutUint32(Len, uint32(len(d)))
	d[2] = Len[0]
	d[3] = Len[1]
	d[4] = Len[2]
	d[5] = Len[3]
	_, err := w.conn.Write(d)
	if err != nil {
		log.Infof("发送错误%v", err)
		return err
	}
	return nil
}

func (w *wogoo) DoSub(stockId string, marketType uint32, bizType uint32, respFunc func(*model.StockSnapshot)) {
	var err = w.sendSub([]*model.SubStockInfo{
		{
			StockId:    stockId,
			MarketType: marketType,
			BizType:    []uint32{bizType},
		},
	})
	if err != nil {
		log.Errorf("发送订阅出错:%+v", err)
	} else {
		if bizType == constants.PushBizTypeSnapshot {
			w.snapshotMap[stockId+"."+strconv.Itoa(int(marketType))] = respFunc
		}
		log.Infof("订阅成功，并且将对应的应答函数添加到snapshotMap, len(snapshotMap)=%v", len(w.snapshotMap))
	}
}
