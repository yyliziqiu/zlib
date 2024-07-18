package zdb

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"github.com/yyliziqiu/zlib/zlog"
)

type DBMigration struct {
	DB       *gorm.DB
	Once     []schema.Tabler
	Cron     []schema.Tabler
	Interval time.Duration
}

func MigrateDBS(ctx context.Context, migrations func() []DBMigration) (err error) {
	for _, migration := range migrations() {
		err = MigrateDB(ctx, migration)
		if err != nil {
			return err
		}
	}
	return nil
}

func MigrateDB(ctx context.Context, migration DBMigration) (err error) {
	db := migration.DB.Set("gorm:table_options", "ENGINE=InnoDB")

	err = migrateDB(db, migration.Once)
	if err != nil {
		return fmt.Errorf("migrate once tables error [%v]", err)
	}

	if len(migration.Cron) == 0 {
		return nil
	}

	err = migrateDB(db, migration.Cron)
	if err != nil {
		return fmt.Errorf("migrate cron tables error [%v]", err)
	}

	go runCronMigrateDB(ctx, migration.Interval, db, migration.Cron)

	return nil
}

func migrateDB(db *gorm.DB, tables []schema.Tabler) error {
	for _, table := range tables {
		err := db.Table(table.TableName()).Migrator().AutoMigrate(&table)
		if err != nil {
			return fmt.Errorf("create table [%s] failed [%v]", table.TableName(), err)
		}
	}
	return nil
}

func runCronMigrateDB(ctx context.Context, interval time.Duration, db *gorm.DB, tables []schema.Tabler) {
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ticker.C:
			err := migrateDB(db, tables)
			if err != nil {
				zlog.Errorf("Migrate DB failed, error: %v.", err)
			}
		case <-ctx.Done():
			return
		}
	}
}

type RecordMigration interface {
	Exist() (bool, error)
	Create() error
}

func MigrateRecords(migrations []RecordMigration) error {
	for _, migration := range migrations {
		exist, err := migration.Exist()
		if err != nil {
			return err
		}
		if exist {
			continue
		}
		err = migration.Create()
		if err != nil {
			return err
		}
	}
	return nil
}
