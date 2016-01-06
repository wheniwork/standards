# REST Scheduler API

While we will accept a framework of your choice, and a language of your choice, we _highly_ encourage using **PHP*** and one of the following frameworks: [Equip](https://github.com/equip/framework), [Radar](https://github.com/radarphp/Radar.Project), or [Proton](https://github.com/alexbilbie/proton).

If you choose to use another language, please be prepared to show fluency in multiple languages (including at least one [C-family language](https://en.wikipedia.org/wiki/List_of_C-family_programming_languages)), and knowledge of multiple design principles ([DRY](https://en.wikipedia.org/wiki/Don%27t_repeat_yourself), [Open/Closed](https://en.wikipedia.org/wiki/Open/closed_principle), etc) and patterns ([Factory](https://en.wikipedia.org/wiki/Factory_method_pattern), [DataMapper](https://en.wikipedia.org/wiki/Data_mapper_pattern), etc).

<sub><sup>* Our core application is written almost entirely in PHP and JS.</sup></sub>

## Requirements

The API must follow REST specification:

- POST should be used to create
- GET should be used to read
- PUT should be used to update (and optionally to create)
- DELETE should be used to delete

Additional methods can be used for expanded functionality.

The API should include the following roles:

- employee (read)
- manager (write)

The `employee` will have much more limited access than a `manager`. The specifics of what each role should be able to do is listed below in [User Stories](#user-stories).

## Data Types

All data structures use the following types:

| type   | description |
| ------ | ----------- |
| int    | a integer number |
| float  | a floating point number |
| string | a string |
| bool   | a boolean |
| id     | a unique identifier |
| fk     | a reference to another id |
| date   | an RFC 2822 formatted date string |

## Data Structures

### User

| field       | type |
| ----------- | ---- |
| id          | id |
| name        | string |
| role        | string |
| email       | string |
| phone       | string |
| created_at  | date |
| updated_at  | date |

The `role` must be either `employee` or `manager`. At least one of `phone` or
`email` must be defined.

### Shift

| field       | type |
| ----------- | ---- |
| id          | id |
| manager_id  | fk |
| employee_id | fk |
| break       | float |
| start_time  | date |
| end_time    | date |
| created_at  | date |
| updated_at  | date |

Both `start_time` and `end_time` are required. Unless defined, the `manager_id`
should always default to the manager that created the shift. Any shift without
an `employee_id` will be visible to all employees.

## User stories

**Please note that this not intended to be a CRUD application.** Only the functionality described by the user stories should be exposed via the API.

- [ ] As an employee, I want to know when I am working, by being able to see all of the shifts assigned to me.
- [ ] As an employee, I want to know who I am working with, by being able to see the employees that are working during the same time period as me.
- [ ] As an employee, I want to know how much I worked, by being able to get a summary of hours worked for each week.
- [ ] As an employee, I want to be able to contact my managers, by seeing manager contact information for my shifts.

- [ ] As a manager, I want to schedule my employees, by creating shifts for any employee.
- [ ] As a manager, I want to see the schedule, by listing shifts within a specific time period.
- [ ] As a manager, I want to be able to change a shift, by updating the time details.
- [ ] As a manager, I want to be able to assign a shift, by changing the employee that will work a shift.
- [ ] As a manager, I want to contact an employee, by seeing employee details.
