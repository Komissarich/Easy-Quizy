<template>
    <div class="play-container">
      <form class="play-form">
        <div class="welcome-header">
        <h1>Введите код для прохождения квиза</h1>
      </div>
        <div class="form-group">
          <label for="quiz_id">Код</label>
          <input 
            type="text" 
            id="quiz_id" 
            v-model="quiz_id"
            placeholder="Введите id игры"
          >
        </div>
        </form>

        <button type="submit" class="login-btn" @click.prevent="startQuiz">
          Играть
        </button>

        
    
        </div>
        <div v-if="errorMessage !== ''" class="error-message">
      {{ errorMessage }} 
    </div>

</template>

<script>
import axios from 'axios'

export default {
    name: 'Play',
  data() {
    return {
      quiz_id: "",
      errorMessage: ""
    }
  },
  
  methods: {

    async startQuiz() {
      console.log("Current quiz_id:", this.quiz_id)
      if (this.quiz_id !== "") {
        this.errorMessage = ''
      try {
        console.log(this.quiz_id)
        let data = await axios.get(`http://localhost:8085/v1/quiz/${encodeURIComponent(this.quiz_id)}`)
        localStorage.setItem("isPlay", "true")
       
        this.$router.push("/play/" + data.data.quiz_id) 
      } catch (error) {
       
        this.errorMessage = 'Викторина с таким id не найдена'
    }
      }}
  }
}
</script>

<style scoped>
.play-container {
  max-width: 500px;
  margin: 0 auto;
  padding: 2rem;
  font-family: 'Arial', sans-serif;
  text-align: center;
}

.error-message {
    max-width: 400px;
  color: #f44336;
  background-color: #ffebee;
  padding: 12px;
  border-radius: 4px;
  margin-bottom: 20px;
  margin: 0 auto;
  border-left: 4px solid #f44336;
  text-align: center;
}
.login-btn {
  width: 100%;
  padding: 0.75rem;
  background-color: #4CAF50;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 1rem;
  cursor: pointer;
  transition: background-color 0.3s;
}
.reg-container {
  max-width: 400px;
  margin: 0 auto;
  padding: 2rem;
  font-family: 'Arial', sans-serif;
  text-align: center;
}
.welcome-header h1 {
  font-size: 2rem;
  color: #2c3e50;
  margin-bottom: 2rem;
  line-height: 1.3;
}

.play-form {
  background: #f8f9fa;
  padding: 2rem;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
}

.form-group {
  margin-bottom: 1.5rem;
  text-align: left;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: 500;
  color: #333;
}

.form-group input {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 1rem;
}

.remember-me {
  display: flex;
  align-items: center;
  margin: 1.5rem 0;
}


</style>