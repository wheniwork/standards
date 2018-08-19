[![Build Status](https://travis-ci.com/ECourant/standards.svg?branch=master)](https://travis-ci.com/ECourant/standards)

# Elliot's Sample Project For When I Work

This project requires Golang 1.10 and PostgreSQL 10. 

The database change scripts currently don't target a specific database instance, and I've just run them on the default `postgres` database instance. If you would like to change what database they are in, you will also need to update the connection string in `config.json` in the Site directory and `database_test_config.json` in the Database directory.

I've also kept track of what needs to be done for this project and some brief notes on how to go about doing that in Issues and Projects. 

If I have time I will also try to create a UI for the REST API using Framework7. I'm not sure how much time I will have to do that though and it will be the last thing I do. 

I've setup Travis CI to make sure that everything works on a system that isn't my own. 

#### Build
```
go get -u github.com/ECourant/standards
cd {path to github.com/ECourant/standards}
go build
```

You will need to tweak the settings in `config.json` in the root directory of the project for your computer. 
You will also want to run `Create Database.sql` on your PostgreSQL instance before running the website. This file can be found in the Database directory.




# REST API Documentation

#### As an employee, I want to know when I am working, by being able to see all of the shifts assigned to me.
One shift shows up in the response that does not have an employee, this is because the employee_id for that record is null, and will show up for all users.
<details><summary>HTTP Request</summary>
<p>

```http request
GET /api/shifts/mine?current_user_id=1

{
    "success": true,
    "results": [
        {
            "id": 114,
            "manager_id": 3,
            "manager_user": {
                "id": 3,
                "name": "Jenny",
                "phone": "1-800-867-5309",
                "role": "manager",
                "created_at": "Sun, Aug 19 11:02:20.537 2018",
                "updated_at": "Sun, Aug 19 11:02:20.537 2018"
            },
            "employee_id": 1,
            "employee_user": {
                "id": 1,
                "name": "Elliot",
                "email": "elliot@elliot.com",
                "role": "employee",
                "created_at": "Sun, Aug 19 11:02:20.537 2018",
                "updated_at": "Sun, Aug 19 11:02:20.537 2018"
            },
            "break": 0,
            "start_time": "Thu, Aug 02 19:31:46.631 2018",
            "end_time": "Thu, Aug 02 20:31:46.631 2018",
            "created_at": "Sun, Aug 19 11:02:25.736 2018",
            "updated_at": "Sun, Aug 19 11:02:25.736 2018"
        },
        {
            "id": 113,
            "manager_id": 3,
            "manager_user": {
                "id": 3,
                "name": "Jenny",
                "phone": "1-800-867-5309",
                "role": "manager",
                "created_at": "Sun, Aug 19 11:02:20.537 2018",
                "updated_at": "Sun, Aug 19 11:02:20.537 2018"
            },
            "break": 0,
            "start_time": "Tue, Aug 07 19:31:46.631 2018",
            "end_time": "Tue, Aug 07 20:31:46.631 2018",
            "created_at": "Sun, Aug 19 11:02:25.697 2018",
            "updated_at": "Sun, Aug 19 11:02:25.697 2018"
        }
    ]
}
```
</p>
</details>

#### As an employee, I want to know who I am working with, by being able to see the employees that are working during the same time period as me.
This will show any shifts that overlap with the shift ID you specified.
<details><summary>HTTP Request</summary>
<p>

```http request
GET /api/shifts/overlapping/1?current_user_id=1

{
    "success": true,
    "results": [
        {
            "id": 2,
            "manager_id": 3,
            "manager_user": {
                "id": 3,
                "name": "Jenny",
                "phone": "1-800-867-5309",
                "role": "manager",
                "created_at": "Sun, Aug 19 11:02:20.537 2018",
                "updated_at": "Sun, Aug 19 11:02:20.537 2018"
            },
            "employee_id": 2,
            "employee_user": {
                "id": 2,
                "name": "Jimmy",
                "email": "jimmy@johns.com",
                "role": "employee",
                "created_at": "Sun, Aug 19 11:02:20.537 2018",
                "updated_at": "Sun, Aug 19 11:02:20.537 2018"
            },
            "break": 0,
            "start_time": "Sun, Aug 19 20:00:00.000 2018",
            "end_time": "Sun, Aug 19 22:00:00.000 2018",
            "created_at": "Sun, Aug 19 11:02:20.537 2018",
            "updated_at": "Sun, Aug 19 11:02:20.537 2018"
        },
        {
            "id": 1,
            "manager_id": 3,
            "manager_user": {
                "id": 3,
                "name": "Jenny",
                "phone": "1-800-867-5309",
                "role": "manager",
                "created_at": "Sun, Aug 19 11:02:20.537 2018",
                "updated_at": "Sun, Aug 19 11:02:20.537 2018"
            },
            "employee_id": 3,
            "employee_user": {
                "id": 3,
                "name": "Jenny",
                "phone": "1-800-867-5309",
                "role": "manager",
                "created_at": "Sun, Aug 19 11:02:20.537 2018",
                "updated_at": "Sun, Aug 19 11:02:20.537 2018"
            },
            "break": 0.45,
            "start_time": "Sun, Aug 19 18:30:00.000 2018",
            "end_time": "Sun, Aug 19 20:30:00.000 2018",
            "created_at": "Sun, Aug 19 11:02:20.537 2018",
            "updated_at": "Sun, Aug 19 11:02:26.222 2018"
        }
    ]
}
```
</p>
</details>

#### As an employee, I want to know how much I worked, by being able to get a summary of hours worked for each week.
This request will show hours scheduled/worked grouped by week. The employee_id is specified in the URL.
<details><summary>HTTP Request</summary>
<p>

```http request
GET /api/summaries/3?current_user_id=1

{
    "success": true,
    "results": [
        {
            "employee_id": 3,
            "employee_user": {
                "id": 3,
                "name": "Jenny",
                "phone": "1-800-867-5309",
                "role": "manager",
                "created_at": "Sun, Aug 19 11:37:31.519 2018",
                "updated_at": "Sun, Aug 19 11:37:31.519 2018"
            },
            "week_start": "Mon, Aug 13 00:00:00.000 2018",
            "week_end": "Mon, Aug 20 00:00:00.000 2018",
            "total_shifts": 2,
            "total_scheduled_time": 3,
            "total_scheduled_time_formatted": "3 Hour(s) 0 Minute(s)",
            "total_worked_time": 0.05,
            "total_worked_time_formatted": "0 Hour(s) 2 Minute(s)",
            "total_break_time": 0,
            "total_break_time_formatted": "0 Hour(s) 0 Minute(s)"
        },
        {
            "employee_id": 3,
            "employee_user": {
                "id": 3,
                "name": "Jenny",
                "phone": "1-800-867-5309",
                "role": "manager",
                "created_at": "Sun, Aug 19 11:37:31.519 2018",
                "updated_at": "Sun, Aug 19 11:37:31.519 2018"
            },
            "week_start": "Mon, Aug 20 00:00:00.000 2018",
            "week_end": "Mon, Aug 27 00:00:00.000 2018",
            "total_shifts": 2,
            "total_scheduled_time": 8,
            "total_scheduled_time_formatted": "8 Hour(s) 0 Minute(s)",
            "total_worked_time": 0,
            "total_worked_time_formatted": "0 Hour(s) 0 Minute(s)",
            "total_break_time": 0,
            "total_break_time_formatted": "0 Hour(s) 0 Minute(s)"
        }
    ]
}
```     
</p>
</details>

#### As an employee, I want to be able to contact my managers, by seeing manager contact information for my shifts.
Manager information is returned as a sub-object with the response for shifts. Field: `manager_user`
<details><summary>HTTP Request</summary>
<p>

```http request
GET /api/shifts/mine?current_user_id=1

{
    "success": true,
    "results": [
        {
            "id": 114,
            "manager_id": 3,
            "manager_user": {
                "id": 3,
                "name": "Jenny",
                "phone": "1-800-867-5309",
                "role": "manager",
                "created_at": "Sun, Aug 19 11:02:20.537 2018",
                "updated_at": "Sun, Aug 19 11:02:20.537 2018"
            },
            "employee_id": 1,
            "employee_user": {
                "id": 1,
                "name": "Elliot",
                "email": "elliot@elliot.com",
                "role": "employee",
                "created_at": "Sun, Aug 19 11:02:20.537 2018",
                "updated_at": "Sun, Aug 19 11:02:20.537 2018"
            },
            "break": 0,
            "start_time": "Thu, Aug 02 19:31:46.631 2018",
            "end_time": "Thu, Aug 02 20:31:46.631 2018",
            "created_at": "Sun, Aug 19 11:02:25.736 2018",
            "updated_at": "Sun, Aug 19 11:02:25.736 2018"
        }
    ]
}
```
</p>
</details>


#### As a manager, I want to schedule my employees, by creating shifts for any employee.
Shifts can be created with a PUT or POST request to `/api/shifts`
If the creation was successful it will return the created object.
<details><summary>HTTP Request</summary>
<p>

```http request
POST /api/shifts?current_user_id=3
{
    "manager_id": 3,
    "employee_id": 1,
    "break": 0,
    "start_time": "Thu, Aug 1 19:31:46.631 2018",
    "end_time": "Thu, Aug 1 20:31:46.631 2018"
}



{
    "success": true,
    "results": {
        "id": 7,
        "manager_id": 3,
        "employee_id": 1,
        "break": 0,
        "start_time": "Wed, Aug 01 19:31:46.631 2018",
        "end_time": "Wed, Aug 01 20:31:46.631 2018",
        "created_at": "Sun, Aug 19 11:50:15.112 2018",
        "updated_at": "Sun, Aug 19 11:50:15.112 2018"
    }
}
```


If the creation was not successful, an error message will be returned. For example; if the start time comes after the end time.
```http request
POST /api/shifts?current_user_id=3
{
    "manager_id": 3,
    "employee_id": 1,
    "break": 0,
    "start_time": "Thu, Aug 3 19:31:46.631 2018",
    "end_time": "Thu, Aug 1 20:31:46.631 2018"
}



{
    "message": "Error, (start_time: Thu, Aug 3 19:31:46.631 2018) must come before (end_time: Thu, Aug 1 20:31:46.631 2018).",
    "success": false
}
```

Or if the shift time's conflict with another shift for that employee.
```http request
POST /api/shifts?current_user_id=3
{
    "manager_id": 3,
    "employee_id": 1,
    "break": 0,
    "start_time": "Thu, Aug 2 19:31:46.631 2018",
    "end_time": "Thu, Aug 2 20:31:46.631 2018"
}



{
    "message": "Error, 1 shift(s) already exist for user ID 1 during the start/end time. Conflicting shift(s): 7.",
    "success": false
}
```
</p>
</details>


#### As a manager, I want to see the schedule, by listing shifts within a specific time period.
There are 4 url params that are used for this request: `date_from`, `date_to`, `date_time_from` and `date_time_to`
Only 1 from and 1 to are allowed per request. The `_time_` params will be converted to `timestamp` in SQL before filtering.
But the normal params will be filtered as `date`. 
<details><summary>HTTP Request</summary>
<p>

```http request
GET /api/shifts?current_user_id=1&date_from=2018-08-01&date_to=2018-08-06

{
    "success": true,
    "results": [
        {
            "id": 7,
            "manager_id": 3,
            "manager_user": {
                "id": 3,
                "name": "Jenny",
                "phone": "1-800-867-5309",
                "role": "manager",
                "created_at": "Sun, Aug 19 11:57:32.489 2018",
                "updated_at": "Sun, Aug 19 11:57:32.489 2018"
            },
            "employee_id": 1,
            "employee_user": {
                "id": 1,
                "name": "Elliot",
                "email": "elliot@elliot.com",
                "role": "employee",
                "created_at": "Sun, Aug 19 11:57:32.489 2018",
                "updated_at": "Sun, Aug 19 11:57:32.489 2018"
            },
            "break": 0,
            "start_time": "Thu, Aug 02 19:31:46.631 2018",
            "end_time": "Thu, Aug 02 20:31:46.631 2018",
            "created_at": "Sun, Aug 19 11:57:32.618 2018",
            "updated_at": "Sun, Aug 19 11:57:32.618 2018"
        }
    ]
}
```
</p>
</details>


#### As a manager, I want to be able to change a shift, by updating the time details.
When updating a shift, the only required information is the shift ID. Only the fields that are being changed need to be included in the request.
However; if a `manager_id` is not specified when updating a shift, the `current_user_id` will be set as the new `manager_id`.
<details><summary>HTTP Request</summary>
<p>

```http request
PUT /api/shifts/1?current_user_id=3
{
    "start_time": "Sun, Aug 19 18:30:00.000 2018",
    "end_time": "Mon, Aug 19 20:30:00.00 2018"
}



{
    "success": true,
    "results": {
        "id": 1,
        "manager_id": 3,
        "employee_id": 3,
        "break": 1,
        "start_time": "Sun, Aug 19 18:30:00.000 2018",
        "end_time": "Sun, Aug 19 20:30:00.000 2018",
        "created_at": "Sun, Aug 19 12:05:11.306 2018",
        "updated_at": "Sun, Aug 19 12:05:11.387 2018"
    }
}
```
</p>
</details>



#### As a manager, I want to be able to assign a shift, by changing the employee that will work a shift.
When updating a shift, the only required information is the shift ID. Only the fields that are being changed need to be included in the request.
However; if a `manager_id` is not specified when updating a shift, the `current_user_id` will be set as the new `manager_id`.
When changing the `employee_id` for a shift, that shift will be validated the same way as when it was created.
If the updated shift would conflict with another shift for that employee (unless the `employee_id` is `NULL`) an error will be returned.
<details><summary>HTTP Request</summary>
<p>

```http request
PUT /api/shifts/1?current_user_id=3
{
    "employee_id": 2
}



{
    "success": true,
    "results": {
        "id": 1,
        "manager_id": 3,
        "employee_id": 2,
        "break": 1,
        "start_time": "Sun, Aug 19 18:30:00.000 2018",
        "end_time": "Sun, Aug 19 20:30:00.000 2018",
        "created_at": "Sun, Aug 19 12:05:11.306 2018",
        "updated_at": "Sun, Aug 19 12:05:11.387 2018"
    }
}
```
</p>
</details>



#### As a manager, I want to contact an employee, by seeing employee details.
`GET` requests should return the field `employee_user` which will include all of the user's details from the users table.
<details><summary>HTTP Request</summary>
<p>

```http request
GET /api/shifts?current_user_id=3&date_from=2018-08-01&date_to=2018-08-06

{
    "success": true,
    "results": [
        {
            "id": 7,
            "manager_id": 3,
            "manager_user": {
                "id": 3,
                "name": "Jenny",
                "phone": "1-800-867-5309",
                "role": "manager",
                "created_at": "Sun, Aug 19 11:57:32.489 2018",
                "updated_at": "Sun, Aug 19 11:57:32.489 2018"
            },
            "employee_id": 1,
            "employee_user": {
                "id": 1,
                "name": "Elliot",
                "email": "elliot@elliot.com",
                "role": "employee",
                "created_at": "Sun, Aug 19 11:57:32.489 2018",
                "updated_at": "Sun, Aug 19 11:57:32.489 2018"
            },
            "break": 0,
            "start_time": "Thu, Aug 02 19:31:46.631 2018",
            "end_time": "Thu, Aug 02 20:31:46.631 2018",
            "created_at": "Sun, Aug 19 11:57:32.618 2018",
            "updated_at": "Sun, Aug 19 11:57:32.618 2018"
        }
    ]
}
```

Or the users endpoint can be queried directly.
```http request
GET /api/users?current_user_id=3

{
    "success": true,
    "results": [
        {
            "id": 1,
            "name": "Elliot",
            "email": "elliot@elliot.com",
            "role": "employee",
            "created_at": "Sun, Aug 19 12:05:11.306 2018",
            "updated_at": "Sun, Aug 19 12:05:11.306 2018"
        },
        {
            "id": 2,
            "name": "Jimmy",
            "email": "jimmy@johns.com",
            "role": "employee",
            "created_at": "Sun, Aug 19 12:05:11.306 2018",
            "updated_at": "Sun, Aug 19 12:05:11.306 2018"
        },
        {
            "id": 3,
            "name": "Jenny",
            "phone": "1-800-867-5309",
            "role": "manager",
            "created_at": "Sun, Aug 19 12:05:11.306 2018",
            "updated_at": "Sun, Aug 19 12:05:11.306 2018"
        },
        {
            "id": 4,
            "name": "Henry",
            "phone": "1-800-123-4561",
            "role": "employee",
            "created_at": "Sun, Aug 19 12:05:11.306 2018",
            "updated_at": "Sun, Aug 19 12:05:11.306 2018"
        }
    ]
}
```
</p>
</details>

# Notes
> Currently the database is storing date/times in central daylight time. I struggled to come up with a good way to convert dates/times to RFC 2822 in Golang, so I opted to instead convert/format them in PostgreSQL. This means that all of the date/time fields are strings in Go but since no operation is being performed on the data itself there, this should be fine. I might change this later and add the formatting to the driver somehow instead? As long as it looks clean and wouldn't effect changes in the future.

> I've taken some code from one of my other website projects that I was working on in Go to make the REST API a bit easier when it came to filtering/sorting/selecting fields to be returned. And the date range handling for Issue #6 I added into that filtering handler.

> I had some issues with nil references in the shifts code. Kept running into this: https://golang.org/doc/faq#nil_error I was able to fix it but the code was so spaghetti I ended up discarding all of the changes because it also ended up introducing a bunch of other bugs. So the verifyShift section works but it's not as pretty as I'd like it to be.
