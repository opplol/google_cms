package controller

import (
	"cms_sheets/lib"
	"cms_sheets/model"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	webpush "github.com/SherClockHolmes/webpush-go"
	"github.com/gin-gonic/gin"
)

type EventController struct{}

func (h EventController) Index(c *gin.Context) {
	eventModel := make([]model.Event, 0)
	file_id := h.valid_google_oauth_and_get_list(c)
	fmt.Printf("%#v\n", file_id)

	h.valid_google_oauth_and_get_data(c, file_id, &eventModel)
	c.HTML(http.StatusOK, "event_list", gin.H{
		"data_model": eventModel,
	})
}

func (h EventController) Subscribe(c *gin.Context) {
	db_conn, err := model.NewGolangAppConnection()
	if err != nil {
		panic("Db connection Fail")
	}
	subscribe_model := model.Subscribe{EndPoint: c.PostForm("endpoint"),
		UserPublicKey: c.PostForm("userPublicKey"),
		UserAuthToken: c.PostForm("userAuthToken")}

	result := db_conn.Create(&subscribe_model)
	fmt.Printf("%#v\n", subscribe_model.ID)
	if result.Error != nil {
		panic("Data save fail")
	}
	c.JSON(http.StatusOK, gin.H{"row_num": subscribe_model.ID})
}

func (h EventController) UnSubscribe(c *gin.Context) {
	db_conn, err := model.NewGolangAppConnection()
	if err != nil {
		panic("Db connection Fail")
	}
	var subscribe model.Subscribe
	db_conn.First(&subscribe, "end_point = ?", c.PostForm("end_point"))
	fmt.Printf("%#v\n", subscribe.ID)
	if subscribe.ID > 0 {
		db_conn.Delete(&subscribe)
	}

	c.JSON(http.StatusOK, gin.H{"row_num": subscribe.ID})
}

func (h EventController) PushMe(c *gin.Context) {
	// Decode subscription
	public_key := os.Getenv("WEB_PUSH_PUBLIC_KEY")
	private_key := os.Getenv("WEB_PUSH_PRIVATE_KEY")
	db_conn, err := model.NewGolangAppConnection()
	if err != nil {
		panic("Db connection Fail")
	}
	var all_subscribe []model.Subscribe
	db_conn.Find(&all_subscribe)
	for _, v := range all_subscribe {
		s := &webpush.Subscription{Endpoint: v.EndPoint, Keys: webpush.Keys{Auth: v.UserAuthToken, P256dh: v.UserPublicKey}}
		json.Unmarshal([]byte(``), s)

		// Send Notification
		resp, err := webpush.SendNotification([]byte("PUSH Test PUSHPUSH!"), s, &webpush.Options{
			Subscriber:      "example@example.com",
			VAPIDPublicKey:  public_key,
			VAPIDPrivateKey: private_key,
			TTL:             30,
		})
		if err != nil {
			panic("Push fail")
		}
		defer resp.Body.Close()

	}
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

func (EventController) valid_google_oauth_and_get_data(c *gin.Context, fileId string, sheetModel *[]model.Event) {
	author := new(lib.GoogleApiAuthSpreadEvent)
	authURL, err := author.CredentialApi(fileId, sheetModel)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, authURL)
	}
}
func (EventController) valid_google_oauth_and_get_list(c *gin.Context) string {
	author := new(lib.GoogleApiAuthDrive)
	result, err := author.CredentialApi("event_list")
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, result.AuthUrl)
	}
	return result.FileId
}
