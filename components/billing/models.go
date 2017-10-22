package billing

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"os"

	"strings"

	"log"

	"donategold.me/centrifugo"
	"donategold.me/components/money"
	"donategold.me/components/widgets"
	"donategold.me/db"
)

var (
	ErrorUnaccepted = errors.New("unaccepted operation")
)

//Про поля можно посмотреть из документации
//https://tech.yandex.ru/money/doc/dg/reference/notification-p2p-incoming-docpage/
//Добавил поле USER, что бы по нему искать
//Лучше Поле Label сделать JSONB
type Bill struct {
	ID               uint
	NotificationType string `json:"notification_type" form:"notification_type"`
	OperationID      string `json:"operation_id" form:"operation_id"`
	Amount           int    `sql:"type:decimal" json:"amount" form:"amount"`
	// копейки отдельно, но учитываться в балансе не будут. Только для проверки хеша
	Kopeiki          int       `sql:"-" json:"-"`
	Currency         int       `json:"-" form:"-"`
	Codepro          bool      `json:"-" form:"-"`
	Datetime         time.Time `json:"datetime" form:"datetime"`
	CreatedAt        time.Time `json:"-" form:"-"`
	Sender           string    `sql:"index" json:"sender" form:"sender"`
	Label            string    `json:"-" form:"-"`
	UserID           uint      `sql:"index"`
	TestNotification bool      `json:"test_notification" form:"test_notification"`
	Unaccepted       bool      `json:"unaccepted" form:"unaccepted"`
	Sha1Hash         string    `sql:"-" json:"sha1_hash" form:"-"`
}

//label example "WidgetName.WidgetID.UserID.[...Other]"
type Label struct {
	UserId   uint
	Widget   string
	WidgetID uint
	Other    []string
}

func parseLabel(label string) (Label, error) {
	var l Label
	splitter := strings.Split(label, ".")
	if len(splitter) < 3 {
		return l, errors.New("No enough label data")
	}
	l.Widget = splitter[0]
	uWID, err := strconv.Atoi(splitter[1])
	if err != nil {
		log.Println(err)
		return l, err
	}
	l.WidgetID = uint(uWID)
	uUID, err := strconv.Atoi(splitter[2])
	if err != nil {
		log.Println(err)
		return l, err
	}
	l.UserId = uint(uUID)
	if len(splitter) > 3 {
		l.Other = splitter[3:]
	}
	return l, nil
}

func (b *Bill) Valid() error {
	if b.Unaccepted {
		return ErrorUnaccepted
	}
	return nil
}

func (b *Bill) Create() error {
	err := b.CompareHash()
	if err != nil {
		return err
	}
	label, err := parseLabel(b.Label)
	if err != nil {
		return err
	}
	err = b.initUser(label)
	if err != nil {
		return err
	}
	if b.TestNotification {
		return nil
	}
	if err = b.Valid(); err != nil {
		return err
	}
	tx := db.DB.Begin()
	tx.Create(b)
	money.AddTrans(tx, b.Amount, b.UserID)
	tx.Commit()
	if tx.Error != nil {
		return tx.Error
	}
	centrifugo.AddMoney(b.UserID, b.Amount)
	err = widgets.PushWidget(label.Widget, label.WidgetID, b.UserID, b.Amount, b.ID)
	if err != nil {
		return err
	}
	return nil
}

//При создании инициализирует пользователя, которому приходят деньги
func (b *Bill) initUser(label Label) error {
	if b.Label == "" && !b.TestNotification {
		return errors.New("Label blank. Who give me money?")
	}
	if label.UserId != 0 {
		b.UserID = label.UserId
	} else {
		return errors.New("UserID not found")
	}
	return nil
}

//Проверка прешдшего хеша с реальным хешом данных
//Если не совпадают, то возвращает ошибку
func (b *Bill) CompareHash() error {
	// копейки добавленны для прохождения проверки
	strOfHash := fmt.Sprintf("%s&%s&%d.%d&%d&%s&%s&%v&%s&%s", b.NotificationType, b.OperationID, b.Amount,
		b.Kopeiki, b.Currency, b.getTimeString(), b.Sender, b.Codepro, os.Getenv("YANDEX_SECRET"), b.Label)
	hashOfBill := sha1.Sum([]byte(strOfHash))
	hash := fmt.Sprintf("%x", hashOfBill)
	if hash != b.Sha1Hash {
		return errors.New("Hash not compare")
	}
	return nil
}

//функция изменения вывода времени. Стандартный RFC3339Nano не выводит наносекунды, если их ноль
//в документации яндекса написана хрень. Там приходит в формате RFC3339, а не в RFC3339Nano
func (b *Bill) getTimeString() string {
	//	_, zoneK := b.Datetime.Zone()
	//	zone := zoneK / 60 / 60
	//	minutes := zoneK % 60 % 60
	//	plus := ""
	//	if zone > 0 {
	//		plus = "+"
	//	}
	//	s := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d.%03d%s%02d:%02d",
	//		b.Datetime.Year(), b.Datetime.Month(), b.Datetime.Day(),
	//		b.Datetime.Hour(), b.Datetime.Minute(), b.Datetime.Nanosecond(), b.Datetime.Second(), plus, zone, minutes)
	s := b.Datetime.Format(time.RFC3339)
	return s

}

func GetTransactionsByUser(user_id int) []Bill {
	var bills []Bill
	db.DB.Where("user_id = ?", user_id).Find(&bills)
	return bills
}

func GetBillFromPostForm(form url.Values) (Bill, error) {
	var bill Bill
	bill.NotificationType = form.Get("notification_type")
	bill.OperationID = form.Get("operation_id")
	//todo: надо разделить по точке на амаунт и копейки
	strsAK := strings.Split(form.Get("amount"), ".")
	if len(strsAK) < 2 {
		return bill, errors.New("cant parse amount")
	}
	bill.Amount, _ = strconv.Atoi(strsAK[0])
	bill.Kopeiki, _ = strconv.Atoi(strsAK[1])
	bill.Currency, _ = strconv.Atoi(form.Get("currency"))
	bill.Codepro, _ = strconv.ParseBool(form.Get("codepro"))
	bill.Datetime, _ = time.Parse(time.RFC3339, form.Get("datetime"))
	bill.Sender = form.Get("sender")
	bill.Label = form.Get("label")
	bill.Sha1Hash = form.Get("sha1_hash")
	return bill, nil
}
