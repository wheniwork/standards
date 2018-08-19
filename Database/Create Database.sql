-- The user_role type will help constrain the role in the DB to only the two roles.
-- It also allows adding more roles easily in the future.
DROP TYPE IF EXISTS user_role CASCADE;
CREATE TYPE user_role AS ENUM ('employee', 'manager');

-- The user type is used to return extra data when querying shifts.
-- This type will be encoded with to_json on the row level.
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
  id         SERIAL    NOT NULL PRIMARY KEY,
  name       TEXT      NOT NULL,
  email      TEXT      NULL,
  phone      TEXT      NULL,
  role       user_role NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT (localtimestamp),
  updated_at TIMESTAMP NOT NULL DEFAULT (localtimestamp),
  -- Make sure that both the phone and email are not blank.
  CHECK ((email IS NOT NULL AND character_length(email) > 0) OR (phone IS NOT NULL AND character_length(phone) > 0))
);
-- Demo data
INSERT INTO public.users (name, email, phone, role)
VALUES ('Elliot', 'elliot@elliot.com', null, 'employee'), --1
       ('Jimmy', 'jimmy@johns.com', null, 'employee'),    --2
       ('Jenny', null, '1-800-867-5309', 'manager'),      --3
       ('Henry', null, '1-800-123-4561', 'employee'); --4

DROP TABLE IF EXISTS public.shifts CASCADE;
CREATE TABLE public.shifts (
  id          SERIAL    NOT NULL PRIMARY KEY,
  manager_id  INT       NOT NULL REFERENCES public.users (id),
  employee_id INT       NULL REFERENCES public.users (id),
  break       FLOAT     NOT NULL DEFAULT 0,
  start_time  TIMESTAMP NOT NULL,
  end_time    TIMESTAMP NOT NULL,
  created_at  TIMESTAMP NOT NULL DEFAULT (localtimestamp),
  updated_at  TIMESTAMP NOT NULL DEFAULT (localtimestamp),
  -- Make sure that the shift starts before it ends (that sounds weird)
  CHECK (start_time < end_time),
  -- Make sure that the break is not longer than the actual shift.
  CHECK (break * INTERVAL '1 Hour' < (end_time - start_time))
);
INSERT INTO public.shifts (manager_id, employee_id, start_time, end_time)
VALUES
       (3, 1, 'Sun, Aug 19 18:00:00.000 2018', 'Mon, Aug 19 20:00:00.00 2018'),
       (3, 2, 'Sun, Aug 19 20:00:00.000 2018', 'Mon, Aug 19 22:00:00.00 2018'),
       (3, 3, 'Sun, Aug 19 22:00:00.000 2018', 'Mon, Aug 20 02:00:00.00 2018'),
       (3, 3, 'Sun, Aug 22 08:00:00.000 2018', 'Mon, Aug 22 14:00:00.00 2018'),
       (3, 1, localtimestamp - INTERVAL '1 Hours', localtimestamp),
       (3, null, localtimestamp - INTERVAL '1 Hours', localtimestamp),
       (3, 1, localtimestamp, localtimestamp + interval '1 hours'),
       (3, 2, localtimestamp - INTERVAL '1 Hours', localtimestamp),
       (3, 3, localtimestamp, localtimestamp + interval '1 hours');

-- This view will be used by the API, it makes it a bit easier to query and to make changes.
-- Also all of the RFC 2822 formatting is done in the view itself.
DROP VIEW IF EXISTS public.vw_users_api;
CREATE VIEW public.vw_users_api AS
  SELECT id, name, email, phone, role, to_char(created_at, 'Dy, Mon DD HH24:MI:SS.MS YYYY') AS created_at, to_char(updated_at, 'Dy, Mon DD HH24:MI:SS.MS YYYY') AS updated_at
  FROM public.users;


DROP VIEW IF EXISTS public.vw_shifts_api;
CREATE VIEW public.vw_shifts_api AS
  SELECT s.id,
         s.manager_id,
         -- The website will unmarshal this employee and the one below for the object.
         -- I feel like doing this helps a bit more when selecting data.
         -- It allows the client app to exclude user data entirely and I don't need to keep track of a lot more columns.
         to_json(row (manager.id,
                     manager.name,
                     manager.email,
                     manager.phone,
                     manager.role,
                     to_char(manager.created_at, 'Dy, Mon DD HH24:MI:SS.MS YYYY'),
                     to_char(manager.updated_at, 'Dy, Mon DD HH24:MI:SS.MS YYYY')) :: public.user) AS manager_user,
         s.employee_id,
         CASE
           WHEN s.employee_id IS NOT NULL THEN to_json(row (employee.id,
                                                           employee.name,
                                                           employee.email,
                                                           employee.phone,
                                                           employee.role,
                                                           to_char(employee.created_at,
                                                                   'Dy, Mon DD HH24:MI:SS.MS YYYY'),
                                                           to_char(employee.updated_at,
                                                                   'Dy, Mon DD HH24:MI:SS.MS YYYY')) :: public.user)
           ELSE NULL END                                                                           AS employee_user,
         s.break,
         to_char(s.start_time, 'Dy, Mon DD HH24:MI:SS.MS YYYY')                                    AS start_time,
         to_char(s.end_time, 'Dy, Mon DD HH24:MI:SS.MS YYYY')                                      AS end_time,
         to_char(s.created_at, 'Dy, Mon DD HH24:MI:SS.MS YYYY')                                    AS created_at,
         to_char(s.updated_at, 'Dy, Mon DD HH24:MI:SS.MS YYYY')                                    AS updated_at
  FROM public.shifts s
         INNER JOIN public.users manager ON manager.id = s.manager_id
         LEFT JOIN public.users employee ON employee.id = s.employee_id;

-- This view will return shifts that overlap with each other based on group_by_id and group_by_employee_id
DROP VIEW IF EXISTS public.vw_shifts_detailed_api;
CREATE VIEW public.vw_shifts_detailed_api AS
  SELECT s.id                                                                                      as group_by_id,
         s.employee_id                                                                             as group_by_employee_id,
         s2.id,
         s2.manager_id,
         to_json(row (manager.id,
                     manager.name,
                     manager.email,
                     manager.phone,
                     manager.role,
                     to_char(manager.created_at, 'Dy, Mon DD HH24:MI:SS.MS YYYY'),
                     to_char(manager.updated_at, 'Dy, Mon DD HH24:MI:SS.MS YYYY')) :: public.user) AS manager_user,
         s2.employee_id,
         CASE
           WHEN s2.employee_id IS NOT NULL THEN to_json(row (employee.id,
                                                            employee.name,
                                                            employee.email,
                                                            employee.phone,
                                                            employee.role,
                                                            to_char(employee.created_at,
                                                                    'Dy, Mon DD HH24:MI:SS.MS YYYY'),
                                                            to_char(employee.updated_at,
                                                                    'Dy, Mon DD HH24:MI:SS.MS YYYY')) :: public.user)
           ELSE NULL END                                                                           AS employee_user,
         s2.break,
         to_char(s2.start_time, 'Dy, Mon DD HH24:MI:SS.MS YYYY')                                   AS start_time,
         to_char(s2.end_time, 'Dy, Mon DD HH24:MI:SS.MS YYYY')                                     AS end_time,
         to_char(s2.created_at, 'Dy, Mon DD HH24:MI:SS.MS YYYY')                                   AS created_at,
         to_char(s2.updated_at, 'Dy, Mon DD HH24:MI:SS.MS YYYY')                                   AS updated_at
  FROM public.shifts s
         INNER JOIN public.shifts s2 ON s2.start_time < s.end_time AND s2.end_time > s.start_time
         INNER JOIN public.users manager ON manager.id = s2.manager_id
         LEFT JOIN public.users employee ON employee.id = s2.employee_id;
;


-- This view is insane.
-- For all of the shifts it looks at the expected hours to be worked.
-- If a shift were to go into another week, such as an overnight shift. Then the hours will be split into those two weeks.
-- If there is a break on a shift, it will be subtracted from the end of the shift. This way if the shift goes into the next week...
--    it will allocate hours from that next week first. This also makes it easier because if you take a one hour break on a shift...
--    that goes into the next week by one hour, it will not show on the next week. This isn't ideal, but without knowing clock in/out...
---   times for a shift (so a break could be determined) its hard to know when a break is.

-- Once a shift has been split, it also takes into account whether that shift is in the past or future, or is currently happening.
-- If a shift is currently happening, it will only indicate the hours worked until the current timestamp.
-- If a shift is in the future it will show 0 hours work, but will still show the hours scheduled.
-- If a shift is in the past it will assume that the start and end of the shift are essentially time clocks. So it'll calculate it as hours worked minus the break.
DROP VIEW IF EXISTS public.vw_shifts_summary_api;
CREATE VIEW public.vw_shifts_summary_api AS
  WITH shifts AS (SELECT employee_id,
                         -- I need to group shifts by week, this will generate a week number looking something like 201833 based of the start of a shift.
                         to_char(start_time, 'YYYYWW')                                                                                                          AS week,
                         -- This will take the week number and then parse it back into a timestamp on its own to show when the week starts and ends
                         to_char(to_date(to_char(start_time, 'YYYYWW'), 'YYYYWW'), 'Dy, Mon DD YYYY')                                             AS week_start,
                         to_char((to_date(to_char(start_time, 'YYYYWW'), 'YYYYWW') + INTERVAL '7 Days'), 'Dy, Mon DD YYYY')                       AS week_end,
                         start_time,
                         end_time,
                         id                                                                                                                                     AS shift_id,
                         -- The LEAST function will help determine whether or not the end of the shift is in the current week as the start.
                         (LEAST((to_date(to_char(start_time, 'YYYYWW'), 'YYYYWW') + INTERVAL '7 Days'), (end_time - (break * INTERVAL '1 Hour'))) - start_time) AS hours_scheduled,
                         -- This will calculate the hours worked if the shift is in progress or is in the past.
                         CASE
                           WHEN start_time < localtimestamp THEN GREATEST((LEAST((to_date(to_char(start_time, 'YYYYWW'), 'YYYYWW') + INTERVAL '7 Days'), (LEAST(end_time, localtimestamp) - (break * INTERVAL '1 Hour'))) - start_time), 0 * INTERVAL '1 Second')
                           ELSE 0 * INTERVAL '1 Second' END                                                                                                     AS hours_worked,
                         break                                                                                                                                  AS breaks,
                         -- The break offset helps calculate if the break needs to be deducted on a shift split into two weeks.
                         CASE
                           WHEN (end_time - (break * INTERVAL '1 Hour')) > (to_date(to_char(start_time, 'YYYYWW'), 'YYYYWW') + INTERVAL '7 Days') THEN 0
                           ELSE break
                             END                                                                                                                                AS break_offset,
                         -- Will return true of the shift goes into the next week.
                         (end_time - (break * INTERVAL '1 Hour')) > (to_date(to_char(start_time, 'YYYYWW'), 'YYYYWW') + INTERVAL '7 Days')                      AS has_overlap
                  FROM public.shifts
                  WHERE employee_id IS NOT NULL),
      -- This CTE table will calculate the second part of shifts that go into the following week.
      shifts_latter AS (SELECT s1.employee_id                                                                                                    AS employee_id,
                               to_char(s2.end_time, 'YYYYWW')                                                                                    AS week,
                               to_char(to_date(to_char(s2.end_time, 'YYYYWW'), 'YYYYWW'), 'Dy, Mon DD YYYY')                       AS week_start,
                               to_char((to_date(to_char(s2.end_time, 'YYYYWW'), 'YYYYWW') + INTERVAL '7 Days'), 'Dy, Mon DD YYYY') AS week_end,
                               s2.shift_id                                                                                                       AS shift_id,
                               CASE
                                 WHEN to_date(to_char(s2.end_time, 'YYYYWW'), 'YYYYWW') < localtimestamp THEN GREATEST(LEAST(s2.end_time, localtimestamp) - to_date(to_char(s2.end_time, 'YYYYWW'), 'YYYYWW') - (s2.breaks * INTERVAL '1 Hour'), 0 * INTERVAL '1 Second')
                                 ELSE 0 * INTERVAL '1 Second'
                                   END                                                                                                           AS hours_worked,
                               s2.end_time - to_date(to_char(s2.end_time, 'YYYYWW'), 'YYYYWW') - (s2.breaks * INTERVAL '1 Hour')                 AS hours_scheduled,
                               s2.breaks                                                                                                         AS breaks
                        FROM shifts s1
                               INNER JOIN shifts s2 ON s1.shift_id = s2.shift_id
                        WHERE s1.has_overlap = true),
      -- Put the first half and second half of shifts together for aggregating data.
      shifts_summarized AS (SELECT shifts.employee_id,
                                   shifts.week,
                                   shifts.week_start,
                                   shifts.week_end,
                                   shifts.shift_id,
                                   shifts.hours_scheduled,
                                   shifts.hours_worked,
                                   shifts.break_offset AS breaks
                            FROM shifts
                            UNION ALL
                            SELECT shifts_latter.employee_id,
                                   shifts_latter.week,
                                   shifts_latter.week_start,
                                   shifts_latter.week_end,
                                   shifts_latter.shift_id,
                                   shifts_latter.hours_scheduled,
                                   shifts_latter.hours_worked,
                                   shifts_latter.breaks
                            FROM shifts_latter)
  SELECT s.employee_id                                                                                                               AS employee_id,
         to_json(row (employee.id,
                     employee.name,
                     employee.email,
                     employee.phone,
                     employee.role,
                     to_char(employee.created_at, 'Dy, Mon DD HH24:MI:SS.MS YYYY'),
                     to_char(employee.updated_at, 'Dy, Mon DD HH24:MI:SS.MS YYYY')) :: public.user)                                  AS employee_user,
         s.week                                                                                                                      AS week,
         s.week_start                                                                                                                AS week_start,
         s.week_end                                                                                                                  AS week_end,
         COUNT(DISTINCT s.shift_id)                                                                                                  AS total_shifts,
         -- Will return the total hours in this week as a decimal
         round(CAST(EXTRACT(epoch FROM SUM(s.hours_scheduled)) / 3600 as numeric), 2)                                                AS total_scheduled_time,
         -- Will return the total hours in this week as a more human readable string for UI.
         to_char(date_part('epoch', SUM(s.hours_scheduled)) * INTERVAL '1 second', 'FMHH24 Hour(s) FMMI ') || 'Minute(s)'            AS total_scheduled_time_formatted,
         round(CAST(EXTRACT(epoch FROM SUM(s.hours_worked)) / 3600 as numeric), 2)                                                   AS total_worked_time,
         to_char(date_part('epoch', SUM(s.hours_worked)) * INTERVAL '1 second', 'FMHH24 Hour(s) FMMI ') || 'Minute(s)'               AS total_worked_time_formatted,
         round(CAST(EXTRACT(epoch FROM SUM(s.breaks) * INTERVAL '1 Hour') / 3600 as numeric), 2)                                     AS total_break_time,
         to_char(date_part('epoch', SUM(s.breaks) * INTERVAL '1 Hour') * INTERVAL '1 second', 'FMHH24 Hour(s) FMMI ') || 'Minute(s)' AS total_break_time_formatted
  FROM shifts_summarized s
         INNER JOIN public.users employee ON employee.id = s.employee_id
  GROUP BY s.employee_id,
           s.week,
           s.week_start,
           s.week_end,
           employee.id;