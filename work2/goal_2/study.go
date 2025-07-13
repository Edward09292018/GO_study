package main

import (
	"errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

// Account 账户表，包含账户ID和余额
type Account struct {
	ID      uint
	Balance float64 // 账户余额
}

// Transaction 转账交易表，记录每笔转账信息
type Transaction struct {
	ID            uint
	FromAccountID uint    // 转出账户ID
	ToAccountID   uint    // 转入账户ID
	Amount        float64 // 转账金额
}

// Transfer 实现从账户A向账户B转账的事务操作
func Transfer(db *gorm.DB, fromAccountID, toAccountID uint, amount float64) error {
	// 开始事务
	return db.Transaction(func(tx *gorm.DB) error {
		// 1. 获取转出账户
		var fromAccount Account
		if err := tx.First(&fromAccount, fromAccountID).Error; err != nil {
			return err
		}

		// 2. 检查余额是否足够
		if fromAccount.Balance < amount {
			return errors.New("insufficient balance")
		}

		// 3. 获取转入账户
		var toAccount Account
		if err := tx.First(&toAccount, toAccountID).Error; err != nil {
			return err
		}

		// 4. 执行转账操作
		fromAccount.Balance -= amount
		toAccount.Balance += amount

		// 更新转出账户余额
		if err := tx.Save(&fromAccount).Error; err != nil {
			return err
		}

		// 更新转入账户余额
		if err := tx.Save(&toAccount).Error; err != nil {
			return err
		}

		// 5. 记录交易信息
		transaction := Transaction{
			FromAccountID: fromAccountID,
			ToAccountID:   toAccountID,
			Amount:        amount,
		}

		if err := tx.Create(&transaction).Error; err != nil {
			return err
		}

		// 返回nil表示事务成功完成
		return nil
	})
}

func main() {
	// 数据库连接DSN
	dsn := "edward:123@tcp(127.0.0.1:3306)/mydb?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 自动迁移表结构
	err = db.AutoMigrate(&Account{}, &Transaction{})
	if err != nil {
		return
	}

	// 初始化账户数据（如果需要）
	//db.Debug().FirstOrCreate(&Account{ID: 1, Balance: 1000})
	db.Debug().FirstOrCreate(&Account{ID: 2, Balance: 500})

	// 执行转账操作
	err = Transfer(db, 1, 2, 100)
	if err != nil {
		log.Printf("Transfer failed: %v", err)
	} else {
		log.Println("Transfer completed successfully")
	}
}
