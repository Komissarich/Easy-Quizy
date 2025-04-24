<template>
    <div class="quiz-detail-container">
      <!-- Лоадер -->
      <div v-if="loading" class="loader">
        <div class="spinner"></div>
        <p>Загружаем данные квиза...</p>
      </div>
  
      <!-- Основной контент -->
      <div v-else-if="quiz" class="quiz-content">
        <!-- Хлебные крошки и кнопка назад -->
        <div class="navigation">
          <button @click="$router.go(-1)" class="back-button">
            ← Назад
          </button>
        </div>
  
        <!-- Заголовок и автор -->
        <div class="quiz-header">
          <h1 class="quiz-title">{{ quiz.title }}</h1>
          <p class="quiz-author">Автор: {{ quiz.author || 'Неизвестен' }}</p>
        </div>
  
        <!-- Основная информация -->
        <div class="quiz-main">
          <!-- Большое изображение -->
          <div class="quiz-image-container">
            <img 
            :src="`${quiz.imageId}`" 
              alt="Обложка квиза"
              class="quiz-image"
            />
          </div>
  
          <!-- Правая колонка с информацией -->
          <div class="quiz-info">
            <!-- Статистика -->
            <div class="stats-section">
              <h3>Статистика квиза</h3>
              <div class="stats-grid">
                <div class="stat-card">
                  <div class="stat-value">{{ stats.success_rate || 0 }}%</div>
                  <div class="stat-label">Успешных прохождений</div>
                </div>
                <div class="stat-card">
                  <div class="stat-value">{{ stats.average_score || 0 }}/10</div>
                  <div class="stat-label">Средняя оценка</div>
                </div>
                <div class="stat-card">
                  <div class="stat-value">{{ quiz.questions_count || 0 }}</div>
                  <div class="stat-label">Вопросов</div>
                </div>
              </div>
            </div>
  
            <!-- Прогресс-бар для визуализации успешности -->
            <div class="progress-section">
              <div class="progress-label">
                <span>Успешность прохождений:</span>
                <span>{{ stats.success_rate || 0 }}%</span>
              </div>
              <div class="progress-bar">
                <div 
                  class="progress-fill"
                  :style="{ width: `${stats.success_rate || 0}%` }"
                ></div>
              </div>
            </div>
  
            <!-- Описание -->
            <div class="description-section">
              <h3>Описание</h3>
              <p>{{ quiz.description || 'Этот квиз пока не имеет описания' }}</p>
            </div>
  
            <!-- Кнопки действий -->
            <div class="actions-section">
              <button 
                @click="toggleFavorite"
                :class="['favorite-button', { 'active': isFavorite }]"
              >
                {{ isFavorite ? '★ В избранном' : '☆ Добавить в избранное' }}
              </button>
              <button 
                @click="startQuiz"
                class="play-button"
              >
                Играть →
              </button>
            </div>
          </div>
        </div>
      </div>
  
      <!-- Если квиз не найден -->
      <div v-else class="not-found">
        <h2>Квиз не найден</h2>
        <p>Возможно, он был удален или вы указали неверный ID</p>
        <button @click="$router.push('/')" class="home-button">
          На главную
        </button>
      </div>
    </div>
  </template>
  
  <script>
  import axios from 'axios'
  import { ref, onMounted } from 'vue'
  import { useRoute, useRouter } from 'vue-router'
  
  export default {
    name: 'ShowQuiz',
    setup() {
      const route = useRoute()
      const router = useRouter()
      const quiz = ref(null)
      const stats = ref({
        success_rate: 0,
        average_score: 0
      })
      const loading = ref(true)
      const isFavorite = ref(false)
      const quiz_id = route.params.id
  
      // Функция загрузки данных квиза
      const fetchQuizData = async () => {
        try {
          // 1. Загрузка основной информации о квизе
          let data = await axios.get(`http://localhost:8085/v1/quiz/${quiz_id}`)
          quiz.value = quizResponse.data
  
          // 2. Загрузка статистики
          const statsResponse = await axios.get(`http://localhost:8085/v1/stats/quiz/${quiz_id}`)
          console.log(statsResponse.data)
          stats.value = statsResponse.data
  
          // 3. Проверка, находится ли квиз в избранном
          const favoriteResponse = await axios.get(`http://localhost:8085/v1/users/favorites/quizzes`, {
            headers: {
              'Authorization': `Bearer ${localStorage.getItem('token')}`
            }
          })
          console.log(favoriteResponse.data)
          isFavorite.value = favoriteResponse.data.isFavorite
  
        } catch (error) {
          console.error('Ошибка загрузки данных:', error)
        } finally {
          loading.value = false
        }
      }
  
      // Переключение избранного
      const toggleFavorite = async () => {
        try {
          if (isFavorite.value) {
            await axios.delete(`http://localhost:8085/v1/favorites/${quiz_id}`, {
              headers: {
                'Authorization': `Bearer ${localStorage.getItem('token')}`
              }
            })
          } else {
            await axios.post(`http://localhost:8085/v1/favorites`, 
              { quiz_id: quiz_id },
              {
                headers: {
                  'Authorization': `Bearer ${localStorage.getItem('token')}`
                }
              }
            )
          }
          isFavorite.value = !isFavorite.value
        } catch (error) {
          console.error('Ошибка обновления избранного:', error)
          alert('Не удалось обновить избранное')
        }
      }
  
      // Начать квиз
      const startQuiz = () => {
        router.push("/play/" + quiz_id)
      }
  
      // Загружаем данные при монтировании компонента
      onMounted(fetchQuizData)
  
      return {
        quiz,
        stats,
        loading,
        isFavorite,
        quiz_id,
        toggleFavorite,
        startQuiz
      }
    }
  }
  </script>
  
  <style scoped>
  .quiz-detail-container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 20px;
    font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
  }
  
  /* Лоадер */
  .loader {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    height: 300px;
  }
  
  .spinner {
    width: 50px;
    height: 50px;
    border: 5px solid #f3f3f3;
    border-top: 5px solid #4a6fa5;
    border-radius: 50%;
    animation: spin 1s linear infinite;
    margin-bottom: 20px;
  }
  
  @keyframes spin {
    0% { transform: rotate(0deg); }
    100% { transform: rotate(360deg); }
  }
  
  /* Навигация */
  .navigation {
    margin-bottom: 20px;
  }
  
  .back-button {
    background: none;
    border: none;
    color: #4a6fa5;
    font-size: 1rem;
    cursor: pointer;
    padding: 5px 10px;
    border-radius: 5px;
    transition: background-color 0.2s;
  }
  
  .back-button:hover {
    background-color: #f0f0f0;
  }
  
  /* Заголовок */
  .quiz-header {
    margin-bottom: 30px;
  }
  
  .quiz-title {
    font-size: 2.2rem;
    color: #333;
    margin-bottom: 5px;
  }
  
  .quiz-author {
    font-size: 1rem;
    color: #666;
  }
  
  /* Основной контент */
  .quiz-main {
    display: flex;
    gap: 40px;
  }
  
  .quiz-image-container {
    flex: 1;
    min-width: 0;
  }
  
  .quiz-image {
    width: 100%;
    max-height: 400px;
    object-fit: cover;
    border-radius: 10px;
    box-shadow: 0 4px 15px rgba(0, 0, 0, 0.1);
  }
  
  .quiz-info {
    flex: 1;
    min-width: 0;
  }
  
  /* Статистика */
  .stats-section {
    margin-bottom: 25px;
  }
  
  .stats-section h3 {
    font-size: 1.2rem;
    margin-bottom: 15px;
    color: #444;
  }
  
  .stats-grid {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 15px;
  }
  
  .stat-card {
    background: #f8f9fa;
    border-radius: 8px;
    padding: 15px;
    text-align: center;
  }
  
  .stat-value {
    font-size: 1.5rem;
    font-weight: bold;
    color: #4a6fa5;
    margin-bottom: 5px;
  }
  
  .stat-label {
    font-size: 0.9rem;
    color: #666;
  }
  
  /* Прогресс-бар */
  .progress-section {
    margin-bottom: 25px;
  }
  
  .progress-label {
    display: flex;
    justify-content: space-between;
    margin-bottom: 8px;
    font-size: 0.9rem;
    color: #555;
  }
  
  .progress-bar {
    height: 10px;
    background: #e9ecef;
    border-radius: 5px;
    overflow: hidden;
  }
  
  .progress-fill {
    height: 100%;
    background: linear-gradient(to right, #4facfe, #00f2fe);
    border-radius: 5px;
    transition: width 0.5s ease;
  }
  
  /* Описание */
  .description-section {
    margin-bottom: 25px;
  }
  
  .description-section h3 {
    font-size: 1.2rem;
    margin-bottom: 10px;
    color: #444;
  }
  
  .description-section p {
    line-height: 1.6;
    color: #555;
  }
  
  /* Кнопки действий */
  .actions-section {
    display: flex;
    gap: 15px;
    margin-top: 30px;
  }
  
  .favorite-button {
    flex: 1;
    padding: 12px;
    background: #f8f9fa;
    border: 1px solid #ddd;
    border-radius: 8px;
    font-size: 1rem;
    color: #555;
    cursor: pointer;
    transition: all 0.2s;
  }
  
  .favorite-button:hover {
    background: #e9ecef;
  }
  
  .favorite-button.active {
    background: #fff3bf;
    border-color: #ffe69c;
    color: #e67700;
  }
  
  .play-button {
    flex: 1;
    padding: 12px;
    background: #4a6fa5;
    border: none;
    border-radius: 8px;
    font-size: 1rem;
    color: white;
    cursor: pointer;
    transition: background-color 0.2s;
  }
  
  .play-button:hover {
    background: #3a5a8f;
  }
  
  /* Не найдено */
  .not-found {
    text-align: center;
    padding: 50px 0;
  }
  
  .not-found h2 {
    font-size: 1.8rem;
    color: #333;
    margin-bottom: 10px;
  }
  
  .not-found p {
    font-size: 1.1rem;
    color: #666;
    margin-bottom: 20px;
  }
  
  .home-button {
    padding: 10px 20px;
    background: #4a6fa5;
    border: none;
    border-radius: 5px;
    color: white;
    font-size: 1rem;
    cursor: pointer;
  }
  
  /* Адаптивность */
  @media (max-width: 768px) {
    .quiz-main {
      flex-direction: column;
    }
    
    .stats-grid {
      grid-template-columns: 1fr;
    }
    
    .actions-section {
      flex-direction: column;
    }
  }
  </style>