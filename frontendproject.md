# Frontend Developer Challenge

Choose a MVVM/Redux framework of your choice. While we encourage using React, Vue or Angular, please use whatever framework and language you are most comfortable with. Our application frontend is written in in Backbone and React.

Please provide a means for someone to checkout your project and run it.

## Requirements

You will be given an access token to an existing WIW account.  Using the WIW API (http://dev.wheniwork.com) you will create the following:

- Handle authentication scenarios with the below.  While we don't need a login page, provide a UX flow for an invalid/unauthorized token.
- An authenticated "page" of the application that will contain the following:
  - A CRUD for Users updating the basic information (email, first_name, last_name, just these fields we dont need the entire user object)
  - A UI to assign and delete positions associated to a user.  Feel free to create positions in the web UI so they are available to add and remove from a user.
  - While we are not asking for 100% test coverage, please demonstrate writing frontend tests against your project.
- The intent of this test is to gauge your abilities with a frontend library/framework.  We would encourage use of a CSS framwork for stying, something like https://getbootstrap.com/docs/4.0/getting-started/introduction/.

