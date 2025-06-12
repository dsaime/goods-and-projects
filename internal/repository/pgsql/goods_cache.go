package pgsql

// GoodsCache взаимодействует с кэшем, читать сохранить
type GoodsCache interface {
	Get(id int) (Good, error)
	Set(Good) error
	Delete(id int) error
}
