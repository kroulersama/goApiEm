package repository

import (
	"net/http"
	"sync"
	"time"

	"dev.gaijin.team/go/golib/logger"
	"gorm.io/gorm"
)

type Repository struct {
	DB  *gorm.DB
	Log *logger.Logger
}

type serviceRepository struct {
	gorm.Model
	mu           sync.Mutex ""
	esrvice_name string     ""
	price        int64      ""
	user_id      int64      ""
	start_date   time.Time  ""
	end_date     time.Time  ""
}

// NewServicserviceRepository - создаёт новую запись сервиса
func (r *Repository) NewServicserviceRepository(w http.ResponseWriter, req *http.Request) *serviceRepository {
	r.Log.Info("Create service")

}

func Create() {
	// todo
}

func Reed() {
	// todo
}

func Update() {
	// todo
}

func Delete() {
	//todo
}

func GetAllCost() {
	//todo
}
