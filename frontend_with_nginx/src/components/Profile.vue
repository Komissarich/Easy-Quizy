<template>
  <div class="profile-container">
    <div class="profile-header">
      <!-- <div class="avatar-container">
        <img :src="user.avatar || defaultAvatar" class="avatar" />
        <button class="edit-avatar-btn" @click="changeAvatar">✎</button>
      </div> -->
      <h1 class="username">{{ user.username }}</h1>
      <p class="user-email">{{ user.email }}</p>
      <!-- <button class="edit-btn" @click="editProfile">Редактировать профиль</button> -->
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
            <div v-for="quiz in userQuizzes" :key="quiz.id" class="quiz-card">
              <h3>{{ quiz.title }}</h3>
              <p>{{ quiz.description }}</p>
            </div>
          </div>
        </div>
        
        <!-- Избранное -->
        <div v-if="activeTab === 'favorites'" class="favorites-section">
          <p v-if="favorites.length === 0" class="empty-message">
            Нет добавленных в избранное тестов
          </p>
          <div v-else class="favorites-list">
            <!-- Список избранных тестов -->
          </div>
        </div>
        
        <!-- Друзья -->
        <div v-if="activeTab === 'friends'" class="friends-section">
          <div class="friends-header">
            <button class="add-friend-btn" @click="showAddFriendModal = true">
              + Добавить друга
            </button>
          </div>
          
          <div v-if="friends.length > 0" class="friends-list">
            <div v-for="friend in friends" :key="friend.id" class="friend-card">
              <img :src="friend.avatar || defaultAvatar" class="friend-avatar" />
              <span class="friend-name">{{ friend.username }}</span>
              <button class="remove-friend-btn" @click="removeFriend(friend.id)">
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

    <!-- Модальное окно добавления друга -->
    <div v-if="showAddFriendModal" class="modal-overlay" @click.self="showAddFriendModal = false">
      <div class="modal-content">
        <h2>Добавить друга</h2>
        <input 
          v-model="friendSearchQuery"
          type="text" 
          placeholder="Введите имя пользователя"
          @input="searchFriends"
        />
        
        <div v-if="searchResults.length > 0" class="search-results">
          <div 
            v-for="user in searchResults" 
            :key="user.id"
            class="user-result"
            @click="addFriend(user)"
          >
            <img :src="user.avatar || defaultAvatar" class="user-avatar" />
            <span>{{ user.username }}</span>
          </div>
        </div>
        
        <button class="close-btn" @click="showAddFriendModal = false">Закрыть</button>
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
      friends: [
        { id: 2, username: 'Мария Петрова', avatar: null },
        { id: 3, username: 'Алексей Смирнов', avatar: null }
      ],
    
        username: localStorage.getItem("username"),
        password: '',
      
        activeTab : ref('tests')
    }
  },
  methods: {
    editProfile() {
        // Логика для редактирования профиля
        console.log('Редактирование профиля')
    },
    searchFriends() {
      // Здесь будет логика поиска друзей
      if (this.friendSearchQuery.length > 2) {
        this.searchResults = [
          { id: 4, username: 'Новый друг', avatar: null }
        ];
      } else {
        this.searchResults = [];
      }
      
  },
  addFriend(user) {
      if (!this.friends.some(f => f.id === user.id)) {
        this.friends.push(user);
      }
      this.showAddFriendModal = false;
      this.friendSearchQuery = '';
      this.searchResults = [];
    },
    removeFriend(friendId) {
      this.friends = this.friends.filter(f => f.id !== friendId);
    }
  },
  setup() {
    
    const route = useRoute()
    const quiz = ref(null)
    const loading = ref(true)

    onMounted(async () => {
      try {
       
       
        let data = await axios.get(`http://localhost:8085/v1/quiz/author/${localStorage.getItem('username')}`,  {
              headers: {
                'Authorization': `Bearer ${localStorage.getItem('token')}`
              },
            })
          console.log(data.data)
      } catch (error) {
        console.error('Ошибка загрузки профиля:', error)
      } finally {
        loading.value = false
      }
    })

    return {
      quiz,
      loading
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
    margin-bottom: 1rem;
  }
  
  .username {
    font-size: 2rem;
    font-weight: bold;
    color: #333;
    margin: 0;
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