<template>
    <div class="question-container">
      <h2>Вопрос {{ questionNumber }} из {{ totalQuestions }}</h2>
      <p class="question-text">{{ question.text }}</p>
      
      <div class="answers-list">
        <div 
          v-for="(answer, index) in question.answers" 
          :key="index" 
          class="answer-option"
        >
          <input
            type="radio"
            :id="'answer-' + index"
            :name="'question-' + questionIndex"
            :checked="selectedAnswer === index"
            @change="$emit('answer-selected', index)"
            class="answer-input"
          >
          <label :for="'answer-' + index" class="answer-label">
            {{ answer.text }}
          </label>
        </div>
      </div>
    </div>
  </template>
  
  <script setup>
  defineProps({
    question: {
      type: Object,
      required: true
    },
    questionIndex: {
      type: Number,
      required: true
    },
    questionNumber: {
      type: Number,
      required: true
    },
    totalQuestions: {
      type: Number,
      required: true
    },
    selectedAnswer: {
      type: Number,
      default: null
    }
  })
  
  defineEmits(['answer-selected'])
  </script>
  
  <style scoped>
  .question-container {
    font-family: Arial, sans-serif;
    padding: 20px;
    border: 1px solid #eee;
    border-radius: 8px;
    margin-bottom: 20px;
  }
  
  .question-text {
    font-size: 1.2rem;
    margin: 15px 0;
  }
  
  .answers-list {
    display: flex;
    flex-direction: column;
    gap: 10px;
  }
  
  .answer-option {
    display: flex;
    align-items: center;
  }
  
  .answer-input {
    margin-right: 10px;
  }
  
  .answer-label {
    cursor: pointer;
    padding: 8px 12px;
    border-radius: 4px;
    transition: background 0.2s;
  }
  
  .answer-label:hover {
    background: #f5f5f5;
  }
  
  input[type="radio"]:checked + .answer-label {
    background: #e3f2fd;
    font-weight: bold;
  }
  </style>