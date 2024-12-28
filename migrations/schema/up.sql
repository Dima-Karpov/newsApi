CREATE TABLE news_list (
  id            SERIAL PRIMARY KEY,                 -- автоматически увеличиваемый идентификатор
  title         VARCHAR(255) NOT NULL,              -- строка для заголовка
  description   TEXT NOT NULL,                      -- строка для описания
  published_at  TIMESTAMP NOT NULL,                 -- время публикации
  link          VARCHAR(255) NOT NULL UNIQUE        -- ссылка на источник
);


ALTER TABLE news_list
DROP COLUMN id;

ALTER TABLE news_list
ADD COLUMN id UUID PRIMARY KEY DEFAULT gen_random_uuid();