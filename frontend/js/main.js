var fakeData = {
  locations: [
    { "id": 1, "name": "Location 1", "address": "123 Road Ln St. Louis, MO 123456"},
    { "id": 2, "name": "Homeless Shelter"}
  ]
};

/* ----------------- COMPONENTS ----------------- */
var Home = {
  template: "#template-home"
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
        app.user = data;
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
    template: "<div>General Settings</div>"
  },
  security: {
    template: "<div>Security Settings</div>"
  }
}

/* ----------------- AUTH GUARD ----------------- */
var requireLogin = function(to, from, next) {
  // If app doesn't exist just go to home page
  if (typeof app === 'undefined') {
    next({ name: "home" });
    return;
  }

  // Allow user to go to destination if logged in
  if (app.user) {
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
/* ------------------- ROUTES ------------------- */
var routes = [
  { path: "/", redirect: "/home" },
  { name: "home", path: "/home", component: Home },
  { name: "locationList", path: "/locations", component: LocationList },
  { name: "locationProfile", path: "/locations/:id", component: LocationProfile },
  { name: "login", path: "/login", component: Login, beforeEnter: requireGuest },
  { name: "register", path: "/register", component: Register, beforeEnter: requireGuest },
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
  created: function() {
    var self = this;
    this.$http.get("login").then(function(data) {
      self.user = data;
    });
  },
  methods: {
    addMessage: function(type, title, message) {
      this.messages.push({ "type": type, "title": title, "message": message });
    }
  }
}).$mount("#app");
