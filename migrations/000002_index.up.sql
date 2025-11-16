CREATE INDEX idx_users_team_id ON users(team_id);
CREATE INDEX idx_users_team_id_active ON users(team_id, is_active);

CREATE INDEX idx_pull_request_author_id ON pull_request(author_id);
CREATE INDEX idx_pull_request_status_id ON pull_request(status_id);
CREATE INDEX idx_pull_request_created_at ON pull_request(created_at DESC);

CREATE INDEX idx_assigned_pr_reviewer_id ON assigned_pr(reviewer_id);
CREATE INDEX idx_assigned_pr_pr_id ON assigned_pr(pr_id);