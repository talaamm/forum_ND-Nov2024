# **Forum Project**

## **Overview**
The **Forum Project** is a web application that enables users to communicate through posts and comments. This project incorporates a category system, user authentication, like/dislike functionality, and filtering mechanisms. It is built with Go, uses SQLite for the database.  

Additionally, the project includes an image upload feature that allows users to create posts containing images alongside text.

---

## **Features**
### **Core Features**
- **User Authentication**
  - Registration with email, username, and password.
  - Login system with session cookies for user authentication.
  - Cookie expiration to manage user sessions.
  - Password encryption (Bonus).
  
- **Posts and Comments**
  - Registered users can create posts and comments.
  - Posts can be associated with categories.

- **Likes and Dislikes**
  - Registered users can like or dislike posts.
  - Total likes/dislikes are visible to all users.

- **Filtering**
  - Posts can be filtered by:
    - Categories (subforums).
    - Created posts (for logged-in users).

### **Extra Features**
- **Image Upload**
  - Registered users can upload images in their posts.
  - Supported formats: JPEG, PNG, GIF.
  - Maximum image size: 20 MB.
  - Error handling for oversized images with proper feedback.

---

## **Technologies Used**
- **Backend:** Go
- **Database:** SQLite
- **Frontend:** HTML, CSS (No frameworks)

---

## **Getting Started**
### **Prerequisites**
- [Go](https://golang.org/) installed on your machine.
- [Docker](https://www.docker.com/) installed.
- Basic understanding of Go and Docker.

### **Setup Instructions**
1. **Clone the Repository**:
   ```bash
   git clone *repo link*
   cd forum
   ```

2. **Database Initialization**:
   - Ensure SQLite is set up.

3. **Run Locally**:
   - Start the server:
     ```bash
     go run .
     ```
   - Open your browser and navigate to `http://localhost:8000`.

---

## **Usage**
- Register an account to create posts or comments.
- Use the navigation bar to filter posts by categories, or your own posts.
- Upload images when creating posts (optional).

---

## **Features to Implement**
- Improve error handling for invalid user inputs.
- Add automated unit tests for backend functionality.
- Expand filtering options based on user feedback.

---

## **Team**
- me: **Tala Amm**
- **Moaz Razem**
- **Amro Khweis** 
- **Noor Halabi**

---

## **License**
This project is for educational purposes as part of the School curriculum and is not licensed for commercial use.
