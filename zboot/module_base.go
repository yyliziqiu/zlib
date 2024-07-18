package zboot

import (
	"context"
	"fmt"

	"github.com/yyliziqiu/zlib/zdb"
	"github.com/yyliziqiu/zlib/zelastic"
	"github.com/yyliziqiu/zlib/zkafka"
	"github.com/yyliziqiu/zlib/zlog"
	"github.com/yyliziqiu/zlib/zredis"
	"github.com/yyliziqiu/zlib/zutil"
)

func BaseInit(config any) InitFunc {
	return func() (err error) {
		val, ok := zutil.StructFieldValue(config, "DB")
		if ok {
			c, ok2 := val.([]zdb.Config)
			if ok2 && len(c) > 0 {
				zlog.Info("Init DB.")
				err = zdb.Init(c...)
				if err != nil {
					return fmt.Errorf("init DB error [%v]", err)
				}
			}
		}

		val, ok = zutil.StructFieldValue(config, "Redis")
		if ok {
			c, ok2 := val.([]zredis.Config)
			if ok2 && len(c) > 0 {
				zlog.Info("Init redis.")
				err = zredis.Init(c...)
				if err != nil {
					return fmt.Errorf("init redis error [%v]", err)
				}
			}
		}

		val, ok = zutil.StructFieldValue(config, "Kafka")
		if ok {
			c, ok2 := val.([]zkafka.Config)
			if ok2 && len(c) > 0 {
				zlog.Info("Init kafka.")
				err = zkafka.Init(c...)
				if err != nil {
					return fmt.Errorf("init kafka error [%v]", err)
				}
			}
		}

		val, ok = zutil.StructFieldValue(config, "Elastic")
		if ok {
			c, ok2 := val.([]zelastic.Config)
			if ok2 && len(c) > 0 {
				zlog.Info("Init elastic.")
				err = zelastic.Init(c...)
				if err != nil {
					return fmt.Errorf("init elastic error [%v]", err)
				}
			}
		}

		return nil
	}
}

func BaseBoot() BootFunc {
	return func(ctx context.Context) error {
		go func() {
			<-ctx.Done()
			zdb.Finally()
			zredis.Finally()
			zkafka.Finally()
			zelastic.Finally()
		}()
		return nil
	}
}
