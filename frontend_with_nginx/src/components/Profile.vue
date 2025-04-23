<template>
    <div class="profile-container">
      <div class="profile-header">
        <h1 class="username" >{{ username }}</h1>
        <button class="edit-btn" @click="editProfile">Редактировать профиль</button>
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
        
        </div>
  
        <div class="tab-content">
          <div v-if="activeTab === 'tests'" class="tests-section">
            <p class="empty-message">Этот пользователь ещё не создал ни одного теста</p>
          </div>
          
          <div v-if="activeTab === 'favorites'" class="favorites-section">
            <!-- Контент для избранного -->
          </div>
          
          <div v-if="activeTab === 'drafts'" class="drafts-section">
            <!-- Контент для черновиков -->
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
        username: localStorage.getItem("username"),
        password: '',
      
        activeTab : ref('tests')
    }
  },
  methods: {
    editProfile() {
        // Логика для редактирования профиля
        console.log('Редактирование профиля')
    }
  },
  
  setup() {
    
    const route = useRoute()
    const quiz = ref(null)
    const loading = ref(true)

    onMounted(async () => {
      try {
        const token = localStorage.getItem('token');
        const data = await axios.post(
            'http://localhost:8085/v1/users/me',
            {
              token: token
            },
            {
              headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json', // Важно явно указать!
              },
            }
          );
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
  </style>