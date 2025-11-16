package model

type Users struct {
	ID           string `gorm:"column:id;type:uuid;primary_key" json:"id"`
	Username     string `gorm:"column:username;type:varchar(50);NOT NULL" json:"username"`
	Password     string `gorm:"column:password;type:varchar(255);NOT NULL" json:"password"`
	FirstName    string `gorm:"column:first_name;type:varchar(100);NOT NULL" json:"firstName"`
	LastName     string `gorm:"column:last_name;type:varchar(100);NOT NULL" json:"lastName"`
	Email        string `gorm:"column:email;type:varchar(100);NOT NULL" json:"email"`
	NumberPhone  string `gorm:"column:number_phone;type:varchar(15)" json:"numberPhone"`
	NumberMobile string `gorm:"column:number_mobile;type:varchar(15)" json:"numberMobile"`
	IDNational   string `gorm:"column:id_national;type:varchar(10)" json:"iDNational"`
	CodePostal   string `gorm:"column:code_postal;type:varchar(10)" json:"codePostal"`
	NameCompany  string `gorm:"column:name_company;type:varchar(100)" json:"nameCompany"`
	ImageProfile string `gorm:"column:image_profile;type:varchar(200)" json:"imageProfile"`
	Gender       string `gorm:"column:gender;type:gender" json:"gender"`
	Address      string `gorm:"column:address;type:varchar(255)" json:"address"`
	Status       string `gorm:"column:status;type:status;NOT NULL" json:"status"`
	CreatedAt    int    `gorm:"column:created_at;type:int4;NOT NULL" json:"createdAt"`
}
