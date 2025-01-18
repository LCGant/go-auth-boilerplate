# Go-Auth-Boilerplate

A lightweight and modular authentication boilerplate built with Go (Golang) and Gin framework. This template includes user registration, login, email verification, and password reset functionality, along with a simple frontend and email integration using Python.

## Features

- **User Registration**: Secure user registration with email and password validation.
- **User Login**: JWT-based authentication with cookies.
- **Email Verification**: Verification links sent to user emails to activate accounts.
- **Password Reset**: Allows users to reset passwords via email.
- **Frontend Templates**: Basic HTML templates for login, registration, and password reset.
- **Database Integration**: MySQL integration using GORM.
- **Email Service**: Email notifications handled via Python Flask.

## Requirements

- **Go**: >= 1.18
- **Python**: >= 3.8
- **MySQL**: >= 5.7
- **Node.js**: (Optional, for advanced frontend development)

**Note:** If you plan to deploy this project in a production environment, ensure the following security measures are implemented:

- **XSRF/CSRF Protection:** Prevent cross-site request forgery attacks by implementing robust token-based protection mechanisms.
- **Secure Cookie Handling:** Use attributes like `HttpOnly`, `Secure`, and `SameSite` to protect cookies from being accessed or transmitted in an insecure manner.
- **TLS/HTTPS:** Enable HTTPS to ensure encrypted communication between the server and clients.
- **Environment Variables:** Store sensitive information, such as database credentials, in environment variables instead of hardcoding them in the source code.

## Getting Started

### Clone the Repository

```bash
$ git clone https://github.com/LCGant/go-auth-boilerplate.git
$ cd go-auth-boilerplate
```

### Configure the Environment

1. Create a `.env` file in the project root with the following variables:

```env
DB_USER=root
DB_PASS=yourpassword
DB_HOST=127.0.0.1
DB_PORT=3306
DB_NAME=mydatabase
SECRET_KEY=your-secret-key
EMAIL_SERVER=http://localhost:5000/send-email
```

2. Update the `config/db.go` file to load configurations dynamically.

### Install Dependencies

```bash
$ go mod tidy
$ pip install flask
```

### Set Up the Database

1. Run the MySQL server.
2. Execute the SQL schema:

```bash
$ mysql -u root -p < schema.sql
```

### Start the Services

1. Run the Go backend:

```bash
$ go run main.go
```

2. Run the Python email server:

```bash
$ cd python_email
$ python email_sender.py
```

3. Access the application at [http://127.0.0.1:8080](http://127.0.0.1:8080).

## Folder Structure

```
├── config
│   └── db.go          # Database connection
├── controllers
│   ├── email_controller.go
│   ├── login_controller.go
│   ├── password_controller.go
│   └── register_controller.go
├── models
│   ├── login.go
│   ├── register.go
│   └── user.go
├── public
│   ├── css
│   │   └── style.css
│   ├── js
│   │   ├── login.js
│   │   ├── register.js
│   │   └── reset-password.js
│   └── pages
│       ├── login.html
│       ├── register.html
│       └── reset-password.html
├── python_email
│   └── email_sender.py # Python-based email server
├── services
│   ├── email_service.go
│   ├── login_service.go
│   ├── password_service.go
│   └── token_service.go
├── go.mod
├── go.sum
├── main.go             # Application entry point
├── schema.sql          # Database schema
└── README.md
```

## API Endpoints

### Authentication

| Method | Endpoint             | Description              |
|--------|----------------------|--------------------------|
| POST   | `/register`          | Register a new user      |
| POST   | `/login`             | Log in a user            |
| POST   | `/reset-password`    | Reset user password      |
| POST   | `/forgot-password`   | Request password reset   |
| GET    | `/verify-token`      | Validate reset token     |
| GET    | `/verify-email`      | Verify email token       |
| GET    | `/verify-email-token`| Activate email           |

## Contributing

1. Fork the repository.
2. Create a new branch:

```bash
$ git checkout -b feature/your-feature-name
```

3. Make your changes and commit them:

```bash
$ git commit -m "Add your message here"
```

4. Push your changes:

```bash
$ git push origin feature/your-feature-name
```

5. Submit a pull request.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Gin Framework](https://github.com/gin-gonic/gin)
- [GORM](https://gorm.io/)
- [Flask](https://flask.palletsprojects.com/)

---

Happy coding!

