(function () {
    //Helpers
    Template7.registerHelper('title_case', function (str) {
        str = str.toLowerCase().split(' ');
        for (let i = 0; i < str.length; i++) {
            str[i] = str[i].charAt(0).toUpperCase() + str[i].slice(1);
        }
        return str.join(' ');
    });

    Template7.registerHelper('role_badge', function (role) {
        if (role == "employee") {
            return "color-blue"
        } else if (role == "manager") {
            return "color-orange"
        }
    });

    //Views
    var homeView = Template7.compile(Dom7("#homeView").html());
    var appView = Template7.compile(Dom7("#appView").html());
    var shiftsFiltersView = Template7.compile(Dom7("#shiftsFiltersView").html());
    var shiftDetailsView = Template7.compile(Dom7("#shiftDetailsView").html());
    var summariesFiltersView = Template7.compile(Dom7("#summariesFiltersView").html());

    var current_user_id = -1;
    var current_user_name = "";
    var current_user_is_manager = false;

    var current_date = new Date();

    // The filter variables will be used to generate urls for API requests.
    var current_shifts_filter = {
        show_only_my_shifts: true,
        date_from: new Date(new Date().setDate(current_date.getDate() - 7)),
        date_to: new Date(),
        page: 1,
        page_size: 10,
        base_url: "http://localhost:8080/api/shifts",
        sort: " start_time",
        message: "Showing shifts for the next 7 days.",
        GetShiftURL: function () {
            var url = this.base_url;
            var params = ["current_user_id=" + current_user_id];

            if (this.show_only_my_shifts) {
                url += "/mine";
            }
            this.message = "Showing shifts from " + parseDateToURLParam(this.date_from) + " to " + parseDateToURLParam(this.date_to);
            params.push("date_from=" + parseDateToURLParam(this.date_from));
            params.push("date_to=" + parseDateToURLParam(this.date_to));
            params.push("page=" + this.page);
            params.push("page_size=" + this.page_size);
            params.push("order=" + encodeURI(this.sort));
            url = url + "?" + params.join("&");
            console.log(url);
            return url;
        },
        Reset: function() {
            this.show_only_my_shifts = true;
            this.date_from = new Date();
            this.date_to = new Date(new Date().setDate(current_date.getDate() + 7));
            this.page = 1;
            this.page_size = 10;
            this.message = "Showing shifts from " + parseDateToURLParam(this.date_from) + " to " + parseDateToURLParam(this.date_to);
        }
    };

    var current_summary_filter = {
        date_from: new Date(new Date().setDate(current_date.getDate() - 7)),
        date_to: new Date(),
        page: 1,
        page_size: 10,
        base_url: "http://localhost:8080/api/summaries",
        sort: "-week_start",
        message: "Showing shifts for the next 7 days.",
        GetSummariesURL: function () {
            var url = this.base_url + "/" + current_user_id;
            var params = ["current_user_id=" + current_user_id];

            this.message = "Showing shifts from " + parseDateToURLParam(this.date_from) + " to " + parseDateToURLParam(this.date_to);
            params.push("date_from=" + parseDateToURLParam(this.date_from));
            params.push("date_to=" + parseDateToURLParam(this.date_to));
            params.push("page=" + this.page);
            params.push("page_size=" + this.page_size);
            params.push("order=" + encodeURI(this.sort));
            url = url + "?" + params.join("&");
            console.log(url);
            return url;
        },
        Reset: function() {
            this.date_from = new Date();
            this.date_to = new Date(new Date().setDate(current_date.getDate() + 7));
            this.page = 1;
            this.page_size = 10;
            this.message = "Showing shifts from " + parseDateToURLParam(this.date_from) + " to " + parseDateToURLParam(this.date_to);
        }
    };

    var current_overlapping_filter = {
        id: -1,
        page: 1,
        page_size: 100,
        base_url: "http://localhost:8080/api/shifts/overlapping",
        GetOverlappingURL: function () {
            var url = this.base_url + "/" + this.id;
            var params = ["current_user_id=" + current_user_id];

            params.push("page=" + this.page);
            params.push("page_size=" + this.page_size);
            url = url + "?" + params.join("&");
            console.log(url);
            return url;
        },
    };



    var shifts = [];
    var summaries = [];
    function parseDateToURLParam(date) {
        return encodeURI((date.getMonth() + 1) + '/' + date.getDate() + '/' + date.getFullYear());
    }

    var app = new Framework7({
        ios: true,
        desktop: false,
        // App root element
        root: '#app',
        // App Name
        name: 'My App',
        // App id
        id: 'com.myapp.test',
        // Enable swipe panel
        panel: {
            swipe: 'left',
        },
        // Add default routes
        routes: [
            {
                path: '/homeView/',
                async: function (routeTo, routeFrom, resolve, reject) {
                    app.request.json("http://localhost:8080/api/users?current_user_id=1&order= name", function (data) {
                        if (data.success) {
                            resolve({
                                template: homeView
                            },
                            {
                                context: data
                            });
                        } else {
                            reject();
                        }
                    });
                },
            },
            {
                path: '/appView/',
                async: function (routeTo, routeFrom, resolve, reject) {
                    shifts = [];
                    summaries = [];
                    app.request.json(current_shifts_filter.GetShiftURL(), function (shiftdata) {
                        if (shiftdata.success) {
                            shifts = shiftdata.results;
                            app.request.json(current_summary_filter.GetSummariesURL(), function (data) {
                                if (data.success) {
                                    summaries = data.results;
                                    resolve({
                                            template: appView
                                        },
                                        {
                                            context: {
                                                current_user_name: current_user_name,
                                                is_manager: current_user_is_manager,
                                                shifts: shifts,
                                                summaries: summaries,
                                                shift_message: current_shifts_filter.message,
                                            }
                                        });
                                } else {
                                    app.dialog.alert(data.message);
                                    reject();
                                }
                            });
                        } else {
                            app.dialog.alert(shiftdata.message);
                            reject();
                        }
                    });
                },
            },
            {
                path: '/shiftsFiltersView/',
                async: function (routeTo, routeFrom, resolve, reject) {
                    resolve({
                        template: shiftsFiltersView
                    },
                    {
                        context: {
                            current_user_name: current_user_name,
                            is_manager: current_user_is_manager,
                        }
                    });
                },
                on: {
                    pageBeforeIn: function (e, page) {
                        var calendarRange = app.calendar.create({
                            inputEl: '#shiftDateRange',
                            dateFormat: 'M dd yyyy',
                            rangePicker: true,
                            closeOnSelect: true,
                            value:[
                                current_shifts_filter.date_from,
                                current_shifts_filter.date_to
                            ]
                        });
                        $("#showOnlyMyShifts")[0].checked = current_shifts_filter.show_only_my_shifts;
                    }
                }
            },
            {
                path: '/summariesFiltersView/',
                async: function (routeTo, routeFrom, resolve, reject) {
                    resolve({
                            template: summariesFiltersView
                        },
                        {
                            context: {
                                current_user_name: current_user_name,
                                is_manager: current_user_is_manager,
                            }
                        });
                },
                on: {
                    pageBeforeIn: function (e, page) {
                        var calendarRange = app.calendar.create({
                            inputEl: '#summaryDateRange',
                            dateFormat: 'M dd yyyy',
                            rangePicker: true,
                            closeOnSelect: true,
                            value:[
                                current_summary_filter.date_from,
                                current_summary_filter.date_to
                            ]
                        });
                    }
                }
            },
            {
                path: '/shiftDetailsView/',
                async: function (routeTo, routeFrom, resolve, reject) {
                    shifts = [];
                    summaries = [];
                    app.request.json(current_overlapping_filter.GetOverlappingURL(), function (data) {
                        if (data.success) {
                            resolve({
                                template: shiftDetailsView
                            },
                            {
                                context: {
                                    current_user_name: current_user_name,
                                    is_manager: current_user_is_manager,
                                    detail: data.results,
                                    shifts: data.results.shifts
                                }
                            });
                        } else {
                            app.dialog.alert(data.message);
                            reject();
                        }
                    });
                },
            },
        ],
        // ... other parameters
    });
    app.router.navigate("/homeView/", {
        animate: false
    });

    $(document).on("click", ".login-select", function () {
        current_user_id = parseInt($(this).attr("user-id"));
        current_user_name = $(this).attr("user-name");
        current_user_is_manager = $(this).attr("user-role") == "manager";
        console.log("New Current User: " + current_user_id);
        current_shifts_filter.Reset();
        current_summary_filter.Reset();
        app.router.navigate("/appView/");
    });

    $(document).on("click", ".filter-shifts", function () {      
        app.router.navigate("/shiftsFiltersView/");
    });

    $(document).on("click", ".filter-summaries", function () {
        app.router.navigate("/summariesFiltersView/");
    });

    $(document).on("click", ".shift-link", function() {
       id = $(this).attr("shift-id");
       current_overlapping_filter.id = parseInt(id);
       app.router.navigate("/shiftDetailsView/");
    });

    $(document).on("click", "#submitShiftFilter", function () {
        current_shifts_filter.show_only_my_shifts = $("#showOnlyMyShifts")[0].checked;
        var shiftDateRange = $("#shiftDateRange").val();
        if (shiftDateRange != "") {
            var date_1 = new Date(shiftDateRange.split(" - ")[0]);
            var date_2 = new Date(shiftDateRange.split(" - ")[1]);
            if (date_1.getTime() < date_2.getTime()) {
                current_shifts_filter.date_from = date_1;
                current_shifts_filter.date_to = date_2;
            } else {
                current_shifts_filter.date_from = date_2;
                current_shifts_filter.date_to = date_1;
            }
        }
        app.router.back("/appView/", {
            ignoreCache: true,
            force: true
        })
    });

    $(document).on("click", "#submitSummaryFilter", function () {
        var shiftDateRange = $("#summaryDateRange").val();
        if (shiftDateRange != "") {
            var date_1 = new Date(shiftDateRange.split(" - ")[0]);
            var date_2 = new Date(shiftDateRange.split(" - ")[1]);
            if (date_1.getTime() < date_2.getTime()) {
                current_summary_filter.date_from = date_1;
                current_summary_filter.date_to = date_2;
            } else {
                current_summary_filter.date_from = date_2;
                current_summary_filter.date_to = date_1;
            }
        }
        app.router.back("/appView/", {
            ignoreCache: true,
            force: true
        });
        setTimeout(function(){
            app.tab.show("#tab-2", false);
        }, 100);
    });
})();