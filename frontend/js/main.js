var fakeData = {
  locations: [
    { "id": 1, "name": "Location 1", "address": "123 Road Ln St. Louis, MO 123456"},
    { "id": 2, "name": "Homeless Shelter"}
  ]
}

/* ----------------- COMPONENTS ----------------- */
var Home = {
  template: "#template-home"
}

var LocationList = {
  template: "#template-locationlist",
  data: function() {
    return {
      locations: fakeData.locations
    }
  }
}

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
}

var Login = {
  template: "#template-login",
  data: function() {
    return {
      formData: {
        email: "",
        password: ""
      }
    }
  },

  methods: {
    login: function() {
      // Login Logic
      console.log(JSON.stringify(this.formData))
    }
  }

}

var Register = {
  template: "#template-register",
  data: function() {
    return {
      formData: {
        email: "",
        password: "",
        confirmPassword: ""
      }
    }
  },

  methods: {
    register: function() {
      // Register Logic
      console.log(JSON.stringify(this.formData));

      if (this.formData.password !== this.formData.confirmPassword) {
        app.addMessage("danger", "Error", "Passwords do not match.");
      }
    }
  }
}

/* ------------------- ROUTES ------------------- */
var routes = [
  { path: "/", redirect: "/home" },
  { name: "home", path: "/home", component: Home },
  { name: "locationList", path: "/locations", component: LocationList },
  { name: "locationProfile", path: "/locations/:id", component: LocationProfile },
  { name: "login", path: "/login", component: Login },
  { name: "register", path: "/register", component: Register }
]

var router = new VueRouter({
  routes // short for routes: routes
})

// Mount to app
var app = new Vue({
  router,
  data: {
    messages: [],
    test: "test"
  },
  methods: {
    addMessage: function(type, title, message) {
      this.messages.push({ "type": type, "title": title, "message": message });
    }
  }
}).$mount("#app")
