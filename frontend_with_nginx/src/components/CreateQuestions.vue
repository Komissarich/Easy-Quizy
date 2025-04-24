<template>
   <div class="test-creator">
    <div class="question-container">
      <div class="questions-list">
        <Question 
          v-for="(q, index) in questions" 
          :key="q.id"
          :questionNumber="index + 1"
          @delete-question="removeQuestion(index)"
          class="question-item"
          v-model="questions[index]"
        />
      </div>
      
     
    </div>
    <button @click="addQuestion" class="add-question-btn">
        + Добавить вопрос
    </button>

    <button @click="SendQuiz" class="send-btn">
        Опубликовать квиз
    </button>
  </div>
  </template>
  
<script setup>
import { inject, ref } from 'vue';
import Question from './Question.vue';
import { useRoute } from 'vue-router'
import axios from 'axios'
const route = useRoute()
const testTitle = ref(route.query.title)
const testDescription = ref(route.query.description)
const testCover = ref(route.query.cover)
  let nextId = 1;
  const questions = ref([
  {
    id: 1,
    question_text: '',
    image_id: "",
    answer: [
      { answer_text: '', is_correct: false },
      { answer_text: '', is_correct: false }
    ]
  }
])
  
const addQuestion = () => {
  questions.value.push({ 
    id: nextId++,
    question_text: '',
    image_id: "",
    answer: [
      { answer_text: '', is_correct: false },
      { answer_text: '', is_correct: false }
    ]
  });
};

  const SendQuiz = async () =>{

  try {
    console.log(questions.value)
    const data = await axios.post(
      
      'http://localhost:8085/v1/quiz',
        {
          name: testTitle.value,
          author: localStorage.getItem("username"),
          image_id: testCover.value,
          description: testDescription.value,
          question: questions.value
        },
        {
          headers: {
            'Content-Type': 'application/json', // Важно явно указать!
          },
        }
      );
      console.log(data.data)
      console.log("Succesfully created quiz")
      } catch (error) {
        console.log(error.status, error)
        // if (error.status == 400) {
        //   this.errorMessage = 'Проверьте почту и пароль; пароль не менее 8 символов'
        // }
      }
  };
  
  const removeQuestion = (index) => {
  if (questions.value.length > 1) {
    questions.value.splice(index, 1)
  } else {
    alert('Должен остаться хотя бы один вопрос')
  }
}
  
  
  </script>
  
  <style scoped>
.test-creator {
  display: flex;
  justify-content: center;
  padding: 20px;
}

.question-container {
  position: relative;
  width: 100%;
  max-width: 1000px;
}

.questions-list {
  width: 80%;
  margin: 0 auto;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.add-question-btn {
  position: absolute;
  top: 1;
  right: 0;
  margin-right: 500px;
  margin-top: 30px;
  padding: 10px 20px;
  background-color: #4CAF50;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s;
  z-index: 10;
}

.send-btn {
    position: absolute;
  top: 1;
  right: 0;
  margin-right: 500px;
  margin-top: 80px;
  padding: 10px 20px;
  background-color: #5cb3ce;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s;
  z-index: 10;
  }
  

.add-question-btn:hover {
  background-color: #45a049;
}



  </style>