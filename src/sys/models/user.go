package models

import (
    "alex/utils"
    "net/url"
    "net/http"
    "sys/core"
    "alex/errors"
    sysErrors "errors"
    "alex/qrcode"
    "time"
)

//通过http.Request获取用户token
func UserToken(req *http.Request) (string, error) {
    c, err := req.Cookie(core.USER_TOKEN)
    if err != nil {
        return "", err
    }
    if c.Value == "" {
        return "", sysErrors.New("token is null")
    }
    return c.Value, nil
}

//fcds.f_users表模型
type FUsers struct {
    ID            uint `gorm:"AUTO_INCREMENT"`
    Username      string `gorm:"type:varchar(50)"`
    Password      string `gorm:"type:varchar(100)"`
    Type          string `gorm:"type:tinyint(2)"`
    Token         string `gorm:"type:varchar(50)"`
    TotpSecret    string `gorm:"type:varchar(100)"`
    TotpUrl       string `gorm:"type:varchar(300)"`
    Mail          string `gorm:"type:varchar(100)"`
    CreateTime    time.Time `gorm:"datetime"`
    UpdateTime    time.Time `gorm:"datetime"`
    LastLoginTime time.Time `gorm:"datetime"`
}


//字段对应中文
func (t FUsers) FieldCn(field string) string {
    labels := map[string]string{
        "username" : "用户名称",
        "type" : "用户类型",
        "password" : "用户密码",
    }
    if label, ok := labels[field]; ok {
        return label
    }
    return field
}




//自定义表名
/*func (t *FUsers) TableName() string {
    return "f_users"
}*/


func NewUserModel() *FUsers {
    return new(FUsers)
}


//判断token是否合法
func GetUserInfo(token string) (*FUsers, error) {
    conn, err := core.Mysql()
    if err != nil {
        return nil, err
    }
    defer core.PoolRollBack(conn)

    user := new(FUsers)
    err = conn.Where("token=?", token).First(&user).Error
    if err != nil {
        return nil, err
    }
    if user.ID == 0 {
        return nil, sysErrors.New("token invalid")
    }

    return user, nil
}

//用户登录校验
func (m *FUsers)Login(post url.Values) (string, errors.IMyError) {
    db, err := core.Mysql()
    if err != nil {
        return "", errors.NewCodeErr(core.ERR_LOGIN_SERVICE, "登录接口错误")
    }
    defer core.PoolRollBack(db)

    username := post.Get("username")
    password := post.Get("password")
    totp := post.Get("totp")
    if username == "" || password == "" || totp == "" {
        return "", errors.NewCodeErr(core.ERR_LOGIN_PARAMS, "登录参数错误")
    }

    user := new(FUsers)
    err = db.Where("username=? AND password=?", username, password).First(&user).Error
    if err != nil {
        return "", errors.NewCodeErr(core.ERR_LOGIN_SERVICE, "登录接口错误", err.Error())
    }
    if user.ID == 0 {
        return "", errors.NewCodeErr(core.ERR_LOGIN_NOTFOUND, "用户不存在")
    }

    //totp检测
    equal, err := utils.TotpCode(totp, user.TotpSecret)
    if err != nil || !equal {
        return "", errors.NewCodeErr(core.ERR_LOGIN_TOTP, "totp码错误")
    }

    //token生成
    user.Token = utils.MD5(utils.StringJoin(user.Username, utils.GetMicroTimer()))
    if err := db.Save(user).Error; err != nil {
        return "", errors.NewCodeErr(core.ERR_LOGIN_SERVICE, "登录接口错误", err.Error())
    }

    return user.Token, nil
}





/********************************************** 接口 begin **********************************************/

func (t *FUsers) TotpImage(req *http.Request) (map[string]interface{}, errors.IMyError) {
    db, err := core.Mysql()
    if err != nil {
        return nil, errors.NewCodeErr(core.ERR_POOL_MYSQL, err)
    }
    defer core.PoolRollBack(db)

    idStr := req.FormValue("id")
    id := utils.MustInt(idStr)
    if id == 0 {
        return nil, errors.NewCodeErr(core.ERR_USER_DETAIL, "id必传")
    }

    user := new(FUsers)
    db.Where("id=?", idStr).First(user)
    if user.ID == 0 {
        return nil, errors.NewCodeErr(core.ERR_USER_DETAIL, "记录不存在")
    }

    pngName := user.Username + ".png"
    fileName := "public/" + user.Username + ".png"
    err = qrcode.Png(user.TotpUrl, fileName, 200, 200)
    if err != nil {
        return nil, errors.NewCodeErr(core.ERR_USER_DETAIL, "生成图片错误", err)
    }
    return map[string]interface{}{"path":pngName}, nil
}


//获取用户详情
func (t *FUsers) Get(req *http.Request) (*UserResult, errors.IMyError) {
    db, err := core.Mysql()
    if err != nil {
        return nil, errors.NewCodeErr(core.ERR_POOL_MYSQL, err)
    }
    defer core.PoolRollBack(db)

    idStr := req.FormValue("id")
    id := utils.MustInt(idStr)
    if id == 0 {
        return nil, errors.NewCodeErr(core.ERR_USER_DETAIL, "id必传")
    }

    user := new(FUsers)
    db.Where("id=?", idStr).First(user)

    if user.ID == 0 {
        return nil, errors.NewCodeErr(core.ERR_USER_DETAIL, "记录不存在")
    }
    return t.UserResult(user), nil

}

//获取用户分页列表
func (t *FUsers) PageList(req *http.Request) (interface{}, errors.IMyError) {
    db, err := core.Mysql()
    if err != nil {
        return nil, errors.NewCodeErr(core.ERR_POOL_MYSQL, err)
    }
    defer core.PoolRollBack(db)

    var list []*FUsers
    res, err := PageList(req, &list, db.Order("id DESC"));
    if err != nil {
        return nil, errors.NewCodeErr(core.ERR_OTHER, err)
    }
    return res, nil
}

//获取管理员类型的中文
func (t *FUsers) GetTypeCn() string {
    if t.Type == "1" {
        return "系统管理员"
    } else {
        return "普通用户"
    }
}

//添加
func (t *FUsers) Create(req *http.Request) (interface{}, errors.IMyError) {
    username := req.FormValue("username")
    if username == "" {
        return nil, errors.NewCodeErr(core.ERR_USER_CREATE, t.FieldCn("username") + "必传")
    }
    password := req.FormValue("password")
    if password == "" {
        return nil, errors.NewCodeErr(core.ERR_USER_CREATE, t.FieldCn("password") + "必传")
    }

    typeStr := req.FormValue("type")
    if typeStr == "" {
        return nil, errors.NewCodeErr(core.ERR_USER_CREATE, "请选择" + t.FieldCn("type"))
    }

    db, err := core.Mysql()
    if err != nil {
        return nil, errors.NewCodeErr(core.ERR_POOL_MYSQL, err)
    }
    defer core.PoolRollBack(db)

    totpUrl, totpSecret, err := utils.TotpUrlAndSecret("fcds", username)
    if err != nil {
        return nil, errors.NewCodeErr(core.ERR_USER_CREATE, err)
    }

    user := new(FUsers)
    user.TotpUrl = totpUrl
    user.TotpSecret = totpSecret
    user.Type = typeStr
    user.Username = username
    user.Password = utils.MD5(password)
    user.CreateTime = time.Now()
    if err := db.Save(user).Error; err != nil {
        return nil, errors.NewCodeErr(core.ERR_POOL_MYSQL, err)
    }

    return map[string]interface{}{"id":user.ID}, nil
}

//修改
func (t *FUsers) Update(req *http.Request) (interface{}, errors.IMyError) {
    id := utils.MustInt(req.FormValue("id"))
    if id == 0 {
        return nil, errors.NewCodeErr(core.ERR_USER_CREATE, "id不能为空")
    }

    db, err := core.Mysql()
    if err != nil {
        return nil, errors.NewCodeErr(core.ERR_POOL_MYSQL, err)
    }
    defer core.PoolRollBack(db)

    user := new(FUsers)

    db.Where("id=?", id).First(user)
    if user.ID == 0 {
        return nil, errors.NewCodeErr(core.ERR_USER_EDIT, "记录不存在")
    }
    username := req.FormValue("username")
    if username == "" {
        return nil, errors.NewCodeErr(core.ERR_USER_CREATE, t.FieldCn("username") + "必传")
    }
    typeStr := req.FormValue("type")
    if typeStr == "" {
        return nil, errors.NewCodeErr(core.ERR_USER_CREATE, "请选择" + t.FieldCn("type"))
    }
    user.Type = typeStr
    user.Username = username
    if password := req.FormValue("password"); password != "" {
        user.Password = utils.MD5(password)
    }
    user.CreateTime = time.Now()
    user.UpdateTime = time.Now()
    if err := db.Save(user).Error; err != nil {
        return nil, errors.NewCodeErr(core.ERR_POOL_MYSQL, err)
    }

    return nil, nil
}

//删除
func (t *FUsers) Del(req *http.Request, userInfo *ContextInfo) (interface{}, errors.IMyError) {
    id := utils.MustInt(req.FormValue("id"))
    if id == 0 {
        return nil, errors.NewCodeErr(core.ERR_USER_DEL, "id不能为空")
    }

    db, err := core.Mysql()
    if err != nil {
        return nil, errors.NewCodeErr(core.ERR_POOL_MYSQL, err)
    }
    defer core.PoolRollBack(db)

    user := new(FUsers)
    db.Where("id=?", id).First(user)
    if user.ID == 0 {
        return nil, errors.NewCodeErr(core.ERR_USER_DEL, "记录不存在")
    }

    if err := validateTotp(req.FormValue("totp"), userInfo); err != nil {
        return nil, errors.NewCodeErr(core.ERR_USER_DEL, err)
    }
    if err := db.Delete(user).Error; err != nil {
        return nil, errors.NewCodeErr(core.ERR_USER_DEL, err)
    }
    return nil, nil
}


/********************************************** 以下和json相关 **********************************************/
type UserResult struct {
    Id         uint `json:"id"`
    Username   string `json:"username"`
    Totp       string `json:"totp"`
    Type       string `json:"type"`
    CreateTime time.Time `json:"createTime"`
}

func (t FUsers) UserResult(user *FUsers) *UserResult {
    return &UserResult{
        Id:user.ID,
        Username:user.Username,
        Totp:user.TotpUrl,
        Type:user.Type,
        CreateTime: time.Time(user.CreateTime),
    }
}

func (t FUsers) UserListResult(users []*FUsers) []*UserResult {
    var lists []*UserResult
    for _, user := range users {
        lists = append(lists, &UserResult{
            Id:user.ID,
            Username:user.Username,
            Totp:user.TotpUrl,
            Type:user.Type,
            CreateTime: user.CreateTime,
        })
    }
    return lists
}


