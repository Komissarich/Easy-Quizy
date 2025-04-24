import { createRouter, createWebHistory } from "vue-router";
import Auth from "./components/Auth.vue";
import Profile from "./components/Profile.vue";
import Play from "./components/Play.vue";
import Register from "./components/Register.vue";
import jwtDecode from 'jwt-decode'
import authState from "./main"
import CreateTest from "./components/CreateTest.vue";
import CreateQuestions from "./components/CreateQuestions.vue";
import Quiz from "./components/Quiz.vue";
import Result from "./components/Result.vue";
import ShowQuiz from "./components/ShowQuiz.vue";
import OtherProfile from "./components/OtherProfile.vue";


const router = createRouter( {
     history: createWebHistory(),
    routes: [
        {path: '/auth', component: Auth, meta: {requireAuth: false}},
        {path: '/register', component: Register, meta: {requireAuth: false}},
     
        {path: '/create-test', component: CreateTest, meta: {requireAuth: true},children: [
            {path: 'questions', component: CreateQuestions} 
        ]},
        {path: '/play', component: Play, meta: {requireAuth: false}, alias: '/' },
        {path: '/play/:quiz_id',component: Quiz,meta: { requireAuth: false }},
        {path: '/play/:quiz_id/result', name: "QuizResult", component: Result,meta: { requireAuth: false }},
        {path: '/profile/me', name: 'Profile', component: Profile, meta: {requireAuth: true}},
        {path: '/profile/:username', name: 'OtherProfile', component: OtherProfile, meta: {requireAuth: true}},
        {path: '/quiz/:quiz_id', name: 'ShowQuiz', component: ShowQuiz, meta: {requireAuth: true}},
    ]
})

router.beforeEach((to, from, next) => {
   
    
    if (to.meta.requireAuth === false) {
        next()
    }
    else {
      
        if (to.meta.requireAuth === true && authState.isLoggedIn === false) {
           
            
            return next({ path: "/auth" })
        }

        const token = localStorage.getItem('token');
    
        if (token) {
          const decoded = jwtDecode(token);
          const isExpired = decoded.exp < Date.now() / 1000;
         
          if (isExpired) {
            console.log("token has expired", isExpired)
            const auth = authState
            auth.logout()
            return next({ path: "/auth" })
          }
         
        }
        
        next()
    }
 
    
  })

export default router;