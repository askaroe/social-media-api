CREATE TABLE posts (
    id bigserial PRIMARY KEY,
    createdAt timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updatedAt timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    image text,
    caption text,
    userId bigserial REFERENCES users(id)
);

CREATE TABLE comments (
    id bigserial PRIMARY KEY,
    createdAt timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updatedAt timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    message text,
    userId bigserial REFERENCES users(id),
    postId bigserial REFERENCES posts(id)
);