<template>
    <div class="result-container">
      <div class="result-header">
        <h1>РЕЗУЛЬТАТЫ</h1>
        <div class="result-summary">
          Правильных ответов: {{ correctCount }} из {{ totalQuestions }}
        </div>
      </div>
  
      <div class="results-list">
        <div v-for="(result, index) in results" :key="index" class="result-item">
          <div class="question">
            <span class="question-number">Вопрос {{ index + 1 }}:</span>
            {{ result.question.questionText }}
          </div>
          <div class="user-answer" :class="{'correct': result.isCorrect, 'incorrect': !result.isCorrect}">
            Ваш ответ: {{ getUserAnswerText(result) }}
            <span v-if="result.isCorrect" class="result-icon">✓</span>
            <span v-else class="result-icon">✗</span>
          </div>
          <div v-if="!result.isCorrect" class="correct-answer">
            Правильный ответ: {{ getCorrectAnswerText(result.question) }}
          </div>
        </div>
      </div>
  
      <!-- Блок оценки квиза -->
      <div v-if="errorMessage !== 'false'" class="rating-section">
        <h3>Оцените этот квиз (ваша статистика отправится с оценкой):</h3>
        <div class="stars">
            <span 
          v-for="star in 5" 
          :key="star"
          class="star"
          :class="{ 'active': star <= (hoverRating !== null ? hoverRating : currentRating) }" 
          @click="rateQuiz(star)"
          @mouseover="hoverRating = star"
          @mouseleave="hoverRating = null"
        >
          ★
        </span>
        </div>

        <p v-if="ratingSubmitted" class="rating-thanks">Спасибо за оценку!</p>
        <button class="add-button" @click="addFavorite"> 
            Добавить в избранное
        </button>
        
        <p v-if="addedFavorite" class="rating-thanks">Добавлено в избранное</p>
      </div>
      <div v-else class="login-prompt">
        <p>Хотите оценить квиз? <a @click="navigateToLogin">Войдите в аккаунт</a></p>
      </div>
  
      <div v-if="errorMessage !== ''" class="error-message">
        {{ errorMessage }}
      </div>
  
      <button @click="returnToQuizzes" class="return-button">
        Вернуться к выбору квиза
      </button>
    </div>
  </template>
  
  <script setup>
  import axios from 'axios'
  import { ref, onMounted } from 'vue'
  import { useRoute, useRouter } from 'vue-router'
import Quiz from './Quiz.vue'
  
  const router = useRouter()
  const userAnswers = ref([])
  const totalQuestions = ref(0)
  const correctCount = ref(0)
  const results = ref([])
  const resultsData = ref(null)
  const errorMessage = ref('')
  const currentRating = ref(0)
  const hoverRating = ref(null)
  const ratingSubmitted = ref(false)
  const addedFavorite = ref(false)
  onMounted(() => {
    
    const savedResults = sessionStorage.getItem('quizResults')
  if (savedResults) {
    resultsData.value = JSON.parse(savedResults)
    sessionStorage.removeItem('quizResults')
    userAnswers.value = resultsData.value.userAnswers
    totalQuestions.value = resultsData.value.totalQuestions
    console.log('Results:', resultsData.value)
  } else {
    console.error('No results found')
    
  }
      processResults()
    
  })

  const rateQuiz = async (rating) => {
 
  if (ratingSubmitted.value === false) {
    currentRating.value = rating
    ratingSubmitted.value = true
    // string quiz_id = 1;
    // map<string, float> players_score = 2;
    // float quiz_rate = 3; post: "/v1/stats/update"
    const username =  localStorage.getItem("username")
    console.log(correctCount.value, totalQuestions.value)
    console.log(parseFloat((3/20) * 100).toFixed(2))
    let data = await axios.post(`http://localhost:8085/v1/stats/update`, 
            {
              quiz_rate: currentRating.value,
              quiz_id: router.currentRoute.value.params.quiz_id,
              player_id: username,
              player_score:parseFloat((correctCount.value/totalQuestions.value) * 100)
            },
            {
              headers: {
                'Authorization': `Bearer ${localStorage.getItem('token')}`
              },
            } 
           
             
           )
         console.log(data.data)
  }
  

}
  const addFavorite = async () => {
    addedFavorite.value = true
    try {
       let data = await axios.post(`http://localhost:8085/v1/users/favorites/quizzes/add`, 
            {
              token: localStorage.getItem('token'),
              quiz_id: router.currentRoute.value.params.quiz_id
            },
            {
              headers: {
              'Authorization': `Bearer ${localStorage.getItem('token')}`,
              'Content-Type': 'application/json', // Важно явно указать!
              },
            } 
           
             
           )
         console.log(data.data)
     } catch (error) {
       console.error('Ошибка добавления в избранное:', error)
     }
  // await axios.post("http://localhost:8080/rate_quiz", JSON.stringify({user_id: localStorage.getItem("username"), quiz_id: router.currentRoute.value.params.quiz_id}))
  }
  const processResults = async () => {
    errorMessage.value = ""
    results.value = userAnswers.value.map(answer => {
      const isCorrect = answer.question.answer[answer.user_answer].isCorrect
      if (isCorrect) correctCount.value++
        return {
            question: answer.question,
            user_answer: answer.user_answer,
            isCorrect
        }
    })
   
    if (localStorage.getItem("auth") === "false") {
        errorMessage.value = "Войдите в аккаунт для сбора статистики и оценки квиза!"
    }

    //await axios.post("http://localhost:8080/player_stat", JSON.stringify({user_id: localStorage.getItem("username"), correctAnswer: correctCount.value, totalquestion: totalQuestions.value}))
   
    
  }
  
  const getUserAnswerText = (result) => {
    return result.question.answer[result.user_answer].answerText
  }
  
  const getCorrectAnswerText = (question) => {
    const correctAnswer = question.answer.find(a => a.isCorrect)
    return correctAnswer.answerText
  }
  
  const returnToQuizzes = () => {
    router.push('/play')
  }
  </script>
  
  <style scoped>

.rating-section {
  margin: 30px 0;
  text-align: center;
  padding: 20px;
  background: #f8f9fa;
  border-radius: 8px;
  
}

.stars {
  font-size: 2rem;
  margin: 10px 0;
}

.star {
  color: #ccc;
  cursor: pointer;
  transition: color 0.2s;
  margin: 0 5px;
}

.star.active {
  color: #ffc107;
}

.star:hover {
  color: #ffc107;
}

.rating-thanks {
  color: #727272;
  font-weight: bold;
}
  .result-container {
    font-family: 'Arial', sans-serif;
    max-width: 800px;
    margin: 0 auto;
    padding: 20px;
  }

  .error-message {
    font-family: 'Arial', sans-serif;
    max-width: 400px;
  color: #0a0a0a;
  background-color: #e2ff92;
  padding: 12px;
  border-radius: 4px;
  margin-bottom: 20px;
  margin: 0 auto;
  border-left: 4px solid #d6ff20;
  text-align: center;
}
  
  .result-header {
    text-align: center;
    margin-bottom: 30px;
  }
  
  .result-header h1 {
    font-size: 2.5rem;
    color: #2c3e50;
    margin-bottom: 10px;
  }
  
  .result-summary {
    font-size: 1.3rem;
    font-weight: bold;
    color: #42b983;
  }
  
  .results-list {
    margin-top: 30px;
  }
  
  .result-item {
    background: white;
    border-radius: 8px;
    padding: 20px;
    margin-bottom: 20px;
    box-shadow: 0 2px 8px rgba(0,0,0,0.1);
  }
  
  .question {
    font-size: 1.1rem;
    margin-bottom: 15px;
  }
  
  .question-number {
    font-weight: bold;
    color: #2c3e50;
  }
  
  .user-answer {
    padding: 10px;
    border-radius: 4px;
    margin-bottom: 5px;
  }
  
  .user-answer.correct {
    background-color: rgba(66, 185, 131, 0.1);
    color: #42b983;
  }
  
  .user-answer.incorrect {
    background-color: rgba(255, 99, 71, 0.1);
    color: #ff6347;
  }
  
  .correct-answer {
    padding: 10px;
    background-color: rgba(66, 185, 131, 0.1);
    border-radius: 4px;
    color: #42b983;
  }
  
  .result-icon {
    margin-left: 10px;
    font-weight: bold;
  }
  
  .return-button {
    display: block;
    width: 100%;
    padding: 12px;
    background-color: #42b983;
    color: white;
    border: none;
    border-radius: 4px;
    font-size: 1rem;
    cursor: pointer;
    margin-top: 30px;
    text-align: center;
  }

  .add-button {
    display: flex;
    width: 50%;
    margin: auto;
    padding: 12px;
    background-color: #72c9fc;
    color: white;
    border: none;
    border-radius: 4px;
    font-size: 1rem;
    cursor: pointer;
    margin-top: 10px;
    align-items: center;
    text-align: center;
    justify-content: center;
  }
  </style>