package initialize

import (
	"database/sql"
	"fmt"
	"go-learning/global"
)

// docker exec -it c930a1205c22 bash
// mysql -uroot -proot123
// use shopdevgo;
// show tables;

// sqlc generate

/*
POSTGRESQL 9.2 vs MYSQL

1. Phình chỉ mục (Index Bloat) & Ghi lặp (Write Amplification):
- POSTGRESQL: khi bạn UPDATE một dòng dữ liệu, hệ thống không ghi đè trực tiếp mà sẽ tạo ra một phiên bản mới của dòng đó ở một vị trí vật lý khác trên đĩa (cơ chế MVCC).
Vì tất cả các indexes của bảng đều trỏ thẳng vào vị trí vật lý (gọi là ctid) của dòng dữ liệu, nên chỉ cần bạn thay đổi một cột nhỏ xíu, vị trí mới được sinh ra,
Postgres sẽ phải cập nhật lại vị trí mới cho TOÀN BỘ các indexes của bảng đó.
(Lưu ý: Postgres vốn có cơ chế tối ưu HOT (Heap-Only Tuples) để tránh cập nhật index, nhưng cơ chế này chỉ hoạt động khi cột được cập nhật không nằm trong index và trang đĩa đó còn chỗ trống.
Ở quy mô ghi dữ liệu khổng lồ của Uber, các bảng thường có rất nhiều index phụ bao phủ hầu hết các cột, khiến cơ chế HOT gần như bị vô hiệu hóa.
Kết quả là đĩa cứng phải ghi liên tục (Write Amplification) và dung lượng index phình to khủng khiếp.)
- MYSQL: các indexes (cả Primary và Secondary) chỉ trỏ về Primary Key. Khi cập nhật dòng, nếu Primary Key không đổi thì các indexes không cần cập nhật theo.

2. Nghẽn mạng khi đồng bộ dữ liệu (Replication Lag):
- POSTGRESQL: dùng Đồng bộ vật lý (Physical Replication) File WAL ghi lại chi tiết từng byte, từng block vật lý thay đổi trên đĩa cứng.
Vì vụ phình index ở trên, số lượng block thay đổi trên đĩa là cực kỳ khổng lồ. Postgres phải gửi toàn bộ các file WAL nặng nề này qua mạng WAN giữa các trung tâm dữ liệu, gây trễ đồng bộ nghiêm trọng.
- MYSQL: dùng Đồng bộ logic (Logical Replication) Binlog của MySQL chỉ gửi đi thông điệp logic ngắn gọn (ví dụ: "Dòng có ID = 5 vừa được cập nhật cột X thành giá trị Y").
Lượng dữ liệu truyền đi siêu nhẹ nên đồng bộ mượt mà hơn rất nhiều.

3. Nâng cấp phiên bản là một "cơn ác mộng" gây downtime kéo dài:
- POSTGRESQL: yêu cầu máy chủ chính (master) và máy chủ phụ (replica) phải chạy cùng một phiên bản lớn (bản 9.3 không thể đồng bộ trực tiếp sang bản 9.2).
Họ không thể nâng cấp cuốn chiếu (rolling upgrade) kiểu "máy này chạy trước, máy kia chạy sau" được. Quy trình nâng cấp bắt buộc phải
tắt master -> chạy công cụ pg_upgrade -> bật lại master -> xóa sạch dữ liệu trên các replica và sao chép lại toàn bộ terabytes dữ liệu mới từ master qua mạng.


POSTGRESQL 17+
- Gộp chỉ mục trùng lặp (B-Tree deduplication) giúp giảm dung lượng index một cách đáng kể
- Cơ chế dọn rác (autovacuum) và quản lý trang đĩa cũng thông minh hơn trước rất nhiều
- Hỗ trợ Logical Replication cực kỳ xịn sò và ổn định. Thay vì gửi cả file log vật lý cồng kềnh,
Postgres giờ cũng có thể chỉ gửi đi các thay đổi logic siêu nhẹ giống như MySQL, giải quyết dứt điểm bài toán nghẽn mạng và lệch dữ liệu giữa các máy chủ
- Nhờ Logical Replication, việc nâng cấp giữa các phiên bản Postgres lớn hiện nay đã có thể thực hiện theo kiểu cuốn chiếu (rolling upgrade).
Bạn có thể dựng một máy replica chạy bản mới, đồng bộ logic từ máy master bản cũ, sau đó chuyển vùng (failover) một cách êm ái với thời gian dừng hệ thống gần như bằng không.
*/

func InitMysqlC() {
	// TODO: Init mysql from file or env
	m := global.Config.Mysql

	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	dsn := "%s:%s@tcp(%s:%v)/%s?charset=utf8mb4&parseTime=True&loc=Local"
	var s = fmt.Sprintf(dsn, m.UserName, m.Password, m.Host, m.Port, m.DBName)

	db, err := sql.Open("mysql", s)
	checkErrorPanic(err, "InitMysqlC initialization error")
	global.Logger.Info("Initializing MySQLC successfully")
	global.Mdbc = db
}
