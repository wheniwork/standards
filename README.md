[![Build Status](https://travis-ci.com/ECourant/standards.svg?branch=master)](https://travis-ci.com/ECourant/standards)

### Elliot's Sample Project For When I Work

This project requires Golang 1.10 and PostgreSQL 10. 
Because of the way I structured the folders I don't believe you can `go get` this project. But at the end I will include a vendor folder which should have any dependency for the actual site.

The database change scripts currently don't target a specific database instance, and I've just run them on the default `postgres` database instance. If you would like to change what database they are in, you will also need to update the connection string in `config.json` in the Site directory and `database_test_config.json` in the Database directory.

I've also kept track of what needs to be done for this project and some brief notes on how to go about doing that in Issues and Projects. 

If I have time I will also try to create a UI for the REST API using Framework7. I'm not sure how much time I will have to do that though and it will be the last thing I do. 



Some small notes:
> Currently the database is storing date/times in central daylight time. I struggled to come up with a good way to convert dates/times to RFC 2822 in Golang, so I opted to instead convert/format them in PostgreSQL. This means that all of the date/time fields are strings in Go but since no operation is being performed on the data itself there, this should be fine. I might change this later and add the formatting to the driver somehow instead? As long as it looks clean and wouldn't effect changes in the future.

> I've taken some code from one of my other website projects that I was working on in Go to make the REST API a bit easier when it came to filtering/sorting/selecting fields to be returned. And the date range handling for Issue #6 I added into that filtering handler.

> I had some issues with nil references in the shifts code. Kept running into this: https://golang.org/doc/faq#nil_error I was able to fix it but the code was so spaghetti I ended up discarding all of the changes because it also ended up introducing a bunch of other bugs. So the verifyShift section works but it's not as pretty as I'd like it to be. I also still need to add some code for making sure that if it's updating a shift it validates that it doesn't conflict properly.
