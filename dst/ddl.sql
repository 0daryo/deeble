CREATE TABLE Users (
  Id STRING(MAX) NOT NULL,
  FirstName STRING(MAX),
  LastName STRING(MAX),
  Email STRING(MAX),
  CreatedAt TIMESTAMP NOT NULL,
  UpdatedAt TIMESTAMP NOT NULL OPTIONS (
    allow_commit_timestamp = true
  ),
) PRIMARY KEY(Id);
