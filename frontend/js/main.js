var fakeData = {
  locations: [
    { "id": 1, "name": "Location 1", "address": "123 Road Ln St. Louis, MO 123456"},
    { "id": 2, "name": "Homeless Shelter"}
  ]
};

/* ----------------- COMPONENTS ----------------- */
var HomeSubComps = {
  welcome: {
    template: "<div>Landing page</div>"
  },
  host: {
    template: "#template-home-host"
  },
  seeker: {
    template: "<div>Seekers Home</div>"
  },
  center: {
    template: "<div>Centers Home</div>"
  },
  organization: {
    template: "<div>Organizations Home</div>"
  }
}

var Home = {
  template: "#template-home",
  components: HomeSubComps,
  computed: {
    currentView: function() {
      if (typeof app === "undefined" || app.user == null) {
        return "welcome";
      }
      else {
        return app.user.type;
      }
    }
  }
};

var LocationList = {
  template: "#template-locationlist",
  data: function() {
    return {
      locations: fakeData.locations
    }
  }
};

var LocationProfile = {
  template: "#template-locationprofile",
  data: function() {
    var self = this;

    var curLocation = fakeData.locations.filter(function(obj) {
      return obj.id == self.$route.params.id;
    })[0];

    return {
      location: curLocation
    }
  }
};

var Login = {
  template: "#template-login",
  data: function() {
    return {
      formData: {
        email: "",
        pass: ""
      }
    }
  },

  methods: {
    login: function() {
      // Example user request
      app.$http.post('login', {}, { params: this.formData }).then(function(data) {
        app.user = JSON.parse(data.body);
        app.addMessage("success", "Success", "User successfully logged in.");
        router.push({ name: "home" })
      }).catch(function(fail) {
        app.addMessage("danger", "Error", "Invalid email or password.")
      });
    }
  }
};

var Register = {
  template: "#template-register",
  data: function() {
    return {
      formData: {
        name: "",
        email: "",
        pass: "",
        confirmPass: "",
        type: "seeker"
      }
    }
  },

  methods: {
    register: function() {
      // Register Logic
      form_valid = true;
      email_re = /^(([^<>()\[\]\\.,;:\s@"]+(\.[^<>()\[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;

      // Make sure name is entered
      if (!this.formData.name) {
        app.addMessage("danger", "Error", "Name field is required.");
        form_valid = false;
      }
      if (!email_re.test(this.formData.email)) {
        app.addMessage("danger", "Error", "Invalid e-mail address.");
        form_valid = false;
      }
      if (this.formData.pass !== this.formData.confirmPass) {
        app.addMessage("danger", "Error", "Passwords do not match.");
        form_valid = false;
      }

      if (form_valid) {
        // Register User
        app.$http.post('user', {}, { params: this.formData }).then(function(data) {
          app.addMessage("success", "Success", "Successfully registered user!");
          router.push({ "name": "login" })
        }).catch(function(data) {
          app.addMessage("danger", "Error", "User with email already exists.");
        });
      }
    }
  },

  computed: {
    namePlaceholder: function() {
      var data = this.formData;
      return (data.type == "seeker" || data.type == "host") ? "First Last" : "Organization Name"
    }
  }
};

var Settings = {
  parent: {
    template: "#template-settings"
  },
  general: {
    template: "#template-settings-general"
  },
  security: {
    template: "<div>Security Settings</div>"
  },
  dependencies: {
    template: "<div>Dependency Settings</div>"
  },
  events: {
    template: "<div>Event Settings</div>"
  },
  newEvent: {
    template: "#template-new-event",
    data: function() {
      return {
        tmpNeed: "",
        needs: [],
        formData: {
          type: "room",
          slots: "",
          divisions: "",
          description: ""
        }
      }
    },
    methods: {
      createEvent: function() {
        this.formData.description = encodeURIComponent(this.formData.description);
        this.formData.needs = this.needs.join(",");
        app.$http.post("event", {}, { params: this.formData }).then(function(data) {
          app.addMessage("success", "Success", "Event successfully created!");
          router.push({ name: "viewEvent", params: { id: JSON.parse(data.body).id } });
        }).catch(function(data) {
          app.addMessage("danger", "Error", "Error creating event");
        });
      },
      addNeed: function() {
        this.needs.push(this.tmpNeed);
        this.tmpNeed = "";
      },
      removeNeed: function(value) {
        var index = this.needs.indexOf(value);
        if (index > -1) {
          this.needs.splice(index, 1);
        }
      }
    }
  },
  viewEvent: {
    template: "#template-view-event",
    data: function() {
      return {
        event: {}
      }
    },
    beforeCreate: function() {
      var self = this;
      // Get data
      Vue.http.get("api/event", { params: this.$route.params }).then(function(data) {
        self.event = JSON.parse(data.body);
      }).catch(function(err) {
        router.push({ name: "home" });
      });
    }
  }
}

/* ----------------- AUTH GUARD ----------------- */
var getUserInfo = function(to, from, next) {
  if (typeof app === "undefined" || app.user == null) {
    Vue.http.get("api/login").then(function(data) {
      app.user = JSON.parse(data.body);
      next();
    }).catch(function(data) {
      next();
    });
  }
  else {
    next();
  }
}

var requireLogin = function(to, from, next) {
  // If app doesn't exist just go to home page
  if (typeof app === 'undefined') {
    Vue.http.get("api/login").then(function(data) {
      app.user = JSON.parse(data.body);
      next();
      return;
    }).catch(function(data) {
      next({ name: "login" });
      return;
    });
  }
  // Allow user to go to destination if logged in
  else if (app.user) {
    next();
  }

  // Redirect to login
  else {
    next({ name: "login" });
  }
}

var requireGuest = function(to, from, next) {
  if (typeof app === "undefined") {
    next();
    return;
  }

  // Allow user to go to destination if logged in
  if (app.user) {
    next({ name: "home" });
  }

  // Redirect to login
  else {
    next();
  }
}

var requireType = function(type) {
  return function(to, from, next) {
    if (app.user.type === type) {
      next();
    }
    else {
      next(false);
    }
  }
}

var logout = function(to, from, next) {
  app.$http.get('logout').then(function() {
    app.user = null;
    app.addMessage("success", "Success", "Successfully logged out");
    next({ name: "home" });
  }).catch(function() {
    app.addMessage("danger", "Error", "Error trying to log user out");
  });
}
/* ------------------- ROUTES ------------------- */
var routes = [
  { path: "/", redirect: "/home" },
  { name: "home", path: "/home", component: Home, beforeEnter: getUserInfo },
  { name: "locationList", path: "/locations", component: LocationList, beforeEnter: getUserInfo },
  { name: "locationProfile", path: "/locations/:id", component: LocationProfile, beforeEnter: getUserInfo },
  { name: "login", path: "/login", component: Login, beforeEnter: requireGuest },
  { name: "register", path: "/register", component: Register, beforeEnter: requireGuest },
  { name: "logout", path: "/logout", beforeEnter: logout },
  { name: "viewEvent", path: "/events/:id", component: Settings.viewEvent, beforeEnter: requireLogin },
  {
    name: "settings",
    path: "/settings",
    component: Settings.parent,
    beforeEnter: requireLogin,
    redirect: "/settings/general",
    children: [
      {
        name: "general",
        path: "general",
        component: Settings.general
      },
      {
        name: "security",
        path: "security",
        component: Settings.security
      },
      {
        name: "dependencies",
        path: "dependencies",
        component: Settings.dependencies,
        beforeEnter: requireType("seeker")
      },
      {
        name: "events",
        path: "events",
        component: Settings.events,
        beforeEnter: requireType("host"),
      },
      {
        name: "newEvent",
        path: "events/new",
        component: Settings.newEvent,
        beforeEnter: requireType("host")
      }
    ]
  }
];

var router = new VueRouter({
  routes // short for routes: routes
});

// Mount to app
var app = new Vue({
  router,
  http: {
    root: '/api'
  },
  data: {
    messages: [],
    user: null
  },
  methods: {
    addMessage: function(type, title, message) {
      this.messages.push({ "type": type, "title": title, "message": message });
    }
  }
}).$mount("#app");
