package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"encoding/json"
	_ "encoding/json"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/kataras/iris/v12"
	"time"
	"xorm.io/xorm"
	_ "xorm.io/xorm"
	"xorm.io/xorm/log"
)
var engine *xorm.Engine
type ResponseBean struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}
type User struct {
	Username string `json:"username" xorm:"VARCHAR(30)"`
	Password string `json:"password" xorm:"VARCHAR(30)"`
}
type Token struct{
	Token string `json:"token"`
}

type UserInformation struct {
	Id       string `json:"id" xorm:"VARCHAR(255)"`
	Username string `json:"username" xorm:"VARCHAR(255)"`
	Age      int `json:"age" xorm:"INT"`
	Gender string `json:"gender" xorm:"VARCHAR(255)"`
	Height float64 `json:"height" xorm:"DOUBLE"`

}
func GenSuccessData(data interface{}, msg string) *ResponseBean {
	return &ResponseBean{200, msg, data}
}

func GenSuccessMsg(msg string) *ResponseBean {
	return &ResponseBean{200, msg, ""}
}

func GenFailedMsg(errMsg string) *ResponseBean {
	return &ResponseBean{400, errMsg, ""}
}




func pad(src []byte) []byte {
   padding := aes.BlockSize - len(src)%aes.BlockSize
   padtext := bytes.Repeat([]byte{byte(padding)}, padding)
   return append(src, padtext...)
}
func unpad(src []byte) ([]byte, error) {
	length := len(src)
	unpadding := int(src[length-1])
	if unpadding > length {
		return nil, errors.New("unpad error. This could happen when incorrect encryption key is used")
	}
	return src[:(length - unpadding)], nil
}
// 加密解密
func encrypt(key []byte, text string) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
			return "", err
		}
	msg := pad([]byte(text))
	ciphertext := make([]byte, aes.BlockSize+len(msg))
	//没有向量，用的空切片
	iv := make([]byte,aes.BlockSize)
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], msg)
	//finalMsg := (base64.StdEncoding.EncodeToString(ciphertext))
	finalMsg := hex.EncodeToString([]byte(ciphertext[aes.BlockSize:]))
	// fmt.Println(hex.EncodeToString([]byte(ciphertext[aes.BlockSize:])))
	return finalMsg, nil

}
func decrypt(key []byte, text string) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	decodedMsg,_ := hex.DecodeString(text)
	iv  :=make([]byte,aes.BlockSize)
	msg := decodedMsg
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(msg,msg)
	unpadMsg, err := unpad(msg)
	if err != nil {
		return "", err
	}
	return string(unpadMsg), nil
}







func main()  {
	key := []byte("0123456789abcdef")
	var err error
	engine, err = xorm.NewEngine("mysql", "root:12345@/mybatis?charset=utf8")
	if err != nil{
		panic(err.Error())
	}
	defer engine.Close()

	engine.ShowSQL(true)
	engine.Logger().SetLevel(log.LOG_DEBUG)


	app := iris.New()
	app.Get("/getRequest", func(context iris.Context){
		path := context.Path()
		app.Logger().Info(path)
		context.JSON(GenSuccessMsg("第一条json"))
	})
	//{
	//	"code": 200,
	//	"msg": "Login success",
	//	"data": {
	//	"time": "2021-07-20 17:47:33",
	//		"token": "c8ecdee5f7668ffcee24486341c779016c6939a394f166c7953602cd191b73c64899309d376838c3dbb654b6636a84c6"
	//}
	//}
	app.Post("/user/login", func(context iris.Context){
		c := &User{}
		if err := context.ReadJSON(c); err != nil{
			context.JSON(GenFailedMsg("Something happen wrong"))
			panic(err.Error())
		}else{
			session := engine.Table("user")
			has, _ := session.Where("username=?",c.Username).And("password=?", c.Password).Get(c)
			if has == true{
				plaintextjson, _ := json.Marshal(c)
				// 生成token
				if encryptText, err := encrypt(key, string(plaintextjson)); err != nil{
					panic(err.Error())
				}else {
					var reData map[string]string
					reData = make(map[string]string)
					reData["token"] = encryptText
					reData["time"] = time.Now().Format("2006-01-02 15:04:05")
					context.JSON(GenSuccessData(reData, "Login success"))
				}
			}
		}
	})
	app.Post("/user/register", func(context iris.Context) {
		c := &User{}
		if err := context.ReadJSON(c); err != nil{
			panic(err.Error())
		}else{
			session := engine.Table("user")
			_, err := session.Insert(c)
			if err != nil {
				panic(err.Error())
			}
			context.JSON(GenSuccessMsg("注册成功"))
		}
	})
	app.Post("/user/information", func(context iris.Context){
		t := &Token{}
		Information := &UserInformation{}
		u := &User{}
		if err := context.ReadJSON(t); err != nil{
			panic(err.Error())
		}else{
			sessionUser := engine.Table("user")
			// 对token解码
			if rawText, err := decrypt(key, t.Token); err != nil{
				panic(err.Error())
			}else{
				// string 转 struct
				if err := json.Unmarshal([]byte(rawText), &u); err != nil{
					panic(err.Error())
				}else{
					has,_ := sessionUser.Where("username=?",u.Username).And("password=?", u.Password).Get(u)
					if has == true{
						sessionInformation := engine.Table("user_information")
						//sessionInformation.Where("username=?", u.Username).Get(Information)
						if get, err := sessionInformation.Where("username=?", u.Username).Get(Information); err != nil{
							panic(err.Error())
						}else {
							if get == true{
								context.JSON(GenSuccessData(Information, "获取信息成功"))
							}
						}
						//sessionInformation.Where("username=?", u.Username).Get(Information)
					}else {
						context.JSON(GenFailedMsg("token不正确"))
					}
				}
			}
		}
	})
	app.Run(iris.Addr(":8080"),iris.WithoutServerError(iris.ErrServerClosed))
}
