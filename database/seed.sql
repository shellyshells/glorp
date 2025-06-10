-- Glorp Database Seed Data
-- This file contains sample data for testing and demonstration

-- Insert admin user (password: AdminPassword123!)
INSERT INTO users (username, email, password_hash, role) VALUES 
('admin', 'admin@forum.com', 'ba3253876aed6bc22d4a6ff53d8406c6ad864195ed144ab5c87621b6c233b548baeae6956df346ec8c17f5ea10f35ee3cbc514797ed7ddd3145464e2a0bab413', 'admin');

-- Insert sample users (password: TestPassword123!)
INSERT INTO users (username, email, password_hash, role) VALUES 
('john_doe', 'john@example.com', 'e7cf3ef4f17c3999a94f2c6f612e8a888e5b1026878e4e19398b23bd38ec221a54474e2f7be8648fcce7e2e2ba1ad5d71c9f05b8fd7bd18c8fdf78e3d59e0e7a', 'user'),
('jane_smith', 'jane@example.com', 'e7cf3ef4f17c3999a94f2c6f612e8a888e5b1026878e4e19398b23bd38ec221a54474e2f7be8648fcce7e2e2ba1ad5d71c9f05b8fd7bd18c8fdf78e3d59e0e7a', 'user'),
('dev_mike', 'mike@example.com', 'e7cf3ef4f17c3999a94f2c6f612e8a888e5b1026878e4e19398b23bd38ec221a54474e2f7be8648fcce7e2e2ba1ad5d71c9f05b8fd7bd18c8fdf78e3d59e0e7a', 'user');

-- Insert default tags
INSERT INTO tags (name) VALUES 
('General'),
('Technology'),
('Programming'),
('Discussion'),
('Help'),
('News'),
('JavaScript'),
('Python'),
('Go'),
('Web Development'),
('Mobile'),
('Gaming'),
('Science'),
('Politics'),
('Sports');

-- Insert sample threads
INSERT INTO threads (title, description, author_id, status) VALUES 
('Welcome to Glorp!', 'This is our first thread. Welcome to the community! Feel free to introduce yourself and share what you''re interested in discussing.', 1, 'open'),
('Best practices for Go development', 'Let''s discuss best practices when developing applications in Go. What patterns do you follow? What libraries do you recommend?', 2, 'open'),
('JavaScript frameworks in 2024', 'What are your thoughts on the current state of JavaScript frameworks? React, Vue, Angular, or something else?', 3, 'open'),
('Career advice for new developers', 'Starting out in tech can be overwhelming. Share your experiences and advice for newcomers to the field.', 4, 'open'),
('Weekly tech news discussion', 'Share and discuss the latest news in technology. What caught your attention this week?', 1, 'open');

-- Link threads to tags
INSERT INTO thread_tags (thread_id, tag_id) VALUES 
(1, 1), (1, 4), -- Welcome thread: General, Discussion
(2, 2), (2, 3), (2, 9), -- Go best practices: Technology, Programming, Go
(3, 2), (3, 3), (3, 7), (3, 10), -- JS frameworks: Technology, Programming, JavaScript, Web Development
(4, 3), (4, 4), -- Career advice: Programming, Discussion
(5, 2), (5, 6); -- Tech news: Technology, News

-- Insert sample messages
INSERT INTO messages (thread_id, author_id, content) VALUES 
(1, 2, 'Thanks for creating this forum! I''m excited to be part of this community.'),
(1, 3, 'Great initiative! Looking forward to some interesting discussions.'),
(1, 4, 'Hello everyone! I''m a software developer with 5 years of experience. Happy to help newcomers!'),

(2, 1, 'One thing I always recommend is to follow the standard Go project layout. It makes code organization much cleaner.'),
(2, 3, 'I agree! Also, always use gofmt and golint. They''re essential for maintaining code quality.'),
(2, 4, 'Don''t forget about proper error handling. It''s one of Go''s strengths when used correctly.'),

(3, 2, 'I''ve been using React for the past few years and I''m quite happy with it. The ecosystem is mature and there''s great tooling.'),
(3, 1, 'Vue.js has been my go-to choice. It''s more approachable for beginners and the documentation is excellent.'),
(3, 4, 'Angular is still relevant for large enterprise applications. The learning curve is steep but it''s very powerful.'),

(4, 1, 'My advice: build projects! Theory is important but hands-on experience is what really matters.'),
(4, 2, 'Network with other developers. Join communities like this one, attend meetups, contribute to open source.'),
(4, 3, 'Don''t be afraid to ask questions. Every senior developer was once a beginner too.'),

(5, 4, 'Did anyone see the latest developments in AI? The pace of innovation is incredible.'),
(5, 2, 'Yes! And the new web standards coming out are also very exciting. WebAssembly is getting more adoption.'),
(5, 3, 'The shift towards edge computing is something I''m watching closely. It''s going to change how we build applications.');

-- Insert sample votes (likes and dislikes)
INSERT INTO votes (message_id, user_id, vote_type) VALUES 
(1, 1, 1), (1, 3, 1), (1, 4, 1), -- Message 1: 3 likes
(2, 1, 1), (2, 2, 1), -- Message 2: 2 likes
(3, 1, 1), (3, 2, 1), (3, 3, 1), -- Message 3: 3 likes
(4, 2, 1), (4, 3, 1), (4, 4, 1), -- Message 4: 3 likes
(5, 1, 1), (5, 2, 1), -- Message 5: 2 likes
(6, 1, 1), (6, 3, 1), -- Message 6: 2 likes
(7, 1, 1), (7, 3, 1), (7, 4, 1), -- Message 7: 3 likes
(8, 2, 1), (8, 4, 1), -- Message 8: 2 likes
(9, 1, 1), (9, 2, 1), -- Message 9: 2 likes
(10, 2, 1), (10, 3, 1), (10, 4, 1), -- Message 10: 3 likes
(11, 3, 1), (11, 4, 1), -- Message 11: 2 likes
(12, 1, 1), (12, 4, 1), -- Message 12: 2 likes
(13, 1, 1), (13, 2, 1), (13, 3, 1), -- Message 13: 3 likes
(14, 1, 1), (14, 3, 1), -- Message 14: 2 likes
(15, 1, 1), (15, 4, 1); -- Message 15: 2 likes