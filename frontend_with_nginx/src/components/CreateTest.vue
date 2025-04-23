<template>
    <div class="create-test-container" v-if="$route.path === '/create-test'">
      <h1>Создание нового теста</h1>
      
      <section class="cover-section">
        <h2>Обложка:</h2>
        <div class="cover-upload" @click="triggerFileInput">
          <img v-if="coverImage" :src="coverImage" class="cover-preview">
          <div v-else class="upload-placeholder">
            Нажмите, чтобы загрузить картинку теста<br>
            <span>(необязательно)</span>
          </div>
        </div>
        <input 
          type="file" 
          ref="fileInput" 
          @change="handleCoverUpload" 
          accept="image/*" 
          style="display: none"
        >
      </section>
  
      <div class="divider"></div>
  
      <section class="test-info">
        <div class="form-group">
          <h2>Название теста:</h2>
          <input 
            type="text" 
            v-model="testTitle" 
            placeholder="Введите название теста" 
            class="title-input"
          >
        </div>
  
        <div class="form-group">
          <h2>Описание теста:</h2>
          <textarea 
            v-model="testDescription" 
            placeholder="Необязательно. Можно добавить позже." 
            class="description-input"
          ></textarea>
        </div>
      </section>
  
      <div class="divider"></div>
  
    
       
  
      <button class="next-btn" @click="goToQuestions">
        Перейти к вопросам
      </button>
    </div>
   <router-view ></router-view>
  </template>
  
  <script setup>
  import { provide, ref } from 'vue'
  import { RouterView } from 'vue-router'
  import { useRouter } from 'vue-router'
  
  const router = useRouter()
  const fileInput = ref(null)
  const coverImage = ref(null)
  const testTitle = ref('')
  const testDescription = ref('')
  const testType = ref('personality')
  
  const triggerFileInput = () => {
    fileInput.value.click()
  }
  
  const handleCoverUpload = (e) => {
    const file = e.target.files[0]
    if (file) {
      const reader = new FileReader()
      reader.onload = (event) => {
        coverImage.value = event.target.result
      }
      reader.readAsDataURL(file)
    }
  }
  provide('testData', {
  title: testTitle.value,
  description: testDescription.value
})


// const uploadImage = async () => {
//   console.log(fileInput.value)
//   authorizeImgur()
//   let accessToken = getAccessToken()
//   //if (!fileInput.value || !accessToken.value) return
  
//   uploading.value = true
//   error.value = ''
  
//   const formData = new FormData()
//   formData.append('image', selectedFile.value)
  
//   try {
//     const response = await axios.post('https://api.imgur.com/3/image', formData, {
//       headers: {
//         'Authorization': `Bearer ${accessToken.value}`,
//         'Content-Type': 'multipart/form-data'
//       }
//     })
    
//     imageUrl.value = response.data.data.link
//   } catch (err) {
//     console.error('Ошибка загрузки:', err)
//     error.value = err.response?.data?.data?.error || 'Не удалось загрузить изображение'
//   } finally {
//     uploading.value = false
//   }
// }

  const goToQuestions = () => {
   // uploadImage(fileInput.value, parseTokenFromRedirect)
    const testData = {
      cover: coverImage.value,
      title: testTitle.value,
      description: testDescription.value
    }
    console.log('Данные теста:', testData)
   // router.push('/auth')
    // Переходим к созданию вопросов
    router.push({path:'/create-test/questions', query: {
      title: testTitle.value,
      description: testDescription.value,
      cover: coverImage.value
    }})
  }
  </script>
  
  <style scoped>
  .create-test-container {
    max-width: 800px;
    margin: 0 auto;
    padding: 2rem;
    font-family: Arial, sans-serif;
  }
  
  h1 {
    font-size: 2rem;
    margin-bottom: 2rem;
    color: #333;
  }
  
  h2 {
    font-size: 1.25rem;
    margin-bottom: 1rem;
    color: #444;
  }
  
  .divider {
    height: 1px;
    background-color: #eee;
    margin: 2rem 0;
  }
  
  .cover-upload {
    width: 100%;
    height: 200px;
    border: 2px dashed #ccc;
    border-radius: 8px;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    overflow: hidden;
  }
  
  .cover-preview {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
  
  .upload-placeholder {
    text-align: center;
    color: #888;
  }
  
  .upload-placeholder span {
    font-size: 0.9rem;
    color: #aaa;
  }
  
  .form-group {
    margin-bottom: 2rem;
  }
  
  .title-input, .description-input {
    width: 100%;
    padding: 0.75rem;
    border: 1px solid #ddd;
    border-radius: 4px;
    font-size: 1rem;
  }
  
  .title-input {
    font-size: 1.5rem;
    font-weight: bold;
  }
  
  .description-input {
    min-height: 100px;
    resize: vertical;
  }
  
  .type-option {
    margin-bottom: 1rem;
  }
  
  .type-option input {
    margin-right: 0.5rem;
  }
  
  .next-btn {
    display: block;
    width: 100%;
    padding: 1rem;
    background-color: #4CAF50;
    color: white;
    border: none;
    border-radius: 4px;
    font-size: 1.1rem;
    cursor: pointer;
    margin-top: 2rem;
  }
  
  .next-btn:hover {
    background-color: #45a049;
  }
  </style>