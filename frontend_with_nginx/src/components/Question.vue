<template>
    <div class="question-editor">
      <div class="question-header">
        <h2>Вопрос №{{ questionNumber }}</h2>
        <div class="question-actions">
          <button class="btn-icon danger" @click="$emit('delete-question')" title="Удалить вопрос">
            <span>×</span>
          </button>
        </div>
      </div>
  
      <div class="question-content">
        <textarea 
          :value="modelValue.question_text"
          @input="updateQuestion('question_text', $event.target.value)"
          placeholder="Введите текст вопроса" 
          class="question-textarea"
        ></textarea>
      </div>
  
      <div class="answers-list">
        <div 
          v-for="(answer, index) in modelValue.answer" 
          :key="index" 
          class="answer-item"
        >
          <div class="answer-content">
            <span class="answer-number">Ответ №{{ index + 1 }}</span>
            <input
              type="text"
              :value="answer.answer_text"
              @input="updateAnswer(index, 'answer_text', $event.target.value)"
              placeholder="Введите текст ответа"
              class="answer-input"
            >
            <input
              type="checkbox"
              :checked="answer.is_correct"
              @change="updateAnswer(index, 'is_correct', $event.target.checked)"
              class="result-input"
            >
            <label class="right-option">
              Правильный ответ?
            </label>
          </div>
          <div class="answer-actions">
            <button class="btn-icon" @click="removeAnswer(index)" title="Удалить ответ">
              <span>×</span>
            </button>
          </div>
        </div>
  
        <button class="add-answer-btn" @click="addAnswer">
          Добавить ответ
        </button>
      </div>
    </div>
  </template>
  
  <script setup>
  const props = defineProps({
    questionNumber: Number,
    modelValue: Object
  })
  
  const emit = defineEmits(['update:modelValue', 'delete-question'])
  
  // Обновление текста вопроса
  const updateQuestion = (field, value) => {
    emit('update:modelValue', {
      ...props.modelValue,
      [field]: value
    })
  }
  
  // Обновление ответов
  const updateAnswer = (answerIndex, field, value) => {
    const updatedAnswers = [...props.modelValue.answer]
    
    if (field === 'is_correct' && value === true) {
    // Сбрасываем все другие чекбоксы
    updatedAnswers.forEach((answer, idx) => {
      if (idx !== answerIndex) {
        answer.is_correct = false
      }
    })
  }
  
  // Обновляем текущий ответ
  updatedAnswers[answerIndex] = {
    ...updatedAnswers[answerIndex],
    [field]: value
  }
    
    emit('update:modelValue', {
      ...props.modelValue,
      answer: updatedAnswers
    })
  }
  
  // Добавление нового ответа
  const addAnswer = () => {
    const updatedAnswers = [
      ...props.modelValue.answer,
      { answer_text: '', is_correct: false }
    ]
    
    emit('update:modelValue', {
      ...props.modelValue,
      answer: updatedAnswers
    })
  }
  
  // Удаление ответа
  const removeAnswer = (index) => {
    if (props.modelValue.answer.length > 2) {
      const updatedAnswers = [...props.modelValue.answer]
      updatedAnswers.splice(index, 1)
      
      emit('update:modelValue', {
        ...props.modelValue,
        answer: updatedAnswers
      })
    } else {
      alert('Должен быть хотя бы 2 ответа')
    }
  }
  </script>
  
  <style scoped>
 .question-editor {
  width: 100%;
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  padding: 1.5rem;
  background-color: #fff;
}

.question-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
}
  
  .question-header h2 {
    margin: 0;
    font-size: 1.5rem;
    color: #333;
  }
  
  .question-actions {
    display: flex;
    gap: 0.5rem;
  }
  
  .btn-icon {
    width: 32px;
    height: 32px;
    border: 1px solid #ddd;
    border-radius: 4px;
    background: none;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.2s;
  }
  
  .btn-icon:hover {
    background-color: #f5f5f5;
  }
  
  .btn-icon.danger:hover {
    background-color: #ffebee;
  }
  
  .question-content {
    margin-bottom: 1.5rem;
  }
  
  .question-textarea {
    width: 96%;
    padding: 1rem;
    border: 1px solid #ddd;
    border-radius: 4px;
    min-height: 80px;
    font-size: 1rem;
    resize: vertical;
  }
  
  .answers-list {
    margin-bottom: 1.5rem;
  }
  
  .answer-item {
    display: flex;
    align-items: center;
    margin-bottom: 1rem;
    padding: 0.75rem;
    border: 1px solid #eee;
    border-radius: 4px;
  }
  
  .answer-content {
    flex-grow: 1;
    display: flex;
    align-items: center;
    gap: 1rem;
  }
  
  .answer-number {
    color: #666;
    min-width: 60px;
  }
  
  .answer-input, .result-input {
    flex-grow: 1;
    padding: 0.5rem;
    border: 1px solid #ddd;
    border-radius: 4px;
  }
  
  .result-input {
    max-width: 150px;
  }
  
  .answer-actions {
    margin-left: 1rem;
  }
  
  .add-answer-btn {
    padding: 0.5rem 1rem;
    background-color: #f5f5f5;
    border: 1px dashed #ccc;
    border-radius: 4px;
    cursor: pointer;
    color: #666;
    transition: all 0.2s;
  }
  
  .add-answer-btn:hover {
    background-color: #e0e0e0;
  }
  


  .question-settings {
    margin-top: 1.5rem;
    padding-top: 1rem;
    border-top: 1px solid #eee;
  }
  
  .setting-option {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    color: #666;
  }
  .right-option {
    display: flex;
    align-items: center;
    margin-left: -9%;
    color: #666;
  }
  </style>