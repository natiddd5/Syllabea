-- Create the users table
CREATE TABLE IF NOT EXISTS users (
                                     id INT AUTO_INCREMENT PRIMARY KEY,
                                     name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    role ENUM('Instructor', 'Manager') NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

-- Create the departments table
CREATE TABLE IF NOT EXISTS departments (
                                           id INT AUTO_INCREMENT PRIMARY KEY,
                                           name VARCHAR(255) UNIQUE NOT NULL
    );

-- Create the courses table
CREATE TABLE IF NOT EXISTS courses (
                                       id INT AUTO_INCREMENT PRIMARY KEY,
                                       name VARCHAR(255) NOT NULL,
    department_id INT NOT NULL,
    FOREIGN KEY (department_id) REFERENCES departments(id) ON DELETE CASCADE
    );

-- Create the syllabi table
CREATE TABLE IF NOT EXISTS syllabi (
                                       id INT AUTO_INCREMENT PRIMARY KEY,
                                       course_id INT NOT NULL,
                                       lecturer_id INT NOT NULL,
                                       status ENUM('Draft', 'Deleted', 'In Review', 'Approved', 'UnsavedDraft') NOT NULL,
    submission_date DATE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    data JSON NOT NULL,
    FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE,
    FOREIGN KEY (lecturer_id) REFERENCES users(id) ON DELETE CASCADE
    );

CREATE TABLE IF NOT EXISTS comments (
                                        id INT AUTO_INCREMENT PRIMARY KEY,
                                        syllabus_id INT NOT NULL,
                                        user_id INT NOT NULL,
                                        content TEXT NOT NULL,
                                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                                        FOREIGN KEY (syllabus_id) REFERENCES syllabi(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
    );



-- Insert sample departments (Hebrew names)
INSERT INTO departments (name)
VALUES
    ('הנדסת תוכנה'),
    ('הנדסת חשמל'),
    ('הנדסת איכות');

-- Insert a sample course (assuming it belongs to department id 1, "הנדסת תוכנה")
INSERT INTO courses (name, department_id)
VALUES ('מערכות מבוזרות', 1);

-- Insert a lecturer into the users table
INSERT INTO users (name, email, role)
VALUES ('מייקל ג''יי מיי', 'michael@example.com', 'Instructor');

-- Insert a syllabus for the course "מערכות מבוזרות" with a JSON null in the data column
INSERT INTO syllabi (course_id, lecturer_id, status, submission_date, data)
VALUES (1, 1, 'Draft', CURDATE(), CAST('null' AS JSON));

