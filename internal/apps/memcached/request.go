package memcached

type UpdateConfig struct {
	Config string `form:"config" json:"config" validate:"required"`
}
