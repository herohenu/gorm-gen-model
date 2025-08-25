package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func main() {
	// 数据库连接配置
	dsn := "root:root@tcp(localhost:3306)/release_atd?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic("数据库连接失败: " + err.Error())
	}

	// 生成器配置
	g := gen.NewGenerator(gen.Config{
		OutPath:           "./po",
		Mode:              gen.WithDefaultQuery,
		FieldNullable:     true,
		FieldWithIndexTag: true,
		FieldWithTypeTag:  true,
		OutFile:           ".go",
	})

	// 关键修改：v0.3.22 版本使用 WithFileNameStrategy 方法设置文件名
	g.WithFileNameStrategy(func(tableName string) string {
		return tableName // 去除 .gen. 后缀
	})

	g.UseDB(db)

	// 自定义类型映射
	dataMap := map[string]func(gorm.ColumnType) string{
		"int":             func(columnType gorm.ColumnType) string { return "int64" },
		"bigint":          func(columnType gorm.ColumnType) string { return "int64" },
		"bigint unsigned": func(columnType gorm.ColumnType) string { return "uint64" },
		"tinyint":         func(columnType gorm.ColumnType) string { return "int" },
		"smallint":        func(columnType gorm.ColumnType) string { return "int" },
	}
	g.WithDataTypeMap(dataMap)

	// 指定需要生成的表
	tables := []string{
		"account_advertiser_detail",
	}

	// 循环处理每个表
	for _, table := range tables {
		g.GenerateModel(table)
	}

	// 执行生成
	g.Execute()
}
