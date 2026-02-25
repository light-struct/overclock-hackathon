import (
    "github.com/pocketbase/pocketbase"
    "github.com/pocketbase/pocketbase/core"
)
func main() {

// блок pocketbase
app := pocketbase.New()

app.OnServe().BindFunc(func(e *core.ServeEvent) error {

    // 1. РЕГИСТРАЦИЯ
    e.Router.POST("/signup", func(c *core.RequestEvent) error {
       email := c.Request.FormValue("email")
       password := c.Request.FormValue("password")
       username := c.Request.FormValue("username")
       role := c.Request.FormValue("role")

       // Валидация роли (защита от того, чтобы кто-то не стал админом просто так)
       if role == "" {
          role = "student" // По умолчанию
       }

       collection, _ := app.FindCollectionByNameOrId("users")
       record := core.NewRecord(collection)

       record.Set("email", email)
       record.Set("password", password)
       record.Set("passwordConfirm", password)
       record.Set("username", username)
       record.Set("role", role) // ЗАПИСЫВАЕМ РОЛЬ

       // Сохраняем пользователя в базу
       if err := app.Save(record); err != nil {
          return c.String(http.StatusBadRequest, "Ошибка регистрации: "+err.Error())
       }

       return c.String(http.StatusCreated, "Пользователь создан!")
    })

    // 2. ЛОГИН
    e.Router.POST("/login", func(c *core.RequestEvent) error {
       email := c.Request.FormValue("email")
       password := c.Request.FormValue("password")

       authRecord, err := app.FindAuthRecordByEmail("users", email)
       if err != nil || !authRecord.ValidatePassword(password) {
          return c.String(http.StatusUnauthorized, "Неверный логин или пароль")
       }

       // Генерация токена
       token, err := authRecord.NewAuthToken()
       if err != nil {
          return c.String(http.StatusInternalServerError, "Ошибка создания токена")
       }

       // возвращаю роль
       userRole := authRecord.GetString("role")

       // возвращаю JSON
       return c.JSON(http.StatusOK, map[string]string{
          "token":    token, // PocketBase выдает токен автоматически в заголовках, но можно и тут
          "username": authRecord.GetString("username"),
          "role":     userRole,
       })
    })

    // 3. ЛОГАУТ
    e.Router.POST("/logout", func(c *core.RequestEvent) error {
       // В JWT-авторизации сервер обычно не хранит состояние сессии.
       // Но если тебе нужно сделать проверку или логику на сервере:
       return c.String(http.StatusOK, "Вы успешно вышли. Удалите токен на клиенте.")
    })

    return nil
})

if err := app.Start(); err != nil {
    log.Fatal(err)
}
}