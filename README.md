Syllabea
Syllabea is a web-based system for managing course syllabi. It is built using the Go programming language with HTML templates and HTMX for dynamic web features. The system helps lecturers and managers create, edit, and view syllabi.

Features
Lecturers can create, edit, and save their syllabi.

Managers can see and manage all syllabi in the system.

The system keeps track of progress through a sidebar that shows how complete the syllabus is.

Uses HTMX for updating parts of the page without reloading the whole page.

Uses Go templates to create simple and fast web pages.

Technologies Used
Go (backend)

HTMX (for dynamic interactions)

HTML and CSS (for the frontend)

MySQL (for the database)

Roles
Lecturer: Can only see and edit their own syllabi.

Manager: Can see and manage all syllabi.

How to Run
Install Go on your system.

Clone this repository.

Set up the MySQL database using the provided SQL script (if available).

Configure the database connection in the code (usually in a config file).

Run the Go server:

go
Copy
Edit
go run main.go
Open your browser and go to http://localhost:8080

Future Plans
Add support for more user roles.

Allow exporting syllabi to PDF.

Improve the design and user experience.
