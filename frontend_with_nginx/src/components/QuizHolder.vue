<template>
    <div class="quizzes-page">
      <h1 class="page-title">Все квизы</h1>
      
      <div v-if="loading" class="loading-spinner">
        <div class="spinner"></div>
      </div>
  
      <div v-else-if="quizzes.length === 0" class="empty-state">
       
        <p>Пока нет доступных квизов</p>
      </div>
  
      <div v-else class="quizzes-grid">
        <div 
          v-for="quiz in quizzes" 
          :key="quiz.id" 
          class="quiz-card"
          @click="$router.push(`/quiz/${quiz.shortID}`)"
        >
          <div class="quiz-card__image-container">
            <img 
                :src="`${quiz.imageId}`" 
              class="quiz-card__image"
              alt="Quiz cover"
            />
            <div class="quiz-card__overlay">
              <span class="quiz-card__questions-count">
                {{ quiz.question.length || 0 }} вопросов
              </span>
             
            </div>
          </div>
          
          <div class="quiz-card__content">
            <div class="quiz-card__header">
              <h3 class="quiz-card__title">{{ quiz.name }}</h3>
              <div class="quiz-card__author">
                <router-link 
                  :to="`/profile/${quiz.author}`" 
                  class="author-link"
                  @click.stop
                >
                  @{{ quiz.author }}
                </router-link>
              </div>
            </div>
            
            <p class="quiz-card__description">
              {{ truncateDescription(quiz.description) }}
            </p>
            
            <div class="quiz-card__footer">
              <div class="quiz-card__tags">
                <span 
                  v-for="tag in quiz.tags?.slice(0, 3)" 
                  :key="tag" 
                  class="tag"
                >
                  {{ tag }}
                </span>
              </div>
              <div class="quiz-card__actions">
                <button 
                  class="play-button"
                  @click.stop="startQuiz(quiz.shortID)"
                >
                  Играть
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </template>
  
  <script>
  import axios from 'axios';
  
  export default {
    name: 'QuizzesPage',
    data() {
      return {
        loading: true,
        quizzes: []
      }
    },
    methods: {
      async fetchQuizzes() {
        try {
          const response = await axios.get('http://localhost:8085/v1/quiz/orderby');
          this.quizzes = response.data.quizzes || [];
        } catch (error) {
          console.error('Ошибка загрузки квизов:', error);
        } finally {
          this.loading = false;
        }
      },
      truncateDescription(desc) {
        if (!desc) return 'Описание отсутствует';
        return desc.length > 100 ? desc.substring(0, 100) + '...' : desc;
      },
      startQuiz(quizId) {
        this.$router.push(`/play/${quizId}`);
      }
    },
    mounted() {
      this.fetchQuizzes();
    }
  }
  </script>
  
  <style scoped>
  .quizzes-page {
    max-width: 1200px;
    margin: 0 auto;
    padding: 2rem;
  }
  
  .page-title {
    font-size: 2rem;
    margin-bottom: 2rem;
    color: #2c3e50;
    text-align: center;
  }
  
  .quizzes-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
    gap: 2rem;
    margin-top: 1.5rem;
  }
  
  .quiz-card {
    background: white;
    border-radius: 12px;
    overflow: hidden;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
    transition: transform 0.2s, box-shadow 0.2s;
    cursor: pointer;
  }
  
  .quiz-card:hover {
    transform: translateY(-5px);
    box-shadow: 0 10px 20px rgba(0, 0, 0, 0.1);
  }
  
  .quiz-card__image-container {
    position: relative;
    height: 200px;
    overflow: hidden;
  }
  
  .quiz-card__image {
    width: 100%;
    height: 100%;
    object-fit: cover;
    transition: transform 0.3s;
  }
  
  .quiz-card:hover .quiz-card__image {
    transform: scale(1.05);
  }
  
  .quiz-card__overlay {
    position: absolute;
    bottom: 0;
    left: 0;
    right: 0;
    padding: 1rem;
    background: linear-gradient(transparent, rgba(0, 0, 0, 0.7));
    color: white;
    display: flex;
    justify-content: space-between;
  }
  
  .quiz-card__content {
    padding: 1.5rem;
  }
  
  .quiz-card__header {
    margin-bottom: 1rem;
  }
  
  .quiz-card__title {
    font-size: 1.25rem;
    margin: 0 0 0.5rem 0;
    color: #2c3e50;
  }
  
  .quiz-card__author {
    font-size: 0.9rem;
    color: #7f8c8d;
  }
  
  .author-link {
    color: #3498db;
    text-decoration: none;
  }
  
  .author-link:hover {
    text-decoration: underline;
  }
  
  .quiz-card__description {
    color: #7f8c8d;
    margin-bottom: 1.5rem;
    line-height: 1.5;
  }
  
  .quiz-card__footer {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
  
  .quiz-card__tags {
    display: flex;
    gap: 0.5rem;
  }
  
  .tag {
    background: #f1f1f1;
    padding: 0.25rem 0.5rem;
    border-radius: 4px;
    font-size: 0.8rem;
    color: #34495e;
  }
  
  .play-button {
    background: #3498db;
    color: white;
    border: none;
    padding: 0.5rem 1rem;
    border-radius: 4px;
    cursor: pointer;
    transition: background 0.2s;
  }
  
  .play-button:hover {
    background: #2980b9;
  }
  
  .loading-spinner {
    display: flex;
    justify-content: center;
    padding: 2rem;
  }
  
  .spinner {
    width: 50px;
    height: 50px;
    border: 5px solid #f3f3f3;
    border-top: 5px solid #3498db;
    border-radius: 50%;
    animation: spin 1s linear infinite;
  }
  
  @keyframes spin {
    0% { transform: rotate(0deg); }
    100% { transform: rotate(360deg); }
  }
  
  .empty-state {
    text-align: center;
    padding: 3rem;
    color: #7f8c8d;
  }
  
  .empty-image {
    width: 200px;
    opacity: 0.7;
    margin-bottom: 1rem;
  }
  </style>