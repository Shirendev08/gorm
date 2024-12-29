<template>
  <div class="home-container">
    <h1 class="welcome-title">Welcome to the Home Page</h1>
    <p class="greeting">Hello!</p>
    <button @click="logout" class="logout-button">Logout</button>
  </div>
</template>

<script setup>
import { useRouter } from 'nuxt/app'
import { useCookie } from '#app'
const router = useRouter()
const authToken = useCookie("accessToken"); // Retrieve the auth token
console.log("authewerw", authToken)
if (!authToken){
  navigateTo("/login")
}
const logout = () => {
  useCookie('authToken').value = null // Clear the auth token from cookies
  router.push("/login") // Redirect to login page after logout
};
</script>

<style scoped>
.home-container {
  max-width: 600px;
  margin: 50px auto;
  text-align: center;
  font-family: Arial, sans-serif;
}

.welcome-title {
  font-size: 2rem;
  color: #333;
}

.greeting {
  font-size: 1.5rem;
  margin: 20px 0;
  color: #555;
}

.logout-button {
  padding: 10px 20px;
  font-size: 1rem;
  color: white;
  background-color: #007bff;
  border: none;
  border-radius: 5px;
  cursor: pointer;
  transition: background-color 0.3s ease;
}

.logout-button:hover {
  background-color: #0056b3;
}
</style>
