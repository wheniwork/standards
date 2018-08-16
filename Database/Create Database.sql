DROP TYPE IF EXISTS user_role CASCADE;
CREATE TYPE user_role AS ENUM ('employee', 'manager');

DROP TYPE IF EXISTS public.user;
CREATE TYPE public.user AS (
  id         INT,
  name       TEXT,
  email      TEXT,
  phone      TEXT,
  role       user_role,
  created_at TEXT,
  updated_at TEXT
);

DROP TABLE IF EXISTS public.users CASCADE;
CREATE TABLE public.users (
  id         SERIAL                   NOT NULL PRIMARY KEY,
  name       TEXT                     NOT NULL,
  email      TEXT                     NULL,
  phone      TEXT                     NULL,
  role       user_role                NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT TIMEZONE('CDT', NOW()),
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT TIMEZONE('CDT', NOW()),
  CHECK ((email IS NOT NULL AND character_length(email) > 0) OR (phone IS NOT NULL AND character_length(phone) > 0))
);
INSERT INTO public.users (name, email, phone, role)
VALUES ('Elliot', 'elliot@elliot.com', null, 'employee'), --1
       ('Jimmy', 'jimmy@johns.com', null, 'employee'),    --2
       ('Jenny', null, '1-800-867-5309', 'manager'),      --3
       ('Henry', null, '1-800-123-4561', 'employee'); --4

DROP TABLE IF EXISTS public.shifts CASCADE;
CREATE TABLE public.shifts (
  id          SERIAL                   NOT NULL PRIMARY KEY,
  manager_id  INT                      NOT NULL REFERENCES public.users (id),
  employee_id INT                      NULL REFERENCES public.users (id),
  break       FLOAT                    NOT NULL DEFAULT 0,
  start_time  TIMESTAMP WITH TIME ZONE NOT NULL,
  end_time    TIMESTAMP WITH TIME ZONE NOT NULL,
  created_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT TIMEZONE('CDT', NOW()),
  updated_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT TIMEZONE('CDT', NOW()),
  CHECK (start_time < end_time)
);
INSERT INTO public.shifts (manager_id, employee_id, start_time, end_time)
VALUES (3, 1, TIMEZONE('CDT', '2018-08-11 8:00AM'), TIMEZONE('CDT', '2018-08-11 2:00PM')),
       (3, 1, TIMEZONE('CDT', NOW()), TIMEZONE('CDT', NOW()) + INTERVAL '2 Hour'),
       (3, 2, TIMEZONE('CDT', NOW()) - INTERVAL '1 Hour', TIMEZONE('CDT', NOW()) + INTERVAL '1 Hour'),
       (3, 3, TIMEZONE('CDT', NOW()) + INTERVAL '1 Hour', TIMEZONE('CDT', NOW()) + INTERVAL '3 Hour');


DROP VIEW IF EXISTS public.vw_users_api;
CREATE VIEW public.vw_users_api AS
  SELECT id,
         name,
         email,
         phone,
         role,
         to_char(created_at, 'Dy, Mon DD HH24:MI:SS.MS OF00 YYYY') AS created_at,
         to_char(updated_at, 'Dy, Mon DD HH24:MI:SS.MS OF00 YYYY') AS updated_at
  FROM public.users;

DROP VIEW IF EXISTS public.vw_shifts_api;
CREATE VIEW public.vw_shifts_api AS
  SELECT s.id,
         s.manager_id,
         to_json(row (manager.id,
                     manager.name,
                     manager.email,
                     manager.phone,
                     manager.role,
                     to_char(manager.created_at, 'Dy, Mon DD HH24:MI:SS.MS OF00 YYYY'),
                     to_char(manager.updated_at, 'Dy, Mon DD HH24:MI:SS.MS OF00 YYYY')) :: public.user) AS manager_user,
         s.employee_id,
         CASE
           WHEN s.employee_id IS NOT NULL THEN to_json(row (employee.id,
                                                           employee.name,
                                                           employee.email,
                                                           employee.phone,
                                                           employee.role,
                                                           to_char(employee.created_at,
                                                                   'Dy, Mon DD HH24:MI:SS.MS OF00 YYYY'),
                                                           to_char(employee.updated_at,
                                                                   'Dy, Mon DD HH24:MI:SS.MS OF00 YYYY')) :: public.user)
           ELSE NULL END                                                                                AS employee_user,
         s.break,
         to_char(s.start_time, 'Dy, Mon DD HH24:MI:SS.MS OF00 YYYY')                                    AS start_time,
         to_char(s.end_time, 'Dy, Mon DD HH24:MI:SS.MS OF00 YYYY')                                      AS end_time,
         to_char(s.created_at, 'Dy, Mon DD HH24:MI:SS.MS OF00 YYYY')                                    AS created_at,
         to_char(s.updated_at, 'Dy, Mon DD HH24:MI:SS.MS OF00 YYYY')                                    AS updated_at
  FROM public.shifts s
         INNER JOIN public.users manager ON manager.id = s.manager_id
         LEFT JOIN public.users employee ON employee.id = s.employee_id;

DROP VIEW IF EXISTS public.vw_shifts_detailed_api;
CREATE VIEW public.vw_shifts_detailed_api AS
  SELECT s.id                                                                                           as group_by_id,
         s.employee_id                                                                                  as group_by_employee_id,
         s2.id,
         s2.manager_id,
         to_json(row (manager.id,
                     manager.name,
                     manager.email,
                     manager.phone,
                     manager.role,
                     to_char(manager.created_at, 'Dy, Mon DD HH24:MI:SS.MS OF00 YYYY'),
                     to_char(manager.updated_at, 'Dy, Mon DD HH24:MI:SS.MS OF00 YYYY')) :: public.user) AS manager_user,
         s2.employee_id,
         CASE
           WHEN s2.employee_id IS NOT NULL THEN to_json(row (employee.id,
                                                            employee.name,
                                                            employee.email,
                                                            employee.phone,
                                                            employee.role,
                                                            to_char(employee.created_at,
                                                                    'Dy, Mon DD HH24:MI:SS.MS OF00 YYYY'),
                                                            to_char(employee.updated_at,
                                                                    'Dy, Mon DD HH24:MI:SS.MS OF00 YYYY')) :: public.user)
           ELSE NULL END                                                                                AS employee_user,
         s2.break,
         to_char(s2.start_time, 'Dy, Mon DD HH24:MI:SS.MS OF00 YYYY')                                   AS start_time,
         to_char(s2.end_time, 'Dy, Mon DD HH24:MI:SS.MS OF00 YYYY')                                     AS end_time,
         to_char(s2.created_at, 'Dy, Mon DD HH24:MI:SS.MS OF00 YYYY')                                   AS created_at,
         to_char(s2.updated_at, 'Dy, Mon DD HH24:MI:SS.MS OF00 YYYY')                                   AS updated_at
  FROM public.shifts s
         INNER JOIN public.shifts s2 ON s2.start_time < s.end_time AND s2.end_time > s.start_time
         INNER JOIN public.users manager ON manager.id = s2.manager_id
         LEFT JOIN public.users employee ON employee.id = s2.employee_id;
;


DROP VIEW IF EXISTS public.vw_shifts_summary_api;
CREATE VIEW public.vw_shifts_summary_api AS
  SELECT id,
         manager_id,
         employee_id,
         break,
         to_char(start_time, 'Dy, Mon DD HH24:MI:SS.MS OF00 YYYY') AS start_time,
         to_char(end_time, 'Dy, Mon DD HH24:MI:SS.MS OF00 YYYY')   AS end_time,
         to_char(created_at, 'Dy, Mon DD HH24:MI:SS.MS OF00 YYYY') AS created_at,
         to_char(updated_at, 'Dy, Mon DD HH24:MI:SS.MS OF00 YYYY') AS updated_at
  FROM public.shifts;