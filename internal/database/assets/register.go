package assets

import (
	"database/sql"
	model "equipmentManager/internal/database/model/tables"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func CrateAssetIndivisual(db *sql.DB, master model.AssetsMaster, asset model.Asset) (bool, error) {
	//個別管理はQuantityを1に固定する
	//フロントからは1が送られるはずだが一応再登録しておく
	asset.Quantity = 1

	tx, err := db.Begin()
	if err != nil {
		log.Println("(個別管理)トランザクション開始失敗:", err)
		return false, err
	}
	masterId, err := createMaster(tx, master)
	if err != nil {
		tx.Rollback()
		log.Println("(個別管理)マスタ登録失敗:", err)
		return false, err
	}
	asset.ItemMasterID = masterId
	err = createAsset(tx, asset)
	if err != nil {
		tx.Rollback()
		log.Println("(個別管理)資産登録失敗:", err)
		return false, err
	}
	err = tx.Commit()
	if err != nil {
		log.Println("(個別管理)トランザクションコミット失敗:", err)
		return false, err
	}
	return true, nil
}

func CreateAssetCollective(db *sql.DB, master model.AssetsMaster, asset model.Asset) (bool, error) {
	tx, err := db.Begin()
	if err != nil {
		log.Println("(全体管理)トランザクション開始失敗:", err)
		return false, err
	}
	res, err := createMaster(tx, master)
	if err != nil {
		tx.Rollback()
		log.Println("(全体管理)マスタ登録失敗:", err)
		return false, err
	}
	asset.ItemMasterID = res
	err = createAsset(tx, asset)
	if err != nil {
		tx.Rollback()
		log.Println("(全体管理)資産登録失敗:", err)
		return false, err
	}
	err = tx.Commit()
	if err != nil {
		log.Println("(全体管理)トランザクションコミット失敗:", err)
		return false, err
	}
	return true, nil
}

// createMaster は資産マスタをデータベースに登録し、登録されたマスタのIDを返す
func createMaster(tx *sql.Tx, master model.AssetsMaster) (int64, error) {
	query := "INSERT INTO assets_masters (name, management_category_id, genre_id, manufacturer, model_number) VALUES (?, ?, ?, ?, ?)"
	res, err := tx.Exec(query, master.Name, master.ManagementCategoryID, master.GenreID.Int64, master.Manufacturer, master.ModelNumber)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// createAsset は資産をデータベースに登録する
// 登録直後は場所＝デフォルト保管場所とする → 貸出時に場所を借りている人に変更するのでクエリパラメータのdefault_locationはlocationを入れる
// 貸出時にownerが決まるので、ownerはnilのまま
func createAsset(tx *sql.Tx, asset model.Asset) error {
	query := `INSERT INTO assets (asset_master_id, quantity, serial_number, status_id, purchase_date, location,default_location ,last_check_date, last_checker, notes)
              VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := tx.Exec(query, asset.ItemMasterID, asset.Quantity, asset.SerialNumber, asset.StatusID, asset.PurchaseDate, asset.Location, asset.Location, asset.LastCheckDate, asset.LastChecker, asset.Notes)
	return err
}
