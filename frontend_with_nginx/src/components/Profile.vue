<template>
  <div class="profile-container">
    <!-- Шапка профиля с расширенной статистикой -->
    <div class="profile-header">
      <div class="user-info">
        <h1 class="username">{{ user.username }}</h1>
        <p class="user-email">{{ user.email }}</p>
      </div>

      <div class="compact-stats">
    <!-- Первая строка -->
    <div class="stats-row">
     
      <div class="stat-item">
        <div class="stat-value">{{ player_stats.avgScore?.toFixed(1) || '0.0' }}</div>
        <div class="stat-label">Средний результат</div>
      </div>
      <div class="stat-item">
        <div class="stat-value">{{ player_stats.bestScore?.toFixed(1) || '0' }}</div>
        <div class="stat-label">Лучший результат</div>
      </div>
      <div class="stat-item">
        <div class="stat-value">{{ player_stats.numSessions || 0 }}</div>
        <div class="stat-label">Пройдено квизов</div>
      </div>
    </div>
    
    <!-- Вторая строка -->
    <div class="stats-row">
      <div class="stat-item">
        <div class="stat-value">{{ author_stats.avgQuizRate?.toFixed(1) || '0.0' }}</div>
        <div class="stat-label">Авторский рейтинг</div>
      </div>
      <div class="stat-item">
        <div class="stat-value">{{ author_stats.bestQuizRate?.toFixed(1) || '0.0' }}</div>
        <div class="stat-label">Лучший квиз</div>
      </div>
      <div class="stat-item">
        <div class="stat-value">{{ (userQuizzes[0]?.quizzes?.length ?? 0) || 0}}</div>
        <div class="stat-label">Создано квизов</div>
      </div>
    </div>
  </div>
</div>

    

    <div class="divider"></div>

    <div class="profile-content">
      <div class="tabs">
        <button 
          class="tab-btn" 
          :class="{ active: activeTab === 'tests' }"
          @click="activeTab = 'tests'"
        >
          Мои тесты
        </button>
        <button 
          class="tab-btn" 
          :class="{ active: activeTab === 'favorites' }"
          @click="activeTab = 'favorites'"
        >
          Избранное
        </button>
        <button 
          class="tab-btn" 
          :class="{ active: activeTab === 'friends' }"
          @click="activeTab = 'friends'"
        >
          Друзья ({{ friends.length }})
        </button>
      </div>

      <div class="tab-content">
        <!-- Мои тесты -->
        <div v-if="activeTab === 'tests'" class="tests-section">
    <p v-if="userQuizzes.length === 0" class="empty-message">
      Этот пользователь ещё не создал ни одного теста
    </p>
    <div v-else class="quiz-list">
      <div 
        v-for="quiz in userQuizzes[0].quizzes" 
        :key="quiz.id" 
        class="quiz-card"
        @click="$router.push(`/quiz/${quiz.shortID}`)"
      >
        <div class="quiz-image-container">
          <img 
            :src="`${quiz.imageId}`" 
            class="quiz-image"
            alt="Quiz cover"
          />
          <div class="quiz-overlay"></div>
        </div>
        <div class="quiz-info">
          <h3 class="quiz-title">{{ quiz.name }}</h3>
          <p class="quiz-description">{{ quiz.description || 'Описание отсутствует' }}</p>
        </div>
      </div>
    </div>
  </div>
        <!-- Избранное -->
        <div v-if="activeTab === 'favorites'" class="favorites-section">
          <p v-if="userQuizzes[1].quizzes.length === 0" class="empty-message">
      Нет добавленных в избранное тестов
    </p>
    <div v-else class="quiz-list">
      <div 
        v-for="quiz in userQuizzes[1].quizzes" 
        :key="quiz.id" 
        class="quiz-card"
        @click="$router.push(`/quiz/${quiz.shortID}`)"
      >
        <div class="quiz-image-container">
          <img 
              :src="`${quiz.imageId}`" 
            class="quiz-image"
            alt="Quiz cover"
          />
          <div class="quiz-overlay"></div>
        </div>
        <div class="quiz-info">
          <h3 class="quiz-title">{{ quiz.name }}</h3>
          
        </div>
        </div>
        </div>
      </div>
        
        <!-- Друзья -->
        <div v-if="activeTab === 'friends'" class="friends-section">
         
          
          <div v-if="friends.length > 0" class="friends-list">
            <div v-for="friend in friends" :key="friend.id" class="friend-card">
              
              <span class="friend-name">{{ friend.username }}</span>
              <button class="remove-friend-btn" @click="removeFriend(friend)">
                ×
              </button>
            </div>
          </div>
          <p v-else class="empty-message">
            У вас пока нет друзей
          </p>
        </div>
      </div>
    </div>

   
  </div>
</template>
<script>
import axios from 'axios'
import { onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'

    export default {
    name: 'Auth',
  data() {
    return {
      activeTab: 'tests',
      showAddFriendModal: false,
     // defaultAvatar: '/default-avatar.png',
      friendSearchQuery: '',
      searchResults: [],
      user: {
        id: 1,
        username: localStorage.getItem("username"),
        email: localStorage.getItem("email"),
        avatar: null
      },
      userQuizzes: [],
      favorites: [],
     
    
        username: localStorage.getItem("username"),
        password: '',
      
        activeTab : ref('tests')
    }
  },
  methods: {

   
  
  
    async removeFriend(friend) {
      console.log("REMOVING ", friend.username)
      this.friends = this.friends.filter(f => f.id !== friend.id);

      let friend_data = await axios.post(`http://localhost:8085/v1/users/friends/remove`, 
      {
        token: localStorage.getItem('token'),
        friend_id: friend.username
      },
        {
          headers: {
            'Authorization': `Bearer ${localStorage.getItem('token')}`
          },
        })
        console.log(friend_data)
    }
  },
  setup() {
    
    const route = useRoute()
    const quiz = ref(null)
    const loading = ref(true)
    const userQuizzes = ref([])
    const friends = ref([  // Теперь это reactive ref
      { id: 1, username: '' },
      { id: 2, username: '' }
    ]);
    const player_stats = ref({
      bestScore: 0,
      avgScore: 0,
      numSessions: 0
    })
    const author_stats = ref({
      avgQuizRate: 0,
      bestQuizRate: 0,
      numQuizzes: 0
    })
      onMounted(async () => {
        try {
      
          let friend_data = await axios.post(`http://localhost:8085/v1/user/friends`, 
              {
                token: localStorage.getItem('token')
              },
               {
                  headers: {
                    'Authorization': `Bearer ${localStorage.getItem('token')}`
                  },
                })
                console.log(friend_data.data) 
               friends.value = friend_data.data.friends
                
        
          let data = await axios.get(`http://localhost:8085/v1/quiz/author/${localStorage.getItem('username')}`,  {
                headers: {
                  'Authorization': `Bearer ${localStorage.getItem('token')}`
                },
              })
           
            userQuizzes.value = data.data.authorQuizzes

            let player_data = await axios.get(`http://localhost:8085/v1/stats/player/${localStorage.getItem('username')}`)  

            console.log("PLAYER STAT",player_data.data)
            player_stats.value = player_data.data.player

            let author_data = await axios.get(`http://localhost:8085/v1/stats/author/${localStorage.getItem('username')}`)  

          console.log("AUTHOR STAT",author_data.data)
          author_stats.value = author_data.data.author

        } catch (error) {
          console.error('Ошибка загрузки профиля:', error, error.data)
        } finally {
          loading.value = false
        }
      })

      return {
        userQuizzes,
        quiz,
        loading,
        author_stats,
        player_stats,
        friends
      }
    }
}

</script>
  
  <style scoped>
  .profile-container {
    max-width: 800px;
    margin: 0 auto;
    padding: 2rem;
    font-family: 'Arial', sans-serif;
  }
  
  .profile-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 0;
  flex-wrap: wrap;
  gap: 20px;
}

.user-info {
  flex: 1;
  min-width: 200px;
}

.username {
  font-size: 2rem;
  margin: 0;
  color: #333;
}

.user-email {
  font-size: 1rem;
  color: #666;
  margin: 5px 0 0 0;
}

.user-stats {
  display: flex;
  gap: 20px;
  flex-wrap: wrap;
}

.compact-stats {
  flex-grow: 1;
  max-width: 500px;
}

.stats-row {
  display: flex;
  justify-content: space-around;
  gap: 0.5rem;
  margin-bottom: 0.5rem;
}

.stat-item {
  flex: 1;
  text-align: center;
  padding: 0.5rem;
  background: #f8f9fa;
  border-radius: 8px;
  min-width: 80px;
}

.stat-value {
  font-weight: bold;
  font-size: 1.1rem;
  color: #4a6fa5;
}

.stat-label {
  font-size: 0.75rem;
  color: #7f8c8d;
  white-space: nowrap;
}

.profile-tabs {
  display: flex;
  border-bottom: 1px solid #eee;
  margin-bottom: 1.5rem;
}

.profile-tabs button {
  padding: 0.75rem 1.5rem;
  background: none;
  border: none;
  cursor: pointer;
  position: relative;
  font-size: 0.9rem;
  color: #7f8c8d;
}

.profile-tabs button.active {
  color: #4a6fa5;
  font-weight: bold;
}

.profile-tabs button.active::after {
  content: '';
  position: absolute;
  bottom: -1px;
  left: 0;
  right: 0;
  height: 2px;
  background: #4a6fa5;
}

/* Адаптивность */
@media (max-width: 768px) {
  .user-profile-header {
    flex-direction: column;
    align-items: flex-start;
  }
}
  
  .stats-row {
    flex-wrap: wrap;
    justify-content: flex-start;
  }
  
  .stat-item {
    min-width: calc(33% - 0.5rem);
  }

.stat-card {
  background: #f8f9fa;
  border-radius: 10px;
  padding: 15px 20px;
  min-width: 120px;
  text-align: center;
}

.stat-value {
  font-size: 1.8rem;
  font-weight: bold;
  color: #4a6fa5;
  margin-bottom: 5px;
}

.stat-label {
  font-size: 0.9rem;
  color: #666;
  margin-bottom: 8px;
}

.rating-stars {
  color: #ccc;
  font-size: 1.2rem;
}

.rating-stars .filled {
  color: #ffc107;
}

.progress-bar {
  height: 6px;
  background: #e9ecef;
  border-radius: 3px;
  margin-top: 8px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: linear-gradient(to right, #4facfe, #00f2fe);
  border-radius: 3px;
}
  
  .edit-btn {
    background: none;
    border: 1px solid #ddd;
    border-radius: 4px;
    padding: 0.5rem 1rem;
    cursor: pointer;
    color: #555;
    transition: all 0.3s;
  }
  
  .edit-btn:hover {
    background-color: #f5f5f5;
  }
  
  .divider {
    height: 1px;
    background-color: #eee;
    margin: 1rem 0;
  }
  
  .quiz-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
  gap: 20px;
  padding: 20px;
}

.quiz-card {
  background: white;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
  cursor: pointer;
  position: relative;
}

.quiz-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 10px 20px rgba(0, 0, 0, 0.15);
}

.quiz-image-container {
  position: relative;
  width: 100%;
  height: 160px;
  overflow: hidden;
}

.quiz-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.5s ease;
}

.quiz-card:hover .quiz-image {
  transform: scale(1.05);
}

.quiz-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: linear-gradient(to top, rgba(0, 0, 0, 0.7), rgba(0, 0, 0, 0.1));
  opacity: 0.7;
  transition: opacity 0.3s ease;
}

.quiz-card:hover .quiz-overlay {
  opacity: 0.5;
}

.quiz-info {
  padding: 16px;
}

.quiz-title {
  margin: 0 0 8px 0;
  font-size: 1.1rem;
  color: #333;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.quiz-description {
  margin: 0;
  font-size: 0.9rem;
  color: #666;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.empty-message {
  text-align: center;
  color: #666;
  padding: 40px;
  font-size: 1.1rem;
}

/* Адаптивность */
@media (max-width: 768px) {
  .quiz-list {
    grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  }
}

@media (max-width: 480px) {
  .quiz-list {
    grid-template-columns: 1fr;
  }
}
  .tabs {
    display: flex;
    border-bottom: 1px solid #eee;
    margin-bottom: 1.5rem;
  }
  
  .tab-btn {
    padding: 0.75rem 1.5rem;
    background: none;
    border: none;
    cursor: pointer;
    font-size: 1rem;
    color: #666;
    position: relative;
    margin-right: 0.5rem;
  }
  
  .tab-btn.active {
    color: #000;
    font-weight: bold;
  }
  
  .tab-btn.active::after {
    content: '';
    position: absolute;
    bottom: -1px;
    left: 0;
    right: 0;
    height: 2px;
    background-color: #4CAF50;
  }
  
  .drafts {
    margin-left: auto;
  }
  
  .empty-message {
    color: #888;
    text-align: center;
    padding: 2rem;
    font-size: 1.1rem;
  }

  .friends-header {
  display: flex;
  justify-content: flex-end;
  margin-bottom: 15px;
}

.add-friend-btn {
  background: #4CAF50;
  color: white;
  border: none;
  padding: 8px 15px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
}

.friends-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
  gap: 15px;
}

.friend-card {
  display: flex;
  align-items: center;
  padding: 10px;
  border: 1px solid #e0e0e0;
  border-radius: 8px;
}

.friend-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  margin-right: 10px;
  object-fit: cover;
}

.friend-name {
  flex-grow: 1;
}

.remove-friend-btn {
  background: none;
  border: none;
  color: #ff5252;
  cursor: pointer;
  font-size: 20px;
  padding: 0 5px;
}

/* Модальное окно */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0,0,0,0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.modal-content {
  background: white;
  padding: 20px;
  border-radius: 8px;
  width: 400px;
  max-width: 90%;
}

.search-results {
  margin: 15px 0;
  max-height: 300px;
  overflow-y: auto;
}

.user-result {
  display: flex;
  align-items: center;
  padding: 10px;
  cursor: pointer;
  border-bottom: 1px solid #eee;
}

.user-result:hover {
  background: #f5f5f5;
}

.user-avatar {
  width: 30px;
  height: 30px;
  border-radius: 50%;
  margin-right: 10px;
  object-fit: cover;
}

.close-btn {
  margin-top: 15px;
  padding: 8px 15px;
  background: #f44336;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}
  </style>