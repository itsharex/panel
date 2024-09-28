package service

import (
	"fmt"
	"net/http"

	"github.com/golang-module/carbon/v2"

	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/internal/data"
	"github.com/TheTNB/panel/internal/http/request"
	"github.com/TheTNB/panel/pkg/types"
)

type MonitorService struct {
	settingRepo biz.SettingRepo
	monitorRepo biz.MonitorRepo
}

func NewMonitorService() *MonitorService {
	return &MonitorService{
		settingRepo: data.NewSettingRepo(),
		monitorRepo: data.NewMonitorRepo(),
	}
}

func (s *MonitorService) GetSetting(w http.ResponseWriter, r *http.Request) {
	setting, err := s.monitorRepo.GetSetting()
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, setting)
}

func (s *MonitorService) UpdateSetting(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.MonitorSetting](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	if err = s.monitorRepo.UpdateSetting(req); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, nil)
}

func (s *MonitorService) Clear(w http.ResponseWriter, r *http.Request) {
	if err := s.monitorRepo.Clear(); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, nil)
}

func (s *MonitorService) List(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.MonitorList](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	monitors, err := s.monitorRepo.List(carbon.CreateFromTimestampMilli(req.Start), carbon.CreateFromTimestampMilli(req.End))
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	var list types.MonitorData
	var bytesSent uint64
	var bytesRecv uint64
	var bytesSent2 uint64
	var bytesRecv2 uint64
	for _, net := range monitors[0].Info.Net {
		if net.Name == "lo" {
			continue
		}
		bytesSent += net.BytesSent
		bytesRecv += net.BytesRecv
	}
	for i, monitor := range monitors {
		// 跳过第一条数据，因为第一条数据的流量为 0
		if i == 0 {
			// MB
			list.Mem.Total = fmt.Sprintf("%.2f", float64(monitor.Info.Mem.Total)/1024/1024)
			list.SWAP.Total = fmt.Sprintf("%.2f", float64(monitor.Info.Swap.Total)/1024/1024)
			continue
		}
		for _, net := range monitor.Info.Net {
			if net.Name == "lo" {
				continue
			}
			bytesSent2 += net.BytesSent
			bytesRecv2 += net.BytesRecv
		}
		list.Times = append(list.Times, monitor.CreatedAt.ToDateTimeString())
		list.Load.Load1 = append(list.Load.Load1, monitor.Info.Load.Load1)
		list.Load.Load5 = append(list.Load.Load5, monitor.Info.Load.Load5)
		list.Load.Load15 = append(list.Load.Load15, monitor.Info.Load.Load15)
		list.CPU.Percent = append(list.CPU.Percent, fmt.Sprintf("%.2f", monitor.Info.Percent[0]))
		list.Mem.Available = append(list.Mem.Available, fmt.Sprintf("%.2f", float64(monitor.Info.Mem.Available)/1024/1024))
		list.Mem.Used = append(list.Mem.Used, fmt.Sprintf("%.2f", float64(monitor.Info.Mem.Used)/1024/1024))
		list.SWAP.Used = append(list.SWAP.Used, fmt.Sprintf("%.2f", float64(monitor.Info.Swap.Used)/1024/1024))
		list.SWAP.Free = append(list.SWAP.Free, fmt.Sprintf("%.2f", float64(monitor.Info.Swap.Free)/1024/1024))
		list.Net.Sent = append(list.Net.Sent, fmt.Sprintf("%.2f", float64(bytesSent2/1024/1024)))
		list.Net.Recv = append(list.Net.Recv, fmt.Sprintf("%.2f", float64(bytesRecv2/1024/1024)))

		// 监控频率为 1 分钟，所以这里除以 60 即可得到每秒的流量
		list.Net.Tx = append(list.Net.Tx, fmt.Sprintf("%.2f", float64(bytesSent2-bytesSent)/60/1024/1024))
		list.Net.Rx = append(list.Net.Rx, fmt.Sprintf("%.2f", float64(bytesRecv2-bytesRecv)/60/1024/1024))

		bytesSent = bytesSent2
		bytesRecv = bytesRecv2
		bytesSent2 = 0
		bytesRecv2 = 0
	}

	Success(w, list)
}