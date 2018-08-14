DROP TYPE IF EXISTS user_role CASCADE;
CREATE TYPE user_role AS ENUM ('employee', 'manager');

DROP TABLE IF EXISTS public.users CASCADE;
CREATE TABLE public.users (
  id         SERIAL    NOT NULL PRIMARY KEY,
  name       TEXT      NOT NULL,
  email      TEXT      NULL,
  phone      TEXT      NULL,
  role       user_role NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT TIMEZONE('utc', NOW()),
  updated_at TIMESTAMP NOT NULL DEFAULT TIMEZONE('utc', NOW()),
  CHECK ((email IS NOT NULL AND character_length(email) > 0) OR (phone IS NOT NULL AND character_length(phone) > 0))
);
INSERT INTO public.users (name, email, phone, role)
VALUES ('Elliot', 'elliot@elliot.com', null, 'employee'), --1
       ('Jimmy', 'jimmy@johns.com', null, 'employee'),    --2
       ('Jenny', null, '1-800-867-5309', 'manager'),      --3
       ('Henry', null, '1-800-123-4561', 'employee'); --4

DROP TABLE IF EXISTS public.shifts;
CREATE TABLE public.shifts (
  id          SERIAL    NOT NULL PRIMARY KEY,
  manager_id  INT       NOT NULL REFERENCES public.users (id),
  employee_id INT       NOT NULL REFERENCES public.users (id),
  break       FLOAT     NOT NULL DEFAULT 0,
  start_time  TIMESTAMP NOT NULL,
  end_time    TIMESTAMP NOT NULL,
  created_at  TIMESTAMP NOT NULL DEFAULT TIMEZONE('utc', NOW()),
  updated_at  TIMESTAMP NOT NULL DEFAULT TIMEZONE('utc', NOW())
);

