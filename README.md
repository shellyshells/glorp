# Glorp - Reddit-inspired Forum

A modern, responsive forum application built with Go, featuring all the essential discussion platform functionalities.

## 🚀 Features

### Mandatory Features (FT-1 to FT-12)

- **FT-1: User Registration** - Secure user registration with validation
- **FT-2: User Login** - JWT-based authentication system
- **FT-3: Thread Creation** - Create discussion threads with tags/categories
- **FT-4: Thread Viewing** - Browse and read discussion threads
- **FT-5: Message Posting** - Post messages in discussion threads
- **FT-6: Like/Dislike System** - Vote on messages (like/dislike)
- **FT-7: Content Management** - Edit/delete own threads and messages
- **FT-8: Message Sorting** - Sort by date or popularity
- **FT-9: Pagination** - Configurable pagination (10, 20, 30, all)
- **FT-10: Tag Filtering** - Filter threads by tags/categories
- **FT-11: Search Functionality** - Search threads by title/tags
- **FT-12: Admin Dashboard** - Complete admin management panel

### Technical Implementation

- **Architecture**: MVC (Model-View-Controller)
- **Backend**: Go (Golang) with Gorilla Mux
- **Database**: SQLite with proper schema design
- **Authentication**: JWT tokens with secure cookies
- **Security**: SHA-512 password hashing
- **Frontend**: HTML templates with JavaScript interactions
- **Styling**: Custom CSS with responsive design

## 📁 Project Structure

```
glorp/
├── README.md
├── go.mod
├── go.sum
├── main.go
├── config/
│   └── database.go
├── models/
│   ├── user.go
│   ├── thread.go
│   ├── message.go
│   ├── vote.go
│   └── tag.go
├── controllers/
│   ├── auth_controller.go
│   ├── thread_controller.go
│   ├── message_controller.go
│   ├── admin_controller.go
│   └── search_controller.go
├── middleware/
│   ├── auth.go
│   └── admin.go
├── utils/
│   ├── hash.go
│   ├── jwt.go
│   └── pagination.go
├── views/
│   ├── layouts/
│   │   └── main.html
│   ├── auth/
│   │   ├── login.html
│   │   └── register.html
│   ├── threads/
│   │   ├── index.html
│   │   ├── show.html
│   │   ├── create.html
│   │   └── edit.html
│   └── admin/
│       └── dashboard.html
├── static/
│   ├── css/
│   │   └── style.css
│   └── js/
│       └── app.js
└── forum.db (created automatically)
```

## 🛠️ Installation & Setup

### Prerequisites

- Go 1.19 or higher
- Git

### Installation Steps

1. **Clone the repository:**
   ```bash
   git clone <your-repository-url>
   cd glorp
   ```

2. **Install dependencies:**
   ```bash
   go mod init glorp
   go get github.com/golang-jwt/jwt/v4
   go get github.com/gorilla/mux
   go get github.com/mattn/go-sqlite3
   go mod tidy
   ```

3. **Run the application:**
   ```bash
   go run main.go
   ```

4. **Access the forum:**
   - Open your browser and navigate to `http://localhost:8080`
   - Default admin credentials: 
     - Username: `admin`
     - Password: `AdminPassword123!`

## 🌐 API Routes

### View Routes (HTML Pages)
- `GET /` - Home page with thread list
- `GET /register` - User registration page
- `GET /login` - User login page
- `GET /threads/{id}` - View specific thread
- `GET /threads/create` - Create thread form (authenticated)
- `GET /threads/{id}/edit` - Edit thread form (authenticated)
- `GET /admin/dashboard` - Admin dashboard (admin only)

### API Routes (JSON Data)
- `POST /api/register` - User registration
- `POST /api/login` - User login
- `POST /api/logout` - User logout
- `GET /api/threads` - Get threads with filters/pagination
- `POST /api/threads` - Create new thread (authenticated)
- `PUT /api/threads/{id}` - Update thread (authenticated)
- `DELETE /api/threads/{id}` - Delete thread (authenticated)
- `POST /api/threads/{id}/messages` - Post message (authenticated)
- `DELETE /api/messages/{id}` - Delete message (authenticated)
- `POST /api/messages/{id}/vote` - Vote on message (authenticated)
- `GET /api/search` - Search threads
- `POST /api/admin/ban/{id}` - Ban/unban user (admin)
- `PUT /api/admin/threads/{id}/status` - Update thread status (admin)

## 👥 Team Composition

- **Developer 1**: [Your Name]
- **Developer 2**: [Partner Name]

## 📊 Project Management Summary

### Project Decomposition
The project was broken down into the following phases:
1. **Setup & Infrastructure** - Database design, project structure
2. **Authentication System** - User registration, login, JWT implementation
3. **Core Forum Features** - Threads, messages, voting system
4. **User Interface** - HTML templates, CSS styling, JavaScript interactions
5. **Admin Features** - Dashboard, moderation tools
6. **Testing & Polish** - Bug fixes, responsive design, final touches

### Task Distribution
- **Backend Development**: Database models, API controllers, authentication
- **Frontend Development**: HTML templates, CSS styling, JavaScript functionality
- **Testing & Integration**: Feature testing, bug fixes, deployment preparation

### Time Management
- **Week 1**: Project setup, database design, authentication system
- **Week 2**: Core forum features, thread and message functionality
- **Week 3**: User interface development, styling, responsiveness
- **Week 4**: Admin features, search functionality, final testing

### Documentation Strategy
- Code documentation through comments and README
- API documentation through route listings
- User guide through intuitive interface design
- Technical documentation in this README file

## 🔧 Key Features Explained

### Security Features
- **Password Security**: SHA-512 hashing with proper validation
- **Authentication**: JWT tokens with HTTP-only cookies
- **Input Validation**: Comprehensive server-side validation
- **Authorization**: Role-based access control (user/admin)

### User Experience
- **Responsive Design**: Mobile-friendly interface
- **Intuitive Navigation**: Clear menu structure and breadcrumbs
- **Real-time Interactions**: AJAX-powered voting and form submissions
- **Search & Filter**: Multiple ways to find relevant content

### Admin Features
- **User Management**: Ban/unban users
- **Content Moderation**: Delete threads/messages, change thread status
- **Statistics Dashboard**: Overview of forum activity
- **Bulk Operations**: Manage multiple items efficiently

## 🎨 Design Philosophy

The forum follows a Reddit-inspired design with:
- Clean, modern interface
- Intuitive voting system
- Hierarchical content organization
- Responsive mobile design
- Accessible color scheme and typography

## 🚀 Future Enhancements

Potential bonus features that could be added:
- **FTB-1**: Image uploads in messages
- **FTB-2**: User profiles with stats and bio
- **FT3**: Friend system with private threads
- Real-time chat functionality
- Email notifications
- Advanced search with filters
- Thread bookmarking
- User reputation system

## 📞 Support

For questions or issues, please contact the development team or check the codebase documentation.

---

**Built with ❤️ using Go, HTML, CSS, and JavaScript**