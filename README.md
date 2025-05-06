Syllabea
--------

Syllabea is a web-based system for managing course syllabi. It is built using the Go programming language with HTML templates and HTMX for dynamic web features. The system helps lecturers and managers create, edit, and view syllabi.

Features
--------

- Lecturers can create, edit, and save their syllabi.
- Managers can see and manage all syllabi in the system.
- The system keeps track of progress through a sidebar that shows how complete the syllabus is.
- Uses HTMX for updating parts of the page without reloading the whole page.
- Uses Go templates to create simple and fast web pages.

Technologies Used
-----------------

- Go (backend)
- HTMX (for dynamic interactions)
- HTML and CSS (for the frontend)
- MySQL (for the database)

Roles
-----

- Lecturer: Can only see and edit their own syllabi.
- Manager: Can see and manage all syllabi.

How to Run
----------

1. Install Go on your system.
2. Clone this repository.
3. Set up the MySQL database using the provided SQL script (if available).
4. Configure the database connection in the code (usually in a config file).
5. Run the Go server:
   go run main.go
6. Open your browser and go to http://localhost:8080

Future Plans
------------

- Add support for more user roles.
- Allow exporting syllabi to PDF.
- Improve the design and user experience.
