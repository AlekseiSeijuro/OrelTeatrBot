-- +goose Up
CREATE TABLE events (
                        id         bigint      PRIMARY KEY,           -- ID из ссылки /e390
                        title      text        NOT NULL,
                        updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE sessions (
                          id            bigint      PRIMARY KEY,        -- ID из ссылки /s3164
                          event_id      bigint      NOT NULL REFERENCES events (id),
                          date_text     text        NOT NULL,           -- дата как на сайте, без разбора
                          actual_now    boolean     NOT NULL,           -- билеты есть прямо сейчас
                          avail_epoch   int         NOT NULL DEFAULT 0, -- +1 при каждом переходе false -> true
                          first_seen_at timestamptz NOT NULL DEFAULT now(),
                          last_seen_at  timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE users (
                       chat_id    bigint      PRIMARY KEY,           -- ID чата в Telegram
                       username   text,
                       is_active  boolean     NOT NULL DEFAULT true,
                       created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE notifications (
                               user_id     bigint      NOT NULL REFERENCES users (chat_id),
                               session_id  bigint      NOT NULL REFERENCES sessions (id),
                               avail_epoch int         NOT NULL,
                               status      text        NOT NULL DEFAULT 'pending'
                                   CHECK (status IN ('pending', 'sent', 'failed')),
                               message_id  bigint,
                               attempts    int         NOT NULL DEFAULT 0,
                               sent_at     timestamptz,
                               PRIMARY KEY (user_id, session_id, avail_epoch)
);

-- Supabase отдаёт таблицы схемы public через публичный REST API.
-- RLS без политик закрывает их для этого API, а бот подключается
-- напрямую ролью postgres и RLS не замечает.
ALTER TABLE events        ENABLE ROW LEVEL SECURITY;
ALTER TABLE sessions      ENABLE ROW LEVEL SECURITY;
ALTER TABLE users         ENABLE ROW LEVEL SECURITY;
ALTER TABLE notifications ENABLE ROW LEVEL SECURITY;

-- +goose Down
DROP TABLE notifications;
DROP TABLE users;
DROP TABLE sessions;
DROP TABLE events;