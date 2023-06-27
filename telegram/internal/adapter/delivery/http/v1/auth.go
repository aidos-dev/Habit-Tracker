package v1

// func (a *AdapterHandler) SignUp(c *gin.Context, username string) {
// 	log.Print("adapter SignUp method called")

// 	type Request struct {
// 		Name string `json:"tg_user_name"`
// 	}

// 	requestData := Request{Name: username}

// 	requestBody, err := json.Marshal(requestData)
// 	if err != nil {
// 		c.String(http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	resp, err := http.Post("http://localhost:8000/telegram/auth/sign-up", "application/json", bytes.NewBuffer(requestBody))
// 	if err != nil {
// 		c.String(http.StatusInternalServerError, err.Error())
// 		return
// 	}
// 	defer resp.Body.Close()

// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		c.String(http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	log.Printf("response body: %v", body)

// 	c.String(http.StatusOK, string(body))

// 	// c.JSON(http.StatusOK, map[string]interface{}{
// 	// 	"tg_user_name": username,
// 	// })
// }

// func (h *Handler) deleteUser(c *gin.Context) {
// 	userId, err := getUserId(c)
// 	if err != nil {
// 		return
// 	}

// 	deletedUserId, err := h.services.User.DeleteUser(userId)
// 	if err != nil {
// 		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error from handler: delete user: %v", err.Error()))
// 		return
// 	}

// 	response := map[string]any{
// 		"Status":         "ok",
// 		"deleted userId": deletedUserId,
// 	}

// 	c.JSON(http.StatusOK, response)
// }

// func prepareUserByClient(c *gin.Context, user models.User, clientType string) models.User {
// 	switch clientType {
// 	case webClient:
// 		user = webUserFormat(c, user)
// 	case telegramClient:
// 		user = telegramUserFormat(c, user)
// 	}

// 	return user
// }

// /*
// webUserFormat prepares user input to be
// registered as a web user, using all required credentials. Also
// it created a telegram userName as a copy of web
// userName just as a placeholder. In the future a web user
// will be able to replace it with real telegram userName
// */
// func webUserFormat(c *gin.Context, user models.User) models.User {
// 	user.TgUsername = fmt.Sprintf("copy_u:%s", user.Username)

// 	return user
// }

// /*
// telegramUserFormat prepares user input to be
// registered as a telegram user, using only
// user name. Otherwise the repository layer will not
// allow to create a user without other credentials
// */
// func telegramUserFormat(c *gin.Context, user models.User) models.User {
// 	var emptyUser models.User // emptyUser created only to return it in case of error

// 	if user.TgUsername == "" {
// 		newErrorResponse(c, http.StatusBadRequest, "error: user name is not specified")
// 		return emptyUser
// 	}

// 	if user.Username == "" {
// 		user.Username = fmt.Sprintf("copy_tg:%s", user.TgUsername)
// 	}

// 	if user.FirstName == "" {
// 		user.FirstName = fmt.Sprintf("c_tg:%s", user.TgUsername)
// 	}

// 	if user.LastName == "" {
// 		user.LastName = fmt.Sprintf("c_tg:%s", user.TgUsername)
// 	}

// 	if user.Email == "" {
// 		user.Email = fmt.Sprintf("c_tg:%s", user.TgUsername)
// 	}

// 	return user
// }
