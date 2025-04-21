<template>
  <div class="quiz-container">
    <h1>{{ quiz.title }}</h1>
 
    
    <div v-if="loading" class="loading">
      Загрузка квиза...
    </div>
    
    
    <template v-else>
      <QuizQuestion
        :question="currentQuestion"
        :question-index="currentQuestionIndex"
        :question-number="currentQuestionIndex + 1"
        :total-questions="quiz.questions.length"
        :selected-answer="selectedAnswer"
        @answer-selected="selectAnswer"
      />
      
      <button 
        class="next-button"
        @click="nextQuestion"
        :disabled="selectedAnswer === null"
      >
        {{ isLastQuestion ? 'Завершить квиз' : 'Следующий вопрос' }}
      </button>
    </template>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import axios from 'axios'
import QuizQuestion from './QuizQuestion.vue'
import router from '@/router'


const route = useRoute()

const quiz = ref({})
const loading = ref(true)
const currentQuestionIndex = ref(0)
const selectedAnswer = ref(null)
const userAnswers = ref([])
// Получаем данные квиза при загрузке компонента
onMounted(async () => {
  try {

     
    let response = await axios.get("http://localhost:8080/get_quiz", {params: {quiz_id:route.params.quiz_id}})
    
    console.log(route.params.quiz_id)
    quiz.value = response.data
  } catch (error) {
    console.error('Ошибка загрузки квиза:', error)
  } finally {
    loading.value = false
  }
})

const props = defineProps({
  quizId: String
})



const currentQuestion = computed(() => {
  return quiz.value.questions?.[currentQuestionIndex.value]
})

const isLastQuestion = computed(() => {
  return currentQuestionIndex.value === quiz.value.questions?.length - 1
})

const selectAnswer = (index) => {
  selectedAnswer.value = index
}

const nextQuestion = () => {
  userAnswers.value.push({question: quiz.value.questions?.[currentQuestionIndex.value], user_answer: selectedAnswer.value})
  if (isLastQuestion.value) {
    sessionStorage.setItem('quizResults', JSON.stringify({
    userAnswers: userAnswers.value,
    totalQuestions: quiz.value.questions.length
  }))
  
  router.push({
    name: 'QuizResult',
    params: {
      quiz_id: route.params.quiz_id
    }
  })
  } else {
    

    currentQuestionIndex.value++
    selectedAnswer.value = null
  }
}

</script>


<style scoped>
.quiz-container {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
}

.loading {
  padding: 40px;
  text-align: center;
}

.next-button {
  display: block;
  width: 100%;
  padding: 12px;
  background-color: #42b983;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 1rem;
  cursor: pointer;
  margin-top: 20px;
}

.next-button:disabled {
  background-color: #cccccc;
  cursor: not-allowed;
}
</style>