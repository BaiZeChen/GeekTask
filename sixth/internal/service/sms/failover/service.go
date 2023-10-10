package failover

import (
	"GeekTask/sixth/internal/domain"
	"GeekTask/sixth/internal/repository"
	"GeekTask/sixth/internal/service/sms"
	"context"
	"errors"
	"sort"
	"time"
)

/*
	由于goroutine过多，就不写测试了（测也测不出啥）
	思想概述：通过平均响应时间，来对各个服务进行优先级排序
*/

func NewFailoverSMSService(svcs []sms.Service, smsRepo repository.SMSRepository, retry, threshold int) *FailoverSMSService {
	if retry <= 0 {
		// 默认重试三次
		retry = 3
	}
	if threshold <= 0 {
		threshold = 1
	}
	var factory []*smsServiceFactory
	for _, svc := range svcs {
		factory = append(factory, &smsServiceFactory{
			svcs: svc,
		})
	}

	failoverSMSService := &FailoverSMSService{
		smsFactory: factory,
		retry:      retry,
		threshold:  threshold,
		ch:         make(chan struct{}),
		ticker:     time.NewTicker(3 * time.Hour),   // 3个小时一次
		monitor:    time.NewTicker(1 * time.Minute), // 1分钟1次
		smsRepo:    smsRepo,
	}
	go failoverSMSService.avgTime()        // 监控平均响应时间
	go failoverSMSService.monitorTimeout() // 用来处理异步数据

	return failoverSMSService
}

type smsServiceFactory struct {
	svcs        sms.Service
	respHistory []int
	avgRespTime int
}

type FailoverSMSService struct {
	smsFactory []*smsServiceFactory
	retry      int
	threshold  int           // 平均响应时间的阈值，超过则进行顺序切换;单位1s
	ch         chan struct{} // 同步信号
	ticker     *time.Ticker
	monitor    *time.Ticker
	smsRepo    repository.SMSRepository
}

func (f *FailoverSMSService) Send(ctx context.Context, biz string, args []string, numbers ...string) error {
	_, ok := ctx.Deadline()
	if !ok {
		// 如果没设置，给个默认3s
		var cancelFunc context.CancelFunc
		ctx, cancelFunc = context.WithTimeout(ctx, 3*time.Second)
		defer cancelFunc()
	}
	factory := f.smsFactory[0] // 第一个代表平均响应时间最低的
	ch := make(chan struct{})
	start := time.Now().Unix()
	go func() {
		defer func() {
			ch <- struct{}{}
		}()
		err := factory.svcs.Send(ctx, biz, args, numbers...)
		if err != nil {
			// 可以先用日志记录一下错误
			// 将这个报错的先排除slice外，然后进行循环，直到找到一个能用的
			for i := 1; i < len(f.smsFactory); i++ {
				factory = f.smsFactory[i]
				err := factory.svcs.Send(ctx, biz, args, numbers...)
				if err != nil {
					// 可以先用日志记录一下错误
					continue
				}
				// 然后将这个能用的与首位交换位置
				f.smsFactory[0], f.smsFactory[i] = f.smsFactory[i], f.smsFactory[0]
				break
			}
		}
	}()

	select {
	case <-ctx.Done():
		end := time.Now().Unix()
		factory.respHistory = append(factory.respHistory, int(end-start))
		// 代表出现超时，需要异步进行入库操作
		go func() {
			extra := domain.SMSCallBackArgs{
				Biz:     biz,
				Args:    args,
				Numbers: numbers,
			}
			err := f.smsRepo.Create(context.Background(), extra)
			if err != nil {
				// 打印一下日志就好，不做过多的处理
				return
			}
		}()
		return errors.New("请求超时")
	case <-ch:
		end := time.Now().Unix()
		factory.respHistory = append(factory.respHistory, int(end-start))
		return nil
	}
}

func (f *FailoverSMSService) avgTime() {
	go func() {
		// 同步操作，只要有改变就更改顺序
		for {
			<-f.ch
			sort.Slice(f.smsFactory, func(i, j int) bool {
				return f.smsFactory[i].avgRespTime > f.smsFactory[j].avgRespTime
			})
		}
	}()
	for {
		select {
		case <-f.ticker.C:
			var sign bool
			for _, factory := range f.smsFactory {
				if len(factory.respHistory) == 0 {
					continue
				}
				sign = true
				count, avg, length := 0, 0, len(factory.respHistory)
				for i := 0; i < length; i++ {
					count += factory.respHistory[i]
				}
				avg = count / length // 取整，小数点不要
				factory.avgRespTime = avg
			}
			if sign {
				f.ch <- struct{}{}
			}
		}
	}
}

func (f *FailoverSMSService) monitorTimeout() {
	for {
		select {
		case <-f.ticker.C:
			list := f.smsRepo.GetRetryList(context.Background())
			if len(list) == 0 {
				continue
			}
			for _, value := range list {
				restry := value
				go func(args domain.SMSCallBackArgs) {
					smsFactoryLen := len(f.smsFactory)
					var success bool
					for i := 0; i < f.retry; i++ {
						SMSServic := f.smsFactory[i%smsFactoryLen]
						err := SMSServic.svcs.Send(context.Background(), args.Biz, args.Args, args.Numbers...)
						if err != nil {
							continue
						}
						success = true
					}
					if success {
						err := f.smsRepo.UpdateStatus(context.Background(), 2)
						if err != nil {
							return
						}
						return
					}
					err := f.smsRepo.UpdateStatus(context.Background(), 3)
					if err != nil {
						return
					}
				}(restry)
			}
		}
	}
}
