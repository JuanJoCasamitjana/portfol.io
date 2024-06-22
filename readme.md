![Portfolio logo](https://github.com/JuanJoCasamitjana/portfol.io/blob/main/web/static/logo_4.0.png)

# Portfol\.io

Portfol\.io is an example project that aims to provide decent complexity for a handful of examples in projects that use Golang and HTMX. Just keep in mind that I'm also learning while developing this project, so if you want a nice guide on how to build your application in the Go style, you are probably better off looking elsewhere. I just selected a handful of technologies that I thought would make it easier to develop without wasting too much time on configuration. In other words, if you plan to do a small application, this is absolutely not the place to look, and if you want to go big, you should take this code with a grain of salt.

# Techlogies, frameworks, dependencies, etc..
Licenses and copyright clauses are housed on [dependencies.md](/dependencies.md)


* **Echo:** for routing and managing the web server. [Echo](https://echo.labstack.com/)
* **Gorm:** for managing the database. [Gorm\.io](https://gorm.io/)
* **Bluemonday:** for HTML sanitization. [Bluemonday](https://github.com/microcosm-cc/bluemonday)
* **Bcrypt:** for password hashing. [bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt)
* **Goeasyi18n:** for internazionalization. [goeasyi18n](https://github.com/eduardolat/goeasyi18n?tab=readme-ov-file)
* **Godotenv:** just because I wanted to use a `.env` file during development. [joho godotenv](https://github.com/joho/godotenv)
* **Gorilla sessions:** for managing sessions and coockies. [sessions](https://github.com/gorilla/sessions)
* **HTMX:** for frontend management. [HTMX](https://htmx.org/)
* **Bootstrap:** for styling. [Bootstrap 4.6](https://getbootstrap.com/docs/4.6/getting-started/introduction/)
* **Summernote:** for html editing on the client. [Summernote](https://summernote.org/)
* **Lumberjack:** for rolling logs. [lumberjack](https://github.com/natefinch/lumberjack)

# Context and features

The idea of the project is to develop a service that allows you to post different kinds of media, such as:
* Articles and other enriched text posts.
* Galleries of images with small descriptions for each image.
* The user's profile can be organized into sections.
* Users may tag different kinds of posts.
* Projects, which are links to git repositories. (Work in Progress)

Users have an account where they can post their content and organize it into sections:
* Profiles have some required information and some optional information.
* Other features have not yet been implemented.

What features are planned for the near future?

* Users may subscribe to other users. (Email may be required for this)
* Users can use tags and queries to filter and look for content.
* Users may upload audio files.
* Users may upload video files.

What features are not planned for the future?

* Recommendation system.
* Chats.


# How to set up
1. Install [go](https://go.dev/).
2. Clone the repository and move into the folder:
    ```bash
    git clone https://github.com/JuanJoCasamitjana/portfol.io.git
    cd portfol.io
    ```
3. Set up the project:
   1. This will download all the required dependencies
   ```bash
   go mod tidy
   ```
   2. You can set up a `.env` file with the following variables:
      1. IMGBB_API_KEY: An api key for the [Imgbb](https://imgbb.com) API
      2. PORT (optional): The port that the application should be started (it is `:8080` by default).
4. To start the project you have 2 options:
   * Install [air](https://github.com/cosmtrek/air) and run `air` in your terminal 
   * Execute `go run ./cmd/main.go` in your terminal


