/* ----------------- COMPONENTS ----------------- */
var Home = { template: "#template-home" }
var LocationList = { template: '#template-locationlist' }
var LocationProfile = { template: '#template-locationprofile' }
var Login = { template: '#template-login' }

/* ------------------- ROUTES ------------------- */
var routes = [
  { path: '/home', component: Home },
  { path: '/locations', component: LocationList },
  { path: '/locations/:id', component: LocationProfile },
  { path: '/login', component: Login }
]

var router = new VueRouter({
  routes // short for routes: routes
})

// Mount to app
var app = new Vue({
  router
}).$mount('#app')
